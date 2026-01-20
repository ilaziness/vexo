import type { ColorSystemOptions } from "@mui/material/styles";

const terminal: ColorSystemOptions = {
  palette: {
    mode: "dark",
    // 主色调：经典Terminal绿色
    primary: {
      main: "#33ff00", // 亮绿色主色
      light: "#66ff33", // 更亮的绿色变体
      dark: "#00cc00", // 稍深的绿色
      contrastText: "#0a0a0a", // 深黑背景
    },
    // 次要色：琥珀色，用于警告和强调
    secondary: {
      main: "#ffb000", // 琥珀色
      light: "#ffcc33", // 更亮的琥珀色
      dark: "#cc8800", // 稍深的琥珀色
      contrastText: "#0a0a0a",
    },
    // 背景色：深黑，避免纯黑以支持scanlines
    background: {
      default: "#0a0a0a", // 深黑主背景
      paper: "#1a1a1a", // 稍浅的面板背景
    },
    // 文本颜色：高对比度绿色主题
    text: {
      primary: "#33ff00", // 亮绿色文本
      secondary: "#1f521f", // 暗绿色次要文本
      disabled: "#0f290f", // 更暗的禁用文本
    },
    // 分割线：暗绿色
    divider: "#1f521f",
    // 动作按钮：使用绿色系
    action: {
      active: "#33ff00",
      hover: "rgba(51, 255, 0, 0.08)",
      selected: "rgba(51, 255, 0, 0.16)",
      disabled: "rgba(31, 82, 31, 0.3)",
      disabledBackground: "rgba(31, 82, 31, 0.12)",
    },
    // 语义化颜色
    success: {
      main: "#33ff00", // 使用主绿色表示成功
      light: "#66ff33",
      dark: "#00cc00",
    },
    error: {
      main: "#ff3333", // 亮红色错误
      light: "#ff6666",
      dark: "#cc0000",
    },
    warning: {
      main: "#ffb000", // 琥珀色警告
      light: "#ffcc33",
      dark: "#cc8800",
    },
    info: {
      main: "#33ff00", // 使用主绿色表示信息
      light: "#66ff33",
      dark: "#00cc00",
    },
  },
};

export default terminal;
