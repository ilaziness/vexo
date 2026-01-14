import React, { useMemo } from "react";
import {
  Box,
  List,
  ListItem,
  LinearProgress,
  Typography,
  IconButton,
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
}

const TransferList: React.FC<TransferListProps> = ({ sessionID }) => {
  const transfersMap = useTransferStore((state) => state.transfers);
  const { removeProgress } = useTransferStore();

  const transfers = useMemo(() => {
    return transfersMap.get(sessionID) || [];
  }, [transfersMap, sessionID]);

  const handleRemove = (id: string) => {
    removeProgress(sessionID, id);
  };

  if (transfers.length === 0) {
    return (
      <Box sx={{ p: 2, textAlign: "center" }}>
        <Typography variant="body2" color="text.secondary">
          暂无传输任务
        </Typography>
      </Box>
    );
  }

  return (
    <Box
      sx={{
        width: "100%",
        height: "100%",
        display: "flex",
        flexDirection: "column",
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
                  <CloudDownload color="secondary" sx={{ flexShrink: 0 }} />
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
                  <ArrowBack sx={{ flexShrink: 0, color: "text.secondary" }} />
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
  );
};

export default TransferList;
