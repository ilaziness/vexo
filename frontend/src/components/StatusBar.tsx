import React from "react";
import { Box, Stack, Tooltip } from "@mui/material";
import MobiledataOffIcon from "@mui/icons-material/MobiledataOff";
import AltRouteIcon from "@mui/icons-material/AltRoute";
import TransferList from "./TransferList";
import SSHTunnel from "./SSHTunnel";

interface StatusBarProps {
  sessionID: string;
  height: string;
}

const StatusBar: React.FC<StatusBarProps> = ({ sessionID, height }) => {
  const [open, setOpen] = React.useState(false);
  const [sshTunnelOpen, setSshTunnelOpen] = React.useState(false);

  const toggleOpen = () => {
    setOpen((prev) => !prev);
  };
  const handleClose = () => {
    setOpen(false);
  };

  const toggleSshTunnelOpen = () => {
    setSshTunnelOpen((prev) => !prev);
  };
  const handleSshTunnelClose = () => {
    setSshTunnelOpen(false);
  };

  return (
    <Box
      sx={{
        height: height,
        px: 2,
        fontSize: 8,
        display: "flex",
        alignItems: "center",
        justifyContent: "flex-end",
        direction: "row",
        bgcolor: "background.default",
      }}
    >
      <Stack direction="row" spacing={2}>
        <Tooltip title="隧道">
          <AltRouteIcon
            fontSize="small"
            onClick={toggleSshTunnelOpen}
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
      </Stack>

      <SSHTunnel
        sessionID={sessionID}
        open={sshTunnelOpen}
        statusBarHeight={height}
        onClose={handleSshTunnelClose}
      />

      <TransferList
        sessionID={sessionID}
        open={open}
        statusBarHeight={height}
        onClose={handleClose}
      />
    </Box>
  );
};

export default StatusBar;
