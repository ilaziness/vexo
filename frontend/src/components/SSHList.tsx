import { Box, Menu, MenuItem, Tab, Tabs } from "@mui/material";
import { tabsClasses } from "@mui/material/Tabs";
import React from "react";
import SSHContainer from "./SSHContainer";
import TerminalIcon from "@mui/icons-material/Terminal";
import { useSSHTabsStore, useReloadSSHTabStore } from "../stores/ssh";
import OpBar from "./OpBar.tsx";
import SSHTabText from "./SSHTabText.tsx";
import { getTabIndex } from "../func/service.ts";

export default function SSHList() {
  const { sshTabs, delTab, pushTab, getByIndex } = useSSHTabsStore();
  const { doTabReload } = useReloadSSHTabStore();
  const [tabValue, setTabValue] = React.useState(sshTabs[0].index); // 当前标签
  const [menuAnchorEl, setMenuAnchorEl] = React.useState<null | HTMLElement>(
    null,
  );
  const [menuTabIndex, setMenuTabIndex] = React.useState<string | null>(null);

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
    const newTabValue = delTab(menuTabIndex, tabValue);
    setTabValue(newTabValue);

    closeContextMenu();
  };

  // duplicate ssh tab
  const handleDuplicateTab = () => {
    closeContextMenu();
    const currentTabInfo = getByIndex(menuTabIndex || "");
    const number = sshTabs.length + 1;
    pushTab({
      index: getTabIndex(),
      name: `新建连接 ${number}`,
      sshInfo: currentTabInfo ? currentTabInfo.sshInfo : undefined,
    });
  };

  // refresh ssh tab
  const handleRefreshTab = () => {
    closeContextMenu();
    if (menuTabIndex) {
      doTabReload(menuTabIndex);
    }
  };

  const tabOnchange = (event: React.SyntheticEvent, newValue: string) => {
    setTabValue(newValue);
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
          height: 45,
          pr: 1,
          "--wails-draggable": "drag",
        }}
      >
        <Tabs
          value={tabValue}
          onChange={tabOnchange}
          aria-label="ssh tabs"
          variant="scrollable"
          scrollButtons="auto"
          sx={{
            minHeight: 45,
            height: 45,
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
                "&.MuiTab-root": {
                  height: 45,
                  minHeight: 45,
                  fontSize: 14,
                  px: 1,
                  py: 0.5,
                  minWidth: 150,
                  width: 180,
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
      <Box sx={{ width: "100%", height: "calc(100% - 45px)" }}>
        {sshTabs.map((item) => (
          <Box
            role="tabpanel"
            key={item.index}
            sx={{
              display: tabValue === item.index ? "flex" : "none",
              width: "100%",
              height: "100%",
            }}
          >
            <SSHContainer tabIndex={item.index} />
          </Box>
        ))}
      </Box>
    </Box>
  );
}
