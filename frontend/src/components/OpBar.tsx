import { Window } from "@wailsio/runtime";
import { Box, IconButton, Stack } from "@mui/material";
import { alpha } from "@mui/material/styles";
import { Close, CropSquare, FilterNone, Remove } from "@mui/icons-material";
import { useState } from "react";

interface OpBarProps {
  onMinimize?: () => void;
  onMaximize?: () => void;
  onClose?: () => void;
  draggable?: boolean;
  showBorder?: boolean;
}

const OpBar = ({
  onClose,
  draggable = true,
  showBorder = true,
}: OpBarProps) => {
  const [isMaximised, setIsMaximised] = useState(false);

  const handleMinimize = () => {
    Window.Minimise().then(() => {});
  };

  const handleMaximize = () => {
    Window.ToggleMaximise().then(() => {
      Window.IsMaximised().then((val: boolean) => {
        setIsMaximised(val);
      });
    });
  };

  const handleClose = () => {
    if (onClose) {
      onClose();
    } else {
      Window.Close().then(() => {});
    }
  };

  return (
    <Box
      sx={{
        "--wails-draggable": draggable ? "drag" : "no-drag",
        height: "100%",
        // width: "100%",
        display: "flex",
        alignItems: "center",
        justifyContent: "flex-end",
        pr: 1,
        borderBottom: showBorder ? 1 : 0,
        borderColor: "divider",
      }}
    >
      <Stack
        direction="row"
        spacing={1}
        sx={{ "--wails-draggable": "no-drag" }}
      >
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

export default OpBar;
