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
import NoteAddIcon from "@mui/icons-material/NoteAdd";
import BookmarksIcon from "@mui/icons-material/Bookmarks";
import BookmarkIcon from "@mui/icons-material/Bookmark";
import LocalLibraryIcon from "@mui/icons-material/LocalLibrary";
import ChevronRightIcon from "@mui/icons-material/ChevronRight";
import ExpandMoreIcon from "@mui/icons-material/ExpandMore";
import CloseIcon from "@mui/icons-material/Close";
import CloudUploadIcon from "@mui/icons-material/CloudUpload";
import HandymanIcon from "@mui/icons-material/Handyman";

import {
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
import { genTabIndex, parseCallServiceError } from "../func/service";
import React, { useState, useEffect } from "react";
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

  // 书签数据
  const [bookmarks, setBookmarks] = useState<BookmarkGroup[]>([]);
  // 书签菜单状态
  const [bookmarkAnchorEl, setBookmarkAnchorEl] = useState<null | HTMLElement>(
    null,
  );
  const [expandedGroups, setExpandedGroups] = useState<{
    [key: string]: boolean;
  }>({});
  const [bookmarkManageOpen, setBookmarkManageOpen] = useState(false);
  // 备份上传相关状态
  const [confirmDialogOpen, setConfirmDialogOpen] = useState(false);
  const [uploading, setUploading] = useState(false);

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
        bookmarkID: bookmark.id,
        host: bookmark.host,
        port: bookmark.port,
        user: bookmark.user,
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

  // 备份上传相关函数
  const handleBackupClick = async () => {
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
      <Stack direction="column" spacing={1} alignItems="center" padding={0.5}>
        <Box sx={{ "--wails-draggable": "drag", cursor: "move" }}>
          <img
            src="/appicon.png"
            alt="Logo"
            style={{ width: "32px", height: "32px" }}
          />
        </Box>

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
              setBookmarkManageOpen(true);
            }}
          >
            <BookmarksIcon />
          </IconButton>
        </Tooltip>

        <ThemeSwitcher />

        <Tooltip title="命令面板">
          <IconButton
            size="small"
            onClick={() => {
              CommandService.ShowWindow().then(() => {});
            }}
          >
            <LocalLibraryIcon />
          </IconButton>
        </Tooltip>

        <Tooltip title="上传备份">
          <IconButton size="small" onClick={handleBackupClick}>
            <CloudUploadIcon />
          </IconButton>
        </Tooltip>

        <Tooltip title="工具">
          <IconButton
            size="small"
            onClick={() => {
              ToolService.ShowWindow().then(() => {});
            }}
          >
            <HandymanIcon />
          </IconButton>
        </Tooltip>

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

        {/* 备份上传确认对话框 */}
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

        {/* 全屏 Loading */}
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
