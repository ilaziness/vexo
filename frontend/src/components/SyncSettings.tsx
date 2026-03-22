import React, { useState, useEffect } from "react";
import {
  Box,
  Typography,
  Paper,
  Stack,
  TextField,
  Button,
  Divider,
  List,
  ListItem,
  ListItemText,
  ListItemButton,
  Chip,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  CircularProgress,
  LinearProgress,
  Alert,
} from "@mui/material";
import {
  CloudUpload,
  CloudDownload,
  Refresh,
  Save,
  Cancel,
  Delete as DeleteIcon,
} from "@mui/icons-material";
import FormRow from "./FormRow";
import { useMessageStore } from "../stores/message";
import {
  UploadSync,
  DownloadSync,
  ListSyncVersions,
  DeleteSyncVersion,
  HealthCheck,
  GetSyncProgress,
  CancelSync,
} from "../../bindings/github.com/ilaziness/vexo/services/syncservice";
import {
  Config,
  SyncProgress,
  LogService,
} from "../../bindings/github.com/ilaziness/vexo/services";
import { parseCallServiceError } from "../func/service";

// SyncConfig 类型从 Config.Sync 推断
type SyncConfig = Config["Sync"];

interface SyncSettingsProps {
  syncConfig: SyncConfig;
  onChange: (config: SyncConfig) => void;
}

interface VersionInfo {
  version_number: number;
  file_size: number;
  created_at: string;
}

const stageLabels: Record<string, string> = {
  preparing: "准备中...",
  packing: "打包数据中...",
  encrypting: "加密数据中...",
  uploading: "上传数据中...",
  downloading: "下载数据中...",
  decrypting: "解密数据中...",
  unpacking: "解压数据中...",
  retrying: "重试中...",
  completed: "完成",
  error: "错误",
};

