import React from "react";
import {
  Box,
  IconButton,
  Slide,
  Stack,
  Tooltip,
  Typography,
} from "@mui/material";
import MobiledataOffIcon from "@mui/icons-material/MobiledataOff";
import SignalCellularAltIcon from "@mui/icons-material/SignalCellularAlt";
import CloseIcon from "@mui/icons-material/Close";
import TransferList from "./TransferList";

interface StatusBarProps {
  sessionID: string;
}

const StatusBar: React.FC<StatusBarProps> = ({ sessionID }) => {
  const [open, setOpen] = React.useState(false);

  const toggleOpen = () => {
    setOpen((prev) => !prev);
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
            onClick={toggleOpen}
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
            onClick={toggleOpen}
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

      <Slide in={open} direction="up" mountOnEnter unmountOnExit>
        <Box
          sx={{
            position: "absolute",
            bottom: "25px",
            left: 0,
            right: 0,
            height: "50%",
            zIndex: 10,
            display: "flex",
            flexDirection: "column",
            bgcolor: "background.default",
            borderTop: 1,
            borderColor: "divider",
            boxShadow: 4,
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
            <IconButton onClick={handleClose} size="small">
              <CloseIcon fontSize="small" />
            </IconButton>
          </Box>
          <Box
            sx={{
              flex: 1,
              display: "flex",
              flexDirection: "column",
              overflow: "hidden",
              p: 0,
            }}
          >
            <TransferList sessionID={sessionID} />
          </Box>
        </Box>
      </Slide>
    </>
  );
};

export default StatusBar;
