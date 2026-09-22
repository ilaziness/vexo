import React, { useEffect, useState, memo } from "react";
import { Box } from "@mui/material";
import type { Terminal as TerminalLib } from "@xterm/xterm";
import type { FitAddon } from "@xterm/addon-fit";
import type { WebglAddon } from "@xterm/addon-webgl";
import type { SearchAddon } from "@xterm/addon-search";
import { Browser } from "@wailsio/runtime";
import {
  LogService,
  SSHService,
  ConfigService,
  AppService,
} from "../../bindings/github.com/ilaziness/vexo/services";
import useTerminalStore from "../stores/terminal";
import Loading from "./Loading";
import TerminalContextMenu from "./TerminalContextMenu";
import { terminalInstances } from "../stores/terminalInstances";
import { sleep } from "../func/service";
import { useSSHTabsStore } from "../stores/ssh";
import { ConnectionStatus } from "../types/ssh";

type XtermModules = {
  Terminal: typeof import("@xterm/xterm").Terminal;
  FitAddon: typeof import("@xterm/addon-fit").FitAddon;
  WebglAddon: typeof import("@xterm/addon-webgl").WebglAddon;
  WebLinksAddon: typeof import("@xterm/addon-web-links").WebLinksAddon;
  Unicode11Addon: typeof import("@xterm/addon-unicode11").Unicode11Addon;
  ImageAddon: typeof import("@xterm/addon-image").ImageAddon;
  LigaturesAddon: typeof import("@xterm/addon-ligatures").LigaturesAddon;
  ClipboardAddon: typeof import("@xterm/addon-clipboard").ClipboardAddon;
  SearchAddon: typeof import("@xterm/addon-search").SearchAddon;
  AttachAddon: typeof import("@xterm/addon-attach").AttachAddon;
};

let xtermModulesPromise: Promise<XtermModules> | null = null;

function loadXtermModules(): Promise<XtermModules> {
  if (!xtermModulesPromise) {
    xtermModulesPromise = Promise.all([
      import("@xterm/xterm"),
      import("@xterm/addon-fit"),
      import("@xterm/addon-webgl"),
      import("@xterm/addon-web-links"),
      import("@xterm/addon-unicode11"),
      import("@xterm/addon-image"),
      import("@xterm/addon-ligatures"),
      import("@xterm/addon-clipboard"),
      import("@xterm/addon-search"),
      import("@xterm/addon-attach"),
      import("@xterm/xterm/css/xterm.css"),
    ])
      .then(
        ([
          xterm,
          fit,
          webgl,
          webLinks,
          unicode11,
          image,
          ligatures,
          clipboard,
          search,
          attach,
        ]) => ({
          Terminal: xterm.Terminal,
          FitAddon: fit.FitAddon,
          WebglAddon: webgl.WebglAddon,
          WebLinksAddon: webLinks.WebLinksAddon,
          Unicode11Addon: unicode11.Unicode11Addon,
          ImageAddon: image.ImageAddon,
          LigaturesAddon: ligatures.LigaturesAddon,
          ClipboardAddon: clipboard.ClipboardAddon,
          SearchAddon: search.SearchAddon,
          AttachAddon: attach.AttachAddon,
        }),
      )
      .catch((err) => {
        // Allow subsequent mounts to retry after a failed dynamic import.
        xtermModulesPromise = null;
        throw err;
      });
  }
  return xtermModulesPromise;
}

const isWebgl2Supported = (() => {
  let isSupported = globalThis.WebGL2RenderingContext ? undefined : false;
  return () => {
    if (isSupported === undefined) {
      const canvas = document.createElement("canvas");
      const gl = canvas.getContext("webgl2", {
        depth: false,
        antialias: false,
      });
      isSupported = gl instanceof globalThis.WebGL2RenderingContext;
    }
    return isSupported;
  };
})();

