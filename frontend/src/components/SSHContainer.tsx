import React from "react";
import { Box, Tab, Tabs } from "@mui/material";
import {
  LogService,
  SSHService,
} from "../../bindings/github.com/ilaziness/vexo/services";
import { SSHLinkInfo } from "../types/ssh";
import Terminal from "./Terminal";
import Sftp from "./Sftp";
import ConnectionForm from "./ConnectionForm";
import { parseCallServiceError } from "../func/service";
import { useSSHTabsStore } from "../stores/ssh";

interface SSHContainerProps {
  tabIndex: string;
}

// SSH 连接容器组件，管理连接状态和错误处理
const SSHContainer: React.FC<SSHContainerProps> = ({ tabIndex }) => {
  const { setName } = useSSHTabsStore();
  const [linkID, setLinkID] = React.useState<string>("");
  const [connectionError, setConnectionError] = React.useState<string>("");
  const [connecting, setConnecting] = React.useState<boolean>(false);
  const [activeTab, setActiveTab] = React.useState(0); // 0 for terminal, 1 for sftp
  const [sftpLoaded, setSftpLoaded] = React.useState(false);
  let tabItems = [
    {
      label: "SSH",
      component: <Terminal linkID={linkID} />,
    },
    {
      label: "SFTP",
      component: <Sftp linkID={linkID} />,
    },
  ];
  const sftpIndex = 1;

  const connect = async (li: SSHLinkInfo) => {
    setConnectionError("");
    setConnecting(true);
    try {
      const ID = await SSHService.Connect(
        li.host,
        li.port,
        li.user,
        li.password || "",
        li.key || "",
      );
      LogService.Debug(`SSH connection established with ID: ${ID}`);
      setLinkID(ID);
      setName(tabIndex, `${li.user}@${li.host}:${li.port}`);
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
    if (newValue === sftpIndex && !sftpLoaded) {
      setSftpLoaded(true);
    }
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
    <Box
      sx={{
        width: "100%",
        height: "100%",
        overflow: "hidden",
        display: "flex",
        flexDirection: "column",
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
            borderBottom: 1,
            borderColor: "divider",
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
                fontSize: 12,
                minWidth: 80,
                width: 80,
                p: 0.5,
              },
            }}
          />
        ))}
      </Tabs>

      <Box sx={{ height: "calc(100% - 30px)" }}>
        {tabItems.map(
          (item, index) =>
            (index != sftpIndex || sftpLoaded) && (
              <Box
                key={index}
                sx={{
                  height: "100%",
                  display: activeTab === index ? "block" : "none",
                }}
              >
                {item.component}
              </Box>
            ),
        )}
      </Box>
    </Box>
  );
};

export default SSHContainer;
