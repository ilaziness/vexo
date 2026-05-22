import type { ColorSystemOptions } from "@mui/material/styles";

/**
 * Amber 主题 - 护眼复古模式
 * 复古终端风格，暖琥珀色减少蓝光刺激
 */
const eyeCare: ColorSystemOptions = {
  palette: {
    mode: "light",
    primary: {
      main: "#92400e", // 深琥珀
      light: "#b45309",
      dark: "#78350f",
      contrastText: "#fef3c7",
    },
    secondary: {
      main: "#d97706", // 金黄
      light: "#fbbf24",
      dark: "#b45309",
      contrastText: "#451a03",
    },
    background: {
      default: "#fef3c7", // 暖琥珀 - 终端区域
      paper: "#fde8a8", // 稍深琥珀 - 工具栏/状态栏区分
    },
    text: {
      primary: "#451a03", // 深棕
      secondary: "#78350f",
      disabled: "#a16207",
    },
    divider: "rgba(69, 26, 3, 0.12)",
    action: {
      active: "#92400e",
      hover: "rgba(146, 64, 14, 0.04)",
      selected: "rgba(146, 64, 14, 0.08)",
      disabled: "rgba(69, 26, 3, 0.26)",
      disabledBackground: "rgba(69, 26, 3, 0.12)",
    },
    success: {
      main: "#65a30d",
      light: "#84cc16",
      dark: "#4d7c0f",
    },
    error: {
      main: "#c2410c",
      light: "#ea580c",
      dark: "#9a3412",
    },
    warning: {
      main: "#d97706",
      light: "#fbbf24",
      dark: "#b45309",
    },
    info: {
      main: "#b45309", // 琥珀棕
      light: "#d97706",
      dark: "#78350f",
    },
  },
};

export default eyeCare;
