import {
  Box,
  IconButton,
  Stack,
  Tooltip,
  Menu,
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
import ChevronRightIcon from "@mui/icons-material/ChevronRight";
import ExpandMoreIcon from "@mui/icons-material/ExpandMore";
import {
  BookmarkService,
  ConfigService,
} from "../../bindings/github.com/ilaziness/vexo/services";
import { useSSHTabsStore } from "../stores/ssh";
import { genTabIndex, parseCallServiceError } from "../func/service";
import { SSHBookmark } from "../../bindings/github.com/ilaziness/vexo/services";
import React, { useState, useEffect } from "react";
import { useMessageStore } from "../stores/message";
import ThemeSwitcher from "./ThemeSwitcher";
import { Events } from "@wailsio/runtime";

interface BookmarkGroup {
  name: string;
  bookmarks: SSHBookmark[];
}

export default function Header() {
  const { sshTabs, pushTab, setCurrentTab } = useSSHTabsStore();
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

  const showSettingWindow = () => {
    ConfigService.ShowWindow().then(() => {});
  };

  const handleAddTab = () => {
    const number = sshTabs.length + 1;
    const newIndex = `${genTabIndex()}`;
    pushTab({
      index: newIndex,
      name: `新建连接 ${number}`,
    });
  };

  useEffect(() => {
    loadBookmarks();
    const unsubscribe = Events.On("eventBookmarkUpdate", (e) => {
      loadBookmarks();
    });

    return () => {
      unsubscribe();
    };
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
    closeBookmarkMenu();
    try {
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
      const newIndex = genTabIndex();
      pushTab({
        index: newIndex,
        name: bookmark.title,
        sshInfo,
      });
      setCurrentTab(newIndex);
    } catch (err) {
      errorMessage(parseCallServiceError(err));
    }
  };

  return (
    <Box
      component={"header"}
      sx={(theme) => ({
        width: "42px",
        backgroundColor: theme.palette.background.paper,
        borderRight: `1px solid ${theme.palette.divider}`,
      })}
    >
      <Stack direction="column" spacing={1} alignItems="center" padding={0.5}>
        <Tooltip title="新连接">
          <IconButton size="small" onClick={handleAddTab}>
            <NoteAddIcon />
          </IconButton>
        </Tooltip>

        <Tooltip title="书签">
          <IconButton size="small" onClick={handleBookmark}>
            <BookmarkIcon />
          </IconButton>
        </Tooltip>

        <Tooltip title="书签管理">
          <IconButton
            size="small"
            onClick={() => {
              BookmarkService.ShowWindow().then(() => {});
            }}
          >
            <BookmarksIcon />
          </IconButton>
        </Tooltip>

        <ThemeSwitcher />

        <Tooltip title="设置">
          <IconButton size="small" onClick={showSettingWindow}>
            <SettingsIcon />
          </IconButton>
        </Tooltip>

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
