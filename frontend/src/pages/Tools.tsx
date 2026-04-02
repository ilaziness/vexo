import { useEffect } from "react";
import {
  Box,
  Card,
  CardContent,
  Grid,
  Typography,
} from "@mui/material";
import NetworkCheckIcon from "@mui/icons-material/NetworkCheck";
import CodeIcon from "@mui/icons-material/Code";
import EditNoteIcon from "@mui/icons-material/EditNote";
import { useNavigate } from "react-router";
import { useToolStore } from "../stores/tool";
import { Tool } from "../types/tool";
import OpBar from "../components/OpBar";

// 图标映射
const iconMap: Record<string, React.ElementType> = {
  NetworkCheck: NetworkCheckIcon,
  Code: CodeIcon,
  RegularExpression: CodeIcon,
  EditNote: EditNoteIcon,
};

interface ToolCardProps {
  readonly tool: Tool;
  readonly onClick: () => void;
}

function ToolCard({ tool, onClick }: ToolCardProps) {
  const IconComponent = iconMap[tool.icon] || CodeIcon;

  return (
    <Card
      onClick={onClick}
      sx={{
        cursor: "pointer",
        height: "100%",
        transition: "all 0.2s ease-in-out",
        "&:hover": {
          transform: "translateY(-4px)",
          boxShadow: (theme) => theme.shadows[8],
        },
      }}
    >
      <CardContent sx={{ textAlign: "center", p: 3 }}>
        <Box
          sx={{
            width: 64,
            height: 64,
            borderRadius: "50%",
            bgcolor: "primary.main",
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
            margin: "0 auto 16px",
          }}
        >
          <IconComponent sx={{ fontSize: 32, color: "white" }} />
        </Box>
        <Typography variant="h6" gutterBottom>
          {tool.name}
        </Typography>
        <Typography variant="body2" color="text.secondary">
          {tool.description}
        </Typography>
        <Typography
          variant="caption"
          color="primary"
          sx={{ display: "block", mt: 1 }}
        >
          {tool.category}
        </Typography>
      </CardContent>
    </Card>
  );
}

export default function Tools() {
  const navigate = useNavigate();
  const { tools, fetchTools } = useToolStore();

  useEffect(() => {
    fetchTools();
  }, [fetchTools]);

  const handleToolClick = (toolId: string) => {
    navigate(`/tools/${toolId}`);
  };

  return (
    <Box
      sx={{
        height: "100vh",
        display: "flex",
        flexDirection: "column",
        bgcolor: "background.default",
      }}
    >
      {/* 操作栏 */}
      <Box sx={{ height: 40 }}>
        <OpBar />
      </Box>

      {/* 标题栏 */}
      <Box
        sx={{
          p: 2,
          borderBottom: 1,
          borderColor: "divider",
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
          "--wails-draggable": "drag",
        }}
      >
        <Typography variant="h5" component="h1">
          运维工具箱
        </Typography>
      </Box>

      {/* 内容区域 */}
      <Box sx={{ flex: 1, overflow: "auto", p: 4 }}>
        <Grid container spacing={3} justifyContent="center">
          {tools.map((tool) => (
            <Grid size={{ xs: 12, sm: 6, md: 4, lg: 3 }} key={tool.id}>
              <ToolCard
                tool={tool}
                onClick={() => handleToolClick(tool.id)}
              />
            </Grid>
          ))}
        </Grid>
      </Box>
    </Box>
  );
}