// Terminal 组件，封装 xterm.js
function Terminal(props: { readonly linkID: string }) {
  const [isInitializing, setIsInitializing] = useState(true);
  const termRef = React.useRef<HTMLDivElement>(null);
  const term = React.useRef<TerminalLib>(null);
  const termFit = React.useRef<FitAddon>(null);
  const termSearch = React.useRef<SearchAddon>(null);
  const webglRef = React.useRef<WebglAddon>(null);
  const resizeTimeout = React.useRef<number | NodeJS.Timeout | null>(null);
  const wsRef = React.useRef<WebSocket>(null);
  const [contextMenu, setContextMenu] = useState<{
    mouseX: number;
    mouseY: number;
  } | null>(null);
  const setConnectionStatus = useSSHTabsStore(
    (state) => state.setConnectionStatus,
  );

  const handleContextMenu = (event: React.MouseEvent) => {
    event.preventDefault();
    setContextMenu(
      contextMenu === null
        ? {
            mouseX: event.clientX + 2,
            mouseY: event.clientY - 6,
          }
        : null,
    );
  };

  const handleClose = () => {
    setContextMenu(null);
  };

  const applyThemeVars = (theme: any) => {
    if (termRef.current) {
      termRef.current.style.setProperty("--term-bg", theme.background);
      termRef.current.style.setProperty("--term-fg", theme.foreground);
    }
  };

  const loadAddon = (mods: XtermModules) => {
    LogService.Debug("loadAddon");
    termFit.current = new mods.FitAddon();
    term.current?.loadAddon(termFit.current);
    term.current?.loadAddon(
      new mods.WebLinksAddon((event, uri) => {
        event.preventDefault();
        Browser.OpenURL(uri);
      }),
    );
    termSearch.current = new mods.SearchAddon();
    term.current?.loadAddon(termSearch.current);
    term.current?.loadAddon(new mods.Unicode11Addon());
    term.current && (term.current.unicode.activeVersion = "11");
    term.current?.loadAddon(new mods.ImageAddon());
    term.current?.loadAddon(new mods.ClipboardAddon());

    if (isWebgl2Supported()) {
      const webglAddon = new mods.WebglAddon();
      webglAddon.onContextLoss(() => {
        console.warn("WebGL context lost. Falling back to DOM rendering.");
        webglAddon.dispose();
      });
      webglRef.current = webglAddon;
      term.current?.loadAddon(webglAddon);
    } else {
      term.current?.loadAddon(new mods.LigaturesAddon());
    }
  };

  const onResize = ({ cols, rows }) => {
    LogService.Debug(`Terminal resized to ${cols}x${rows}`);
    SSHService.Resize(props.linkID, cols, rows);
  };

  const initTerminal = async (mountedRef: { current: boolean }) => {
    if (term.current) {
      return;
    }
    LogService.Debug("Initializing terminal for link ID: " + props.linkID);
    try {
      const [config, mods] = await Promise.all([
        ConfigService.ReadConfig(),
        loadXtermModules(),
      ]);
      if (!mountedRef.current) return;
      const settings = config?.Terminal || useTerminalStore.getState();
      LogService.Debug(`Terminal setting ${JSON.stringify(settings)}`);
      if (term.current) return;
      if (!termRef.current) {
        LogService.Error(
          `Terminal DOM node missing for link ID: ${props.linkID}`,
        );
        setIsInitializing(false);
        return;
      }
      // 获取当前终端主题
      const terminalTheme = useTerminalStore.getState().getCurrentTheme();
      applyThemeVars(terminalTheme);
      term.current = new mods.Terminal({
        allowProposedApi: true,
        cursorBlink: true,
        cursorStyle: "block",
        fontFamily: settings.fontFamily,
        fontSize: settings.fontSize,
        lineHeight: settings.lineHeight,
        rightClickSelectsWord: true,
        theme: terminalTheme,
      });
      loadAddon(mods);
      term.current.open(termRef.current);
      terminalInstances.set(props.linkID, term.current);
      await sleep(50);
      if (!mountedRef.current) return;
      termFit.current?.fit();

      const cols = term.current?.cols || 80;
      const rows = term.current?.rows || 24;
      const wsAddr = await AppService.GetWSAddr();
      if (!mountedRef.current) return;
      const wsUrl = `ws://${wsAddr}/ws/terminal?id=${props.linkID}&cols=${cols}&rows=${rows}`;
      LogService.Debug(
        `Connecting to WebSocket at ${wsUrl} for terminal ${props.linkID}`,
      );
      const ws = new WebSocket(wsUrl);
      ws.binaryType = "arraybuffer";
      wsRef.current = ws;
      const attachAddon = new mods.AttachAddon(ws);
      term.current?.loadAddon(attachAddon);

      ws.onopen = async () => {
        if (!mountedRef.current) return;
        LogService.Debug("WebSocket connected for terminal " + props.linkID);
        setIsInitializing(false);
        setConnectionStatus(props.linkID, ConnectionStatus.Connected);

        term.current?.focus();
        term.current?.onResize(onResize);
        sleep(1000).then(() => {
          if (!mountedRef.current) return;
          LogService.Debug("Fitting terminal after WebSocket connection");
          termFit.current?.fit();
        });
      };

      ws.onerror = (error) => {
        if (!mountedRef.current) return;
        // WebSocket error event is an Event object, use message property if available
        const errorMessage =
          "message" in error
            ? (error as any).message
            : "Unknown WebSocket error";
        LogService.Error(
          `WebSocket error for terminal ${props.linkID}: ${errorMessage}`,
        );
        term.current?.write(`\r\n*** WebSocket error: ${errorMessage} ***\r\n`);
      };

      ws.onclose = () => {
        if (!mountedRef.current) return;
        LogService.Debug(`WebSocket closed for terminal ${props.linkID}`);
        setConnectionStatus(props.linkID, ConnectionStatus.Disconnected);
        term.current?.write(`\r\n*** SSH connection closed ***\r\n`);
      };
    } catch (err) {
      const message = err instanceof Error ? err.message : String(err);
      LogService.Error(
        `Failed to initialize terminal ${props.linkID}: ${message}`,
      );
      if (!mountedRef.current) return;
      if (term.current) {
        term.current.write(
          `\r\n*** Failed to initialize terminal: ${message} ***\r\n`,
        );
      }
      setIsInitializing(false);
      setConnectionStatus(props.linkID, ConnectionStatus.Disconnected);
    }
  };

  // 监听终端主题变化，动态更新终端主题
  useEffect(() => {
    const unsubscribe = useTerminalStore.subscribe((state, prevState) => {
      if (state.theme !== prevState.theme && term.current) {
        const newTheme = state.getCurrentTheme();
        term.current.options.theme = newTheme;
        applyThemeVars(newTheme);
        term.current.refresh(0, term.current.rows - 1);
        LogService.Debug(`Terminal theme updated to: ${state.theme}`);
      }
    });

    return unsubscribe;
  }, []);

  useEffect(() => {
    const mountedRef = { current: true };

    initTerminal(mountedRef).then(() => {
      if (!mountedRef.current) return;
    });

    const observer = new ResizeObserver(() => {
      if (resizeTimeout.current) {
        clearTimeout(resizeTimeout.current);
      }
      resizeTimeout.current = globalThis.setTimeout(() => {
        termFit.current?.fit();
      }, 100);
    });
    if (termRef.current) {
      observer.observe(termRef.current);
    }

    return () => {
      mountedRef.current = false;
      LogService.Debug(`Terminal component unmounting ${props.linkID}`);
      terminalInstances.remove(props.linkID);
      try {
        wsRef.current?.close();
        webglRef.current?.dispose();
        term.current?.dispose();
        term.current = null;
      } catch (e) {
        console.error("Error disposing terminal:", e);
      }
      if (resizeTimeout.current) {
        clearTimeout(resizeTimeout.current);
      }
      observer.disconnect();
    };
  }, [props.linkID]);

  return (
    <Box
      sx={{
        width: "100%",
        height: "100%",
        display: "flex",
      }}
    >
      <Box
        ref={termRef}
        onContextMenu={handleContextMenu}
        sx={{
          flex: 1,
          width: "100%",
          height: "100%",
          backgroundColor: "var(--term-bg)",
          paddingLeft: "3px",
        }}
      />
      <TerminalContextMenu
        contextMenu={contextMenu}
        onClose={handleClose}
        linkID={props.linkID}
      />
      {isInitializing && (
        <Box
          sx={{
            width: "100%",
            height: "100%",
            position: "absolute",
            top: 0,
            left: 0,
          }}
        >
          <Loading message="Initializing terminal..." />
        </Box>
      )}
    </Box>
  );
}

export default memo(Terminal);
