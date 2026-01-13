import React, { useState, useEffect } from "react";
import {
  Box,
  TextField,
  Button,
  Typography,
  Paper,
  Stack,
  Divider,
  MenuItem,
  Select,
  FormControl,
  IconButton,
  InputAdornment,
} from "@mui/material";
import MoreHorizIcon from "@mui/icons-material/MoreHoriz";
import { SSHBookmark } from "../../bindings/github.com/ilaziness/vexo/services";
import { CommonService } from "../../bindings/github.com/ilaziness/vexo/services";
import { useMessageStore } from "../stores/message";
import FormRow from "./FormRow";

interface BookmarkFormProps {
  bookmark: SSHBookmark | null;
  groupNames: string[]; // 添加分组列表
  onSave: (bookmark: SSHBookmark) => Promise<void>;
  onTestConnection: (bookmark: SSHBookmark) => Promise<void>;
  onSaveAndConnect: (bookmark: SSHBookmark) => Promise<void>;
}

const BookmarkForm: React.FC<BookmarkFormProps> = ({
  bookmark,
  groupNames,
  onSave,
  onTestConnection,
  onSaveAndConnect,
}) => {
  const { errorMessage } = useMessageStore();

  const [formData, setFormData] = useState<SSHBookmark>({
    id: "",
    title: "",
    group_name: "default",
    host: "",
    port: 22,
    private_key: "",
    private_key_password: "",
    user: "",
    password: "",
  });

  const [isLoading, setIsLoading] = useState(false);

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
        private_key_password: "",
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

  const validateForm = (): boolean => {
    if (!formData.title.trim()) {
      errorMessage("书签名称不能为空");
      return false;
    }
    if (!formData.host.trim()) {
      errorMessage("主机地址不能为空");
      return false;
    }
    if (!formData.port || formData.port <= 0) {
      errorMessage("端口必须是有效的正整数");
      return false;
    }
    if (!formData.user.trim()) {
      errorMessage("用户名不能为空");
      return false;
    }
    if (!formData.password.trim() && !formData.private_key.trim()) {
      errorMessage("密码和密钥文件至少填写一项");
      return false;
    }
    return true;
  };

  const handleSave = async () => {
    if (validateForm()) {
      setIsLoading(true);
      try {
        await onSave(formData);
      } finally {
        setIsLoading(false);
      }
    }
  };

  const handleTestConnection = async () => {
    if (validateForm()) {
      setIsLoading(true);
      try {
        await onTestConnection(formData);
      } finally {
        setIsLoading(false);
      }
    }
  };

  const handleSaveAndConnect = async () => {
    if (validateForm()) {
      setIsLoading(true);
      try {
        await onSaveAndConnect(formData);
      } finally {
        setIsLoading(false);
      }
    }
  };

  // 处理分组选择变化
  const handleGroupChange = (event: any) => {
    setFormData((prev) => ({
      ...prev,
      group_name: event.target.value,
    }));
  };

  // 处理选择文件
  const handleSelectFile = async () => {
    try {
      const selectedPath = await CommonService.SelectFile();
      if (selectedPath) {
        setFormData((prev) => ({
          ...prev,
          private_key: selectedPath,
        }));
      }
    } catch (error) {
      console.error("选择文件失败:", error);
    }
  };

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
      <Box sx={{ flex: 1, overflowY: "auto", p: 2 }}>
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
                <FormRow label="书签名称" labelWidth={120}>
                  <TextField
                    fullWidth
                    size="small"
                    name="title"
                    value={formData.title}
                    onChange={handleChange}
                    placeholder="请输入书签名称"
                  />
                </FormRow>
                <FormRow label="分组" labelWidth={120}>
                  <FormControl fullWidth size="small">
                    <Select
                      value={formData.group_name}
                      onChange={handleGroupChange}
                      displayEmpty
                    >
                      {groupNames.length === 0 && (
                        <MenuItem value="default">default</MenuItem>
                      )}
                      {groupNames.map((groupName) => (
                        <MenuItem key={groupName} value={groupName}>
                          {groupName}
                        </MenuItem>
                      ))}
                    </Select>
                  </FormControl>
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
                <FormRow label="主机地址" labelWidth={120}>
                  <TextField
                    fullWidth
                    size="small"
                    name="host"
                    value={formData.host}
                    onChange={handleChange}
                    placeholder="例如: 192.168.1.100 或 example.com"
                  />
                </FormRow>
                <FormRow label="端口" labelWidth={120}>
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
                <FormRow label="用户名" labelWidth={120}>
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
                <FormRow label="密码" labelWidth={120}>
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
                <FormRow label="密钥文件" labelWidth={120}>
                  <TextField
                    fullWidth
                    size="small"
                    name="private_key"
                    value={formData.private_key}
                    placeholder="私钥文件路径（可选）"
                    slotProps={{
                      input: {
                        readOnly: true,
                        endAdornment: (
                          <InputAdornment position="end">
                            <IconButton
                              edge="end"
                              onClick={handleSelectFile}
                              size="small"
                              title="选择文件"
                            >
                              <MoreHorizIcon />
                            </IconButton>
                          </InputAdornment>
                        ),
                      },
                    }}
                  />
                </FormRow>
                <FormRow label="密钥密码" labelWidth={120}>
                  <TextField
                    fullWidth
                    size="small"
                    name="private_key_password"
                    type="password"
                    value={formData.private_key_password}
                    onChange={handleChange}
                    placeholder="私钥密码（可选）"
                  />
                </FormRow>
              </Stack>
            </Box>

            <Divider />

            {/* 操作按钮 */}
            <Box
              sx={{
                display: "flex",
                justifyContent: "flex-end",
                gap: 1,
                pt: 1,
              }}
            >
              <Button
                variant="outlined"
                color="secondary"
                onClick={handleTestConnection}
                loading={isLoading}
              >
                测试连接
              </Button>
              <Button
                variant="outlined"
                onClick={handleSave}
                loading={isLoading}
              >
                保存
              </Button>
              <Button
                variant="contained"
                onClick={handleSaveAndConnect}
                loading={isLoading}
              >
                保存并连接
              </Button>
            </Box>
          </Stack>
        </Paper>
      </Box>
    </Box>
  );
};

export default BookmarkForm;
