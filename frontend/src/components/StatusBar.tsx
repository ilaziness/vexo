import React from "react";
import { Box, Stack, Tooltip } from "@mui/material";
import MobiledataOffIcon from "@mui/icons-material/MobiledataOff";
import TransferList from "./TransferList";

interface StatusBarProps {
  sessionID: string;
  height: string;
}

const StatusBar: React.FC<StatusBarProps> = ({ sessionID, height }) => {
  const [open, setOpen] = React.useState(false);
  const toggleOpen = () => {
    setOpen((prev) => !prev);
  };
  const handleClose = () => {
    setOpen(false);
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
