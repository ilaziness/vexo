import React, { useEffect, useState } from "react";
import { Box } from "@mui/material";
import "@xterm/xterm/css/xterm.css";
import "../styles/terminal.css";
import { Terminal as TerminalLib } from "@xterm/xterm";
import { FitAddon } from "@xterm/addon-fit";
import { WebglAddon } from "@xterm/addon-webgl";
import { WebLinksAddon } from "@xterm/addon-web-links";
import { Unicode11Addon } from "@xterm/addon-unicode11";
import { ImageAddon } from "@xterm/addon-image";
import { LigaturesAddon } from "@xterm/addon-ligatures";
import { SearchAddon } from "@xterm/addon-search";
import { Events } from "@wailsio/runtime";
import {
  LogService,
  SSHService,
  ConfigService,
} from "../../bindings/github.com/ilaziness/vexo/services";
import useTerminalStore from "../stores/terminal";
import { decodeBase64, encodeBase64 } from "../func/decode";
import Loading from "./Loading";
import { parseCallServiceError, sleep } from "../func/service";

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
    term.current?.loadAddon(new WebLinksAddon());
    termSerach.current = new SearchAddon();
    term.current?.loadAddon(termSerach.current);
  };

  const loadAddonAfterOpen = () => {
    LogService.Debug("loadAddonAfterOpen");
    term.current?.loadAddon(new Unicode11Addon());
    term.current && (term.current.unicode.activeVersion = "11");
    term.current?.loadAddon(new ImageAddon());

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

  const initTerminal = async () => {
    if (term.current) {
      return;
    }
    LogService.Debug("Initializing terminal for link ID: " + props.linkID);
    const config = await ConfigService.ReadConfig();
    const settings = config?.Terminal || useTerminalStore.getState();
    LogService.Debug(`Terminal setting ${JSON.stringify(settings)}`);
    term.current = new TerminalLib({
      allowProposedApi: true,
      cursorBlink: true,
      cursorStyle: "bar",
      fontFamily: settings.fontFamily,
      fontSize: settings.fontSize,
      lineHeight: settings.lineHeight,
      theme: {
        background: "#1e1e1e",
      },
    });
    if (termRef.current) {
      loadAddonBeforeOpen();
      term.current.open(termRef.current);
      loadAddonAfterOpen();
      
      // 等待字体加载完成，防止因默认字体与配置字体宽度不一致导致 fit 计算出的行列数偏差
      try {
        await document.fonts.ready;
      } catch (e) {
        console.warn("Error waiting for fonts to load:", e);
      }
      await sleep(100);
      termFit.current?.fit();

      try {
        // 启动 SSH 连接
        LogService.Debug("start ssh session");
        await SSHService.Start(
          props.linkID,
          term.current?.cols || 80,
          term.current?.rows || 24,
        );
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
      setIsInitializing(false);

      // Handle user input
      term.current?.onData(handleInputData);
      // Handle terminal resize
      term.current?.onResize(onResize);

      // 再次 fit 以确保在连接建立期间如果有布局变化能及时更新，并触发 onResize 同步给后端
      // 此时 onResize 已注册，如果尺寸有变化会自动通知后端
      termFit.current?.fit();
    }
  };

  const onResize = ({ cols, rows }) => {
    LogService.Debug(`Terminal resized to ${cols}x${rows}`);
    SSHService.Resize(props.linkID, cols, rows);
  };

  useEffect(() => {
    // 注册事件监听器
    sshOutputHandler.current = (event: any) => {
      const dataObj = event.data;
      if (dataObj.id !== props.linkID) return;
      term.current?.write(decodeBase64(dataObj.data));
    };
    const unsubscribe = Events.On("sshOutput", sshOutputHandler.current);
    initTerminal();

    const observer = new ResizeObserver(() => {
      if (resizeTimeout.current) {
        clearTimeout(resizeTimeout.current);
      }
      resizeTimeout.current = window.setTimeout(() => {
        termFit.current?.fit();
      }, 200);
    });
    if (termRef.current) {
      observer.observe(termRef.current);
    }

    return () => {
      LogService.Info(
        `Terminal component unmounting, closing SSH connection ${props.linkID}`,
      );
      unsubscribe();
      SSHService.CloseByID(props.linkID)
        .then(() => {})
        .catch((err) => {
          //LogService.Error(err.message);
        });
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
        display: "flex",
        flexDirection: "column",
        overflow: "hidden",
        position: "relative",
        bgcolor: "#1e1e1e",
      }}
    >
      <Box
        ref={termRef}
        sx={{
          width: "100%",
          height: "100%",
          minWidth: "100%",
          minHeight: "100%",
        }}
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
