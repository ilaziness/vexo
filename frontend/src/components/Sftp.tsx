import React, { useEffect, useState } from "react";
import {
  Box,
  Button,
  Chip,
  CircularProgress,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  IconButton,
  LinearProgress,
  ListItemIcon,
  ListItemText,
  Menu,
  MenuItem,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  TextField,
  Typography,
} from "@mui/material";
import {
  CloudDownload as DownloadIcon,
  CreateNewFolder as CreateFolderIcon,
  Delete as DeleteIcon,
  Description as FileIcon,
  DriveFileRenameOutline as RenameIcon,
  Folder as FolderIcon,
  FolderOpen as FolderOpenIcon,
  NoteAdd as CreateFileIcon,
  UploadFile as UploadIcon,
} from "@mui/icons-material";
import MoreVertIcon from "@mui/icons-material/MoreVert";
import {
  LogService,
  SftpService,
  SSHService,
} from "../../bindings/github.com/ilaziness/vexo/services";
import { formatFileSize, parseCallServiceError } from "../func/service";
import { sortFileList } from "../func/ftp";
import { useMessageStore } from "../stores/common";
import SftpNavbar from "./SftpNavbar";
import { FileInfo } from "../func/types";

interface SftpProps {
  linkID: string;
}

const Sftp: React.FC<SftpProps> = ({ linkID }) => {
  const [sftpLoaded, setSftpLoaded] = useState(false);
  const [currentPath, setCurrentPath] = useState<string>("/tmp"); // 临时默认值，将在useEffect中更新为实际home目录
  const [fileList, setFileList] = useState<FileInfo[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [showHiddenFiles, setShowHiddenFiles] = useState<boolean>(false);
  const { errorMessage: showMessageError } = useMessageStore();
  const [contextMenu, setContextMenu] = useState<{
    mouseX: number;
    mouseY: number;
    file: FileInfo | null;
  } | null>(null);
  const [blankContextMenu, setBlankContextMenu] = useState<{
    mouseX: number;
    mouseY: number;
  } | null>(null);

  const [deleteConfirmOpen, setDeleteConfirmOpen] = useState(false);
  const [itemToDelete, setItemToDelete] = useState<FileInfo | null>(null);
  const [renamingItem, setRenamingItem] = useState<FileInfo | null>(null);
  const [renamingName, setRenamingName] = useState("");
  const [fullScreenLoading, setFullScreenLoading] = useState(false);
  const [uploadType, setUploadType] = useState<"file" | "directory" | null>(
    null,
  );

  // 初始化SFTP连接
  useEffect(() => {
    const initSftp = async () => {
      if (sftpLoaded) return;
      try {
        await SSHService.StartSftp(linkID);
        LogService.Debug("SFTP connection established");
        // 获取默认的home目录
        const homePath = await SftpService.GetWd(linkID);
        setCurrentPath(homePath);
        await refreshFileList(homePath);
        setSftpLoaded(true);
      } catch (err: any) {
        showMessageError(parseCallServiceError(err));
        LogService.Error(`Failed to initialize SFTP: ${err.message || err}`);
      }
    };

    const handleFocus = () => {
      initSftp().then(() => {});
    };

    const handleVisibilityChange = () => {
      if (!document.hidden) {
        initSftp().then(() => {});
      }
    };

    // 添加事件监听
    window.addEventListener('focus', handleFocus);
    document.addEventListener('visibilitychange', handleVisibilityChange);

    // 组件卸载时关闭SFTP连接
    return () => {
      window.removeEventListener('focus', handleFocus);
      document.removeEventListener('visibilitychange', handleVisibilityChange);
      SftpService.Close();
    };
  }, [linkID]);

  const refreshFileList = async (path?: string) => {
    const targetPath = path || currentPath;
    setLoading(true);
    try {
      const files = await SftpService.ListFiles(linkID, targetPath, showHiddenFiles);
      setFileList(sortFileList(files));
    } catch (err: any) {
      showMessageError(parseCallServiceError(err));
      LogService.Error(`Failed to list files: ${err.message || err}`);
    } finally {
      setLoading(false);
    }
  };

  const handleItemClick = async (file: FileInfo) => {
    if (file.isDir) {
      const newPath =
        currentPath === "/" ? `/${file.name}` : `${currentPath}/${file.name}`;
      setCurrentPath(newPath);
      await refreshFileList(newPath);
    }
  };

  const handleContextMenu = (event: React.MouseEvent, file: FileInfo) => {
    event.preventDefault();
    setContextMenu({
      mouseX: event.clientX,
      mouseY: event.clientY,
      file,
    });
    setBlankContextMenu(null); // Close blank context menu if open
  };

  const handleBlankContextMenu = (event: React.MouseEvent) => {
    event.preventDefault();
    setBlankContextMenu({
      mouseX: event.clientX,
      mouseY: event.clientY,
    });
    setContextMenu(null); // Close file context menu if open
  };

  const handleCloseMenu = () => {
    setContextMenu(null);
  };

  const handleCloseBlankMenu = () => {
    setBlankContextMenu(null);
  };

  const handleDownload = async () => {
    if (!contextMenu?.file) return;

    try {
      setFullScreenLoading(true);
      const remotePath =
        currentPath === "/"
          ? `/${contextMenu.file.name}`
          : `${currentPath}/${contextMenu.file.name}`;

      if (contextMenu.file.isDir) {
        // 下载目录
        await SftpService.DownloadDirectory(linkID, "", remotePath);
      } else {
        // 下载文件
        await SftpService.DownloadFileDialog(linkID, remotePath);
      }

      LogService.Info(`Downloaded: ${contextMenu.file.name}`);
    } catch (err: any) {
      showMessageError(parseCallServiceError(err));
      LogService.Error(`Failed to download: ${err.message || err}`);
    } finally {
      setFullScreenLoading(false);
      handleCloseMenu();
    }
  };

  const handleUpload = async (type: "file" | "directory") => {
    try {
      setFullScreenLoading(true);
      setUploadType(type);

      const remotePath = currentPath;

      if (type === "file") {
        await SftpService.UploadFileDialog(linkID, remotePath);
      } else {
        await SftpService.UploadDirectory(linkID, "", remotePath);
      }

      await refreshFileList();
      LogService.Info(`Upload ${type} completed`);
    } catch (err: any) {
      showMessageError(parseCallServiceError(err));
      LogService.Error(`Failed to upload ${type}: ${err.message || err}`);
    } finally {
      setFullScreenLoading(false);
      setUploadType(null);
      handleCloseBlankMenu();
    }
  };

  const handleDelete = async () => {
    if (!itemToDelete) return;

    try {
      const path =
        currentPath === "/"
          ? `/${itemToDelete.name}`
          : `${currentPath}/${itemToDelete.name}`;
      await SftpService.DeleteFile(linkID, path);
      await refreshFileList();
      LogService.Info(`Deleted: ${itemToDelete.name}`);
    } catch (err: any) {
      showMessageError(parseCallServiceError(err));
      LogService.Error(`Failed to delete: ${err.message || err}`);
    } finally {
      setDeleteConfirmOpen(false);
      setItemToDelete(null);
    }
  };

  const handleRenameClick = () => {
    if (!contextMenu?.file) return;
    setRenamingItem(contextMenu.file);
    setRenamingName(contextMenu.file.name);
    handleCloseMenu();
  };

  const handleRename = async () => {
    if (!renamingItem || !renamingName.trim()) return;

    try {
      setFullScreenLoading(true);
      const oldPath =
        currentPath === "/"
          ? `/${renamingItem.name}`
          : `${currentPath}/${renamingItem.name}`;
      const newPath =
        currentPath === "/"
          ? `/${renamingName.trim()}`
          : `${currentPath}/${renamingName.trim()}`;
      await SftpService.RenameFile(linkID, oldPath, newPath);
      await refreshFileList();
      LogService.Info(
        `Renamed: ${renamingItem.name} to ${renamingName.trim()}`,
      );
    } catch (err: any) {
      showMessageError(parseCallServiceError(err));
      LogService.Error(`Failed to rename: ${err.message || err}`);
    } finally {
      setFullScreenLoading(false);
      setRenamingItem(null);
      setRenamingName("");
    }
  };

  const handleCreateFile = async () => {
    try {
      setFullScreenLoading(true);
      const newFileName = prompt("请输入新文件名:");
      if (!newFileName) {
        LogService.Info("Create file cancelled by user");
        return;
      }

      const filePath =
        currentPath === "/"
          ? `/${newFileName}`
          : `${currentPath}/${newFileName}`;
      await SftpService.CreateFile(linkID, filePath);
      await refreshFileList();
      LogService.Info(`Created file: ${newFileName}`);
    } catch (err: any) {
      showMessageError(parseCallServiceError(err));
      LogService.Error(`Failed to create file: ${err.message || err}`);
    } finally {
      setFullScreenLoading(false);
      handleCloseBlankMenu();
    }
  };

  const handleCreateDirectory = async () => {
    try {
      setFullScreenLoading(true);
      const newDirName = prompt("请输入新文件夹名:");
      if (!newDirName) {
        LogService.Info("Create directory cancelled by user");
        return;
      }

      const dirPath =
        currentPath === "/" ? `/${newDirName}` : `${currentPath}/${newDirName}`;
      await SftpService.CreateDirectory(linkID, dirPath);
      await refreshFileList();
      LogService.Info(`Created directory: ${newDirName}`);
    } catch (err: any) {
      showMessageError(parseCallServiceError(err));
      LogService.Error(`Failed to create directory: ${err.message || err}`);
    } finally {
      setFullScreenLoading(false);
      handleCloseBlankMenu();
    }
  };

  const navigateToParent = async () => {
    if (currentPath === "/") return;
    const parentPath = currentPath.substring(0, currentPath.lastIndexOf("/"));
    const newPath = parentPath || "/";
    setCurrentPath(newPath);
    await refreshFileList(newPath);
  };

  return (
    <Box sx={{ height: "100%", display: "flex", flexDirection: "column" }}>
      <SftpNavbar
        currentPath={currentPath}
        onPathChange={async (path: string) => {
          setCurrentPath(path);
          await refreshFileList(path);
        }}
        onRefresh={async () => {
          await refreshFileList();
        }}
        onNavigateToParent={navigateToParent}
        disableParentButton={currentPath === "/"}
        showHiddenFiles={showHiddenFiles}
        onToggleShowHidden={() => {
          setShowHiddenFiles(!showHiddenFiles);
          refreshFileList().then(() => {}); // 切换显示隐藏文件后立即刷新列表
        }}
      />

      {loading ? (
        <Box
          sx={{
            display: "flex",
            justifyContent: "center",
            alignItems: "center",
            flex: 1,
          }}
        >
          <CircularProgress />
        </Box>
      ) : (
        <TableContainer
          component={Paper}
          onContextMenu={(e) => handleBlankContextMenu(e)}
          sx={{ flex: 1, overflow: "auto", boxShadow: "none", borderRadius: 0 }}
        >
          <Table stickyHeader>
            <TableHead>
              <TableRow>
                <TableCell>类型</TableCell>
                <TableCell>名称</TableCell>
                <TableCell align="center">大小</TableCell>
                <TableCell align="center">修改时间</TableCell>
                <TableCell align="center">操作</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {fileList.map((file, index) => (
                <TableRow
                  hover
                  key={index}
                  sx={{
                    cursor: file.isDir ? "pointer" : "default",
                  }}
                  onContextMenu={(e) => handleContextMenu(e, file)}
                  onDoubleClick={() => handleItemClick(file)}
                >
                  <TableCell component="th" scope="row">
                    {file.isDir ? (
                      <FolderIcon sx={{ color: "#FFB74D" }} />
                    ) : (
                      <FileIcon sx={{ color: "#81C784" }} />
                    )}
                  </TableCell>
                  <TableCell>
                    {renamingItem?.name === file.name ? (
                      <TextField
                        value={renamingName}
                        onChange={(e) => setRenamingName(e.target.value)}
                        onBlur={handleRename}
                        onKeyDown={(e) => e.key === "Enter" && handleRename()}
                        autoFocus
                        size="small"
                        sx={{ width: "100%" }}
                      />
                    ) : (
                      <Box sx={{ display: "flex", alignItems: "center" }}>
                        <Typography noWrap sx={{ flex: 1 }}>
                          {file.name}
                        </Typography>
                        {file.isDir && (
                          <Chip
                            label="目录"
                            size="small"
                            sx={{
                              mx: 1,
                              height: "20px",
                              fontSize: "0.65rem",
                            }}
                          />
                        )}
                        <Typography
                          noWrap
                          variant="subtitle2"
                          color="textSecondary"
                        >
                          {file.mode}
                        </Typography>
                      </Box>
                    )}
                  </TableCell>
                  <TableCell align="center">
                    {file.isDir ? "-" : formatFileSize(file.size)}
                  </TableCell>
                  <TableCell align="center">
                    {new Date(file.modTime).toLocaleString()}
                  </TableCell>
                  <TableCell align="center">
                    <IconButton onClick={(e) => handleContextMenu(e, file)}>
                      <MoreVertIcon />
                    </IconButton>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      )}

      {/* 文件/文件夹右键菜单 */}
      <Menu
        open={contextMenu !== null}
        onClose={handleCloseMenu}
        anchorReference="anchorPosition"
        anchorPosition={
          contextMenu !== null
            ? { top: contextMenu.mouseY, left: contextMenu.mouseX }
            : undefined
        }
      >
        <MenuItem onClick={handleDownload}>
          <ListItemIcon>
            <DownloadIcon fontSize="small" />
          </ListItemIcon>
          <ListItemText>下载</ListItemText>
        </MenuItem>
        <MenuItem
          onClick={() => {
            if (contextMenu?.file) {
              setItemToDelete(contextMenu.file);
              setDeleteConfirmOpen(true);
            }
            handleCloseMenu();
          }}
        >
          <ListItemIcon>
            <DeleteIcon fontSize="small" />
          </ListItemIcon>
          <ListItemText>删除</ListItemText>
        </MenuItem>
        <MenuItem onClick={handleRenameClick}>
          <ListItemIcon>
            <RenameIcon fontSize="small" />
          </ListItemIcon>
          <ListItemText>重命名</ListItemText>
        </MenuItem>
      </Menu>

      {/* 空白区域右键菜单 */}
      <Menu
        open={blankContextMenu !== null}
        onClose={handleCloseBlankMenu}
        anchorReference="anchorPosition"
        anchorPosition={
          blankContextMenu !== null
            ? { top: blankContextMenu.mouseY, left: blankContextMenu.mouseX }
            : undefined
        }
      >
        <MenuItem onClick={handleCreateFile}>
          <ListItemIcon>
            <CreateFileIcon fontSize="small" />
          </ListItemIcon>
          <ListItemText>新建文件</ListItemText>
        </MenuItem>
        <MenuItem onClick={handleCreateDirectory}>
          <ListItemIcon>
            <CreateFolderIcon fontSize="small" />
          </ListItemIcon>
          <ListItemText>新建文件夹</ListItemText>
        </MenuItem>
        <MenuItem onClick={() => handleUpload("file")}>
          <ListItemIcon>
            <UploadIcon fontSize="small" />
          </ListItemIcon>
          <ListItemText>上传文件</ListItemText>
        </MenuItem>
        <MenuItem onClick={() => handleUpload("directory")}>
          <ListItemIcon>
            <FolderOpenIcon fontSize="small" />
          </ListItemIcon>
          <ListItemText>上传文件夹</ListItemText>
        </MenuItem>
      </Menu>

      {/* 删除确认对话框 */}
      <Dialog
        open={deleteConfirmOpen}
        onClose={() => setDeleteConfirmOpen(false)}
      >
        <DialogTitle>确认删除</DialogTitle>
        <DialogContent>
          <Typography>
            确定要删除 "{itemToDelete?.name}" 吗？此操作不可撤销。
          </Typography>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setDeleteConfirmOpen(false)}>取消</Button>
          <Button onClick={handleDelete} color="error">
            删除
          </Button>
        </DialogActions>
      </Dialog>

      {/* 全屏加载效果 */}
      {fullScreenLoading && (
        <Box
          sx={{
            position: "absolute",
            top: 0,
            left: 0,
            width: "100%",
            height: "100%",
            backgroundColor: "rgba(0, 0, 0, 0.5)",
            display: "flex",
            justifyContent: "center",
            alignItems: "center",
            zIndex: 9999,
          }}
        >
          <Box sx={{ width: "80%", maxWidth: 400 }}>
            <LinearProgress />
            <Typography
              variant="h6"
              color="white"
              align="center"
              sx={{ mt: 2 }}
            >
              {uploadType === "file"
                ? "正在上传文件..."
                : uploadType === "directory"
                  ? "正在上传文件夹..."
                  : renamingItem
                    ? "正在重命名..."
                    : "请稍候..."}
            </Typography>
          </Box>
        </Box>
      )}
    </Box>
  );
};

export default Sftp;
