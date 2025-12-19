import { Box, IconButton, Tooltip, Menu, MenuItem } from "@mui/material";
import React from "react";
import SSHContainer from "./SSHContainer";
import { Tab } from "@mui/material";
import TabContext from "@mui/lab/TabContext";
import TabList from "@mui/lab/TabList";
import AddIcon from "@mui/icons-material/Add";
import TerminalIcon from "@mui/icons-material/Terminal";
import { useSSHTabsStore } from "../stores/ssh";
import { getTabIndex } from "../func/service";

export default function SSHList() {
  const { sshTabs, pushTab, delTab } = useSSHTabsStore();
  const [tabValue, setTabValue] = React.useState(sshTabs[0].index); // 当前标签
  const [menuAnchorEl, setMenuAnchorEl] = React.useState<null | HTMLElement>(
    null
  );
  const [menuTabIndex, setMenuTabIndex] = React.useState<string | null>(null);

  const handleContextMenu = (
    e: React.MouseEvent<HTMLElement>,
    index: string
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
        display: "flex",
        flexDirection: "column",
        flex: 1,
        minHeight: 0,
      }}
    >
      <TabContext value={tabValue}>
        <Box
          sx={{
            borderBottom: 1,
            borderColor: "divider",
            display: "flex",
            alignItems: "center",
            direction: "row",
          }}
        >
          <TabList onChange={tabOnchange} sx={{ minHeight: 35, height: 35 }}>
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
                    height: 35,
                    minHeight: 20,
                    fontSize: 12,
                    px: 1,
                    py: 0.5,
                    minWidth: 150,
                    width: 180,
                  },
                }}
              />
            ))}
          </TabList>
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
          <Tooltip title="新建连接">
            <IconButton size="small" onClick={addTab}>
              <AddIcon />
            </IconButton>
          </Tooltip>
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
      </TabContext>
    </Box>
  );
}
