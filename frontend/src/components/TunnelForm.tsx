import React, { useState } from "react";
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Button,
  CircularProgress,
  Box,
} from "@mui/material";
import { useMessageStore } from "../stores/message";
import { SSHTunnelService } from "../../bindings/github.com/ilaziness/vexo/services";
import { parseCallServiceError } from "../func/service";

interface TunnelFormProps {
  open: boolean;
  tunnelType: "local" | "remote" | "dynamic";
  sessionID: string;
  onClose: () => void;
  onSuccess: () => void;
}

const TunnelForm: React.FC<TunnelFormProps> = ({
  open,
  tunnelType,
  sessionID,
  onClose,
  onSuccess,
}) => {
  const { errorMessage, successMessage } = useMessageStore();

  const [formData, setFormData] = useState({
    localPort: "",
    remoteAddr: "",
    // 动态转发可能只需要本地端口
  });
  const [loading, setLoading] = useState(false);

  // 获取隧道类型显示名称
  const getTunnelTypeName = () => {
    switch (tunnelType) {
      case "local":
        return "本地转发";
      case "remote":
        return "远程转发";
      case "dynamic":
        return "动态端口转发";
      default:
        return "隧道";
    }
  };

  // 验证远程地址格式 (host:port)
  const validateRemoteAddr = (addr: string): boolean => {
    if (!addr.trim()) return false;

    // 基本格式验证: host:port
    const parts = addr.split(":");
    if (parts.length !== 2) return false;

    const [host, portStr] = parts;
    if (!host.trim()) return false;

    const port = Number.parseInt(portStr, 10);
    if (Number.isNaN(port) || port < 1 || port > 65535) return false;

    return true;
  };

  // 验证本地端口
  const validateLocalPort = (portStr: string): boolean => {
    if (!portStr.trim()) return false;

    const port = Number.parseInt(portStr, 10);
    return !Number.isNaN(port) && port >= 1 && port <= 65535;
  };

  // 表单验证
  const isFormValid = () => {
    if (!validateLocalPort(formData.localPort)) return false;

    // 动态转发可能不需要远程地址
    if (tunnelType === "dynamic") {
      return true; // 动态转发可能只需要本地端口
    }

    return validateRemoteAddr(formData.remoteAddr);
  };

  // 处理表单提交
  const handleSubmit = async () => {
    if (!isFormValid()) return;

    setLoading(true);
    try {
      const localPort = Number.parseInt(formData.localPort, 10);

      // 目前只实现本地转发
      if (tunnelType === "local") {
        await SSHTunnelService.StartLocal(
          sessionID,
          localPort,
          formData.remoteAddr,
        );
        successMessage(`${getTunnelTypeName()}启动成功`);
      } else {
        // 远程转发和动态转发预留
        throw new Error(`${getTunnelTypeName()}功能暂未实现`);
      }

      onSuccess();
      onClose();
      setFormData({ localPort: "", remoteAddr: "" });
    } catch (err) {
      const errorMsg =
        err instanceof Error ? parseCallServiceError(err) : "启动隧道失败";
      errorMessage(`${getTunnelTypeName()}失败: ${errorMsg}`);
    } finally {
      setLoading(false);
    }
  };

  // 渲染表单字段
  const renderFormFields = () => {
    let localPortHelperText = "1-65535 之间的端口号";
    if (!formData.localPort.trim()) {
      localPortHelperText = "请输入本地端口";
    } else if (validateLocalPort(formData.localPort) === false) {
      localPortHelperText = "端口必须是 1-65535 之间的数字";
    }

    switch (tunnelType) {
      case "local":
        // 提取本地端口 helperText 逻辑
        return (
          <>
            <TextField
              label="本地端口"
              value={formData.localPort}
              onChange={(e) =>
                setFormData({ ...formData, localPort: e.target.value })
              }
              placeholder="例如: 8080"
              type="number"
              slotProps={{
                htmlInput: {
                  min: 1,
                  max: 65535,
                },
              }}
              helperText={localPortHelperText}
              error={
                formData.localPort.trim() !== "" &&
                validateLocalPort(formData.localPort) === false
              }
              fullWidth
              required
            />
            {/*
              提取远程地址 helperText 逻辑
            */}
            {(() => {
              let remoteAddrHelperText = "格式: host:port";
              if (formData.remoteAddr.trim() === "") {
                remoteAddrHelperText = "请输入远程地址";
              } else if (validateRemoteAddr(formData.remoteAddr) === false) {
                remoteAddrHelperText = "格式不正确，应为 host:port";
              }
              return (
                <TextField
                  label="远程地址"
                  value={formData.remoteAddr}
                  onChange={(e) =>
                    setFormData({ ...formData, remoteAddr: e.target.value })
                  }
                  placeholder="例如: 127.0.0.1:80 或 example.com:443"
                  helperText={remoteAddrHelperText}
                  error={
                    formData.remoteAddr.trim() !== "" &&
                    !validateRemoteAddr(formData.remoteAddr)
                  }
                  fullWidth
                  required
                />
              );
            })()}
          </>
        );

      case "remote":
        return (
          <Box sx={{ textAlign: "center", py: 4 }}>远程转发功能开发中...</Box>
        );

      case "dynamic":
        return (
          <Box sx={{ textAlign: "center", py: 4 }}>
            动态端口转发功能开发中...
          </Box>
        );

      default:
        return null;
    }
  };

  return (
    <Dialog open={open} onClose={onClose} maxWidth="sm" fullWidth>
      <DialogTitle>创建{getTunnelTypeName()}</DialogTitle>
      <DialogContent>
        <Box sx={{ display: "flex", flexDirection: "column", gap: 2, mt: 1 }}>
          {renderFormFields()}
        </Box>
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose} disabled={loading}>
          取消
        </Button>
        <Button
          onClick={handleSubmit}
          variant="contained"
          disabled={!isFormValid() || loading}
          startIcon={loading ? <CircularProgress size={16} /> : null}
        >
          {loading ? "启动中..." : "启动隧道"}
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default TunnelForm;
