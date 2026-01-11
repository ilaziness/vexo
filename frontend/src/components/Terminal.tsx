import React, { useEffect } from "react";
import { Box } from "@mui/material";
import "@xterm/xterm/css/xterm.css";
import "../styles/terminal.css";
import { Terminal as TerminalLib } from "@xterm/xterm";
import { FitAddon } from "@xterm/addon-fit";
import { Events } from "@wailsio/runtime";
import {
  LogService,
  SSHService,
  ConfigService,
} from "../../bindings/github.com/ilaziness/vexo/services";
import useTerminalStore from "../stores/terminal";
import { decodeBase64, encodeBase64 } from "../func/decode";

// Terminal 组件，封装 xterm.js
export default function Terminal(props: { linkID: string }) {
  const termRef = React.useRef<HTMLDivElement>(null);
  const term = React.useRef<TerminalLib>(null);
  const fit = React.useRef<FitAddon>(null);
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

  const initTerminal = async () => {
    if (term.current) {
      return;
    }
    LogService.Debug("Initializing terminal for link ID: " + props.linkID);
    const config = await ConfigService.ReadConfig();
    const settings = config?.Terminal || useTerminalStore.getState();
    term.current = new TerminalLib({
      cursorBlink: true,
      cursorStyle: "bar",
      fontFamily: settings.fontFamily,
      fontSize: settings.fontSize,
      lineHeight: settings.lineHeight,
    });
    fit.current = new FitAddon();
    term.current.loadAddon(fit.current);
    if (termRef.current) {
      term.current.open(termRef.current);

      // 使用 requestAnimationFrame 确保 DOM 完全渲染后再执行 fit
      // 这样可以避免字符间距计算错误的问题
      requestAnimationFrame(() => {
        // 双重 RAF 确保布局计算完成
        requestAnimationFrame(() => {
          fit.current?.fit();
          term.current?.focus();

          // 在 fit 之后立即刷新终端渲染
          term.current?.refresh(0, term.current.rows - 1);

          // 启动 SSH 连接
          SSHService.Start(
            props.linkID,
            term.current?.cols || 80,
            term.current?.rows || 24,
          )
            .then(() => {
              LogService.Info(
                `SSH connection started for link ID: ${props.linkID}`,
              );
            })
            .catch((err) => {
              LogService.Error(
                `Failed to start SSH connection for link ID: ${props.linkID}, error: ${err}`,
              );
              term.current?.write(`\r\nConnection error: ${err}\r\n`);
            });
        });
      });

      // Handle user input
      term.current?.onData(handleInputData);
      // Handle terminal resize
      term.current?.onResize(onResize);
      window.addEventListener("resize", () => {
        fit.current?.fit();
      });
    }
  };

  const onResize = ({ cols, rows }) => {
    LogService.Debug(`Terminal resized to ${cols}x${rows}`);
    if (resizeTimeout.current) {
      clearTimeout(resizeTimeout.current);
    }
    resizeTimeout.current = window.setTimeout(() => {
      SSHService.Resize(props.linkID, cols, rows);
      resizeTimeout.current = null;
    }, 200);
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
      term.current?.dispose();
      window.removeEventListener("resize", () => {
        fit.current?.fit();
      });
    };
  }, [props.linkID]);

  return (
    <Box
      sx={{
        width: "100%",
        height: "100%",
        display: "flex",
        flexDirection: "column",
        overflow: "hidden",
        border: "1px solid red",
      }}
    >
      <Box sx={{ width: "100%", height: "calc(100% - 20px)" }}>
        <Box ref={termRef} sx={{ width: "100%", height: "100%" }} />
      </Box>
      <Box
        sx={{
          height: 20,
          px: 0.5,
          py: 0.1,
          fontSize: 8,
          display: "flex",
          alignItems: "center",
        }}
      >
        status bar
      </Box>
    </Box>
  );
}
