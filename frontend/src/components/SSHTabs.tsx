import { Box, Menu, MenuItem, Tab, Tabs } from "@mui/material";
import { tabsClasses } from "@mui/material/Tabs";
import React from "react";
import SSHContainer from "./SSHContainer.tsx";
import { Events } from "@wailsio/runtime";
import TerminalIcon from "@mui/icons-material/Terminal";
import { useSSHTabsStore, useReloadSSHTabStore } from "../stores/ssh.ts";
import { useTransferStore } from "../stores/transfer.ts";
import OpBar from "./OpBar.tsx";
import SSHTabText from "./SSHTabText.tsx";
import { genTabIndex } from "../func/service.ts";
import { ProgressData } from "../../bindings/github.com/ilaziness/vexo/services/models.ts";
import { LogService } from "../../bindings/github.com/ilaziness/vexo/services/index.ts";

const tabHeight = "40px";

export default function SSHTabs() {
  const { sshTabs, currentTab, delTab, pushTab, getByIndex, setCurrentTab } =
    useSSHTabsStore();
  const { doTabReload } = useReloadSSHTabStore();
  const { addProgress } = useTransferStore();
  const [menuAnchorEl, setMenuAnchorEl] = React.useState<null | HTMLElement>(
    null,
  );
  const [menuTabIndex, setMenuTabIndex] = React.useState<string | null>(null);

  React.useEffect(() => {
    const unsubscribe = Events.On("eventProgress", (event: any) => {
      const eventData = event.data as ProgressData;
      LogService.Debug(`eventProgress: ${JSON.stringify(event)}`);
      addProgress(eventData);
    });

    return () => {
      unsubscribe();
    };
  }, [addProgress]);

  const handleContextMenu = (
    e: React.MouseEvent<HTMLElement>,
    tabIndex: string,
  ) => {
    e.preventDefault();
    setMenuAnchorEl(e.currentTarget);
    setMenuTabIndex(tabIndex);
  };

  const closeContextMenu = () => {
    setMenuAnchorEl(null);
    setMenuTabIndex(null);
  };

  const handleCloseTab = () => {
    if (menuTabIndex === null) return;
    delTab(menuTabIndex, currentTab);
    closeContextMenu();
  };

  // duplicate ssh tab
  const handleDuplicateTab = () => {
    closeContextMenu();
    const currentTabInfo = getByIndex(menuTabIndex || "");
    const number = sshTabs.length + 1;
    const newIndex = genTabIndex();
    pushTab({
      index: newIndex,
      name: `新建连接 ${number}`,
      sshInfo: currentTabInfo ? currentTabInfo.sshInfo : undefined,
    });
    setCurrentTab(newIndex);
  };

  // refresh ssh tab
  const handleRefreshTab = () => {
    closeContextMenu();
    if (menuTabIndex) {
      doTabReload(menuTabIndex);
    }
  };

  const tabOnchange = (event: React.SyntheticEvent, newValue: string) => {
    setCurrentTab(newValue);
  };

  return (
    <Box
      sx={{
        height: "100%",
        width: "calc(100% - 42px)",
        display: "flex",
        flexDirection: "column",
        flex: 1,
      }}
    >
      <Box
        sx={{
          borderBottom: 1,
          borderColor: "divider",
          display: "flex",
          alignItems: "center",
          width: "100%",
          height: tabHeight,
          pr: 1,
          "--wails-draggable": "drag",
        }}
      >
        <Tabs
          value={currentTab}
          onChange={tabOnchange}
          aria-label="ssh tabs"
          variant="scrollable"
          scrollButtons="auto"
          sx={{
            minHeight: tabHeight,
            height: tabHeight,
            flex: 1,
            [`& .${tabsClasses.scrollButtons}`]: {
              "&.Mui-disabled": { opacity: 0.5 },
            },
            [`& .${tabsClasses.scrollButtons}:first-of-type`]: {
              borderRight: 1,
              borderColor: "#FFF",
            },
            [`& .${tabsClasses.scrollButtons}:last-of-type`]: {
              borderLeft: 1,
              borderColor: "#FFF",
            },
          }}
        >
          {sshTabs.map((item) => (
            <Tab
              key={item.index}
              value={item.index}
              label={<SSHTabText text={item.name} />}
              icon={<TerminalIcon sx={{ fontSize: 14 }} />}
              iconPosition="start"
              onContextMenu={(e) => handleContextMenu(e, item.index)}
              sx={{
                textTransform: "none",
                "&.MuiTab-root": {
                  height: tabHeight,
                  minHeight: tabHeight,
                  p: 0.5,
                  minWidth: 150,
                  width: 150,
                  borderRight: 2,
                  borderColor: "divider",
                },
              }}
            />
          ))}
        </Tabs>
        <OpBar />
        <Menu
          open={Boolean(menuAnchorEl)}
          anchorEl={menuAnchorEl}
          onClose={closeContextMenu}
          anchorOrigin={{ vertical: "bottom", horizontal: "center" }}
          transformOrigin={{ vertical: "top", horizontal: "center" }}
        >
          <MenuItem onClick={handleCloseTab}>关闭</MenuItem>
          <MenuItem onClick={handleDuplicateTab}>复制</MenuItem>
          <MenuItem onClick={handleRefreshTab}>刷新</MenuItem>
        </Menu>
      </Box>
      {/* ssh tab body */}
      <Box
        sx={{
          width: "100%",
          height: `calc(100% - ${tabHeight})`,
          position: "relative",
        }}
      >
        {sshTabs.map((item) => (
          <Box
            role="tabpanel"
            key={item.index}
            sx={{
              display: "flex",
              width: "100%",
              height: "100%",
              position: "absolute",
              top: 0,
              left: currentTab === item.index ? 0 : "-9999rem",
            }}
          >
            <SSHContainer tabIndex={item.index} />
          </Box>
        ))}
      </Box>
    </Box>
  );
}
