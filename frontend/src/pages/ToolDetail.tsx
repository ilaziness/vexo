import React, { Suspense } from "react";
import { useParams } from "react-router";
import { Box, CircularProgress, Typography } from "@mui/material";
import ToolLayout from "../components/ToolLayout";

const toolComponents: Record<
  string,
  React.LazyExoticComponent<React.ComponentType>
> = {
  "port-check": React.lazy(() => import("../components/PortCheckTool")),
  encoder: React.lazy(() => import("../components/EncoderTool")),
  regex: React.lazy(() => import("../components/RegexTool")),
  hash: React.lazy(() => import("../components/HashTool")),
  timestamp: React.lazy(() => import("../components/TimestampTool")),
};

export default function ToolDetail() {
  const { toolId } = useParams<{ toolId: string }>();
  const ToolComponent = toolId ? toolComponents[toolId] : null;

  return (
    <ToolLayout>
      {ToolComponent ? (
        <Suspense
          fallback={
            <Box sx={{ display: "flex", justifyContent: "center", mt: 8 }}>
              <CircularProgress size={28} />
            </Box>
          }
        >
          <ToolComponent />
        </Suspense>
      ) : (
        <Box sx={{ textAlign: "center", mt: 8 }}>
          <Typography variant="h6" color="text.secondary">
            未知工具: {toolId}
          </Typography>
        </Box>
      )}
    </ToolLayout>
  );
}
