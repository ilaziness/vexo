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
  Dialog,
  AppBar,
  Toolbar,
  Typography,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  Backdrop,
} from "@mui/material";
import SettingsIcon from "@mui/icons-material/Settings";
import AddBoxIcon from "@mui/icons-material/AddBox";
import BookmarksIcon from "@mui/icons-material/Bookmarks";
import BookmarkIcon from "@mui/icons-material/Bookmark";
import ChevronRightIcon from "@mui/icons-material/ChevronRight";
import ExpandMoreIcon from "@mui/icons-material/ExpandMore";
import CloseIcon from "@mui/icons-material/Close";
import CloudUploadIcon from "@mui/icons-material/CloudUpload";
import HandymanIcon from "@mui/icons-material/Handyman";
import TerminalIcon from "@mui/icons-material/Terminal";
import AddToQueueIcon from "@mui/icons-material/AddToQueue";
import AutoAwesomeIcon from '@mui/icons-material/AutoAwesome';
import {
  AppService,
  BookmarkService,
  ConfigService,
  SSHBookmark,
  CommandService,
  LogService,
  ToolService,
} from "../../bindings/github.com/ilaziness/vexo/services";
import { UploadSync } from "../../bindings/github.com/ilaziness/vexo/services/syncservice";
import { ReadConfig } from "../../bindings/github.com/ilaziness/vexo/services/configservice";
import { useSSHTabsStore } from "../stores/ssh";
import { ConnectionStatus } from "../types/ssh";
import { useAIAssistantStore } from "../stores/aiAssistant";
import { genTabIndex, parseCallServiceError } from "../func/service";
import React, { useState, useEffect, useMemo, useCallback } from "react";
import { useMessageStore } from "../stores/message";
import ThemeSwitcher from "./ThemeSwitcher";
import { Events } from "@wailsio/runtime";
import BookmarkManager from "./Bookmark";
import Loading from "./Loading";

interface BookmarkGroup {
  name: string;
  bookmarks: SSHBookmark[];
}

