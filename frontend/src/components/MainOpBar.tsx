import { Window } from "@wailsio/runtime";
import { Box, IconButton, Stack } from "@mui/material";
import { alpha } from "@mui/material/styles";
import { Close, CropSquare, FilterNone, Remove } from "@mui/icons-material";
import { useState } from "react";
import { AppService } from "../../bindings/github.com/ilaziness/vexo/services";

const MainOpBar = () => {
  const [isMaximised, setIsMaximised] = useState(false);
  const handleMinimize = () => {
    AppService.MainWindowMin().then(() => {});
  };

  const handleMaximize = () => {
    AppService.MainWindowMax().then(() => {
      Window.IsMaximised().then((val: boolean) => {
        setIsMaximised(val);
      });
    });
  };

  const handleClose = () => {
    AppService.MainWindowClose().then(() => {});
  };

  return (
    <Box sx={{ "--wails-draggable": "no-drag" }}>
      <Stack direction="row" spacing={1}>
        <IconButton
          size="small"
          onClick={handleMinimize}
          sx={(theme) => ({
            bgcolor: alpha(theme.palette.action.active, 0.08),
            "&:hover": { bgcolor: alpha(theme.palette.action.active, 0.12) },
          })}
        >
          <Remove fontSize="small" sx={{ color: "#f7b424" }} />
        </IconButton>
        <IconButton
          size="small"
          onClick={handleMaximize}
          sx={(theme) => ({
            bgcolor: alpha(theme.palette.action.active, 0.08),
            "&:hover": { bgcolor: alpha(theme.palette.action.active, 0.12) },
          })}
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
          sx={(theme) => ({
            bgcolor: alpha(theme.palette.action.active, 0.08),
            "&:hover": { bgcolor: alpha(theme.palette.action.active, 0.12) },
          })}
        >
          <Close fontSize="small" sx={{ color: "#FF5F57" }} />
        </IconButton>
      </Stack>
    </Box>
  );
};

export default MainOpBar;
