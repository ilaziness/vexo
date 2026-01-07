import React, { useEffect } from "react";
import { Box } from "@mui/material";
import "@xterm/xterm/css/xterm.css";
import { Terminal as TerminalLib } from "@xterm/xterm";
import { FitAddon } from "@xterm/addon-fit";
import { Events } from "@wailsio/runtime";
import {
  LogService,
  SSHService,
} from "../../bindings/github.com/ilaziness/vexo/services";
import useTerminalStore from "../stores/terminal";
import { decodeBase64 } from "../func/service";

// Terminal 组件，封装 xterm.js
export default function Terminal(props: { linkID: string }) {
  const termRef = React.useRef<HTMLDivElement>(null);
  const term = React.useRef<TerminalLib>(null);
  const fit = React.useRef<FitAddon>(null);
  const resizeTimeout = React.useRef<number | null>(null);
  const sshOutputHandler = React.useRef<(event: any) => void>(null);

  // Handle output from backend
  React.useEffect(() => {
    sshOutputHandler.current = (event: any) => {
      const dataObj = event.data;
      if (dataObj.id !== props.linkID) return;
      term.current?.write(decodeBase64(dataObj.data));
    };

    // 注册事件监听器
    const unsubscribe = Events.On("sshOutput", sshOutputHandler.current);

    // 清理函数
    return () => {
      unsubscribe();
    };
  }, [props.linkID]); // 依赖 linkID 确保为每个终端创建独立的监听器

  // 使用 useCallback 确保 onData 回调函数的引用保持稳定
  const handleData = React.useCallback(
    (data: string) => {
      LogService.Debug(`Terminal input data: ${data}`);
      Events.Emit("sshInput", {
        id: props.linkID,
        data: data,
      });
    },
    [props.linkID],
  );

  const initTerminal = () => {
    if (term.current) {
      return;
    }
    LogService.Debug("Initializing terminal for link ID: " + props.linkID);
    const settings = useTerminalStore.getState();
    term.current = new TerminalLib({
      cursorBlink: true,
      cursorStyle: "block",
      fontFamily: settings.fontFamily,
      fontSize: settings.fontSize,
      lineHeight: settings.lineHeight,
      letterSpacing: 0.1,
    });
    fit.current = new FitAddon();
    term.current.loadAddon(fit.current);
    if (termRef.current) {
      term.current.open(termRef.current);
      term.current.focus();
      fit.current.fit();
      // Handle user input
      term.current?.onData(handleData);
      SSHService.Start(props.linkID)
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
    }
    term.current.refresh(0, 0);
    // Handle terminal resize
    term.current?.onResize(onResize);
    window.addEventListener("resize", () => {
      fit.current?.fit();
    });
  };

  const onResize = ({ cols, rows }) => {
    if (resizeTimeout.current) {
      clearTimeout(resizeTimeout.current);
    }
    resizeTimeout.current = window.setTimeout(() => {
      SSHService.Resize(props.linkID, cols, rows);
      resizeTimeout.current = null;
    }, 200);
  };

  useEffect(() => {
    initTerminal();

    return () => {
      LogService.Info(
        `Terminal component unmounting, closing SSH connection ${props.linkID}`,
      );
      SSHService.CloseByID(props.linkID);
      term.current?.dispose();
      window.removeEventListener("resize", () => {
        fit.current?.fit();
      });
    };
  }, [props.linkID]);

  return (
    <Box
      sx={{
        height: "100%",
        width: "100%",
        display: "flex",
        overflow: "hidden",
      }}
    >
      <div ref={termRef} style={{ flex: 1, minHeight: 0 }} />
    </Box>
  );
}
