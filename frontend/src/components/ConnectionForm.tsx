import React, { useState } from "react";
import {
  Button,
  FormControl,
  TextField,
  Typography,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Select,
  MenuItem,
  InputLabel,
  Box,
} from "@mui/material";
import {
  LogService,
  SSHService,
  BookmarkService,
} from "../../bindings/github.com/ilaziness/vexo/services";
import { SSHLinkInfo } from "../types/ssh";
import { useMessageStore } from "../stores/message";
import { parseCallServiceError } from "../func/service";
import {
  SSHBookmark,
  BookmarkGroup,
} from "../../bindings/github.com/ilaziness/vexo/services/models";
import { generateRandomId } from "../func/id";

// 输入连接信息表单
interface ConnectionFormProps {
  onConnect: (info: SSHLinkInfo) => void;
  error?: string;
  connecting?: boolean;
}

const ConnectionForm: React.FC<ConnectionFormProps> = ({
  onConnect,
  error,
  connecting,
}) => {
  const [host, setHost] = useState("172.20.10.4");
  const [port, setPort] = useState<string>("22");
  const [user, setUser] = useState("test");
  const [password, setPassword] = useState("test");
  const [key, setKey] = useState("");
  const [keyPassword, setKeyPassword] = useState("");
  const [saveDialogOpen, setSaveDialogOpen] = useState(false);
  const [selectedGroup, setSelectedGroup] = useState("default");
  const [groups, setGroups] = useState<string[]>([]);

  const { errorMessage } = useMessageStore();

  const onSelectKeyFile = async () => {
    try {
      const file = await SSHService.SelectKeyFile();
      LogService.Debug(`select key file ${file}`);
      setKey(file);
    } catch (err: any) {
      errorMessage(parseCallServiceError(err));
    }
  };

  const submit = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    // validate required fields: host, port, user
    if (!host.trim() || !port.toString().trim() || !user.trim()) {
      errorMessage("Host、Port 和 Username 为必填项");
      return;
    }
    const p = Number(port);
    if (Number.isNaN(p) || p <= 0) {
      errorMessage("Port 必须是有效的数字");
      return;
    }
    // password 和 key 二选一
    if (!password.trim() && !key.trim()) {
      errorMessage("必须提供 Password 或 Private Key 其一");
      return;
    }

    onConnect({
      host,
      port: p,
      user,
      password,
      key,
      keyPassword: key ? keyPassword : undefined,
    });
  };

  const handleSaveClick = () => {
    // 加载书签分组
    BookmarkService.ListBookmarks()
      .then((bookmarkGroups) => {
        const groupNames = bookmarkGroups
          .filter((group): group is BookmarkGroup => group !== null)
          .map((group) => group.name);
        setGroups(groupNames);
        setSelectedGroup(groupNames[0] || "default");
        setSaveDialogOpen(true);
      })
      .catch((err) => {
        errorMessage(parseCallServiceError(err));
      });
  };

  const handleSaveBookmark = () => {
    // 创建新的SSHBookmark对象
    const newBookmark: SSHBookmark = {
      id: generateRandomId(),
      title: `${host}:${port}`,
      group_name: selectedGroup,
      host,
      port: Number(port),
      private_key: key,
      private_key_password: keyPassword,
      user,
      password,
    };

    // 保存到书签
    BookmarkService.SaveBookmark(newBookmark)
      .then(() => {
        setSaveDialogOpen(false);
        // 也可以选择连接
        onConnect({
          host,
          port: Number(port),
          user,
          password,
          key,
          keyPassword: key ? keyPassword : undefined,
        });
      })
      .catch((err) => {
        errorMessage(parseCallServiceError(err));
      });
  };

  return (
    <>
      <Box
        component="form"
        sx={{
          display: "flex",
          flexDirection: "column",
          pt: 3,
          alignItems: "center",
          flex: 1,
        }}
        onSubmit={submit}
      >
        <FormControl sx={{ width: 300 }}>
          <TextField
            required
            label="Host"
            value={host}
            onChange={(e) => setHost(e.target.value)}
            variant="outlined"
            margin="normal"
            size="small"
            sx={{ m: 0.8 }}
          />
          <TextField
            required
            label="Port"
            value={port}
            onChange={(e) => setPort(e.target.value)}
            variant="outlined"
            margin="normal"
            size="small"
            sx={{ m: 0.8 }}
          />
          <TextField
            required
            label="Username"
            value={user}
            onChange={(e) => setUser(e.target.value)}
            variant="outlined"
            margin="normal"
            size="small"
            sx={{ m: 0.8 }}
          />
          <TextField
            label="Password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            variant="outlined"
            margin="normal"
            type="password"
            size="small"
            sx={{ m: 0.8 }}
          />
          <Box sx={{ display: "flex", gap: 1, m: 0.8 }}>
            <TextField
              label="Key file"
              value={key}
              variant="outlined"
              size="small"
              slotProps={{
                input: { readOnly: true },
              }}
              placeholder="please select key file"
              sx={{ flex: 1 }}
            />
            <Button onClick={onSelectKeyFile} variant="outlined">
              Select
            </Button>
          </Box>
          {key && (
            <TextField
              label="Key Password"
              value={keyPassword}
              onChange={(e) => setKeyPassword(e.target.value)}
              variant="outlined"
              margin="normal"
              type="password"
              size="small"
              sx={{ m: 0.8 }}
              placeholder="private key password (optional)"
            />
          )}
        </FormControl>
        {error && (
          <Typography color="error" sx={{ mt: 1, mb: 1 }}>
            {error}
          </Typography>
        )}
        <Box sx={{ display: "flex", gap: 1, mt: 1 }}>
          <Button
            variant="contained"
            type="submit"
            size="large"
            loading={connecting}
          >
            连接
          </Button>
          <Button
            variant="outlined"
            size="large"
            loading={connecting}
            onClick={handleSaveClick}
          >
            保存
          </Button>
        </Box>
      </Box>

      {/* 保存书签对话框 */}
      <Dialog open={saveDialogOpen} onClose={() => setSaveDialogOpen(false)}>
        <DialogTitle>保存到书签</DialogTitle>
        <DialogContent>
          <Box sx={{ mt: 2, minWidth: 400 }}>
            <InputLabel id="bookmark-group-select-label">选择分组</InputLabel>
            <Select
              labelId="bookmark-group-select-label"
              value={selectedGroup}
              label="选择分组"
              onChange={(e) => setSelectedGroup(e.target.value)}
              fullWidth
            >
              {groups.map((group) => (
                <MenuItem key={group} value={group}>
                  {group}
                </MenuItem>
              ))}
            </Select>
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setSaveDialogOpen(false)}>取消</Button>
          <Button onClick={handleSaveBookmark}>保存</Button>
        </DialogActions>
      </Dialog>
    </>
  );
};

export default ConnectionForm;
