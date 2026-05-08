import { useEffect } from "react";
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogContentText,
  DialogActions,
  Button,
  Box,
} from "@mui/material";
import SettingsIcon from "@mui/icons-material/Settings";
import AutoAwesomeIcon from '@mui/icons-material/AutoAwesome';
import { useAIConfigStore } from "../stores/aiConfig";
import { ConfigService } from "../../bindings/github.com/ilaziness/vexo/services";

interface AIConfigDialogProps {
  open: boolean;
  onClose: () => void;
}

export default function AIConfigDialog({ open, onClose }: AIConfigDialogProps) {
  const { providers, loadProviders } = useAIConfigStore();

  // 加载提供商列表
  useEffect(() => {
    if (open) {
      loadProviders();
    }
  }, [open, loadProviders]);

  const handleGoToSettings = () => {
    onClose();
    ConfigService.ShowWindow();
  };

  return (
    <Dialog open={open} onClose={onClose} maxWidth="sm" fullWidth>
      <DialogTitle>
        <Box sx={{ display: "flex", alignItems: "center", gap: 1 }}>
          <AutoAwesomeIcon color="primary" />
          AI 功能未启用
        </Box>
      </DialogTitle>
      <DialogContent>
        <DialogContentText>
          要使用 AI 命令助手功能，您需要先配置 AI 服务。
        </DialogContentText>
        <DialogContentText sx={{ mt: 1 }} color="text.secondary">
          支持的提供商：
        </DialogContentText>
        <Box
          component="ul"
          sx={{
            mt: 0.5,
            pl: 2,
            color: "text.secondary",
            fontSize: "0.875rem",
          }}
        >
          {providers.map((provider) => (
            <li key={provider.name}>
              {provider.label}
              {!provider.needs_api_key && "（无需 API Key）"}
            </li>
          ))}
        </Box>
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose}>取消</Button>
        <Button
          onClick={handleGoToSettings}
          variant="contained"
          startIcon={<SettingsIcon />}
        >
          前往设置
        </Button>
      </DialogActions>
    </Dialog>
  );
}
