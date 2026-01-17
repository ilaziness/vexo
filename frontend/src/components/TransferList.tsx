import React, { useMemo } from "react";
import {
  Box,
  List,
  ListItem,
  LinearProgress,
  Typography,
  IconButton,
  Slide,
} from "@mui/material";
import {
  CloudUpload,
  CloudDownload,
  Close,
  ArrowForward,
  ArrowBack,
} from "@mui/icons-material";
import { useTransferStore } from "../stores/transfer";
import { formatFileSize } from "../func/service";
import { ProgressData } from "../../bindings/github.com/ilaziness/vexo/services/models";

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
  const transfersMap = useTransferStore((state) => state.transfers);
  const { removeProgress } = useTransferStore();

  const transfers = useMemo(() => {
    return transfersMap.get(sessionID) || [];
  }, [transfersMap, sessionID]);

  const handleRemove = (id: string) => {
    removeProgress(sessionID, id);
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
            <IconButton onClick={onClose} size="small">
              <Close fontSize="small" />
            </IconButton>
          </Box>
          {transfers.length === 0 ? (
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
                {transfers.map((transfer: ProgressData) => {
                  const progress =
                    transfer.TotalSize > 0
                      ? (transfer.Rate / transfer.TotalSize) * 100
                      : 0;
                  const isUpload =
                    transfer.TransferType.toLowerCase().includes("upload");

                  return (
                    <ListItem
                      key={transfer.ID}
                      sx={{
                        border: 1,
                        borderColor: "divider",
                        borderRadius: 1,
                        mb: 1,
                        "&:hover": {
                          backgroundColor: "action.hover",
                          borderColor: "primary.main",
                        },
                        transition: "all 0.2s ease-in-out",
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
                          <CloudUpload color="primary" sx={{ flexShrink: 0 }} />
                        ) : (
                          <CloudDownload
                            color="secondary"
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
                            sx={{ flexShrink: 0, color: "text.secondary" }}
                          />
                        ) : (
                          <ArrowBack
                            sx={{ flexShrink: 0, color: "text.secondary" }}
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
                            sx={{ height: 6, borderRadius: 3 }}
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
                          {progress.toFixed(1)}%
                        </Typography>

                        {/* 清除图标 */}
                        <IconButton
                          size="small"
                          onClick={() => handleRemove(transfer.ID)}
                          sx={{ flexShrink: 0 }}
                        >
                          <Close fontSize="small" />
                        </IconButton>
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
