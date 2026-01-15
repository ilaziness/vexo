import type { ColorSystemOptions } from "@mui/material/styles";

const eyeCare: ColorSystemOptions = {
  palette: {
    mode: "light",
    primary: {
      main: "#5d4037",
      light: "#8b6b61",
      dark: "#321911",
      contrastText: "#ffffff",
    },
    secondary: {
      main: "#795548",
      light: "#a98274",
      dark: "#4b2c20",
      contrastText: "#ffffff",
    },
    background: { default: "#fdf6e3", paper: "#eee8d5" },
    text: { primary: "#423629", secondary: "#635240", disabled: "#9c8c74" },
    divider: "rgba(66, 54, 41, 0.12)",
    action: {
      active: "#423629",
      hover: "rgba(66, 54, 41, 0.04)",
      selected: "rgba(66, 54, 41, 0.08)",
      disabled: "rgba(66, 54, 41, 0.26)",
      disabledBackground: "rgba(66, 54, 41, 0.12)",
    },
    success: { main: "#556b2f" },
    error: { main: "#a0522d" },
    warning: { main: "#cd853f" },
    info: { main: "#5d4037" },
  },
};

export default eyeCare;
