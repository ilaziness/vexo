import type { ColorSystemOptions } from "@mui/material/styles";

const deep: ColorSystemOptions = {
  palette: {
    mode: "dark",
    primary: {
      main: "#40a9ff",
      light: "#69c0ff",
      dark: "#096dd9",
      contrastText: "#ffffff",
    },
    secondary: {
      main: "#73d13d",
      light: "#95de64",
      dark: "#389e0d",
      contrastText: "#ffffff",
    },
    background: { default: "#000000", paper: "#141414" },
    text: {
      primary: "rgba(255, 255, 255, 0.85)",
      secondary: "rgba(255, 255, 255, 0.45)",
      disabled: "rgba(255, 255, 255, 0.25)",
    },
    divider: "#303030",
    action: {
      active: "#ffffff",
      hover: "rgba(255, 255, 255, 0.08)",
      selected: "rgba(255, 255, 255, 0.12)",
      disabled: "rgba(255, 255, 255, 0.3)",
      disabledBackground: "rgba(255, 255, 255, 0.12)",
    },
    success: { main: "#52c41a" },
    error: { main: "#ff4d4f" },
    warning: { main: "#faad14" },
    info: { main: "#1890ff" },
  },
};

export default deep;
