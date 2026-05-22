import type { ColorSystemOptions } from "@mui/material/styles";

/**
 * Midnight 主题 - 暗色夜空模式
 * 深夜工作模式，夜空感深蓝灰，专业终端风格
 */
const dark: ColorSystemOptions = {
  palette: {
    mode: "dark",
    primary: {
      main: "#60a5fa", // 明亮蓝
      light: "#93c5fd",
      dark: "#3b82f6",
      contrastText: "#0f172a",
    },
    secondary: {
      main: "#22d3ee", // 霓虹青
      light: "#67e8f9",
      dark: "#06b6d4",
      contrastText: "#0f172a",
    },
    background: {
      default: "#0f172a", // 深 slate - 终端区域
      paper: "#1a2234", // 稍浅蓝黑 - 工具栏/状态栏区分
    },
    text: {
      primary: "#e2e8f0", // 柔和白
      secondary: "#94a3b8",
      disabled: "#64748b",
    },
    divider: "rgba(226, 232, 240, 0.12)",
    action: {
      active: "#60a5fa",
      hover: "rgba(96, 165, 250, 0.08)",
      selected: "rgba(96, 165, 250, 0.16)",
      disabled: "rgba(226, 232, 240, 0.3)",
      disabledBackground: "rgba(226, 232, 240, 0.12)",
    },
    success: {
      main: "#34d399",
      light: "#6ee7b7",
      dark: "#10b981",
    },
    error: {
      main: "#f87171",
      light: "#fca5a5",
      dark: "#ef4444",
    },
    warning: {
      main: "#fbbf24",
      light: "#fcd34d",
      dark: "#f59e0b",
    },
    info: {
      main: "#60a5fa",
      light: "#93c5fd",
      dark: "#3b82f6",
    },
  },
};

export default dark;
