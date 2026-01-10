import React, { useEffect, useState } from "react";
import {
  Box,
  Button,
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
  const [currentPath, setCurrentPath] = useState<string>("/tmp"); // 临时默认值,将在useEffect中更新为实际home目录
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
  const [fileToDelete, setFileToDelete] = useState<FileInfo | null>(null);
  const [renameDialogOpen, setRenameDialogOpen] = useState(false);
  const [renamingItem, setRenamingItem] = useState<FileInfo | null>(null);
  const [renamingName, setRenamingName] = useState("");
  const [fullScreenLoading, setFullScreenLoading] = useState(false);
  const [createFileDialogOpen, setCreateFileDialogOpen] = useState(false);
  const [createDirDialogOpen, setCreateDirDialogOpen] = useState(false);
  const [newFileName, setNewFileName] = useState("");

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
    initSftp().then(() => {});

    // 组件卸载时关闭SFTP连接
    return () => {
      SftpService.Close();
    };
  }, [linkID]);

  const refreshFileList = async (path?: string) => {
    const targetPath = path || currentPath;
    setLoading(true);
    try {
      const files = await SftpService.ListFiles(
        linkID,
        targetPath,
        showHiddenFiles,
      );
      console.log(files);
      setFileList(sortFileList(files));
    } catch (err: any) {
      showMessageError(parseCallServiceError(err));
      LogService.Error(`Failed to list files: ${err.message || err}`);
      throw err; // 重新抛出错误,让调用者知道失败了
    } finally {
      setLoading(false);
    }
  };

  const handleItemClick = async (file: FileInfo) => {
    if (file.isDir) {
      const newPath =
        currentPath === "/" ? `/${file.name}` : `${currentPath}/${file.name}`;
      try {
        setFullScreenLoading(true);
        await refreshFileList(newPath);
        setCurrentPath(newPath);
      } catch (err: any) {
        // 错误已在refreshFileList中处理
      } finally {
        setFullScreenLoading(false);
      }
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

    const remotePath =
      currentPath === "/"
        ? `/${contextMenu.file.name}`
        : `${currentPath}/${contextMenu.file.name}`;
    const fileName = contextMenu.file.name;
    const isDir = contextMenu.file.isDir;

    handleCloseMenu();

    try {
      if (isDir) {
        // 下载目录
        await SftpService.DownloadDirectory(linkID, "", remotePath);
      } else {
        // 下载文件
        await SftpService.DownloadFileDialog(linkID, remotePath);
      }

      LogService.Info(`Downloaded: ${fileName}`);
    } catch (err: any) {
      showMessageError(parseCallServiceError(err));
      LogService.Error(`Failed to download: ${err.message || err}`);
    }
  };

  const handleUpload = async (type: "file" | "directory") => {
    handleCloseBlankMenu();

    try {
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
    }
  };

  const handleDelete = async () => {
    if (!fileToDelete) return;

    try {
      setFullScreenLoading(true);
      const path =
        currentPath === "/"
          ? `/${fileToDelete.name}`
          : `${currentPath}/${fileToDelete.name}`;
      await SftpService.DeleteFile(linkID, path);
      await refreshFileList();
      LogService.Info(`Deleted: ${fileToDelete.name}`);
    } catch (err: any) {
      showMessageError(parseCallServiceError(err));
      LogService.Error(`Failed to delete: ${err.message || err}`);
    } finally {
      setFullScreenLoading(false);
      setDeleteConfirmOpen(false);
      setFileToDelete(null);
    }
  };

  const handleRenameClick = () => {
    if (!contextMenu?.file) return;
    setRenamingItem(contextMenu.file);
    setRenamingName(contextMenu.file.name);
    setRenameDialogOpen(true);
    handleCloseMenu();
  };

  const handleRename = async () => {
    if (!renamingItem || !renamingName.trim()) return;

    try {
      setFullScreenLoading(true);
      setRenameDialogOpen(false);
      const oldPath =
        currentPath === "/"
          ? `/${renamingItem.name}`
          : `${currentPath}/${renamingItem.name}`;
      const newPath =
        currentPath === "/"
          ? `/${renamingName.trim()}`
          : `${currentPath}/${renamingName.trim()}`;
      if (oldPath === newPath) {
        return;
      }
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

  const handleCreateFile = () => {
    setCreateFileDialogOpen(true);
    setNewFileName("");
    handleCloseBlankMenu();
  };

  const handleCreateFileConfirm = async () => {
    if (!newFileName.trim()) {
      return;
    }

    try {
      setFullScreenLoading(true);
      setCreateFileDialogOpen(false);
      const filePath =
        currentPath === "/"
          ? `/${newFileName.trim()}`
          : `${currentPath}/${newFileName.trim()}`;
      await SftpService.CreateFile(linkID, filePath);
      await refreshFileList();
      LogService.Info(`Created file: ${newFileName.trim()}`);
    } catch (err: any) {
      showMessageError(parseCallServiceError(err));
      LogService.Error(`Failed to create file: ${err.message || err}`);
    } finally {
      setFullScreenLoading(false);
      setNewFileName("");
    }
  };

  const handleCreateDirectory = () => {
    setCreateDirDialogOpen(true);
    setNewFileName("");
    handleCloseBlankMenu();
  };

  const handleCreateDirectoryConfirm = async () => {
    if (!newFileName.trim()) {
      return;
    }

    try {
      setFullScreenLoading(true);
      setCreateDirDialogOpen(false);
      const dirPath =
        currentPath === "/"
          ? `/${newFileName.trim()}`
          : `${currentPath}/${newFileName.trim()}`;
      await SftpService.CreateDirectory(linkID, dirPath);
      await refreshFileList();
      LogService.Info(`Created directory: ${newFileName.trim()}`);
    } catch (err: any) {
      showMessageError(parseCallServiceError(err));
      LogService.Error(`Failed to create directory: ${err.message || err}`);
    } finally {
      setFullScreenLoading(false);
      setNewFileName("");
    }
  };

  const navigateToParent = async () => {
    if (currentPath === "/") return;
    const parentPath = currentPath.substring(0, currentPath.lastIndexOf("/"));
    const newPath = parentPath || "/";
    try {
      setFullScreenLoading(true);
      await refreshFileList(newPath);
      setCurrentPath(newPath);
    } catch (err: any) {
      // 错误已在refreshFileList中处理
    } finally {
      setFullScreenLoading(false);
    }
  };

  // sftpLoaded为false时显示初始加载动画
  if (!sftpLoaded) {
    return (
      <Box
        sx={{
          height: "100%",
          display: "flex",
          justifyContent: "center",
          alignItems: "center",
        }}
      >
        <CircularProgress />
      </Box>
    );
  }

  return (
    <Box
      sx={{
        height: "100%",
        display: "flex",
        flexDirection: "column",
        position: "relative",
      }}
    >
      <SftpNavbar
        currentPath={currentPath}
        onPathChange={async (path: string) => {
          try {
            setFullScreenLoading(true);
            await refreshFileList(path);
            // 只有加载成功才更新路径
            setCurrentPath(path);
          } catch (err: any) {
            // 错误已在refreshFileList中处理,这里不需要额外处理
          } finally {
            setFullScreenLoading(false);
          }
        }}
        onRefresh={async () => {
          try {
            setFullScreenLoading(true);
            await refreshFileList();
          } finally {
            setFullScreenLoading(false);
          }
        }}
        onNavigateToParent={navigateToParent}
        disableParentButton={currentPath === "/"}
        showHiddenFiles={showHiddenFiles}
        onToggleShowHidden={() => {
          const newShowHiddenFiles = !showHiddenFiles;
          setShowHiddenFiles(newShowHiddenFiles);
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
          sx={{ flex: 1, overflow: "auto", boxShadow: "none" }}
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
                    <Box sx={{ display: "flex", alignItems: "center" }}>
                      <Typography noWrap sx={{ flex: 1 }}>
                        {file.name}
                      </Typography>
                      <Typography
                        noWrap
                        variant="subtitle2"
                        color="textSecondary"
                      >
                        {file.mode}
                      </Typography>
                    </Box>
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
              setFileToDelete(contextMenu.file);
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
            确定要删除 "{fileToDelete?.name}" 吗?此操作不可撤销。
          </Typography>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setDeleteConfirmOpen(false)}>取消</Button>
          <Button onClick={handleDelete} color="error">
            删除
          </Button>
        </DialogActions>
      </Dialog>

      {/* 新建文件对话框 */}
      <Dialog
        open={createFileDialogOpen}
        onClose={() => {
          setCreateFileDialogOpen(false);
          setNewFileName("");
        }}
      >
        <DialogTitle>新建文件</DialogTitle>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            label="文件名"
            type="text"
            fullWidth
            variant="outlined"
            value={newFileName}
            onChange={(e) => setNewFileName(e.target.value)}
            onKeyDown={(e) => {
              if (e.key === "Enter") {
                handleCreateFileConfirm();
              }
            }}
          />
        </DialogContent>
        <DialogActions>
          <Button
            onClick={() => {
              setCreateFileDialogOpen(false);
              setNewFileName("");
            }}
          >
            取消
          </Button>
          <Button onClick={handleCreateFileConfirm} variant="contained">
            创建
          </Button>
        </DialogActions>
      </Dialog>

      {/* 新建文件夹对话框 */}
      <Dialog
        open={createDirDialogOpen}
        onClose={() => {
          setCreateDirDialogOpen(false);
          setNewFileName("");
        }}
      >
        <DialogTitle>新建文件夹</DialogTitle>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            label="文件夹名"
            type="text"
            fullWidth
            variant="outlined"
            value={newFileName}
            onChange={(e) => setNewFileName(e.target.value)}
            onKeyDown={(e) => {
              if (e.key === "Enter") {
                handleCreateDirectoryConfirm();
              }
            }}
          />
        </DialogContent>
        <DialogActions>
          <Button
            onClick={() => {
              setCreateDirDialogOpen(false);
              setNewFileName("");
            }}
          >
            取消
          </Button>
          <Button onClick={handleCreateDirectoryConfirm} variant="contained">
            创建
          </Button>
        </DialogActions>
      </Dialog>

      {/* 重命名对话框 */}
      <Dialog
        open={renameDialogOpen}
        onClose={() => {
          setRenameDialogOpen(false);
          setRenamingItem(null);
          setRenamingName("");
        }}
      >
        <DialogTitle>重命名</DialogTitle>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            label="新名称"
            type="text"
            fullWidth
            variant="outlined"
            value={renamingName}
            onChange={(e) => setRenamingName(e.target.value)}
            onKeyDown={(e) => {
              if (e.key === "Enter") {
                handleRename();
              }
            }}
          />
        </DialogContent>
        <DialogActions>
          <Button
            onClick={() => {
              setRenameDialogOpen(false);
              setRenamingItem(null);
              setRenamingName("");
            }}
          >
            取消
          </Button>
          <Button onClick={handleRename} variant="contained">
            重命名
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
              请稍候...
            </Typography>
          </Box>
        </Box>
      )}
    </Box>
  );
};

export default Sftp;
