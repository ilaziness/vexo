import React, { useState, useEffect } from "react";
import {
  Box,
  Typography,
  Paper,
  TextField,
  Button,
  Stack,
  FormControlLabel,
  Switch,
} from "@mui/material";
import FormRow from "../FormRow";
import { useMessageStore } from "../../stores/message";
import {
  Config,
  ConfigService,
} from "../../../bindings/github.com/ilaziness/vexo/services";
import { parseCallServiceError } from "../../func/service";

interface SSHSettingsProps {
  config: Config["SSH"];
}

const defaultSSHConfig: Config["SSH"] = {
  serverAliveInterval: 30,
  dialTimeoutSec: 30,
  autoReconnect: true,
};

const SSHSettings: React.FC<SSHSettingsProps> = ({ config }) => {
  const [localConfig, setLocalConfig] = useState({
    ...defaultSSHConfig,
    ...config,
  });
  const [saving, setSaving] = useState(false);
  const { errorMessage, successMessage } = useMessageStore();

  useEffect(() => {
    setLocalConfig({ ...defaultSSHConfig, ...config });
  }, [config]);

  const handleChange = <K extends keyof Config["SSH"]>(
    field: K,
    value: Config["SSH"][K],
  ) => {
    setLocalConfig((prev) => ({ ...prev, [field]: value }));
  };

  const validateConfig = (): string[] => {
    const errors: string[] = [];
    if (localConfig.serverAliveInterval < 0) {
      errors.push("心跳间隔不能为负数（0 表示关闭）");
    }
    if (!localConfig.dialTimeoutSec || localConfig.dialTimeoutSec < 1) {
      errors.push("拨号超时必须大于等于 1 秒");
    }
    return errors;
  };

  const handleSave = async () => {
    const errors = validateConfig();
    if (errors.length > 0) {
      errorMessage(errors.join("\n"));
      return;
    }

    setSaving(true);
    try {
      await ConfigService.SaveSSHConfig(localConfig);
      successMessage("连接配置保存成功");
    } catch (error) {
      console.error("Failed to save SSH config:", error);
      errorMessage(parseCallServiceError(error) || "连接配置保存失败");
    } finally {
      setSaving(false);
    }
  };

  return (
    <Box>
      <Typography variant="h5" gutterBottom sx={{ mb: 3, fontWeight: 600 }}>
        连接设置
      </Typography>
      <Paper sx={{ p: 2 }} elevation={1}>
        <Stack spacing={1.5}>
          <FormRow label="心跳间隔（秒）">
            <TextField
              fullWidth
              size="small"
              type="number"
              slotProps={{
                htmlInput: { min: 0, step: 1 },
              }}
              value={localConfig.serverAliveInterval ?? 0}
              onChange={(e) => {
                const value =
                  e.target.value === ""
                    ? 0
                    : Number.parseInt(e.target.value, 10);
                handleChange(
                  "serverAliveInterval",
                  Number.isNaN(value) ? 0 : value,
                );
              }}
              helperText="对应 ServerAliveInterval，0 表示关闭保活"
            />
          </FormRow>
          <FormRow label="拨号超时（秒）">
            <TextField
              fullWidth
              size="small"
              type="number"
              slotProps={{
                htmlInput: { min: 1, step: 1 },
              }}
              value={localConfig.dialTimeoutSec ?? 30}
              onChange={(e) => {
                const value =
                  e.target.value === ""
                    ? 0
                    : Number.parseInt(e.target.value, 10);
                handleChange(
                  "dialTimeoutSec",
                  Number.isNaN(value) ? 0 : value,
                );
              }}
            />
          </FormRow>
          <FormRow label="自动重连">
            <FormControlLabel
              control={
                <Switch
                  checked={Boolean(localConfig.autoReconnect)}
                  onChange={(e) =>
                    handleChange("autoReconnect", e.target.checked)
                  }
                />
              }
              label={localConfig.autoReconnect ? "已开启" : "已关闭"}
            />
          </FormRow>
        </Stack>
        <Box sx={{ display: "flex", justifyContent: "flex-end", mt: 3 }}>
          <Button variant="contained" onClick={handleSave} loading={saving}>
            保存
          </Button>
        </Box>
      </Paper>
    </Box>
  );
};

export default SSHSettings;
