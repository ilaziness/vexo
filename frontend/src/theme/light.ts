import type { ColorSystemOptions } from "@mui/material/styles";
import { lightTokens as t } from "./paletteTokens";

/**
 * Daylight 主题 - 亮色 IDE 风格
 * 终端/聊天区灰底，Chrome 白底，用户蓝 / AI 绿
 */
const light: ColorSystemOptions = {
  palette: {
    mode: "light",
    primary: {
      main: t.primary,
      light: t.primaryLight,
      dark: t.primaryDark,
      contrastText: "#FFFFFF",
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
      disabled: "#94A3B8",
    },
    grey: {
      100: t.aiBubble,
    },
    divider: t.divider,
    action: {
      active: t.primary,
      hover: "rgba(37, 99, 235, 0.06)",
      selected: "rgba(37, 99, 235, 0.10)",
      disabled: "rgba(15, 23, 42, 0.26)",
      disabledBackground: "rgba(15, 23, 42, 0.12)",
    },
    success: {
      main: t.success,
      light: "#34D399",
      dark: "#047857",
    },
    error: {
      main: "#EF4444",
      light: "#F87171",
      dark: "#DC2626",
    },
    warning: {
      main: "#F59E0B",
      light: "#FBBF24",
      dark: "#D97706",
    },
    info: {
      main: t.primary,
      light: t.primaryLight,
      dark: t.primaryDark,
    },
  },
};

export default light;
