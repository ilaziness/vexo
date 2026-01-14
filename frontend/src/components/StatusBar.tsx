import React from "react";
import {
  IconButton,
  Dialog,
  DialogTitle,
  DialogContent,
  Slide,
  Icon,
  Box,
  Stack,
} from "@mui/material";
import { CloudUpload } from "@mui/icons-material";
import TransferList from "./TransferList";

interface StatusBarProps {
  sessionID: string;
}

const StatusBar: React.FC<StatusBarProps> = ({ sessionID }) => {
  const [open, setOpen] = React.useState(false);

  const handleClick = () => {
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
  };

  return (
    <>
      <Stack direction="row" spacing={2}>
        <Icon
          onClick={handleClick}
          sx={{
            fontSize: 12,
            color: "text.secondary",
            "&:hover": {
              color: "primary.main",
              cursor: "pointer",
            },
          }}
        >
          <CloudUpload fontSize="small" />
        </Icon>
        <Icon
          onClick={handleClick}
          sx={{
            fontSize: 15,
            color: "text.secondary",
            "&:hover": {
              color: "primary.main",
            },
          }}
        >
          <CloudUpload fontSize="small" />
        </Icon>
      </Stack>

      <Dialog
        open={open}
        onClose={handleClose}
        fullScreen
        slots={{
          transition: Slide,
        }}
        slotProps={{
          transition: {
            direction: "up",
          },
          paper: {
            sx: {
              width: "100%",
              height: "50%",
            },
          },
        }}
      >
        <DialogTitle>传输列表</DialogTitle>
        <DialogContent
          sx={{
            flex: 1,
            display: "flex",
            flexDirection: "column",
            p: 0,
          }}
        >
          <TransferList sessionID={sessionID} />
        </DialogContent>
      </Dialog>
    </>
  );
};

export default StatusBar;
