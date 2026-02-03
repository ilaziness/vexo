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
import { Events, Browser } from "@wailsio/runtime";
import {
  LogService,
  SSHService,
  ConfigService,
} from "../../bindings/github.com/ilaziness/vexo/services";
import useTerminalStore from "../stores/terminal";
import { decodeBase64, encodeBase64 } from "../func/decode";
import Loading from "./Loading";
import { parseCallServiceError, sleep } from "../func/service";
import TerminalContextMenu from "./TerminalContextMenu";
import { terminalInstances } from "../stores/terminalInstances";

const isWebgl2Supported = (() => {
  let isSupported = window.WebGL2RenderingContext ? undefined : false;
  return () => {
    if (isSupported === undefined) {
      const canvas = document.createElement("canvas");
      const gl = canvas.getContext("webgl2", {
        depth: false,
        antialias: false,
      });
      isSupported = gl instanceof window.WebGL2RenderingContext;
    }
    return isSupported;
  };
})();

// Terminal 组件，封装 xterm.js
export default function Terminal(props: { linkID: string }) {
  const [isInitializing, setIsInitializing] = useState(true);
  const termRef = React.useRef<HTMLDivElement>(null);
  const term = React.useRef<TerminalLib>(null);
  const termFit = React.useRef<FitAddon>(null);
  const termSerach = React.useRef<SearchAddon>(null);
  const resizeTimeout = React.useRef<number | null>(null);
  const sshOutputHandler = React.useRef<(event: any) => void>(null);
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

  // 使用 useCallback 确保 onData 回调函数的引用保持稳定
  const handleInputData = React.useCallback(
    (data: string) => {
      LogService.Debug(`Terminal input data: ${data}`);
      Events.Emit("sshInput", {
        id: props.linkID,
        data: encodeBase64(data),
      });
    },
    [props.linkID],
  );

  const loadAddonBeforeOpen = () => {
    LogService.Debug("loadAddonBeforeOpen");
    termFit.current = new FitAddon();
    term.current?.loadAddon(termFit.current);
    term.current?.loadAddon(
      new WebLinksAddon((event, uri) => {
        event.preventDefault();
        Browser.OpenURL(uri);
      }),
    );
    termSerach.current = new SearchAddon();
    term.current?.loadAddon(termSerach.current);
  };

  const loadAddonAfterOpen = () => {
    LogService.Debug("loadAddonAfterOpen");
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
      loadAddonBeforeOpen();
      term.current.open(termRef.current);
      terminalInstances.set(props.linkID, term.current);
      loadAddonAfterOpen();
      if (!mountedRef.current) return;

      try {
        // 启动 SSH 连接
        LogService.Debug(
          `start ssh session ${term.current?.cols}x${term.current?.rows}`,
        );
        await SSHService.Start(
          props.linkID,
          term.current?.cols || 80,
          term.current?.rows || 24,
        );
        if (!mountedRef.current) return;
        LogService.Info(`SSH connection started for link ID: ${props.linkID}`);

        term.current?.focus();
      } catch (err) {
        LogService.Error(
          `Failed to start SSH connection for link ID: ${props.linkID}, error: ${err}`,
        );
        term.current?.write(
          `Connection error: ${parseCallServiceError(err)}\r\n`,
        );
      }

      // Handle user input
      term.current?.onData(handleInputData);
      // Handle terminal resize
      term.current?.onResize(onResize);
      termFit.current?.fit();
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
    // 注册事件监听器
    sshOutputHandler.current = (event: any) => {
      const dataObj = event.data;
      if (dataObj.id !== props.linkID) return;
      term.current?.write(decodeBase64(dataObj.data));
    };
    const unsubscribe = Events.On("sshOutput", sshOutputHandler.current);
    const unsbuscribeClose = Events.On("sshClose", (event: any) => {
      const dataObj = event.data;
      if (dataObj.id !== props.linkID) return;
      LogService.Debug(`SSH connection closed for link ID: ${props.linkID}`);
      term.current?.write(`\r\n*** SSH connection closed ***\r\n`);
    });

    initTerminal(mountedRef).then(() => {
      if (!mountedRef.current) return;
      setIsInitializing(false);
    });

    const observer = new ResizeObserver(() => {
      if (resizeTimeout.current) {
        clearTimeout(resizeTimeout.current);
      }
      resizeTimeout.current = window.setTimeout(() => {
        termFit.current?.fit();
      }, 100);
    });
    if (termRef.current) {
      observer.observe(termRef.current);
    }

    return () => {
      mountedRef.current = false;
      LogService.Debug(`Terminal component unmounting ${props.linkID}`);
      unsubscribe();
      unsbuscribeClose();
      terminalInstances.remove(props.linkID);
      try {
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
        minWidth: "100%",
        minHeight: "100%",
      }}
    >
      <Box
        ref={termRef}
        className={styles.terminalWrapper}
        onContextMenu={handleContextMenu}
        sx={{
          width: "100%",
          height: "100%",
          minWidth: "100%",
          minHeight: "100%",
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
