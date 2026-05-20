import { useEffect, useState } from "react";
import {
  Box,
  Typography,
  Switch,
  FormControlLabel,
  TextField,
  Select,
  MenuItem,
  Button,
  FormControl,
  InputLabel,
  Alert,
  CircularProgress,
  Divider,
  IconButton,
  InputAdornment,
} from "@mui/material";
import VisibilityIcon from "@mui/icons-material/Visibility";
import VisibilityOffIcon from "@mui/icons-material/VisibilityOff";
import { useAIConfigStore } from "../stores/aiConfig";

export default function AISettings() {
  const [showApiKey, setShowApiKey] = useState(false);

  const {
    config,
    isLoading,
    error,
    providers,
    loadConfig,
    saveConfig,
    resetConfig,
    loadProviders,
    updatePartialConfig,
    clearError,
  } = useAIConfigStore();

  useEffect(() => {
    loadConfig();
    loadProviders();
  }, [loadConfig, loadProviders]);

  // 处理启用/禁用切换
  const handleEnabledChange = (checked: boolean) => {
    updatePartialConfig({ enabled: checked });
  };

  // 处理提供商切换
  const handleProviderChange = (provider: string) => {
    updatePartialConfig({ provider });
  };

  // 处理保存配置
  const handleSave = async () => {
    if (config) {
      await saveConfig(config);
    }
  };

  // 处理恢复默认
  const handleReset = async () => {
    if (window.confirm("确定要恢复默认 AI 配置吗？")) {
      await resetConfig();
    }
  };



  if (!config) {
    return (
      <Box sx={{ display: "flex", justifyContent: "center", p: 4 }}>
        <CircularProgress />
      </Box>
    );
  }

  const currentProvider = providers.find((p) => p.name === config.provider);
  const needsApiKey = currentProvider?.needs_api_key ?? true;
  const needsEndpoint = currentProvider?.needs_endpoint ?? false;

  return (
    <Box sx={{ p: 3, maxWidth: 800 }}>
      <Typography variant="h5" gutterBottom>
        AI 设置
      </Typography>

      {error && (
        <Alert severity="error" sx={{ mb: 2 }} onClose={clearError}>
          {error}
        </Alert>
      )}

      {/* 启用/禁用开关 */}
      <FormControlLabel
        control={
          <Switch
            checked={config.enabled}
            onChange={(e) => handleEnabledChange(e.target.checked)}
          />
        }
        label="启用 AI 功能"
        sx={{ mb: 2, display: "block" }}
      />

      <Divider sx={{ my: 2 }} />

      {/* 提供商选择 */}
      <FormControl fullWidth sx={{ mb: 2 }}>
        <InputLabel id="provider-select-label">模型提供商</InputLabel>
        <Select
          labelId="provider-select-label"
          label="模型提供商"
          value={config.provider}
          onChange={(e) => handleProviderChange(e.target.value as string)}
          disabled={!config.enabled}
        >
          {providers.map((option) => (
            <MenuItem key={option.name} value={option.name}>
              {option.label}
            </MenuItem>
          ))}
        </Select>
      </FormControl>

      {/* 端点配置 */}
      {needsEndpoint && (
        <Box sx={{ mb: 2 }}>
          <TextField
            fullWidth
            label="API 端点地址"
            value={config.endpoint}
            onChange={(e) =>
              updatePartialConfig({ endpoint: e.target.value })
            }
            disabled={!config.enabled}
            sx={{ mb: 1 }}
            helperText="Ollama 填 http://localhost:11434；其他兼容/OpenAI 兼容提供商填写对应 API 地址"
          />
        </Box>
      )}

      {/* API Key 输入 */}
      {needsApiKey && (
        <TextField
          fullWidth
          label="API Key"
          type={showApiKey ? "text" : "password"}
          value={config.api_key || ""}
          onChange={(e) => updatePartialConfig({ api_key: e.target.value })}
          disabled={!config.enabled}
          sx={{ mb: 2 }}
          slotProps={{
            input: {
              endAdornment: (
                <InputAdornment position="end">
                  <IconButton
                    onClick={() => setShowApiKey(!showApiKey)}
                    edge="end"
                  >
                    {showApiKey ? <VisibilityOffIcon /> : <VisibilityIcon />}
                  </IconButton>
                </InputAdornment>
              ),
            }
          }}
        />
      )}

      {/* 模型输入 */}
      <TextField
        fullWidth
        label="模型"
        value={config.model}
        onChange={(e) => updatePartialConfig({ model: e.target.value.trim() })}
        disabled={!config.enabled}
        sx={{ mb: 2 }}
        helperText="请输入模型名称，如：gpt-4o、llama3.2、gemini-2.5-flash 等"
        required
        error={!config.model}
      />

      <Divider sx={{ my: 2 }} />

      {/* 高级选项 */}
      <Typography variant="h6" gutterBottom>
        高级选项
      </Typography>

      <TextField
        fullWidth
        label="Temperature"
        type="number"
        slotProps={{ htmlInput: { min: 0, max: 1, step: 0.1 } }}
        value={config.temperature}
        onChange={(e) =>
          updatePartialConfig({ temperature: parseFloat(e.target.value) })
        }
        disabled={!config.enabled}
        sx={{ mb: 2 }}
        helperText="数值越低，输出越确定（0.0-1.0）"
      />

      <TextField
        fullWidth
        label="Max Tokens"
        type="number"
        slotProps={{ htmlInput: { min: 100, max: 4000, step: 100 } }}
        value={config.max_tokens}
        onChange={(e) =>
          updatePartialConfig({ max_tokens: parseInt(e.target.value) })
        }
        disabled={!config.enabled}
        sx={{ mb: 2 }}
        helperText="生成文本的最大长度"
      />

      
      <Divider sx={{ my: 2 }} />

      {/* 操作按钮 */}
      <Box sx={{ display: "flex", gap: 2, mt: 3 }}>
        <Button
          variant="contained"
          onClick={handleSave}
          disabled={isLoading || !config.enabled}
        >
          {isLoading ? <CircularProgress size={20} /> : "保存配置"}
        </Button>

        <Button
          variant="outlined"
          color="secondary"
          onClick={handleReset}
          disabled={isLoading}
        >
          恢复默认设置
        </Button>
      </Box>
    </Box>
  );
}
