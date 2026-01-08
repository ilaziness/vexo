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
    });
    fit.current = new FitAddon();
    term.current.loadAddon(fit.current);
    if (termRef.current) {
      term.current.open(termRef.current);
      term.current.focus();
      fit.current.fit();
      // Handle user input
      term.current?.onData(handleInputData);
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
    }
    term.current.refresh(0, 0);
    // Handle terminal resize
    term.current?.onResize(onResize);
    window.addEventListener("resize", () => {
      fit.current?.fit();
    });
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
          LogService.Error(err.message);
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
        height: "100%",
        width: "100%",
        display: "flex",
        overflow: "hidden",
      }}
    >
      <Box
        ref={termRef}
        sx={{ width: "100%", height: "100%", border: "1px solid red" }}
      />
    </Box>
  );
}
