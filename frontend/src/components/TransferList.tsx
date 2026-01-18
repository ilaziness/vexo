import React, { useMemo } from "react";
import {
  Box,
  List,
  ListItem,
  LinearProgress,
  Typography,
  IconButton,
  Slide,
  Tooltip,
} from "@mui/material";
import {
  CloudUpload,
  CloudDownload,
  Close,
  ArrowForward,
  ArrowBack,
  ClearAll,
  Cancel,
} from "@mui/icons-material";
import { useTransferStore } from "../stores/transfer";
import { formatFileSize } from "../func/service";
import { ProgressData } from "../../bindings/github.com/ilaziness/vexo/services/models";
import { SftpService } from "../../bindings/github.com/ilaziness/vexo/services";

interface TransferListProps {
  sessionID: string;
  open: boolean;
  statusBarHeight: string;
  onClose: () => void;
}

const TransferList: React.FC<TransferListProps> = ({
  sessionID,
  open,
  statusBarHeight,
  onClose,
}) => {
  const {
    transfers: transfersMap,
    removeProgress,
    clearSession,
  } = useTransferStore();
  const transfersList = useMemo(() => {
    return transfersMap.get(sessionID) || [];
  }, [transfersMap, sessionID]);
  const handleRemove = (id: string) => {
    removeProgress(sessionID, id);
  };
  const handleClear = () => {
    clearSession(sessionID);
  };

  return (
    <Box
      sx={{
        position: "absolute",
        bottom: `${statusBarHeight}`,
        left: 0,
        right: 0,
        height: "50%",
        zIndex: 10,
        overflow: "hidden",
        pointerEvents: open ? "auto" : "none",
      }}
    >
      <Slide in={open} direction="up" mountOnEnter unmountOnExit>
        <Box
          sx={{
            height: "100%",
            display: "flex",
            flexDirection: "column",
            bgcolor: "background.default",
            borderTop: 1,
            borderColor: "divider",
          }}
        >
          <Box
            sx={{
              p: 1,
              px: 2,
              borderBottom: 1,
              borderColor: "divider",
              display: "flex",
              justifyContent: "space-between",
              alignItems: "center",
              height: "40px",
            }}
          >
            <Typography variant="subtitle2" sx={{ fontWeight: "bold" }}>
              传输列表
            </Typography>
            <Box sx={{ display: "flex", gap: 1 }}>
              <Tooltip title="清空列表">
                <IconButton onClick={handleClear} size="small">
                  <ClearAll fontSize="small" />
                </IconButton>
              </Tooltip>
              <Tooltip title="关闭">
                <IconButton onClick={onClose} size="small">
                  <Close fontSize="small" />
                </IconButton>
              </Tooltip>
            </Box>
          </Box>
          {transfersList.length === 0 ? (
            <Box sx={{ p: 2, textAlign: "center" }}>
              <Typography variant="body2" color="text.secondary">
                暂无传输任务
              </Typography>
            </Box>
          ) : (
            <Box
              sx={{
                flex: 1,
                display: "flex",
                flexDirection: "column",
                overflow: "hidden",
                p: 0,
              }}
            >
              <List sx={{ flex: 1, overflow: "auto" }}>
                {transfersList.map((transfer: ProgressData) => {
                  const progress = transfer.Done ? 100 : transfer.Rate;
                  const isUpload =
                    transfer.TransferType.toLowerCase().includes("upload");
                  const isCompleted = transfer.Done;
                  const hasError =
                    transfer.Error && transfer.Error.trim() !== "";

                  return (
                    <ListItem
                      key={transfer.ID}
                      sx={{
                        border: 1,
                        borderColor: hasError ? "error.main" : "divider",
                        borderRadius: 1,
                        mb: 1,
                        "&:hover": {
                          backgroundColor: "action.hover",
                          borderColor: hasError ? "error.main" : "primary.main",
                        },
                        transition: "all 0.2s ease-in-out",
                      }}
                    >
                      <Box
                        sx={{
                          display: "flex",
                          flexDirection: "column",
                          width: "100%",
                          gap: 1,
                        }}
                      >
                        <Box
                          sx={{
                            display: "flex",
                            alignItems: "center",
                            width: "100%",
                            gap: 1,
                            flexWrap: "wrap",
                          }}
                        >
                          {/* 上传/下载图标 */}
                          {isUpload ? (
                            <CloudUpload
                              color={hasError ? "error" : "primary"}
                              sx={{ flexShrink: 0 }}
                            />
                          ) : (
                            <CloudDownload
                              color={hasError ? "error" : "secondary"}
                              sx={{ flexShrink: 0 }}
                            />
                          )}

                          {/* 本地文件全路径 */}
                          <Typography
                            variant="body2"
                            sx={{
                              flex: 1,
                              minWidth: 0,
                              overflow: "hidden",
                              textOverflow: "ellipsis",
                              whiteSpace: "nowrap",
                            }}
                            title={transfer.LocalFile}
                          >
                            {transfer.LocalFile}
                          </Typography>

                          {/* 箭头 */}
                          {isUpload ? (
                            <ArrowForward
                              sx={{
                                flexShrink: 0,
                                color: hasError
                                  ? "error.main"
                                  : "primary.light",
                              }}
                            />
                          ) : (
                            <ArrowBack
                              sx={{
                                flexShrink: 0,
                                color: hasError
                                  ? "error.main"
                                  : "primary.light",
                              }}
                            />
                          )}

                          {/* 远程全路径 */}
                          <Typography
                            variant="body2"
                            sx={{
                              flex: 1,
                              minWidth: 0,
                              overflow: "hidden",
                              textOverflow: "ellipsis",
                              whiteSpace: "nowrap",
                            }}
                            title={transfer.RemoteFile}
                          >
                            {transfer.RemoteFile}
                          </Typography>

                          {/* 进度条 */}
                          <Box sx={{ flex: 2, minWidth: 100, maxWidth: 200 }}>
                            <LinearProgress
                              variant="determinate"
                              value={progress}
                              sx={{
                                height: 6,
                                borderRadius: 3,
                                "& .MuiLinearProgress-bar": {
                                  backgroundColor: hasError
                                    ? "error.main"
                                    : isCompleted
                                      ? "success.main"
                                      : undefined,
                                },
                              }}
                            />
                          </Box>

                          {/* 文件大小 */}
                          <Typography
                            variant="caption"
                            sx={{ flexShrink: 0, minWidth: "fit-content" }}
                          >
                            {formatFileSize(transfer.TotalSize)}
                          </Typography>

                          {/* 上传百分比 */}
                          <Typography
                            variant="caption"
                            sx={{ flexShrink: 0, minWidth: "fit-content" }}
                          >
                            {progress.toFixed(2)}%
                          </Typography>

                          {/* 取消或清除图标 */}
                          {!transfer.Done ? (
                            <Tooltip title="取消传输">
                              <IconButton
                                size="small"
                                onClick={() =>
                                  SftpService.CancelTransfer(transfer.ID)
                                }
                                sx={{ flexShrink: 0 }}
                              >
                                <Cancel fontSize="small" />
                              </IconButton>
                            </Tooltip>
                          ) : (
                            <Tooltip title="删除">
                              <IconButton
                                size="small"
                                onClick={() => handleRemove(transfer.ID)}
                                sx={{ flexShrink: 0 }}
                              >
                                <Close fontSize="small" />
                              </IconButton>
                            </Tooltip>
                          )}
                        </Box>

                        {/* 错误信息 */}
                        {hasError && (
                          <Typography
                            variant="caption"
                            color="error"
                            sx={{
                              px: 1,
                              wordBreak: "break-word",
                            }}
                          >
                            错误: {transfer.Error}
                          </Typography>
                        )}
                      </Box>
                    </ListItem>
                  );
                })}
              </List>
            </Box>
          )}
        </Box>
      </Slide>
    </Box>
  );
};

export default TransferList;
