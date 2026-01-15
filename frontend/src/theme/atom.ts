import type { ColorSystemOptions } from "@mui/material/styles";

const atom: ColorSystemOptions = {
  palette: {
    mode: "dark",
    primary: {
      main: "#61afef",
      light: "#98d5fd",
      dark: "#2889d6",
      contrastText: "#000000",
    },
    secondary: {
      main: "#c678dd",
      light: "#e4aef5",
      dark: "#9447aa",
      contrastText: "#ffffff",
    },
    background: { default: "#282c34", paper: "#21252b" },
    text: { primary: "#abb2bf", secondary: "#8a919d", disabled: "#5c6370" },
    divider: "rgba(255, 255, 255, 0.12)",
    action: {
      active: "#abb2bf",
      hover: "rgba(171, 178, 191, 0.08)",
      selected: "rgba(171, 178, 191, 0.16)",
      disabled: "rgba(171, 178, 191, 0.3)",
      disabledBackground: "rgba(171, 178, 191, 0.12)",
    },
    success: { main: "#98c379" },
    error: { main: "#e06c75" },
    warning: { main: "#e5c07b" },
    info: { main: "#61afef" },
  },
};

export default atom;
