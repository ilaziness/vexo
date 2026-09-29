import React, { useEffect, useMemo, useRef, memo } from "react";
import { Box, Button, Tab, Tabs, Typography } from "@mui/material";
import {
  LogService,
  SSHService,
  BookmarkService,
} from "../../bindings/github.com/ilaziness/vexo/services";
import { ConnectionStatus, SSHLinkInfo } from "../types/ssh";
import Terminal from "./Terminal";
import Sftp from "./Sftp";
import ConnectionForm from "./ConnectionForm";
import Loading from "./Loading";
import { formatSSHConnectError } from "../func/service";
import { useSSHTabsStore, useReloadSSHTabStore } from "../stores/ssh";
import { SSH_STATUS_BAR_HEIGHT } from "../func/aiSidebar";
import StatusBar from "./StatusBar";

interface SSHContainerProps {
  tabIndex: string;
}

const tabHeight = "30px";
const statusBarHeight = `${SSH_STATUS_BAR_HEIGHT}px`;

// SSH 连接容器组件，管理连接状态和错误处理
const SSHTabBody: React.FC<SSHContainerProps> = ({ tabIndex }) => {
  const setName = useSSHTabsStore((state) => state.setName);
  const setSSHInfo = useSSHTabsStore((state) => state.setSSHInfo);
  const getByIndex = useSSHTabsStore((state) => state.getByIndex);
  const reloadTab = useReloadSSHTabStore((state) => state.reloadTab);
  const setTabConnectionStatus = useSSHTabsStore(
    (state) => state.setTabConnectionStatus,
  );
  const [linkID, setLinkID] = React.useState<string>("");
  const [connectionError, setConnectionError] = React.useState<string>("");
  const [connecting, setConnecting] = React.useState<boolean>(false);
  const [activeTab, setActiveTab] = React.useState(0); // 0 for terminal, 1 for sftp
  const [sftpLoaded, setSftpLoaded] = React.useState(false);
  const [isReloading, setIsReloading] = React.useState<boolean>(false);
  const [lastSSHInfo, setLastSSHInfo] = React.useState<SSHLinkInfo | null>(
    null,
  );
  const [reconnectFailed, setReconnectFailed] = React.useState(false);
  const reloadingRef = useRef(false);
  const tabInfo = useMemo(() => getByIndex(tabIndex), [tabIndex]);
  const tabItems = useMemo(
    () => [
      {
        label: "SSH",
        component: <Terminal linkID={linkID} />,
      },
      {
        label: "SFTP",
        component: <Sftp linkID={linkID} />,
      },
    ],
    [linkID],
  );
  const sftpIndex = 1;

  // connect ssh server
  const connect = async (li: SSHLinkInfo) => {
    setConnectionError("");
    setReconnectFailed(false);
    setConnecting(true);
    setLastSSHInfo({ ...li, linkID: undefined });
    try {
      LogService.Debug(`SSHLinkInfo ${JSON.stringify(li)}`);
      setTabConnectionStatus(tabIndex, ConnectionStatus.Connecting);
      let nextLinkID = "";
      if (li.bookmarkID != "" && li.bookmarkID != undefined) {
        nextLinkID = await BookmarkService.ConnectBookmarkByID(li.bookmarkID);
      } else {
        nextLinkID = await SSHService.Connect(
          li.host,
          li.port,
          li.user,
          li.password || "",
          li.key || "",
          li.keyPassword || "",
          li.proxyJumpID || "",
          !!li.useAgent,
        );
      }
      LogService.Debug(`SSH connection established with ID: ${nextLinkID}`);
      setLinkID(nextLinkID);
      setName(tabIndex, `${li.user}@${li.host}:${li.port}`);
      setSSHInfo(tabIndex, { ...li, linkID: nextLinkID });
      // Session is Exec-ready immediately; do not wait for terminal WebSocket.
      setTabConnectionStatus(tabIndex, ConnectionStatus.Connected);
    } catch (err: any) {
      const msg = formatSSHConnectError(err);
      LogService.Error(`Connection failed: ${err?.message || err}`).then(
        () => {},
      );
      // 旧会话已在刷新时关掉，失败后不能再挂到这个 linkID 上。
      setLinkID("");
      setConnectionError(msg);
      if (reloadingRef.current) {
        setReconnectFailed(true);
      }
      setTabConnectionStatus(tabIndex, ConnectionStatus.Disconnected);
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
        if (linkID != "") {
          try {
            await SSHService.CloseByID(linkID);
          } catch (e) {
            console.warn("CloseByID error during reload:", e);
          }
        }
        // 重置状态
        setActiveTab(0);
        setSftpLoaded(false);
        // 如果有保存的连接信息，重新连接
        if (lastSSHInfo) {
          reloadingRef.current = true;
          setIsReloading(true);
          try {
            await connect(lastSSHInfo);
          } finally {
            reloadingRef.current = false;
          }
        }
      }
    };
    doReload();
  }, [reloadTab]);

  useEffect(() => {
    return () => {
      if (linkID) {
        LogService.Debug(
          `SSHContainer unmounting, closing connection ${linkID}`,
        );
        SSHService.CloseByID(linkID).catch((err) => {
          console.warn("Error closing connection on unmount:", err);
        });
      }
    };
  }, [linkID]);

  useEffect(() => {
    if (tabInfo?.sshInfo) {
      setIsReloading(true);
      connect(tabInfo.sshInfo).then(() => {});
    }
  }, []);

  if (isReloading) {
    return <Loading message="Loading SSH connection..." />;
  }

  if (linkID === "") {
    if (reconnectFailed && connectionError) {
      return (
        <Box
          sx={{
            width: "100%",
            height: "100%",
            display: "flex",
            flexDirection: "column",
            alignItems: "center",
            justifyContent: "center",
            gap: 2,
            px: 3,
          }}
        >
          <Typography color="error" align="center">
            {connectionError}
          </Typography>
          <Button
            variant="contained"
            disabled={!lastSSHInfo || connecting}
            onClick={() => {
              if (!lastSSHInfo) return;
              reloadingRef.current = true;
              setIsReloading(true);
              connect(lastSSHInfo).finally(() => {
                reloadingRef.current = false;
              });
            }}
          >
            重新连接
          </Button>
        </Box>
      );
    }
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
            minHeight: tabHeight,
            height: tabHeight,
            borderBottom: 1,
            borderColor: "divider",
          },
        }}
      >
        {tabItems.map((item) => (
          <Tab
            key={item.label}
            label={item.label}
            sx={{
              "&.MuiTab-root": {
                height: tabHeight,
                minHeight: tabHeight,
                fontSize: 12,
                minWidth: 80,
                width: 80,
                p: 0.5,
              },
            }}
          />
        ))}
      </Tabs>

      {/* 终端区域 */}
      <Box
        sx={{
          height: `calc(100% - ${tabHeight} - ${statusBarHeight})`,
          width: "100%",
          display: "flex",
          flexDirection: "row",
          overflow: "hidden",
        }}
      >
        {/* SSH/SFTP 内容区域 */}
        <Box
          sx={{
            width: "100%",
            height: "100%",
            position: "relative",
          }}
        >
          {tabItems.map(
            (item, index) =>
              (index != sftpIndex || sftpLoaded) && (
                <Box
                  key={item.label}
                  sx={{
                    width: "100%",
                    height: "100%",
                    position: "absolute",
                    top: 0,
                    left: 0,
                    visibility: activeTab === index ? "visible" : "hidden",
                    pointerEvents: activeTab === index ? "auto" : "none",
                    zIndex: activeTab === index ? 1 : 0,
                  }}
                >
                  {item.component}
                </Box>
              ),
          )}
        </Box>
      </Box>

      {/*status bar*/}
      <StatusBar sessionID={linkID} height={statusBarHeight} />
    </Box>
  );
};

export default memo(SSHTabBody);
