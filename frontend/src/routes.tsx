import type { ComponentType } from "react";
import { Box, Typography, Button } from "@mui/material";
import { createHashRouter, useRouteError, isRouteErrorResponse } from "react-router";
import App from "./pages/App.tsx";
import SubMainWindow from "./pages/SubMainWindow";
import Loading from "./components/Loading";

/** Shown while a lazy route chunk is loading (e.g. setting/command/tools windows). */
function RouteLoadingFallback() {
  return <Loading message="Loading..." />;
}

function RouteErrorFallback() {
  const error = useRouteError();
  const message = isRouteErrorResponse(error)
    ? `${error.status} ${error.statusText}`
    : error instanceof Error
      ? error.message
      : "页面加载失败";

  return (
    <Box
      sx={{
        height: "100%",
        display: "flex",
        flexDirection: "column",
        alignItems: "center",
        justifyContent: "center",
        gap: 2,
        p: 3,
      }}
    >
      <Typography variant="h6">页面加载失败</Typography>
      <Typography variant="body2" color="text.secondary">
        {message}
      </Typography>
      <Button variant="outlined" onClick={() => window.location.reload()}>
        重新加载
      </Button>
    </Box>
  );
}

const lazyRoute = (importer: () => Promise<{ default: ComponentType }>) => ({
  HydrateFallback: RouteLoadingFallback,
  ErrorBoundary: RouteErrorFallback,
  lazy: async () => {
    const { default: Component } = await importer();
    return { Component };
  },
});

export default createHashRouter([
  {
    path: "/",
    Component: App,
  },
  {
    path: "/setting",
    ...lazyRoute(() => import("./pages/Setting")),
  },
  {
    path: "/command",
    ...lazyRoute(() => import("./pages/Command")),
  },
  {
    path: "/tools",
    ...lazyRoute(() => import("./pages/Tools")),
  },
  {
    path: "/tools/:toolId",
    ...lazyRoute(() => import("./pages/ToolDetail")),
  },
  {
    // Secondary main window: keep eager to avoid a blank webview on open.
    path: "/submainwindow",
    Component: SubMainWindow,
  },
]);