export default function Header() {
  const { sshTabs, pushTab, setCurrentTab } = useSSHTabsStore();
  const { errorMessage, successMessage } = useMessageStore();
  const { toggleSidebarOpen } = useAIAssistantStore();

  const [bookmarks, setBookmarks] = useState<BookmarkGroup[]>([]);
  const [bookmarkAnchorEl, setBookmarkAnchorEl] = useState<null | HTMLElement>(
    null,
  );
  const [expandedGroups, setExpandedGroups] = useState<{
    [key: string]: boolean;
  }>({});
  const [bookmarkManageOpen, setBookmarkManageOpen] = useState(false);
  const [confirmDialogOpen, setConfirmDialogOpen] = useState(false);
  const [uploading, setUploading] = useState(false);

  const showSettingWindow = useCallback(() => {
    ConfigService.ShowWindow();
  }, []);

  const handleAIAssistant = useCallback(() => {
    toggleSidebarOpen();
  }, [toggleSidebarOpen]);

  const handleAddTab = useCallback(() => {
    const number = sshTabs.length + 1;
    const newIndex = `${genTabIndex()}`;
    pushTab({
      index: newIndex,
      name: `新建连接 ${number}`,
    });
  }, [sshTabs.length, pushTab]);

  const handleNewMainWindow = useCallback(() => {
    AppService.NewMainWindow();
  }, []);

  const handleBookmark = useCallback((event: React.MouseEvent<HTMLElement>) => {
    setBookmarkAnchorEl(event.currentTarget);
  }, []);

  const handleBackupClick = useCallback(async () => {
    try {
      const config = await ReadConfig();
      const syncConfig = config?.Sync;
      if (
        !syncConfig?.serverUrl ||
        !syncConfig?.syncId ||
        !syncConfig?.userKey
      ) {
        errorMessage("请先前往设置页面配置同步参数");
        return;
      }
      setConfirmDialogOpen(true);
    } catch (error) {
      LogService.Error("Failed to read config: " + String(error));
      errorMessage("读取配置失败");
    }
  }, [errorMessage]);

  const menuItems = useMemo(
    () => [
      { title: "新建空白标签", icon: <AddBoxIcon />, onClick: handleAddTab },
      {
        title: "新建窗口",
        icon: <AddToQueueIcon />,
        onClick: handleNewMainWindow,
      },
      { title: "书签", icon: <BookmarkIcon />, onClick: handleBookmark },
      {
        title: "书签管理",
        icon: <BookmarksIcon />,
        onClick: () => setBookmarkManageOpen(true),
      },
      {
        title: "命令面板",
        icon: <TerminalIcon />,
        onClick: () => CommandService.ShowWindow(),
      },
      {
        title: "AI助手",
        icon: <AutoAwesomeIcon />,
        onClick: handleAIAssistant,
      },
      {
        title: "上传备份",
        icon: <CloudUploadIcon />,
        onClick: handleBackupClick,
      },
      {
        title: "工具",
        icon: <HandymanIcon />,
        onClick: () => ToolService.ShowWindow(),
      },
      { title: "设置", icon: <SettingsIcon />, onClick: showSettingWindow },
    ],
    [
      handleAddTab,
      handleNewMainWindow,
      handleBookmark,
      handleBackupClick,
      showSettingWindow,
      handleAIAssistant,
    ],
  );

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
        bookmarkID: bookmark.id,
        host: bookmark.host,
        port: bookmark.port,
        user: bookmark.user,
        proxyJumpID: bookmark.proxy_jump_id,
      };
      const newIndex = genTabIndex();
      pushTab({
        index: newIndex,
        name: bookmark.title,
        sshInfo,
        connectionStatus: ConnectionStatus.Connecting,
      });
      setCurrentTab(newIndex);
    } catch (err) {
      errorMessage(parseCallServiceError(err));
    }
  };

  const handleCancelBackup = () => {
    setConfirmDialogOpen(false);
  };

  const handleConfirmBackup = () => {
    setConfirmDialogOpen(false);
    handleUpload();
  };

  const handleUpload = async () => {
    setUploading(true);
    try {
      await UploadSync();
      successMessage("数据上传成功");
    } catch (error) {
      LogService.Error("Upload failed: " + String(error));
      errorMessage("数据上传失败: " + parseCallServiceError(error));
    } finally {
      setUploading(false);
    }
  };

  return (
    <Box
      component={"header"}
      sx={(theme) => ({
        width: "42px",
        overflow: "hidden",
        backgroundColor: theme.palette.background.paper,
        borderRight: `1px solid ${theme.palette.divider}`,
      })}
    >
      <Stack direction="column" spacing={1} sx={{ alignItems: "center", padding: 0.5 }}>
        <Box sx={{ "--wails-draggable": "drag", cursor: "move" }}>
          <img
            src="/appicon.png"
            alt="Logo"
            style={{ width: "32px", height: "32px" }}
          />
        </Box>

        {menuItems.map((item, index) => (
          <Tooltip key={index} title={item.title}>
            <IconButton size="small" onClick={item.onClick}>
              {item.icon}
            </IconButton>
          </Tooltip>
        ))}

        <ThemeSwitcher />

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
        <Dialog
          fullScreen
          open={bookmarkManageOpen}
          onClose={() => setBookmarkManageOpen(false)}
        >
          <AppBar
            position="static"
            color="default"
            elevation={1}
            sx={{ "--wails-draggable": "drag" }}
          >
            <Toolbar>
              <IconButton
                edge="start"
                onClick={() => setBookmarkManageOpen(false)}
                size="large"
              >
                <CloseIcon />
              </IconButton>
              <Typography sx={{ ml: 1 }} variant="h6">
                书签管理
              </Typography>
            </Toolbar>
          </AppBar>
          <Box sx={{ height: "calc(100% - 64px)" }}>
            <BookmarkManager
              onRequestClose={() => setBookmarkManageOpen(false)}
            />
          </Box>
        </Dialog>

        <Dialog
          open={confirmDialogOpen}
          onClose={handleCancelBackup}
          maxWidth="xs"
          fullWidth
        >
          <DialogTitle>确认上传</DialogTitle>
          <DialogContent>
            <Typography>确定要上传备份吗？</Typography>
          </DialogContent>
          <DialogActions>
            <Button onClick={handleCancelBackup}>取消</Button>
            <Button variant="contained" onClick={handleConfirmBackup}>
              确定
            </Button>
          </DialogActions>
        </Dialog>

        <Backdrop
          open={uploading}
          sx={(theme) => ({
            zIndex: theme.zIndex.modal + 1,
            backgroundColor: "rgba(0, 0, 0, 0.5)",
          })}
        >
          <Loading message="正在上传备份..." size={60} />
        </Backdrop>
      </Stack>
    </Box>
  );
}
