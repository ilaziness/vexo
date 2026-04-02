import { useParams } from "react-router";
import ToolLayout from "../components/ToolLayout";
import PortCheckTool from "../components/PortCheckTool";
import EncoderTool from "../components/EncoderTool";
import RegexTool from "../components/RegexTool";
import HashTool from "../components/HashTool";
import TimestampTool from "../components/TimestampTool";
import { Box, Typography } from "@mui/material";

const toolComponents: Record<string, React.ComponentType> = {
  "port-check": PortCheckTool,
  encoder: EncoderTool,
  regex: RegexTool,
  hash: HashTool,
  timestamp: TimestampTool,
};

export default function ToolDetail() {
  const { toolId } = useParams<{ toolId: string }>();
  const ToolComponent = toolId ? toolComponents[toolId] : null;

  return (
    <ToolLayout>
      {ToolComponent ? (
        <ToolComponent />
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
