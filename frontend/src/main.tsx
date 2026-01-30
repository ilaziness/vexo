import React from "react";
import ReactDOM from "react-dom/client";
import { ThemeProvider } from "@mui/material/styles";
import CssBaseline from "@mui/material/CssBaseline";
import { RouterProvider } from "react-router/dom";
import routes from "./routes.ts";
import theme from "./theme";
import { ReadConfig } from "../bindings/github.com/ilaziness/vexo/services/configservice";
import { Config } from "../bindings/github.com/ilaziness/vexo/services";
import useTerminalStore from "./stores/terminal";
import type { AppTheme } from "./types/ssh";

(async () => {
  try {
    const cfg = await ReadConfig();
    if (cfg) {
      // Terminal settings
      const t = (cfg as Config).Terminal;
      const updates: Partial<{
        fontFamily: string;
        fontSize: number;
        lineHeight: number;
      }> = {};

      if (t && typeof t.fontFamily === "string" && t.fontFamily)
        updates.fontFamily = t.fontFamily;
      if (t && typeof t.fontSize === "number" && t.fontSize > 0)
        updates.fontSize = t.fontSize;
      if (t && typeof t.lineHeight === "number" && t.lineHeight > 0)
        updates.lineHeight = t.lineHeight;

      if (Object.keys(updates).length > 0) {
        useTerminalStore.setState(updates as any);
      }

      // Sync terminal theme with global theme if present
      const globalTheme = (cfg as Config).General?.Theme;
      if (globalTheme && typeof globalTheme === "string") {
        useTerminalStore
          .getState()
          .syncWithGlobalTheme(globalTheme as AppTheme);
      }
    }
  } catch (err) {
    console.error("ReadConfig failed:", err);
  } finally {
    ReactDOM.createRoot(document.getElementById("root") as HTMLElement).render(
      <React.StrictMode>
        <ThemeProvider theme={theme} defaultMode="dark">
          <CssBaseline />
          <RouterProvider router={routes} />
        </ThemeProvider>
      </React.StrictMode>,
    );
  }
})();
