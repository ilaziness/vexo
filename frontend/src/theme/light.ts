import type { ColorSystemOptions } from "@mui/material/styles";

const light: ColorSystemOptions = {
  palette: {
    mode: "light",
    // 主色调：清爽的蓝色，代表技术感和专业性
    primary: {
      main: "#1976d2", // 蓝色主色
      light: "#42a5f5",
      dark: "#1565c0",
      contrastText: "#ffffff",
    },
    // 次要色：柔和的青色
    secondary: {
      main: "#0288d1",
      light: "#03a9f4",
      dark: "#01579b",
      contrastText: "#ffffff",
    },
    // 背景色：简洁的浅色系
    background: {
      default: "#f0f0f0", // 主背景 - 更柔和的灰色，减少眩光
      paper: "#fefefe", // 卡片、面板背景 - 略带灰调的白色
    },
    // 文本颜色
    text: {
      primary: "rgba(0, 0, 0, 0.87)",
      secondary: "rgba(0, 0, 0, 0.6)",
      disabled: "rgba(0, 0, 0, 0.38)",
    },
    // 分割线
    divider: "rgba(0, 0, 0, 0.12)",
    // 动作按钮
    action: {
      active: "rgba(0, 0, 0, 0.54)",
      hover: "rgba(0, 0, 0, 0.04)",
      selected: "rgba(0, 0, 0, 0.08)",
      disabled: "rgba(0, 0, 0, 0.26)",
      disabledBackground: "rgba(0, 0, 0, 0.12)",
    },
    // 成功、错误、警告、信息
    success: {
      main: "#2e7d32",
      light: "#4caf50",
      dark: "#1b5e20",
    },
    error: {
      main: "#d32f2f",
      light: "#ef5350",
      dark: "#c62828",
    },
    warning: {
      main: "#ed6c02",
      light: "#ff9800",
      dark: "#e65100",
    },
    info: {
      main: "#0288d1",
      light: "#03a9f4",
      dark: "#01579b",
    },
  },
};

export default light;
