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
import useTerminalStore from "../stores/terminal";
import { AppTheme } from "../types/ssh";
import { SetTheme } from "../../bindings/github.com/ilaziness/vexo/services/configservice";
import { LogService } from "../../bindings/github.com/ilaziness/vexo/services";

// 主题选项配置
const THEME_OPTIONS: {
  value: AppTheme;
  label: string;
  mode: "light" | "dark";
}[] = [
  { value: "light", label: "亮白", mode: "light" },
  { value: "dark", label: "暗色", mode: "dark" },
  { value: "blueDark", label: "蓝夜", mode: "dark" },
  { value: "atom", label: "极黑", mode: "dark" },
  { value: "deep", label: "深邃", mode: "dark" },
  { value: "eyeCare", label: "暖棕", mode: "light" },
];

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

  const changeColorMode = (val: AppTheme) => {
    console.log("changeColorMode", val);

    // 保存主题配置
    SetTheme(val).catch((err) => {
      LogService.Error("Failed to save theme:" + err.message);
    });

    // 同步终端主题
    useTerminalStore.getState().syncWithGlobalTheme(val);

    const themeConfig = THEME_OPTIONS.find((t) => t.value === val);
    if (!themeConfig) return;

    setMode(themeConfig.mode);
    if (themeConfig.mode === "light") {
      setColorScheme({ light: val });
    } else {
      setColorScheme({ dark: val });
    }
  };

  const handleColorModeSelect = (val: AppTheme) => {
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
            mode === "dark" && !colorScheme
              ? "dark"
              : mode === "dark"
                ? darkColorScheme
                : lightColorScheme
          }
          onChange={(e) => handleColorModeSelect(e.target.value as AppTheme)}
        >
          {THEME_OPTIONS.map((option) => (
            <MenuItem key={option.value}>
              <FormControlLabel
                value={option.value}
                control={
                  <Radio
                    checked={mode !== "system" && colorScheme === option.value}
                    size="small"
                  />
                }
                label={option.label}
                sx={{ width: "100%" }}
              />
            </MenuItem>
          ))}
        </RadioGroup>
      </Menu>
    </>
  );
}
