//import React from "react";
import ReactDOM from "react-dom/client";
import { createTheme, ThemeProvider } from "@mui/material/styles";
import CssBaseline from "@mui/material/CssBaseline";
import { RouterProvider } from "react-router/dom";
import routes from "./routes.ts";

const darkTheme = createTheme({
  palette: {
    mode: "dark",
  },
});

ReactDOM.createRoot(document.getElementById("root") as HTMLElement).render(
  //<React.StrictMode>
  <ThemeProvider theme={darkTheme}>
    <CssBaseline />
    <RouterProvider router={routes} />
  </ThemeProvider>,
  //</React.StrictMode>
);
