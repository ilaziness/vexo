import React from "react";
import {
  Dialog,
  DialogTitle,
  DialogContent,
  Slide,
  Stack,
  Tooltip,
} from "@mui/material";
import MobiledataOffIcon from "@mui/icons-material/MobiledataOff";
import SignalCellularAltIcon from "@mui/icons-material/SignalCellularAlt";
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
        <Tooltip title="传输列表">
          <MobiledataOffIcon
            fontSize="small"
            onClick={handleClick}
            sx={{
              color: "text.secondary",
              "&:hover": {
                color: "primary.main",
                cursor: "pointer",
              },
            }}
          />
        </Tooltip>
        <Tooltip title="传输列表">
          <SignalCellularAltIcon
            fontSize="small"
            onClick={handleClick}
            sx={{
              color: "text.secondary",
              "&:hover": {
                color: "primary.main",
                cursor: "pointer",
              },
            }}
          />
        </Tooltip>
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
