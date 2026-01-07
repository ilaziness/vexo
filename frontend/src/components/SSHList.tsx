import {
  Box,
  IconButton,
  Menu,
  MenuItem,
  Tab,
  Tabs,
  Tooltip,
} from "@mui/material";
import { tabsClasses } from "@mui/material/Tabs";
import React from "react";
import SSHContainer from "./SSHContainer";
import AddIcon from "@mui/icons-material/Add";
import TerminalIcon from "@mui/icons-material/Terminal";
import { useSSHTabsStore } from "../stores/ssh";
import { getTabIndex } from "../func/service";

export default function SSHList() {
  const { sshTabs, pushTab, delTab } = useSSHTabsStore();
  const [tabValue, setTabValue] = React.useState(sshTabs[0].index); // 当前标签
  const [menuAnchorEl, setMenuAnchorEl] = React.useState<null | HTMLElement>(
    null,
  );
  const [menuTabIndex, setMenuTabIndex] = React.useState<string | null>(null);

  const handleContextMenu = (
    e: React.MouseEvent<HTMLElement>,
    index: string,
  ) => {
    e.preventDefault();
    setMenuAnchorEl(e.currentTarget);
    setMenuTabIndex(index);
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

  const handleDuplicateTab = () => {
    closeContextMenu();
  };

  const handleRefreshTab = () => {
    closeContextMenu();
  };

  const addTab = () => {
    const number = sshTabs.length + 1;
    const newIndex = `${getTabIndex()}`;
    pushTab({
      index: newIndex,
      name: `新建连接 ${number}`,
    });
    setTabValue(newIndex);
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
        minHeight: 0,
      }}
    >
      <Box
        sx={{
          borderBottom: 1,
          borderColor: "divider",
          display: "flex",
          alignItems: "center",
          direction: "row",
          width: "100%",
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
              label={item.name}
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
        <Tooltip title="新建连接">
          <IconButton size="small" onClick={addTab}>
            <AddIcon />
          </IconButton>
        </Tooltip>
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
      {sshTabs.map((item) => (
        <Box
          role="tabpanel"
          key={item.index}
          sx={{
            flex: 1,
            padding: 0,
            display: tabValue === item.index ? "flex" : "none",
            minHeight: 0,
          }}
        >
          <SSHContainer />
        </Box>
      ))}
    </Box>
  );
}
