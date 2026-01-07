import React from "react";
import { Paper, Tab, Tabs } from "@mui/material";
import {
  LogService,
  SSHService,
} from "../../bindings/github.com/ilaziness/vexo/services";
import { SSHLinkInfo } from "../types/ssh";
import Terminal from "./Terminal";
import Sftp from "./Sftp";
import ConnectionForm from "./ConnectionForm";
import { parseCallServiceError } from "../func/service";

// SSH 连接容器组件，管理连接状态和错误处理
export default function SSHContainer() {
  const [linkID, setLinkID] = React.useState<string>("");
  const [connectionError, setConnectionError] = React.useState<string>("");
  const [connecting, setConnecting] = React.useState<boolean>(false);
  const [activeTab, setActiveTab] = React.useState(0); // 0 for terminal, 1 for sftp

  const tabItems = [
    {
      label: "ssh",
      component: <Terminal linkID={linkID} />,
    },
    {
      label: "sftp",
      component: <Sftp linkID={linkID} />,
    },
  ];

  const connect = async (li: SSHLinkInfo) => {
    setConnectionError("");
    setConnecting(true);
    try {
      const ID = await SSHService.Connect(
        li.host,
        li.port || 22,
        li.user || "",
        li.password || "",
        li.key || "",
      );
      LogService.Debug(`SSH connection established with ID: ${ID}`);
      setLinkID(ID);
    } catch (err: any) {
      const msg = "Connection failed";
      LogService.Error(`${msg}: ${err.message || err}`).then(() => {});
      setConnectionError(parseCallServiceError(err));
    } finally {
      setConnecting(false);
    }
  };

  const handleTabChange = (event: React.SyntheticEvent, newValue: number) => {
    setActiveTab(newValue);
  };

  if (linkID === "") {
    return (
      <ConnectionForm
        onConnect={connect}
        error={connectionError}
        connecting={connecting}
      />
    );
  }

  return (
    <Paper
      sx={{
        height: "100%",
        overflow: "hidden",
        display: "flex",
        flexDirection: "column",
        flex: 1,
        p: 1,
        pt: 0,
      }}
    >
      <Tabs
        value={activeTab}
        onChange={handleTabChange}
        aria-label="ssh tabs"
        sx={{
          "&.MuiTabs-root": {
            minHeight: 30,
            height: 30,
          },
        }}
      >
        {tabItems.map((item, index) => (
          <Tab
            key={index}
            label={item.label}
            sx={{
              "&.MuiTab-root": {
                height: 30,
                minHeight: 30,
                fontSize: 10,
                minWidth: 80,
                width: 80,
                p: 0.5,
                borderBottom: 1,
                borderColor: "divider",
              },
            }}
          />
        ))}
      </Tabs>

      {tabItems.map((item, index) => (
        <div
          key={index}
          style={{
            height: "100%",
            display: activeTab === index ? "block" : "none",
          }}
        >
          {item.component}
        </div>
      ))}
    </Paper>
  );
}
