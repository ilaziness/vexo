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
  List,
  ListItem,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  Collapse,
} from "@mui/material";
import SettingsIcon from "@mui/icons-material/Settings";
import NoteAddIcon from "@mui/icons-material/NoteAdd";
import BookmarksIcon from "@mui/icons-material/Bookmarks";
import BookmarkIcon from "@mui/icons-material/Bookmark";
import PaletteIcon from "@mui/icons-material/Palette";
import ChevronRightIcon from "@mui/icons-material/ChevronRight";
import ExpandMoreIcon from "@mui/icons-material/ExpandMore";
import {
  BookmarkService,
  ConfigService,
} from "../../bindings/github.com/ilaziness/vexo/services";
import { useSSHTabsStore } from "../stores/ssh";
import { getTabIndex } from "../func/service";
import { SSHBookmark } from "../../bindings/github.com/ilaziness/vexo/services";
import React, { useState, useEffect } from "react";
import { useMessageStore } from "../stores/common";

interface BookmarkGroup {
  name: string;
  bookmarks: SSHBookmark[];
}

export default function Header() {
  const { pushTab, setCurrentTab } = useSSHTabsStore();
  const { mode, setMode } = useColorScheme();
  const { errorMessage } = useMessageStore();

  // 书签数据
  const [bookmarks, setBookmarks] = useState<BookmarkGroup[]>([]);
  // 书签菜单状态
  const [bookmarkAnchorEl, setBookmarkAnchorEl] = useState<null | HTMLElement>(
    null,
  );
  const [expandedGroups, setExpandedGroups] = useState<{
    [key: string]: boolean;
  }>({});

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

  useEffect(() => {
    loadBookmarks();
  }, []);

  const loadBookmarks = async () => {
    try {
      const config = await BookmarkService.ListBookmarks();
      if (config && Array.isArray(config)) {
        const validBookmarks = config.filter(
          (group): group is BookmarkGroup =>
            group !== null && group !== undefined,
        );
        setBookmarks(validBookmarks);
      } else {
        setBookmarks([]);
      }
    } catch (error) {
      console.error("Failed to load bookmarks:", error);
      errorMessage(`Failed to load bookmarks:${error}`);
    }
  };

  const handleBookmark = (event: React.MouseEvent<HTMLElement>) => {
    setBookmarkAnchorEl(event.currentTarget);
  };

  const closeBookmarkMenu = () => {
    setBookmarkAnchorEl(null);
    setExpandedGroups({});
  };

  const toggleGroup = (groupName: string) => {
    setExpandedGroups((prev) => ({
      ...prev,
      [groupName]: !prev[groupName],
    }));
  };

  const handleBookmarkSelect = async (bookmarkID: string) => {
    const bookmark = await BookmarkService.GetBookmarkByID(bookmarkID);
    if (!bookmark) {
      errorMessage("bookmark not found");
      return;
    }
    if (!bookmark.password && !bookmark.private_key) {
      errorMessage("password and private key file are empty");
      return;
    }
    const sshInfo = {
      host: bookmark.host,
      port: bookmark.port,
      user: bookmark.user,
      password: bookmark.password || undefined,
      key: bookmark.private_key || undefined,
      keyPassword: bookmark.private_key_password || undefined,
    };
    const newIndex = getTabIndex().toString();
    pushTab({
      index: newIndex,
      name: bookmark.title,
      sshInfo,
    });
    setCurrentTab(newIndex);
    closeBookmarkMenu();
  };

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

        {/* 书签菜单 */}
        <Menu
          anchorEl={bookmarkAnchorEl}
          open={Boolean(bookmarkAnchorEl)}
          onClose={closeBookmarkMenu}
          slotProps={{
            paper: {
              elevation: 8,
              sx: {
                overflow: "visible",
                mt: 1,
                minWidth: 200,
              },
            },
          }}
        >
          <List dense>
            {bookmarks.map((group) => (
              <div key={group.name}>
                <ListItem disablePadding>
                  <ListItemButton onClick={() => toggleGroup(group.name)}>
                    <ListItemIcon sx={{ minWidth: 32 }}>
                      {expandedGroups[group.name] ? (
                        <ExpandMoreIcon fontSize="small" />
                      ) : (
                        <ChevronRightIcon fontSize="small" />
                      )}
                    </ListItemIcon>
                    <ListItemText primary={group.name} />
                  </ListItemButton>
                </ListItem>
                <Collapse
                  in={expandedGroups[group.name]}
                  timeout="auto"
                  unmountOnExit
                >
                  <List component="div" disablePadding dense>
                    {group.bookmarks.map((bookmark) => (
                      <ListItem key={bookmark.id} disablePadding>
                        <ListItemButton
                          sx={{ pl: 4 }}
                          onClick={() => handleBookmarkSelect(bookmark.id)}
                        >
                          <ListItemIcon sx={{ minWidth: 24 }}>
                            <BookmarkIcon fontSize="small" />
                          </ListItemIcon>
                          <ListItemText primary={bookmark.title} />
                        </ListItemButton>
                      </ListItem>
                    ))}
                  </List>
                </Collapse>
              </div>
            ))}
          </List>
        </Menu>
      </Stack>
    </Box>
  );
}
