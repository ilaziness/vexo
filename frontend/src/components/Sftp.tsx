import React, { useState, useEffect } from "react";
import {
  Box,
  List,
  ListItem,
  ListItemIcon,
  ListItemText,
  ListItemButton,
  Menu,
  MenuItem,
  IconButton,
  Typography,
  CircularProgress,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  TextField,
} from "@mui/material";
import {
  Folder as FolderIcon,
  Description as FileIcon,
  Refresh as RefreshIcon,
  CloudDownload as DownloadIcon,
  UploadFile as UploadIcon,
  Delete as DeleteIcon,
  DriveFileRenameOutline as RenameIcon,
} from "@mui/icons-material";
import { SftpService, LogService } from "../../bindings/github.com/ilaziness/vexo/services";
import { parseCallServiceError } from "../func/service";

interface FileInfo {
  name: string;
  size: number;
  mode: number;
  modTime: string;
  isDir: boolean;
}

interface SftpProps {
  linkID: string;
}

const Sftp: React.FC<SftpProps> = ({ linkID }) => {
  const [currentPath, setCurrentPath] = useState<string>("/");
  const [fileList, setFileList] = useState<FileInfo[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [contextMenu, setContextMenu] = useState<{
    mouseX: number;
    mouseY: number;
    file: FileInfo | null;
  } | null>(null);
  
  // Dialog states
  const [renameDialogOpen, setRenameDialogOpen] = useState(false);
  const [newName, setNewName] = useState("");
  const [uploadDialogOpen, setUploadDialogOpen] = useState(false);
  const [selectedLocalPath, setSelectedLocalPath] = useState("");
  const [remoteFileName, setRemoteFileName] = useState("");

  // 初始化SFTP连接
  useEffect(() => {
    const initSftp = async () => {
      try {
        await SftpService.Connect(linkID);
        await refreshFileList();
      } catch (err: any) {
        setError(parseCallServiceError(err));
        LogService.Error(`Failed to initialize SFTP: ${err.message || err}`);
      }
    };

    initSftp();

    // 组件卸载时关闭SFTP连接
    return () => {
      SftpService.Close();
    };
  }, [linkID]);

  const refreshFileList = async () => {
    setLoading(true);
    setError(null);
    try {
      const files = await SftpService.ListFiles(currentPath);
      setFileList(files);
    } catch (err: any) {
      setError(parseCallServiceError(err));
      LogService.Error(`Failed to list files: ${err.message || err}`);
    } finally {
      setLoading(false);
    }
  };

  const handleItemClick = async (file: FileInfo) => {
    if (file.isDir) {
      const newPath = currentPath === "/" ? `/${file.name}` : `${currentPath}/${file.name}`;
      setCurrentPath(newPath);
      await refreshFileList();
    }
  };

  const handleContextMenu = (event: React.MouseEvent, file: FileInfo) => {
    event.preventDefault();
    setContextMenu({
      mouseX: event.clientX,
      mouseY: event.clientY,
      file,
    });
  };

  const handleCloseMenu = () => {
    setContextMenu(null);
  };

  const handleDownload = async () => {
    if (!contextMenu?.file) return;

    try {
      // 在实际应用中，这里需要一个本地文件保存对话框
      // 为了简化，我们使用一个基本的提示
      const localPath = prompt("请输入本地保存路径:", `./downloads/${contextMenu.file.name}`);
      if (!localPath) {
        LogService.Info("Download cancelled by user");
        return;
      }
      
      const remotePath = currentPath === "/" ? `/${contextMenu.file.name}` : `${currentPath}/${contextMenu.file.name}`;
      await SftpService.DownloadFile(remotePath, localPath);
      LogService.Info(`Downloaded file: ${contextMenu.file.name} to ${localPath}`);
    } catch (err: any) {
      setError(parseCallServiceError(err));
      LogService.Error(`Failed to download file: ${err.message || err}`);
    } finally {
      handleCloseMenu();
    }
  };

  const handleSelectUploadFile = () => {
    // 在实际应用中，这会打开一个文件选择对话框
    // 为了简化，我们使用一个提示来获取本地文件路径
    const localPath = prompt("请输入本地文件路径:");
    if (localPath) {
      setSelectedLocalPath(localPath);
      setRemoteFileName(contextMenu?.file?.name || localPath.split('/').pop() || 'uploaded_file');
      setUploadDialogOpen(true);
    }
    handleCloseMenu();
  };

  const handleUpload = async () => {
    if (!selectedLocalPath || !remoteFileName) return;
    
    try {
      const remotePath = currentPath === "/" ? `/${remoteFileName}` : `${currentPath}/${remoteFileName}`;
      await SftpService.UploadFile(selectedLocalPath, remotePath);
      await refreshFileList();
      LogService.Info(`Uploaded file: ${remoteFileName}`);
    } catch (err: any) {
      setError(parseCallServiceError(err));
      LogService.Error(`Failed to upload file: ${err.message || err}`);
    } finally {
      setUploadDialogOpen(false);
      setSelectedLocalPath("");
      setRemoteFileName("");
    }
  };

  const handleDelete = async () => {
    if (!contextMenu?.file) return;

    try {
      const path = currentPath === "/" ? `/${contextMenu.file.name}` : `${currentPath}/${contextMenu.file.name}`;
      await SftpService.DeleteFile(path);
      await refreshFileList();
      LogService.Info(`Deleted file: ${contextMenu.file.name}`);
    } catch (err: any) {
      setError(parseCallServiceError(err));
      LogService.Error(`Failed to delete file: ${err.message || err}`);
    } finally {
      handleCloseMenu();
    }
  };

  const handleRenameClick = () => {
    if (!contextMenu?.file) return;
    setNewName(contextMenu.file.name);
    setRenameDialogOpen(true);
    handleCloseMenu();
  };

  const handleRename = async () => {
    if (!contextMenu?.file || !newName.trim()) return;

    try {
      const oldPath = currentPath === "/" ? `/${contextMenu.file.name}` : `${currentPath}/${contextMenu.file.name}`;
      const newPath = currentPath === "/" ? `/${newName.trim()}` : `${currentPath}/${newName.trim()}`;
      await SftpService.RenameFile(oldPath, newPath);
      await refreshFileList();
      LogService.Info(`Renamed file: ${contextMenu.file.name} to ${newName.trim()}`);
    } catch (err: any) {
      setError(parseCallServiceError(err));
      LogService.Error(`Failed to rename file: ${err.message || err}`);
    } finally {
      setRenameDialogOpen(false);
      setNewName("");
    }
  };

  const navigateToParent = () => {
    if (currentPath === "/") return;
    const parentPath = currentPath.substring(0, currentPath.lastIndexOf("/"));
    setCurrentPath(parentPath || "/");
  };

  return (
    <Box sx={{ height: "100%", display: "flex", flexDirection: "column" }}>
      <Box sx={{ p: 1, borderBottom: 1, borderColor: "divider", display: "flex", alignItems: "center" }}>
        <IconButton size="small" onClick={navigateToParent} disabled={currentPath === "/"}>
          <Typography variant="body2">..</Typography>
        </IconButton>
        <Typography variant="body2" noWrap sx={{ flex: 1, ml: 1 }}>
          {currentPath}
        </Typography>
        <IconButton size="small" onClick={refreshFileList}>
          <RefreshIcon fontSize="small" />
        </IconButton>
      </Box>

      {error && (
        <Box sx={{ p: 2 }}>
          <Typography color="error">Error: {error}</Typography>
        </Box>
      )}

      {loading ? (
        <Box sx={{ display: "flex", justifyContent: "center", alignItems: "center", flex: 1 }}>
          <CircularProgress />
        </Box>
      ) : (
        <List sx={{ flex: 1, overflow: "auto" }}>
          {fileList.map((file, index) => (
            <ListItem
              key={index}
              secondaryAction={
                <IconButton edge="end" onClick={(e) => handleContextMenu(e, file)}>
                  <i className="material-icons">more_vert</i>
                </IconButton>
              }
              disablePadding
            >
              <ListItemButton onContextMenu={(e) => handleContextMenu(e, file)} onClick={() => handleItemClick(file)}>
                <ListItemIcon>
                  {file.isDir ? <FolderIcon /> : <FileIcon />}
                </ListItemIcon>
                <ListItemText
                  primary={file.name}
                  secondary={`${file.size} bytes | ${new Date(file.modTime).toLocaleString()}`}
                />
              </ListItemButton>
            </ListItem>
          ))}
        </List>
      )}

      {/* 右键菜单 */}
      <Menu
        open={contextMenu !== null}
        onClose={handleCloseMenu}
        anchorReference="anchorPosition"
        anchorPosition={
          contextMenu !== null ? { top: contextMenu.mouseY, left: contextMenu.mouseX } : undefined
        }
      >
        <MenuItem onClick={handleDownload}>
          <ListItemIcon>
            <DownloadIcon fontSize="small" />
          </ListItemIcon>
          <ListItemText>下载</ListItemText>
        </MenuItem>
        <MenuItem onClick={handleSelectUploadFile}>
          <ListItemIcon>
            <UploadIcon fontSize="small" />
          </ListItemIcon>
          <ListItemText>上传</ListItemText>
        </MenuItem>
        <MenuItem onClick={handleRenameClick}>
          <ListItemIcon>
            <RenameIcon fontSize="small" />
          </ListItemIcon>
          <ListItemText>重命名</ListItemText>
        </MenuItem>
        <MenuItem onClick={handleDelete}>
          <ListItemIcon>
            <DeleteIcon fontSize="small" />
          </ListItemIcon>
          <ListItemText>删除</ListItemText>
        </MenuItem>
      </Menu>
      
      {/* 重命名对话框 */}
      <Dialog open={renameDialogOpen} onClose={() => setRenameDialogOpen(false)}>
        <DialogTitle>重命名文件</DialogTitle>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            label="新名称"
            fullWidth
            variant="outlined"
            value={newName}
            onChange={(e) => setNewName(e.target.value)}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setRenameDialogOpen(false)}>取消</Button>
          <Button onClick={handleRename}>确定</Button>
        </DialogActions>
      </Dialog>
      
      {/* 上传对话框 */}
      <Dialog open={uploadDialogOpen} onClose={() => setUploadDialogOpen(false)}>
        <DialogTitle>确认上传</DialogTitle>
        <DialogContent>
          <Typography>确定要上传文件 "{selectedLocalPath.split('/').pop() || selectedLocalPath}" 到 "{currentPath}"?</Typography>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setUploadDialogOpen(false)}>取消</Button>
          <Button onClick={handleUpload}>上传</Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

export default Sftp;