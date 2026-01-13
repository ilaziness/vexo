import {
  Box,
  IconButton,
  Stack,
  Tooltip,
  useColorScheme,
  Menu,
  MenuItem,
  FormControlLabel,
  Radio,
  RadioGroup,
} from "@mui/material";
import SettingsIcon from "@mui/icons-material/Settings";
import NoteAddIcon from "@mui/icons-material/NoteAdd";
import BookmarksIcon from "@mui/icons-material/Bookmarks";
import BookmarkIcon from "@mui/icons-material/Bookmark";
import PaletteIcon from "@mui/icons-material/Palette";
import {
  BookmarkService,
  ConfigService,
} from "../../bindings/github.com/ilaziness/vexo/services";
import { useSSHTabsStore } from "../stores/ssh";
import { getTabIndex } from "../func/service";
import React, { useState } from "react";

export default function Header() {
  const { pushTab } = useSSHTabsStore();
  const { mode, setMode } = useColorScheme();

  // 颜色模式选择菜单相关状态
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const open = Boolean(anchorEl);

  const showSettingWindow = () => {
    ConfigService.ShowWindow().then(() => {});
  };

  const handleAddTab = () => {
    const number = useSSHTabsStore.getState().sshTabs.length + 1;
    const newIndex = `${getTabIndex()}`;
    pushTab({
      index: newIndex,
      name: `新建连接 ${number}`,
    });
  };

  const changeColorMode = (mode: "system" | "light" | "dark") => {
    console.log("changeColorMode", mode);
    setMode(mode);
  };

  const handleColorModeMenuOpen = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleColorModeMenuClose = () => {
    setAnchorEl(null);
  };

  const handleColorModeSelect = (mode: "system" | "light" | "dark") => {
    changeColorMode(mode);
    handleColorModeMenuClose();
  };

  const handleBookmark = () => {};

  const menuIcon = [
    {
      title: "新连接",
      icon: <NoteAddIcon />,
      clickHandler: handleAddTab,
    },
    {
      title: "书签",
      icon: <BookmarkIcon />,
      clickHandler: handleBookmark,
    },
    {
      title: "书签管理",
      icon: <BookmarksIcon />,
      clickHandler: () => {
        BookmarkService.ShowWindow().then(() => {});
      },
    },
    {
      title: "颜色模式",
      icon: <PaletteIcon />,
      clickHandler: handleColorModeMenuOpen,
    },
    {
      title: "设置",
      icon: <SettingsIcon />,
      clickHandler: showSettingWindow,
    },
  ];

  return (
    <Box
      component={"header"}
      sx={(theme) => ({
        width: "42px",
        backgroundColor: theme.palette.background.paper,
        borderRight: `1px solid ${theme.palette.divider}`,
        transition: "background-color 0.3s ease",
      })}
    >
      <Stack direction="column" spacing={1} alignItems="center" padding={0.5}>
        {menuIcon.map((item, key) => (
          <Tooltip title={item.title} key={key}>
            <IconButton size="small" onClick={item.clickHandler}>
              {item.icon}
            </IconButton>
          </Tooltip>
        ))}

        {/* 颜色模式选择菜单 */}
        <Menu
          anchorEl={anchorEl}
          open={open}
          onClose={handleColorModeMenuClose}
          onClick={(e) => e.stopPropagation()}
          slotProps={{
            paper: {
              elevation: 8,
              sx: {
                overflow: "visible",
                mt: 1,
              },
            },
          }}
        >
          <RadioGroup
            value={mode}
            onChange={(e) =>
              handleColorModeSelect(
                e.target.value as "system" | "light" | "dark",
              )
            }
          >
            <MenuItem>
              <FormControlLabel
                value="system"
                control={<Radio checked={mode === "system"} size="small" />}
                label="系统"
                sx={{ width: "100%" }}
              />
            </MenuItem>
            <MenuItem>
              <FormControlLabel
                value="light"
                control={<Radio checked={mode === "light"} size="small" />}
                label="明亮"
                sx={{ width: "100%" }}
              />
            </MenuItem>
            <MenuItem>
              <FormControlLabel
                value="dark"
                control={<Radio checked={mode === "dark"} size="small" />}
                label="暗黑"
                sx={{ width: "100%" }}
              />
            </MenuItem>
          </RadioGroup>
        </Menu>
      </Stack>
    </Box>
  );
}
