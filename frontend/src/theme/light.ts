import type { ColorSystemOptions } from "@mui/material/styles";

/**
 * Daylight 主题 - 亮色专业模式
 * 适合白天使用，清爽专业，类似现代IDE明亮模式
 */
const light: ColorSystemOptions = {
  palette: {
    mode: "light",
    primary: {
      main: "#2563eb", // 专业蓝
      light: "#60a5fa",
      dark: "#1d4ed8",
      contrastText: "#ffffff",
    },
    secondary: {
      main: "#0891b2", // 终端青
      light: "#22d3ee",
      dark: "#0e7490",
      contrastText: "#ffffff",
    },
    background: {
      default: "#f8f9fa", // 轻微暖白 - 终端区域
      paper: "#eef2f6", // 稍深灰蓝 - 工具栏/状态栏区分
    },
    text: {
      primary: "#1f2937", // 深灰，非纯黑
      secondary: "#6b7280",
      disabled: "#9ca3af",
    },
    divider: "rgba(31, 41, 55, 0.12)",
    action: {
      active: "#2563eb",
      hover: "rgba(37, 99, 235, 0.04)",
      selected: "rgba(37, 99, 235, 0.08)",
      disabled: "rgba(31, 41, 55, 0.26)",
      disabledBackground: "rgba(31, 41, 55, 0.12)",
    },
    success: {
      main: "#10b981",
      light: "#34d399",
      dark: "#059669",
    },
    error: {
      main: "#ef4444",
      light: "#f87171",
      dark: "#dc2626",
    },
    warning: {
      main: "#f59e0b",
      light: "#fbbf24",
      dark: "#d97706",
    },
    info: {
      main: "#2563eb",
      light: "#60a5fa",
      dark: "#1d4ed8",
    },
  },
};

export default light;
