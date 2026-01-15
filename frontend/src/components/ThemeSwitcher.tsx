import {
  IconButton,
  Tooltip,
  Menu,
  MenuItem,
  FormControlLabel,
  Radio,
  RadioGroup,
} from "@mui/material";
import { useColorScheme } from "@mui/material/styles";
import PaletteIcon from "@mui/icons-material/Palette";
import React, { useState } from "react";

export default function ThemeSwitcher() {
  const {
    mode,
    setMode,
    colorScheme,
    setColorScheme,
    lightColorScheme,
    darkColorScheme,
  } = useColorScheme();

  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const open = Boolean(anchorEl);

  const handleColorModeMenuOpen = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleColorModeMenuClose = () => {
    setAnchorEl(null);
  };

  const changeColorMode = (
    val: "system" | "light" | "dark" | "atom" | "deep" | "eyeCare",
  ) => {
    console.log("changeColorMode", val);
    if (val === "system") {
      setMode("system");
      return;
    }

    if (val === "light" || val === "eyeCare") {
      setMode("light");
      setColorScheme({ light: val });
      return;
    }

    setMode("dark");
    setColorScheme({ dark: val });
  };

  const handleColorModeSelect = (
    val: "system" | "light" | "dark" | "atom" | "deep" | "eyeCare",
  ) => {
    changeColorMode(val);
    handleColorModeMenuClose();
  };

  return (
    <>
      <Tooltip title="颜色模式">
        <IconButton size="small" onClick={handleColorModeMenuOpen}>
          <PaletteIcon />
        </IconButton>
      </Tooltip>
      <Menu
        anchorEl={anchorEl}
        open={open}
        onClose={handleColorModeMenuClose}
        onClick={(e) => e.stopPropagation()}
        slotProps={{
          paper: {
            elevation: 8,
            sx: {
              overflow: "visible",
              mt: 1,
            },
          },
        }}
      >
        <RadioGroup
          value={
            mode === "system"
              ? "system"
              : mode === "dark"
                ? darkColorScheme
                : lightColorScheme
          }
          onChange={(e) => handleColorModeSelect(e.target.value as any)}
        >
          <MenuItem>
            <FormControlLabel
              value="system"
              control={<Radio checked={mode === "system"} size="small" />}
              label="系统"
              sx={{ width: "100%" }}
            />
          </MenuItem>
          <MenuItem>
            <FormControlLabel
              value="light"
              control={
                <Radio
                  checked={mode !== "system" && colorScheme === "light"}
                  size="small"
                />
              }
              label="明亮"
              sx={{ width: "100%" }}
            />
          </MenuItem>
          <MenuItem>
            <FormControlLabel
              value="dark"
              control={
                <Radio
                  checked={mode !== "system" && colorScheme === "dark"}
                  size="small"
                />
              }
              label="暗黑"
              sx={{ width: "100%" }}
            />
          </MenuItem>
          <MenuItem>
            <FormControlLabel
              value="atom"
              control={
                <Radio
                  checked={mode !== "system" && colorScheme === "atom"}
                  size="small"
                />
              }
              label="极客黑"
              sx={{ width: "100%" }}
            />
          </MenuItem>
          <MenuItem>
            <FormControlLabel
              value="deep"
              control={
                <Radio
                  checked={mode !== "system" && colorScheme === "deep"}
                  size="small"
                />
              }
              label="深邃夜"
              sx={{ width: "100%" }}
            />
          </MenuItem>
          <MenuItem>
            <FormControlLabel
              value="eyeCare"
              control={
                <Radio
                  checked={mode !== "system" && colorScheme === "eyeCare"}
                  size="small"
                />
              }
              label="护眼"
              sx={{ width: "100%" }}
            />
          </MenuItem>
        </RadioGroup>
      </Menu>
    </>
  );
}
