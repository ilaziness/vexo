import React, { useState, useEffect } from "react";
import {
  Box,
  TextField,
  Button,
  Typography,
  Paper,
  Stack,
  Divider,
} from "@mui/material";
import { SSHBookmark } from "../../bindings/github.com/ilaziness/vexo/services";

interface BookmarkFormProps {
  bookmark: SSHBookmark | null;
  onSave: (bookmark: SSHBookmark) => void;
  onTestConnection: (bookmark: SSHBookmark) => void;
  onSaveAndConnect: (bookmark: SSHBookmark) => void;
}

const BookmarkForm: React.FC<BookmarkFormProps> = ({
  bookmark,
  onSave,
  onTestConnection,
  onSaveAndConnect,
}) => {
  const [formData, setFormData] = useState<SSHBookmark>({
    id: "",
    title: "",
    group_name: "default",
    host: "",
    port: 22,
    private_key: "",
    user: "",
    password: "",
  });

  useEffect(() => {
    if (bookmark) {
      setFormData({ ...bookmark });
    } else {
      setFormData({
        id: "",
        title: "",
        group_name: "default",
        host: "",
        port: 22,
        private_key: "",
        user: "",
        password: "",
      });
    }
  }, [bookmark]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData((prev) => ({
      ...prev,
      [name]: name === "port" ? Number(value) : value,
    }));
  };

  const handleSave = () => {
    if (formData) {
      onSave(formData);
    }
  };

  const handleTestConnection = () => {
    if (formData) {
      onTestConnection(formData);
    }
  };

  const handleSaveAndConnect = () => {
    if (formData) {
      onSaveAndConnect(formData);
    }
  };

  // 横向表单项组件
  const FormRow: React.FC<{
    label: string;
    children: React.ReactNode;
  }> = ({ label, children }) => (
    <Box
      sx={{
        display: "flex",
        alignItems: "center",
        mb: 2,
        "&:last-child": { mb: 0 },
      }}
    >
      <Typography
        sx={{
          width: 120,
          flexShrink: 0,
          fontWeight: 500,
          fontSize: "0.95rem",
        }}
      >
        {label}
      </Typography>
      <Box sx={{ flex: 1 }}>{children}</Box>
    </Box>
  );

  if (!bookmark) {
    return (
      <Box
        sx={{
          height: "100%",
          display: "flex",
          flexDirection: "column",
          alignItems: "center",
          justifyContent: "center",
          p: 3,
        }}
      >
        <Typography variant="h6" color="text.secondary" sx={{ mb: 1 }}>
          请选择一个书签进行编辑
        </Typography>
        <Typography variant="body2" color="text.disabled">
          从左侧列表中选择或创建一个新书签
        </Typography>
      </Box>
    );
  }

  return (
    <Box
      sx={{
        height: "100%",
        display: "flex",
        flexDirection: "column",
      }}
    >
      {/* 标题栏 */}
      <Box
        sx={{
          p: 2,
          pb: 1.5,
          borderBottom: 1,
          borderColor: "divider",
        }}
      >
        <Typography variant="h6" sx={{ fontWeight: 600, fontSize: "1.1rem" }}>
          书签详情
        </Typography>
      </Box>

      {/* 表单内容 */}
      <Box sx={{ flex: 1, overflowY: "auto", p: 3 }}>
        <Paper sx={{ p: 3 }} elevation={1}>
          <Stack spacing={2.5}>
            {/* 基本信息 */}
            <Box>
              <Typography
                variant="subtitle2"
                sx={{ mb: 2, fontWeight: 600, color: "primary.main" }}
              >
                基本信息
              </Typography>
              <Stack spacing={2}>
                <FormRow label="书签名称">
                  <TextField
                    fullWidth
                    size="small"
                    name="title"
                    value={formData.title}
                    onChange={handleChange}
                    placeholder="请输入书签名称"
                  />
                </FormRow>
                <FormRow label="分组">
                  <TextField
                    fullWidth
                    size="small"
                    name="group_name"
                    value={formData.group_name}
                    onChange={handleChange}
                    placeholder="请输入分组名称"
                  />
                </FormRow>
              </Stack>
            </Box>

            <Divider />

            {/* 连接信息 */}
            <Box>
              <Typography
                variant="subtitle2"
                sx={{ mb: 2, fontWeight: 600, color: "primary.main" }}
              >
                连接信息
              </Typography>
              <Stack spacing={2}>
                <FormRow label="主机地址">
                  <TextField
                    fullWidth
                    size="small"
                    name="host"
                    value={formData.host}
                    onChange={handleChange}
                    placeholder="例如: 192.168.1.100 或 example.com"
                  />
                </FormRow>
                <FormRow label="端口">
                  <TextField
                    fullWidth
                    size="small"
                    name="port"
                    type="number"
                    value={formData.port}
                    onChange={handleChange}
                    placeholder="默认: 22"
                  />
                </FormRow>
                <FormRow label="用户名">
                  <TextField
                    fullWidth
                    size="small"
                    name="user"
                    value={formData.user}
                    onChange={handleChange}
                    placeholder="SSH 登录用户名"
                  />
                </FormRow>
              </Stack>
            </Box>

            <Divider />

            {/* 认证信息 */}
            <Box>
              <Typography
                variant="subtitle2"
                sx={{ mb: 2, fontWeight: 600, color: "primary.main" }}
              >
                认证信息
              </Typography>
              <Stack spacing={2}>
                <FormRow label="密码">
                  <TextField
                    fullWidth
                    size="small"
                    name="password"
                    type="password"
                    value={formData.password}
                    onChange={handleChange}
                    placeholder="密码认证（可选）"
                  />
                </FormRow>
                <FormRow label="密钥路径">
                  <TextField
                    fullWidth
                    size="small"
                    name="private_key"
                    value={formData.private_key}
                    onChange={handleChange}
                    placeholder="私钥文件路径（可选）"
                  />
                </FormRow>
              </Stack>
            </Box>
          </Stack>
        </Paper>
      </Box>

      {/* 底部操作按钮 */}
      <Box
        sx={{
          borderTop: 1,
          borderColor: "divider",
          p: 2,
          display: "flex",
          justifyContent: "flex-end",
          gap: 1,
          backgroundColor: "background.paper",
        }}
      >
        <Button variant="outlined" color="secondary" onClick={handleTestConnection}>
          测试连接
        </Button>
        <Button variant="outlined" onClick={handleSave}>
          保存
        </Button>
        <Button variant="contained" onClick={handleSaveAndConnect}>
          保存并连接
        </Button>
      </Box>
    </Box>
  );
};

export default BookmarkForm;
