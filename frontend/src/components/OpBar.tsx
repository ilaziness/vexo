"use client";

import { Window } from "@wailsio/runtime";
import { Box, IconButton, Stack } from "@mui/material";
import { Close, CropSquare, FilterNone, Remove } from "@mui/icons-material";
import { useState } from "react";

const OpBar = () => {
  const [isMaximised, setIsMaximised] = useState(false);
  const handleMinimize = () => {
    Window.Minimise().then(() => {});
  };

  const handleMaximize = () => {
    Window.ToggleMaximise().then(() => {});
    setIsMaximised(!isMaximised);
  };

  const handleClose = () => {
    Window.Close().then(() => {});
  };

  return (
    <Box sx={{ "--wails-draggable": "no-drag" }}>
      <Stack direction="row" spacing={1}>
        <IconButton
          size="small"
          onClick={handleMinimize}
          sx={{
            bgcolor: "rgba(255,255,255, 0.1)",
            "&:hover": { bgcolor: "rgba(255,255,255, 0.2)" },
          }}
        >
          <Remove fontSize="small" sx={{ color: "#FFBD2E" }} />
        </IconButton>
        <IconButton
          size="small"
          onClick={handleMaximize}
          sx={{
            bgcolor: "rgba(255,255,255, 0.1)",
            "&:hover": { bgcolor: "rgba(255,255,255, 0.2)" },
          }}
        >
          {isMaximised ? (
            <FilterNone fontSize="small" sx={{ color: "#28C940" }} />
          ) : (
            <CropSquare fontSize="small" sx={{ color: "#28C940" }} />
          )}
        </IconButton>
        <IconButton
          size="small"
          onClick={handleClose}
          sx={{
            bgcolor: "rgba(255,255,255, 0.1)",
            "&:hover": { bgcolor: "rgba(255,255,255, 0.2)" },
          }}
        >
          <Close fontSize="small" sx={{ color: "#FF5F57" }} />
        </IconButton>
      </Stack>
    </Box>
  );
};

export default OpBar;
