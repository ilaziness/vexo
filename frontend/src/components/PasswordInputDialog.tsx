import React, { useEffect, useState } from "react";
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Button,
  Typography,
} from "@mui/material";
import { Events } from "@wailsio/runtime";
import { BookmarkService } from "../../bindings/github.com/ilaziness/vexo/services";

const PasswordInputDialog: React.FC = () => {
  const [open, setOpen] = useState(false);
  const [password, setPassword] = useState("");
  const [message, setMessage] = useState("");

  useEffect(() => {
    // 监听密码输入事件
    const unsubscribe = Events.On("eventInput_Password", (event: any) => {
      console.log("Received eventInput_Password event:", event);
      setMessage(event.data || "需要密码来加密/解密密码");
      setOpen(true);
      setPassword("");
    });

    return () => {
      unsubscribe();
    };
  }, []);

  const handleClose = () => {
    setOpen(false);
    setPassword("");
    setMessage("");
  };

  const handleSubmit = async () => {
    if (!password.trim()) {
      setMessage("密码不能为空");
      return;
    }

    try {
      // 调用SetUserPassword保存密码，这也会通过channel通知后端
      await BookmarkService.SetUserPassword(password);
      console.log("Password set successfully");
      handleClose();
    } catch (error) {
      console.error("Failed to set password:", error);
      setMessage("设置密码失败: " + error);
    }
  };

  const handleKeyDown = (event: React.KeyboardEvent) => {
    if (event.key === "Enter") {
      handleSubmit();
    }
  };

  return (
    <Dialog open={open} onClose={handleClose} maxWidth="sm" fullWidth>
      <DialogTitle>输入密码</DialogTitle>
      <DialogContent>
        {message && (
          <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
            {message}
          </Typography>
        )}
        <TextField
          autoFocus
          margin="dense"
          label="密码"
          type="password"
          fullWidth
          variant="outlined"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          onKeyDown={handleKeyDown}
          placeholder="请输入密码"
        />
      </DialogContent>
      <DialogActions>
        <Button onClick={handleClose}>取消</Button>
        <Button onClick={handleSubmit} variant="contained">
          确定
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default PasswordInputDialog;
