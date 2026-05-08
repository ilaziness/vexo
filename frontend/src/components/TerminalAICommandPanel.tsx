import { useState, useRef, useEffect } from "react";
import {
  Box,
  Typography,
  TextField,
  IconButton,
  Paper,
  Alert,
  Chip,
  Tooltip,
  CircularProgress,
  Fade,
  InputAdornment,
} from "@mui/material";
import SendIcon from "@mui/icons-material/Send";
import CloseIcon from "@mui/icons-material/Close";
import ClearAllIcon from '@mui/icons-material/ClearAll';
import ContentCopyIcon from "@mui/icons-material/ContentCopy";
import PlayArrowIcon from "@mui/icons-material/PlayArrow";
import EditIcon from "@mui/icons-material/Edit";
import AutoAwesomeIcon from '@mui/icons-material/AutoAwesome';
import WarningIcon from "@mui/icons-material/Warning";
import { useAICommandStore, AIMessage } from "../stores/aiCommand";

interface TerminalAICommandPanelProps {
  sessionID: string;
  onClose: () => void;
}

export default function TerminalAICommandPanel({
  sessionID,
  onClose,
}: TerminalAICommandPanelProps) {
  const [input, setInput] = useState("");
  const [editingCommand, setEditingCommand] = useState<string | null>(null);
  const messagesEndRef = useRef<HTMLDivElement>(null);

  const {
    messages,
    isLoading,
    error,
    generateCommand,
    sendToTerminal,
    copyToClipboard,
    clearHistory,
    resetError,
  } = useAICommandStore();

  // 滚动到最新消息
  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages, isLoading]);

  // 发送消息
  const handleSend = async () => {
    if (!input.trim()) return;
    await generateCommand(sessionID);
    setInput("");
  };

  // 处理按键
  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  };

  // 运行命令
  const handleRun = (command: string) => {
    sendToTerminal(command, [sessionID]);
  };

  // 复制命令
  const handleCopy = async (command: string) => {
    await copyToClipboard(command);
  };

  // 编辑命令
  const handleEdit = (command: string) => {
    setEditingCommand(command);
  };

  // 保存编辑后的命令
  const handleSaveEdit = () => {
    if (editingCommand) {
      handleRun(editingCommand);
      setEditingCommand(null);
    }
  };

  // 获取安全警告颜色
  const getWarningColor = (level?: string) => {
    switch (level) {
      case "high":
        return "error";
      case "medium":
        return "warning";
      case "low":
        return "info";
      default:
        return "info";
    }
  };

  // 渲染消息
  const renderMessage = (message: AIMessage) => {
    const isUser = message.role === "user";

    return (
      <Box
        key={message.id}
        sx={{
          display: "flex",
          flexDirection: isUser ? "row-reverse" : "row",
          gap: 1,
          mb: 2,
        }}
      >
        <Box
          sx={{
            width: 32,
            height: 32,
            borderRadius: "50%",
            bgcolor: isUser ? "primary.main" : "secondary.main",
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
            flexShrink: 0,
          }}
        >
          {isUser ? (
            <Typography variant="caption" sx={{ color: "white" }}>
              我
            </Typography>
          ) : (
            <AutoAwesomeIcon sx={{ color: "white", fontSize: 18 }} />
          )}
        </Box>

        <Box
          sx={{
            maxWidth: "80%",
            bgcolor: isUser ? "primary.light" : "grey.100",
            color: isUser ? "white" : "text.primary",
            p: 1.5,
            borderRadius: 2,
            wordBreak: "break-word",
          }}
        >
          <Typography variant="body2" sx={{ whiteSpace: "pre-wrap" }}>
            {message.content}
          </Typography>

          {/* 命令展示和操作按钮 */}
          {message.command && (
            <Box sx={{ mt: 1 }}>
              <Paper
                variant="outlined"
                sx={{
                  p: 1,
                  bgcolor: "background.default",
                  fontFamily: "monospace",
                  fontSize: "0.875rem",
                  position: "relative",
                }}
              >
                {editingCommand === message.command.command ? (
                  <TextField
                    fullWidth
                    size="small"
                    value={editingCommand}
                    onChange={(e) => setEditingCommand(e.target.value)}
                    onKeyDown={(e) => {
                      if (e.key === "Enter") handleSaveEdit();
                    }}
                    autoFocus
                  />
                ) : (
                  <Typography
                    variant="body2"
                    sx={{ fontFamily: "monospace", pr: 8 }}
                  >
                    $ {message.command.command}
                  </Typography>
                )}

                <Box
                  sx={{
                    position: "absolute",
                    top: 4,
                    right: 4,
                    display: "flex",
                    gap: 0.5,
                  }}
                >
                  <Tooltip title="复制">
                    <IconButton
                      size="small"
                      onClick={() => handleCopy(message.command!.command)}
                    >
                      <ContentCopyIcon fontSize="small" />
                    </IconButton>
                  </Tooltip>
                  <Tooltip title="运行">
                    <IconButton
                      size="small"
                      color="primary"
                      onClick={() => handleRun(message.command!.command)}
                    >
                      <PlayArrowIcon fontSize="small" />
                    </IconButton>
                  </Tooltip>
                  <Tooltip title="编辑">
                    <IconButton
                      size="small"
                      onClick={() => handleEdit(message.command!.command)}
                    >
                      <EditIcon fontSize="small" />
                    </IconButton>
                  </Tooltip>
                </Box>
              </Paper>

              {/* 安全警告 */}
              {message.command.safety_warning && (
                <Alert
                  severity={getWarningColor(
                    message.command.safety_warning.level
                  )}
                  icon={<WarningIcon />}
                  sx={{ mt: 1 }}
                >
                  <Typography variant="caption" sx={{ fontWeight: "bold" }}>
                    {message.command.safety_warning.message}
                  </Typography>
                  <Typography variant="caption" sx={{ display: "block" }}>
                    {message.command.safety_warning.risk_detail}
                  </Typography>
                  <Typography variant="caption" sx={{ display: "block" }}>
                    建议: {message.command.safety_warning.suggestion}
                  </Typography>
                </Alert>
              )}

              {/* 替代方案 */}
              {message.command.alternatives &&
                message.command.alternatives.length > 0 && (
                  <Box sx={{ mt: 1 }}>
                    <Typography variant="caption" color="text.secondary">
                      替代方案:
                    </Typography>
                    {message.command.alternatives.map((alt, index) => (
                      <Chip
                        key={index}
                        label={alt}
                        size="small"
                        sx={{ ml: 0.5, mt: 0.5 }}
                        onClick={() => handleRun(alt)}
                      />
                    ))}
                  </Box>
                )}
            </Box>
          )}
        </Box>
      </Box>
    );
  };

  return (
    <Fade in>
      <Box
        sx={{
          width: "40%",
          height: "100%",
          display: "flex",
          flexDirection: "column",
          bgcolor: "background.paper",
          borderLeft: 1,
          borderColor: "divider",
        }}
      >
        {/* 头部 */}
        <Box
          sx={{
            p: 2,
            borderBottom: 1,
            borderColor: "divider",
            display: "flex",
            alignItems: "center",
            justifyContent: "space-between",
          }}
        >
          <Box sx={{ display: "flex", alignItems: "center", gap: 1 }}>
            <AutoAwesomeIcon color="primary" />
            <Typography variant="h6">AI 命令助手</Typography>
          </Box>
          <Box>
            <Tooltip title="清除历史">
              <IconButton size="small" onClick={clearHistory} sx={{ mr: 1 }}>
                <ClearAllIcon fontSize="small" />
              </IconButton>
            </Tooltip>
            <Tooltip title="关闭">
            <IconButton onClick={onClose}>
              <CloseIcon />
            </IconButton>
            </Tooltip>
          </Box>
        </Box>

        {/* 错误提示 */}
        {error && (
          <Alert severity="error" sx={{ m: 2 }} onClose={resetError}>
            {error}
          </Alert>
        )}

        {/* 消息区域 */}
        <Box sx={{ flex: 1, overflow: "auto", p: 2 }}>
          {messages.length === 0 && (
            <Box sx={{ textAlign: "center", py: 4, color: "text.secondary" }}>
              <AutoAwesomeIcon sx={{ fontSize: 48, mb: 2, opacity: 0.5 }} />
              <Typography variant="body2">
                描述你想要的操作，AI 将生成命令
              </Typography>
              <Box sx={{ mt: 2, textAlign: "left" }}>
                <Typography variant="caption" color="text.secondary">
                  示例:
                </Typography>
                <Typography variant="caption" sx={{ display: "block" }}>
                  • 查找大于100MB的文件
                </Typography>
                <Typography variant="caption" sx={{ display: "block" }}>
                  • 查看端口占用
                </Typography>
                <Typography variant="caption" sx={{ display: "block" }}>
                  • 压缩日志文件
                </Typography>
              </Box>
            </Box>
          )}

          {messages.map(renderMessage)}

          {isLoading && (
            <Box sx={{ display: "flex", justifyContent: "center", p: 2 }}>
              <CircularProgress size={24} />
            </Box>
          )}

          <div ref={messagesEndRef} />
        </Box>

        {/* 输入区域 */}
        <Box sx={{ p: 2, borderTop: 1, borderColor: "divider" }}>
          <TextField
            fullWidth
            multiline
            maxRows={3}
            placeholder="输入你的问题..."
            value={input}
            onChange={(e) => setInput(e.target.value)}
            onKeyDown={handleKeyDown}
            disabled={isLoading}
            slotProps={{
              input: {
                endAdornment: (
                  <InputAdornment position="end">
                    <IconButton
                      onClick={handleSend}
                      disabled={!input.trim() || isLoading}
                      color="primary"
                    >
                      <SendIcon />
                    </IconButton>
                  </InputAdornment>
                ),
              }
            }}
          />
        </Box>
      </Box>
    </Fade>
  );
}
