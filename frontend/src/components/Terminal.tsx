import React, { useEffect, useState } from "react";
import { Box } from "@mui/material";
import "@xterm/xterm/css/xterm.css";
import styles from "../styles/Terminal.module.css";
import { Terminal as TerminalLib } from "@xterm/xterm";
import { FitAddon } from "@xterm/addon-fit";
import { WebglAddon } from "@xterm/addon-webgl";
import { WebLinksAddon } from "@xterm/addon-web-links";
import { Unicode11Addon } from "@xterm/addon-unicode11";
import { ImageAddon } from "@xterm/addon-image";
import { LigaturesAddon } from "@xterm/addon-ligatures";
import { ClipboardAddon } from "@xterm/addon-clipboard";
import { SearchAddon } from "@xterm/addon-search";
import { AttachAddon } from "@xterm/addon-attach";
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
export default function Terminal(props: { readonly linkID: string }) {
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

  const loadAddon = () => {
    LogService.Debug("loadAddon");
    termFit.current = new FitAddon();
    term.current?.loadAddon(termFit.current);
    term.current?.loadAddon(
      new WebLinksAddon((event, uri) => {
        event.preventDefault();
        Browser.OpenURL(uri);
      }),
    );
    termSearch.current = new SearchAddon();
    term.current?.loadAddon(termSearch.current);
    term.current?.loadAddon(new Unicode11Addon());
    term.current && (term.current.unicode.activeVersion = "11");
    term.current?.loadAddon(new ImageAddon());
    term.current?.loadAddon(new ClipboardAddon());

    if (isWebgl2Supported()) {
      const webglAddon = new WebglAddon();
      webglAddon.onContextLoss((e) => {
        console.warn("WebGL context lost. Falling back to DOM rendering.");
        webglAddon.dispose();
      });
      webglRef.current = webglAddon;
      term.current?.loadAddon(webglAddon);
    } else {
      term.current?.loadAddon(new LigaturesAddon());
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
    const config = await ConfigService.ReadConfig();
    if (!mountedRef.current) return;
    const settings = config?.Terminal || useTerminalStore.getState();
    LogService.Debug(`Terminal setting ${JSON.stringify(settings)}`);
    if (term.current) return;
    // 获取当前终端主题
    const terminalTheme = useTerminalStore.getState().getCurrentTheme();
    applyThemeVars(terminalTheme);
    term.current = new TerminalLib({
      allowProposedApi: true,
      cursorBlink: true,
      cursorStyle: "bar",
      fontFamily: settings.fontFamily,
      fontSize: settings.fontSize,
      lineHeight: settings.lineHeight,
      rightClickSelectsWord: true,
      theme: terminalTheme,
    });
    if (termRef.current) {
      loadAddon();
      term.current.open(termRef.current);
      terminalInstances.set(props.linkID, term.current);
      await sleep(50);
      termFit.current?.fit();

      const cols = term.current?.cols || 80;
      const rows = term.current?.rows || 24;
      const wsAddr = await AppService.GetWSAddr();
      const wsUrl = `ws://${wsAddr}/ws/terminal?id=${props.linkID}&cols=${cols}&rows=${rows}`;
      LogService.Debug(
        `Connecting to WebSocket at ${wsUrl} for terminal ${props.linkID}`,
      );
      const ws = new WebSocket(wsUrl);
      wsRef.current = ws;
      const attachAddon = new AttachAddon(ws);
      term.current?.loadAddon(attachAddon);

      ws.onopen = async () => {
        LogService.Debug("WebSocket connected for terminal " + props.linkID);
        setIsInitializing(false);

        if (!mountedRef.current) return;
        term.current?.focus();
        term.current?.onResize(onResize);
        sleep(1000).then(() => {
          LogService.Debug("Fitting terminal after WebSocket connection");
          termFit.current?.fit();
        });
      };

      ws.onerror = (error) => {
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
        LogService.Debug(`WebSocket closed for terminal ${props.linkID}`);
        term.current?.write(`\r\n*** SSH connection closed ***\r\n`);
      };
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
      }}
    >
      <Box
        ref={termRef}
        className={styles.terminalWrapper}
        onContextMenu={handleContextMenu}
        sx={{
          width: "100%",
          height: "100%",
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
