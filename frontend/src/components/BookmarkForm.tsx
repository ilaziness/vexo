import React, { useState, useEffect } from "react";
import { Box, TextField, Button, Grid, Typography } from "@mui/material";
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

  if (!formData) {
    return (
      <Box
        sx={{
          p: 2,
          display: "flex",
          justifyContent: "center",
          alignItems: "center",
          height: "100%",
        }}
      >
        <Typography variant="h6" color="textSecondary">
          请选择一个书签进行编辑
        </Typography>
      </Box>
    );
  }

  return (
    <Box sx={{ p: 2 }}>
      <Grid container spacing={2}>
        <Grid size={12}>
          <TextField
            fullWidth
            label="书签名称"
            name="title"
            value={formData.title}
            onChange={handleChange}
            variant="outlined"
          />
        </Grid>
        <Grid size={12}>
          <TextField
            fullWidth
            label="分组"
            name="group_name"
            value={formData.group_name}
            onChange={handleChange}
            variant="outlined"
          />
        </Grid>
        <Grid size={12}>
          <TextField
            fullWidth
            label="主机地址"
            name="host"
            value={formData.host}
            onChange={handleChange}
            variant="outlined"
          />
        </Grid>
        <Grid size={6}>
          <TextField
            fullWidth
            label="端口"
            name="port"
            type="number"
            value={formData.port}
            onChange={handleChange}
            variant="outlined"
          />
        </Grid>
        <Grid size={6}>
          <TextField
            fullWidth
            label="用户名"
            name="user"
            value={formData.user}
            onChange={handleChange}
            variant="outlined"
          />
        </Grid>
        <Grid size={12}>
          <TextField
            fullWidth
            label="密码"
            name="password"
            type="password"
            value={formData.password}
            onChange={handleChange}
            variant="outlined"
          />
        </Grid>
        <Grid size={12}>
          <TextField
            fullWidth
            label="密钥路径"
            name="private_key"
            value={formData.private_key}
            onChange={handleChange}
            variant="outlined"
          />
        </Grid>
      </Grid>
      <Box
        sx={{ mt: 3, display: "flex", justifyContent: "flex-start", gap: 2 }}
      >
        <Button variant="contained" color="primary" onClick={handleSave}>
          保存
        </Button>
        <Button
          variant="outlined"
          color="secondary"
          onClick={handleTestConnection}
        >
          测试连接
        </Button>
        <Button
          variant="contained"
          color="success"
          onClick={handleSaveAndConnect}
        >
          保存并连接
        </Button>
      </Box>
    </Box>
  );
};

export default BookmarkForm;
