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

// 提取终端设置应用函数，降低认知复杂度
function applyTerminalSettings(t: Config["Terminal"]): void {
  const updates: Partial<{
    fontFamily: string;
    fontSize: number;
    lineHeight: number;
  }> = {};

  if (t?.fontFamily) {
    updates.fontFamily = t.fontFamily;
  }
  if (t?.fontSize && t.fontSize > 0) {
    updates.fontSize = t.fontSize;
  }
  if (t?.lineHeight && t.lineHeight > 0) {
    updates.lineHeight = t.lineHeight;
  }

  if (Object.keys(updates).length > 0) {
    useTerminalStore.setState(updates as any);
  }
}

// 提取主题同步函数
function syncTerminalTheme(globalTheme: string | undefined): void {
  if (globalTheme) {
    useTerminalStore.getState().syncWithGlobalTheme(globalTheme as AppTheme);
  }
}

try {
  const config = await ReadConfig();
  if (config) {
    applyTerminalSettings(config.Terminal);
    syncTerminalTheme(config.General?.Theme);
  }
} catch (err) {
  console.error("ReadConfig failed:", err);
}

ReactDOM.createRoot(document.getElementById("root") as HTMLElement).render(
  <React.StrictMode>
    <ThemeProvider theme={theme} defaultMode="dark">
      <CssBaseline />
      <RouterProvider router={routes} />
    </ThemeProvider>
  </React.StrictMode>,
);
