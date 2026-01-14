import React from "react";
import {
  Box,
  List,
  ListItem,
  ListItemText,
  LinearProgress,
  Typography,
  IconButton,
} from "@mui/material";
import { CloudUpload, CloudDownload, Close } from "@mui/icons-material";
import { useTransferStore } from "../stores/transfer";
import { ProgressData } from "../../bindings/github.com/ilaziness/vexo/services/models";

interface TransferListProps {
  sessionID: string;
}

const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return "0 B";
  const k = 1024;
  const sizes = ["B", "KB", "MB", "GB"];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + " " + sizes[i];
};

const TransferList: React.FC<TransferListProps> = ({ sessionID }) => {
  const { getTransfersBySession, removeProgress } = useTransferStore();
  const transfers = getTransfersBySession(sessionID);

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
    <Box sx={{ width: "100%", height: "100%", display: "flex", flexDirection: "column" }}>
      <List sx={{ flex: 1, overflow: "auto" }}>
        {transfers.map((transfer: ProgressData) => {
          const progress = transfer.TotalSize > 0 ? (transfer.Rate / transfer.TotalSize) * 100 : 0;
          const isUpload = transfer.TransferType.toLowerCase().includes("upload");

          return (
            <ListItem
              key={transfer.ID}
              sx={{
                border: 1,
                borderColor: "divider",
                borderRadius: 1,
                mb: 1,
                flexDirection: "column",
                alignItems: "stretch",
              }}
            >
              <Box sx={{ display: "flex", alignItems: "center", width: "100%" }}>
                {isUpload ? (
                  <CloudUpload color="primary" sx={{ mr: 1 }} />
                ) : (
                  <CloudDownload color="secondary" sx={{ mr: 1 }} />
                )}
                <ListItemText
                  primary={
                    <Typography variant="body2" sx={{ fontWeight: "bold" }}>
                      {isUpload ? "上传" : "下载"}: {transfer.LocalFile.split("/").pop() || transfer.LocalFile}
                    </Typography>
                  }
                  secondary={
                    <Typography variant="caption" color="text.secondary">
                      {transfer.RemoteFile}
                    </Typography>
                  }
                />
                <IconButton size="small" onClick={() => handleRemove(transfer.ID)}>
                  <Close fontSize="small" />
                </IconButton>
              </Box>
              <Box sx={{ mt: 1, width: "100%" }}>
                <Box sx={{ display: "flex", justifyContent: "space-between", mb: 0.5 }}>
                  <Typography variant="caption">
                    {formatFileSize(transfer.Rate)} / {formatFileSize(transfer.TotalSize)}
                  </Typography>
                  <Typography variant="caption">
                    {progress.toFixed(1)}%
                  </Typography>
                </Box>
                <LinearProgress variant="determinate" value={progress} />
              </Box>
            </ListItem>
          );
        })}
      </List>
    </Box>
  );
};

export default TransferList;
