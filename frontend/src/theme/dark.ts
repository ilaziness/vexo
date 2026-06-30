import type { ColorSystemOptions } from "@mui/material/styles";
import { darkTokens as t } from "./paletteTokens";

/**
 * Midnight 主题 - GitHub Dark 风格
 * 终端深底，Chrome 抬升，用户蓝 / AI 绿
 */
const dark: ColorSystemOptions = {
  palette: {
    mode: "dark",
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
      contrastText: "#0D1117",
    },
    background: {
      default: t.bgDefault,
      paper: t.bgPaper,
    },
    text: {
      primary: t.textPrimary,
      secondary: t.textSecondary,
      disabled: "#6E7681",
    },
    grey: {
      800: t.aiBubble,
    },
    divider: t.divider,
    action: {
      active: t.primary,
      hover: "rgba(88, 166, 255, 0.10)",
      selected: "rgba(88, 166, 255, 0.18)",
      disabled: "rgba(230, 237, 243, 0.30)",
      disabledBackground: "rgba(230, 237, 243, 0.12)",
    },
    success: {
      main: t.success,
      light: "#56D364",
      dark: "#238636",
    },
    error: {
      main: "#F85149",
      light: "#FF7B72",
      dark: "#DA3633",
    },
    warning: {
      main: "#D29922",
      light: "#E3B341",
      dark: "#BB8009",
    },
    info: {
      main: t.primary,
      light: t.primaryLight,
      dark: t.primaryDark,
    },
  },
};

export default dark;
