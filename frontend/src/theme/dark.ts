import type { ColorSystemOptions } from "@mui/material/styles";

const dark: ColorSystemOptions = {
  palette: {
    mode: "dark",
    // 主色调：柔和的蓝色，护眼且专业
    primary: {
      main: "#90caf9", // 偏浅的蓝色，在暗色背景下更清晰
      light: "#bbdefb",
      dark: "#42a5f5",
      contrastText: "rgba(0, 0, 0, 0.87)",
    },
    // 次要色：青色系
    secondary: {
      main: "#26c6da",
      light: "#4dd0e1",
      dark: "#0097a7",
      contrastText: "rgba(0, 0, 0, 0.87)",
    },
    // 背景色：深色系，避免纯黑，更护眼
    background: {
      default: "#0a0e27", // 深蓝黑色主背景
      paper: "#1a1f3a", // 稍浅的卡片背景
    },
    // 文本颜色：高对比度，易读
    text: {
      primary: "#e3f2fd", // 浅蓝白色，更柔和
      secondary: "rgba(227, 242, 253, 0.7)",
      disabled: "rgba(227, 242, 253, 0.5)",
    },
    // 分割线
    divider: "rgba(227, 242, 253, 0.12)",
    // 动作按钮
    action: {
      active: "#90caf9",
      hover: "rgba(144, 202, 249, 0.08)",
      selected: "rgba(144, 202, 249, 0.16)",
      disabled: "rgba(227, 242, 253, 0.3)",
      disabledBackground: "rgba(227, 242, 253, 0.12)",
    },
    // 成功、错误、警告、信息
    success: {
      main: "#66bb6a",
      light: "#81c784",
      dark: "#388e3c",
    },
    error: {
      main: "#f44336",
      light: "#e57373",
      dark: "#d32f2f",
    },
    warning: {
      main: "#ffa726",
      light: "#ffb74d",
      dark: "#f57c00",
    },
    info: {
      main: "#29b6f6",
      light: "#4fc3f7",
      dark: "#0288d1",
    },
  },
};

export default dark;
