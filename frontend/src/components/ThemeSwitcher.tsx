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

// 主题类型定义
type ThemeOption = "system" | "light" | "dark" | "blueDark" | "atom" | "deep" | "eyeCare";

// 主题选项配置
const THEME_OPTIONS: { value: ThemeOption; label: string }[] = [
  { value: "system", label: "系统" },
  { value: "light", label: "明亮" },
  { value: "dark", label: "经典暗色" },
  { value: "blueDark", label: "蓝色暗色" },
  { value: "atom", label: "代码风格" },
  { value: "deep", label: "深色主题" },
  { value: "eyeCare", label: "护眼模式" },
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

  const changeColorMode = (val: ThemeOption) => {
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

  const handleColorModeSelect = (val: ThemeOption) => {
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
              : mode === "dark" && !colorScheme
                ? "dark"
                : mode === "dark"
                  ? darkColorScheme
                  : lightColorScheme
          }
          onChange={(e) => handleColorModeSelect(e.target.value as ThemeOption)}
        >
          {THEME_OPTIONS.map((option) => (
            <MenuItem key={option.value}>
              <FormControlLabel
                value={option.value}
                control={
                  <Radio
                    checked={
                      option.value === "system"
                        ? mode === "system"
                        : mode !== "system" && colorScheme === option.value
                    }
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
