import React, { useEffect } from "react";
import { Box, Tab, Tabs } from "@mui/material";
import {
  LogService,
  SSHService,
} from "../../bindings/github.com/ilaziness/vexo/services";
import { SSHLinkInfo } from "../types/ssh";
import Terminal from "./Terminal";
import Sftp from "./Sftp";
import ConnectionForm from "./ConnectionForm";
import Loading from "./Loading";
import StatusBar from "./StatusBar";
import { parseCallServiceError } from "../func/service";
import { useSSHTabsStore, useReloadSSHTabStore } from "../stores/ssh";

interface SSHContainerProps {
  tabIndex: string;
}

// SSH 连接容器组件，管理连接状态和错误处理
const SSHContainer: React.FC<SSHContainerProps> = ({ tabIndex }) => {
  const { setName, setSSHInfo, getByIndex } = useSSHTabsStore();
  const { reloadTab } = useReloadSSHTabStore();
  const [linkID, setLinkID] = React.useState<string>("");
  const [connectionError, setConnectionError] = React.useState<string>("");
  const [connecting, setConnecting] = React.useState<boolean>(false);
  const [activeTab, setActiveTab] = React.useState(0); // 0 for terminal, 1 for sftp
  const [sftpLoaded, setSftpLoaded] = React.useState(false);
  const [isReloading, setIsReloading] = React.useState<boolean>(false);
  const [lastSSHInfo, setLastSSHInfo] = React.useState<SSHLinkInfo | null>(
    null,
  );
  const tabInfo = getByIndex(tabIndex);
  const tabItems = [
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

  // connect ssh server
  const connect = async (li: SSHLinkInfo) => {
    setConnectionError("");
    setConnecting(true);
    setLastSSHInfo(li);
    try {
      LogService.Debug(`SSHLinkInfo ${JSON.stringify(li)}`);
      const linkID = await SSHService.Connect(
        li.host,
        li.port,
        li.user,
        li.password || "",
        li.key || "",
        li.keyPassword || "",
      );
      LogService.Debug(`SSH connection established with ID: ${linkID}`);
      setLinkID(linkID);
      setName(tabIndex, `${li.user}@${li.host}:${li.port}`);
      setSSHInfo(tabIndex, li);
    } catch (err: any) {
      const msg = "Connection failed";
      LogService.Error(`${msg}: ${err.message || err}`).then(() => {});
      setConnectionError(parseCallServiceError(err));
    } finally {
      setConnecting(false);
      setIsReloading(false);
    }
  };

  const handleTabChange = (event: React.SyntheticEvent, newValue: number) => {
    setActiveTab(newValue);
    if (newValue === sftpIndex && !sftpLoaded) {
      setSftpLoaded(true);
    }
  };

  useEffect(() => {
    const doReload = async () => {
      if (reloadTab.index === tabIndex) {
        LogService.Debug(`reload tab ${reloadTab.index} - ${tabIndex}`);
        if (linkID != "") await SSHService.CloseByID(linkID);
        // 重置状态
        setActiveTab(0);
        setSftpLoaded(false);
        // 如果有保存的连接信息，重新连接
        if (lastSSHInfo) {
          setIsReloading(true);
          await connect(lastSSHInfo);
        }
      }
    };
    doReload();
  }, [reloadTab]);

  useEffect(() => {
    if (tabInfo?.sshInfo) {
      setIsReloading(true);
      connect(tabInfo.sshInfo).then(() => {});
    }
  }, [tabIndex]);

  if (isReloading) {
    return <Loading message="Loading SSH connection..." />;
  }

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
        {/*ssh sftp*/}
        <Box sx={{ height: "calc(100% - 20px)" }}>
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
        {/*status bar*/}
        <Box
          sx={{
            height: 20,
            px: 0.5,
            py: 1,
            fontSize: 8,
            display: "flex",
            alignItems: "center",
            justifyContent: "flex-end",
            direction: "row",
          }}
        >
          <StatusBar sessionID={linkID} />
        </Box>
      </Box>
    </Box>
  );
};

export default SSHContainer;
