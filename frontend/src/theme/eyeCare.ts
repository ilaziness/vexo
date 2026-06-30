import type { ColorSystemOptions } from "@mui/material/styles";
import { eyeCareTokens as t } from "./paletteTokens";

/**
 * Amber 主题 - 护眼暖色模式
 * 暖琥珀底，teal AI 头像，assistant 气泡带细边框
 */
const eyeCare: ColorSystemOptions = {
  palette: {
    mode: "light",
    primary: {
      main: t.primary,
      light: t.primaryLight,
      dark: t.primaryDark,
      contrastText: "#FFFBEB",
    },
    secondary: {
      main: t.secondary,
      light: t.secondaryLight,
      dark: t.secondaryDark,
      contrastText: "#FFFFFF",
    },
    background: {
      default: t.bgDefault,
      paper: t.bgPaper,
    },
    text: {
      primary: t.textPrimary,
      secondary: t.textSecondary,
      disabled: "#A16207",
    },
    grey: {
      100: t.aiBubble,
    },
    divider: t.divider,
    action: {
      active: t.primary,
      hover: "rgba(180, 83, 9, 0.06)",
      selected: "rgba(180, 83, 9, 0.10)",
      disabled: "rgba(69, 26, 3, 0.26)",
      disabledBackground: "rgba(69, 26, 3, 0.12)",
    },
    success: {
      main: t.success,
      light: "#84CC16",
      dark: "#4D7C0F",
    },
    error: {
      main: "#C2410C",
      light: "#EA580C",
      dark: "#9A3412",
    },
    warning: {
      main: "#D97706",
      light: "#FBBF24",
      dark: "#B45309",
    },
    info: {
      main: t.primary,
      light: t.primaryLight,
      dark: t.primaryDark,
    },
  },
};

export default eyeCare;
