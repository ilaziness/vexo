import type { ColorSystemOptions } from "@mui/material/styles";

const light: ColorSystemOptions = {
  palette: {
    mode: "light",
    primary: {
      main: "#6442d6", // 紫色主色调
      light: "#9675dd", // 更亮的紫色变体
      dark: "#4a2da0",  // 更深的紫色变体
      contrastText: "#ffffff",
    },
    secondary: {
      main: "#7c4dff", // 紫色系的次要色，稍微亮一些
      light: "#b794f6", // 更亮的紫色
      dark: "#5e35b1",  // 更深的紫色
      contrastText: "#ffffff",
    },
    background: {
      default: "#fafafa",
      paper: "#ffffff",
    },
    text: {
      primary: "rgba(0, 0, 0, 0.87)",
      secondary: "rgba(0, 0, 0, 0.6)",
      disabled: "rgba(0, 0, 0, 0.38)",
    },
    divider: "rgba(0, 0, 0, 0.12)",
    action: {
      active: "rgba(0, 0, 0, 0.54)",
      hover: "rgba(0, 0, 0, 0.04)",
      selected: "rgba(0, 0, 0, 0.08)",
      disabled: "rgba(0, 0, 0, 0.26)",
      disabledBackground: "rgba(0, 0, 0, 0.12)",
    },
    success: {
      main: "#4caf50",
      light: "#81c784",
      dark: "#388e3c",
    },
    error: {
      main: "#f44336",
      light: "#e57373",
      dark: "#d32f2f",
    },
    warning: {
      main: "#ff9800",
      light: "#ffb74d",
      dark: "#f57c00",
    },
    info: {
      main: "#2196f3",
      light: "#64b5f6",
      dark: "#1976d2",
    },
  },
};

export default light;
