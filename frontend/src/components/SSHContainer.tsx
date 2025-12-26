import React, { useState } from "react";
import {
  FormControl,
  Paper,
  TextField,
  Button,
  Typography,
  Tabs,
  Tab,
} from "@mui/material";
import {
  LogService,
  SSHService,
} from "../../bindings/github.com/ilaziness/vexo/services";
import { SSHLinkInfo } from "../types/ssh";
import Terminal from "./Terminal";
import Sftp from "./Sftp";
import { useMessageStore } from "../stores/common";
import { parseCallServiceError } from "../func/service";

// 输入连接信息表单
const NewConnection = (props: {
  onConnect: (info: SSHLinkInfo) => void;
  error?: string;
  connecting?: boolean;
}) => {
  const [host, setHost] = useState("");
  const [port, setPort] = useState<string>("22");
  const [user, setUser] = useState("");
  const [password, setPassword] = useState("");
  const [key, setKey] = useState("");

  const { errorMessage } = useMessageStore();

  const onSelectKeyFile = async () => {
    try{
      const file = await SSHService.SelectKeyFile()
      LogService.Debug(`select key file ${file}`)
      setKey(file)
    } catch(err: any) {
      errorMessage(parseCallServiceError(err))
    }
  }

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

    props.onConnect({
      host,
      port: p,
      user,
      password,
      key,
    });
  };

  return (
    <>
      <Paper
        component="form"
        sx={{
          display: "flex",
          flexDirection: "column",
          pt: 2,
          alignItems: "center",
          flex: 1,
          height: "100%",
        }}
        onSubmit={submit}
      >
        <FormControl>
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
          <Button onClick={onSelectKeyFile} sx={{m: 0.8}}>
            选择密钥文件
          </Button>
        </FormControl>
        {props.error && (
          <Typography color="error" sx={{ mt: 1, mb: 1 }}>
            {props.error}
          </Typography>
        )}
        <Button variant="contained" type="submit" loading={props.connecting}>
          {props.connecting ? "Connecting..." : "Connect"}
        </Button>
      </Paper>
    </>
  );
};

// SSH 连接容器组件，管理连接状态和错误处理
export default function SSHContainer() {
  "use memo";
  const [linkID, setLinkID] = React.useState<string>("");
  const [connectionError, setConnectionError] = React.useState<string>("");
  const [connecting, setConnecting] = React.useState<boolean>(false);
  const [activeTab, setActiveTab] = React.useState(0); // 0 for terminal, 1 for sftp

  const connect = async (li: SSHLinkInfo) => {
    setConnectionError("");
    setConnecting(true);
    try {
      const ID = await SSHService.Connect(
        li.host,
        li.port || 22,
        li.user || "",
        li.password || "",
        li.key || ""
      );
      LogService.Debug(`SSH connection established with ID: ${ID}`);
      setLinkID(ID);
    } catch (err: any) {
      const msg = "Connection failed";
      LogService.Error(`${msg}: ${err.message || err}`).then(() => {});
      setConnectionError(parseCallServiceError(err));
    } finally {
      setConnecting(false);
    }
  };

  const handleTabChange = (event: React.SyntheticEvent, newValue: number) => {
    setActiveTab(newValue);
  };

  if (linkID === "") {
    return (
      <NewConnection
        onConnect={connect}
        error={connectionError}
        connecting={connecting}
      />
    );
  }

  return (
    <Paper sx={{ height: "100%", display: "flex", flexDirection: "column", flex: 1 }}>
      <Tabs value={activeTab} onChange={handleTabChange} aria-label="ssh tabs">
        <Tab label="Terminal" />
        <Tab label="SFTP" />
      </Tabs>
      {activeTab === 0 && <Terminal linkID={linkID} />}
      {activeTab === 1 && <Sftp linkID={linkID} />}
    </Paper>
  );
}