const SyncSettings: React.FC<SyncSettingsProps> = ({
  syncConfig,
  onChange,
}) => {
  const [loading, setLoading] = useState(false);
  const [versions, setVersions] = useState<VersionInfo[]>([]);
  const [healthStatus, setHealthStatus] = useState<string>("");
  const [restoreDialogOpen, setRestoreDialogOpen] = useState(false);
  const [selectedVersion, setSelectedVersion] = useState<number | null>(null);
  const [progress, setProgress] = useState<SyncProgress | null>(null);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [versionToDelete, setVersionToDelete] = useState<number | null>(null);
  // 恢复弹框无限滚动状态
  const [restoreVersions, setRestoreVersions] = useState<VersionInfo[]>([]);
  const [restorePage, setRestorePage] = useState(0);
  const [restoreHasMore, setRestoreHasMore] = useState(true);
  const [restoreLoading, setRestoreLoading] = useState(false);
  // 重启提示弹框状态
  const [restartDialogOpen, setRestartDialogOpen] = useState(false);
  const [checkingHealth, setCheckingHealth] = useState(false);
  const { successMessage, errorMessage } = useMessageStore();

  // 本地状态用于表单编辑
  const [localConfig, setLocalConfig] = useState<SyncConfig>({
    serverUrl: syncConfig.serverUrl || "",
    syncId: syncConfig.syncId || "",
    userKey: syncConfig.userKey || "",
  });

  // 当外部配置变化时更新本地状态
  useEffect(() => {
    setLocalConfig({
      serverUrl: syncConfig.serverUrl || "",
      syncId: syncConfig.syncId || "",
      userKey: syncConfig.userKey || "",
    });
  }, [syncConfig]);

  const isConfigured =
    localConfig.serverUrl && localConfig.syncId && localConfig.userKey;

  useEffect(() => {
    if (isConfigured) {
      checkHealth();
      loadVersions();
    }
  }, [isConfigured]);

  // 进度轮询
  useEffect(() => {
    if (!loading) {
      setProgress(null);
      return;
    }

    const interval = setInterval(async () => {
      try {
        const p = await GetSyncProgress();
        setProgress(p);

        // 如果已完成或出错，停止轮询
        if (p.isCompleted || p.error) {
          clearInterval(interval);
        }
      } catch (error) {
        LogService.Warn("Failed to get progress: " + String(error));
      }
    }, 500);

    return () => clearInterval(interval);
  }, [loading]);

  const checkHealth = async () => {
    setCheckingHealth(true);
    try {
      await HealthCheck();
      setHealthStatus("connected");
    } catch (error) {
      LogService.Error("Health check failed: " + String(error));
      setHealthStatus("disconnected");
    } finally {
      setCheckingHealth(false);
    }
  };

  const loadVersions = async () => {
    if (!isConfigured) return;
    try {
      // 设置页面只加载前10个版本
      const result = await ListSyncVersions(10, 0);
      setVersions(result || []);
    } catch (error) {
      LogService.Error("Failed to load versions: " + String(error));
      errorMessage("加载历史版本失败: " + parseCallServiceError(error));
    }
  };

  const handleSaveConfig = () => {
    // 保存时添加协议前缀
    const configToSave: SyncConfig = {
      ...localConfig,
      serverUrl: localConfig.serverUrl,
    };
    onChange(configToSave);
    successMessage("同步配置已保存");
    checkHealth();
  };

  const handleUpload = async () => {
    setLoading(true);
    setProgress(null);
    try {
      await UploadSync();
      successMessage("数据上传成功");
      await loadVersions();
    } catch (error) {
      LogService.Error("Upload failed: " + String(error));
      errorMessage("数据上传失败: " + parseCallServiceError(error));
    } finally {
      setLoading(false);
    }
  };

  const handleDownload = async () => {
    if (!selectedVersion) {
      errorMessage("请选择要恢复的版本");
      return;
    }
    setLoading(true);
    setProgress(null);
    try {
      await DownloadSync(selectedVersion);
      // 恢复成功，显示重启提示
      setRestartDialogOpen(true);
    } catch (error) {
      LogService.Error("Download failed: " + String(error));
      // 恢复失败也需要重启（数据库已关闭）
      setRestartDialogOpen(true);
    } finally {
      setLoading(false);
      setRestoreDialogOpen(false);
    }
  };

  const handleRestartDialogClose = () => {
    setRestartDialogOpen(false);
  };

  // 加载恢复弹框的版本列表（支持分页）
  const loadRestoreVersions = async (page: number) => {
    if (!isConfigured || restoreLoading) return;
    setRestoreLoading(true);
    try {
      const limit = 50;
      const offset = page * limit;
      const result = await ListSyncVersions(limit, offset);
      const newVersions = result || [];
      if (page === 0) {
        setRestoreVersions(newVersions);
      } else {
        setRestoreVersions((prev) => [...prev, ...newVersions]);
      }
      setRestoreHasMore(newVersions.length === limit);
      setRestorePage(page);
    } catch (error) {
      LogService.Error("Failed to load restore versions: " + String(error));
      errorMessage("加载版本列表失败: " + parseCallServiceError(error));
    } finally {
      setRestoreLoading(false);
    }
  };

  // 打开恢复弹框时加载初始数据
  const handleOpenRestoreDialog = () => {
    setRestoreDialogOpen(true);
    setRestorePage(0);
    setRestoreHasMore(true);
    setSelectedVersion(null);
    loadRestoreVersions(0);
  };

  // 处理滚动加载更多
  const handleRestoreScroll = (event: React.UIEvent<HTMLDivElement>) => {
    const { scrollTop, scrollHeight, clientHeight } = event.currentTarget;
    // 滚动到底部时加载更多
    if (
      scrollHeight - scrollTop - clientHeight < 50 &&
      restoreHasMore &&
      !restoreLoading
    ) {
      loadRestoreVersions(restorePage + 1);
    }
  };

  const handleCancel = async () => {
    try {
      await CancelSync();
      successMessage("同步已取消");
    } catch (error) {
      LogService.Debug("Cancel failed: " + String(error));
    }
  };

  const handleDeleteConfirm = async () => {
    if (!versionToDelete) return;
    try {
      await DeleteSyncVersion(versionToDelete);
      successMessage("版本删除成功");
      await loadVersions();
    } catch (error) {
      LogService.Error("Delete version failed: " + String(error));
      errorMessage("删除版本失败: " + parseCallServiceError(error));
    } finally {
      setDeleteDialogOpen(false);
      setVersionToDelete(null);
    }
  };

  const handleDeleteCancel = () => {
    setDeleteDialogOpen(false);
    setVersionToDelete(null);
  };

  const formatFileSize = (bytes: number) => {
    if (bytes === 0) return "0 B";
    const k = 1024;
    const sizes = ["B", "KB", "MB", "GB"];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return (
      Number.parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + " " + sizes[i]
    );
  };

  const formatDate = (dateStr: string) => {
    return new Date(dateStr).toLocaleString("zh-CN");
  };

  const getStageLabel = (stage: string) => {
    return stageLabels[stage] || stage;
  };

  return (
    <Box>
      <Typography variant="h5" gutterBottom sx={{ mb: 3, fontWeight: 600 }}>
        数据同步
      </Typography>

      <Paper sx={{ p: 3, mb: 3 }} elevation={1}>
        <Typography variant="h6" gutterBottom>
          同步配置
        </Typography>
        <Stack spacing={2}>
          <FormRow label="服务器地址">
            <TextField
              fullWidth
              size="small"
              value={localConfig.serverUrl || ""}
              onChange={(e) =>
                setLocalConfig({
                  ...localConfig,
                  serverUrl: e.target.value.trim(),
                })
              }
              placeholder="例如: http://localhost:8080 或 https://sync.example.com"
            />
          </FormRow>
          <FormRow label="同步 ID">
            <TextField
              fullWidth
              size="small"
              value={localConfig.syncId || ""}
              onChange={(e) =>
                setLocalConfig({
                  ...localConfig,
                  syncId: e.target.value.trim(),
                })
              }
              placeholder="请输入同步 ID"
            />
          </FormRow>
          <FormRow label="用户密钥">
            <TextField
              fullWidth
              size="small"
              type="password"
              value={localConfig.userKey || ""}
              onChange={(e) =>
                setLocalConfig({
                  ...localConfig,
                  userKey: e.target.value.trim(),
                })
              }
              placeholder="请输入用户密钥"
            />
          </FormRow>
          <Box sx={{ display: "flex", gap: 2, mt: 2, alignItems: "center" }}>
            <Button
              variant="contained"
              startIcon={<Save />}
              onClick={handleSaveConfig}
            >
              保存配置
            </Button>
            {healthStatus && (
              <Chip
                label={healthStatus === "connected" ? "已连接" : "未连接"}
                color={healthStatus === "connected" ? "success" : "error"}
                size="small"
              />
            )}
            <Button
              variant="outlined"
              size="small"
              onClick={checkHealth}
              disabled={checkingHealth}
              startIcon={
                checkingHealth ? <CircularProgress size={16} /> : <Refresh />
              }
            >
              刷新
            </Button>
          </Box>
        </Stack>
      </Paper>

      {isConfigured && (
        <Paper sx={{ p: 3 }} elevation={1}>
          <Typography variant="h6" gutterBottom>
            同步操作
          </Typography>

          {/* 进度显示 */}
          {loading && progress && (
            <Box sx={{ mb: 3 }}>
              <Box sx={{ display: "flex", alignItems: "center", mb: 1 }}>
                <Typography variant="body2" sx={{ flex: 1 }}>
                  {getStageLabel(progress.stage)}
                </Typography>
                <Typography variant="body2" color="text.secondary">
                  {progress.progress.toFixed(1)}%
                </Typography>
              </Box>
              <LinearProgress
                variant="determinate"
                value={progress.progress}
                sx={{ mb: 1 }}
              />
              <Box sx={{ display: "flex", justifyContent: "space-between" }}>
                <Typography variant="caption" color="text.secondary">
                  {formatFileSize(progress.doneBytes)} /{" "}
                  {formatFileSize(progress.totalBytes)}
                </Typography>
                {progress.speed > 0 && (
                  <Typography variant="caption" color="text.secondary">
                    {formatFileSize(progress.speed)}/s
                  </Typography>
                )}
              </Box>
              {progress.error && (
                <Alert severity="error" sx={{ mt: 1 }}>
                  {progress.error}
                </Alert>
              )}
            </Box>
          )}

          <Box sx={{ display: "flex", gap: 2, mb: 3 }}>
            <Button
              variant="contained"
              startIcon={
                loading ? <CircularProgress size={20} /> : <CloudUpload />
              }
              onClick={handleUpload}
              disabled={loading}
            >
              上传数据
            </Button>
            <Button
              variant="outlined"
              startIcon={
                loading ? <CircularProgress size={20} /> : <CloudDownload />
              }
              onClick={handleOpenRestoreDialog}
              disabled={loading || versions.length === 0}
            >
              恢复数据
            </Button>
            <Button
              variant="outlined"
              startIcon={<Refresh />}
              onClick={loadVersions}
              disabled={loading}
            >
              刷新版本
            </Button>
            {loading && (
              <Button
                variant="outlined"
                color="error"
                startIcon={<Cancel />}
                onClick={handleCancel}
              >
                取消
              </Button>
            )}
          </Box>

          <Divider sx={{ my: 2 }} />

          <Typography variant="subtitle1" gutterBottom>
            历史版本
          </Typography>
          {versions.length === 0 ? (
            <Typography color="text.secondary">暂无历史版本</Typography>
          ) : (
            <List dense>
              {versions.map((version) => (
                <ListItem
                  key={version.version_number}
                  divider
                  secondaryAction={
                    <Button
                      size="small"
                      color="error"
                      startIcon={<DeleteIcon />}
                      onClick={() => {
                        setVersionToDelete(version.version_number);
                        setDeleteDialogOpen(true);
                      }}
                    >
                      删除
                    </Button>
                  }
                >
                  <ListItemText
                    primary={`版本 ${version.version_number}`}
                    secondary={`${formatFileSize(version.file_size)} • ${formatDate(
                      version.created_at,
                    )}`}
                  />
                </ListItem>
              ))}
            </List>
          )}
        </Paper>
      )}

      {/* 恢复对话框 */}
      <Dialog
        open={restoreDialogOpen}
        onClose={() => setRestoreDialogOpen(false)}
        maxWidth="sm"
        fullWidth
      >
        <DialogTitle>恢复数据</DialogTitle>
        <DialogContent>
          <Typography gutterBottom>请选择要恢复的版本:</Typography>
          <Box
            sx={{
              maxHeight: 400,
              overflow: "auto",
            }}
            onScroll={handleRestoreScroll}
          >
            <List>
              {restoreVersions.map((version) => (
                <ListItemButton
                  key={version.version_number}
                  selected={selectedVersion === version.version_number}
                  onClick={() => setSelectedVersion(version.version_number)}
                >
                  <ListItemText
                    primary={`版本 ${version.version_number}`}
                    secondary={`${formatFileSize(version.file_size)} • ${formatDate(
                      version.created_at,
                    )}`}
                  />
                </ListItemButton>
              ))}
            </List>
            {restoreLoading && (
              <Box sx={{ display: "flex", justifyContent: "center", py: 2 }}>
                <CircularProgress size={24} />
              </Box>
            )}
            {!restoreHasMore && restoreVersions.length > 0 && (
              <Typography
                variant="caption"
                color="text.secondary"
                sx={{ display: "block", textAlign: "center", py: 1 }}
              >
                已加载全部版本
              </Typography>
            )}
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setRestoreDialogOpen(false)}>取消</Button>
          <Button
            variant="contained"
            onClick={handleDownload}
            disabled={!selectedVersion || loading}
            startIcon={
              loading ? <CircularProgress size={20} /> : <CloudDownload />
            }
          >
            恢复
          </Button>
        </DialogActions>
      </Dialog>

      {/* 删除确认对话框 */}
      <Dialog
        open={deleteDialogOpen}
        onClose={handleDeleteCancel}
        maxWidth="xs"
        fullWidth
      >
        <DialogTitle>确认删除</DialogTitle>
        <DialogContent>
          <Typography>
            确定要删除版本 {versionToDelete} 吗？此操作不可恢复。
          </Typography>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleDeleteCancel}>取消</Button>
          <Button
            variant="contained"
            color="error"
            onClick={handleDeleteConfirm}
            disabled={loading}
          >
            删除
          </Button>
        </DialogActions>
      </Dialog>

      {/* 重启提示对话框 */}
      <Dialog
        open={restartDialogOpen}
        onClose={handleRestartDialogClose}
        maxWidth="xs"
        fullWidth
      >
        <DialogTitle>需要重启</DialogTitle>
        <DialogContent>
          <Typography>数据恢复操作已完成，请重启应用。</Typography>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleRestartDialogClose} variant="contained">
            知道了
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

export default SyncSettings;
