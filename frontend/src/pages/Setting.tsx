import React, { useState, useEffect } from "react";
import {
  ReadConfig,
  SaveConfig,
} from "../../bindings/github.com/ilaziness/vexo/services/configservice";
import {
  Config,
  ConfigService,
} from "../../bindings/github.com/ilaziness/vexo/services";
import {
  Box,
  List,
  ListItem,
  ListItemButton,
  ListItemText,
  Divider,
  TextField,
  Button,
  Typography,
  Container,
  Paper,
  Stack,
  IconButton,
} from "@mui/material";
import { MoreHoriz } from "@mui/icons-material";
import { SelectDirectory } from "../../bindings/github.com/ilaziness/vexo/services/commonservice";
import Message from "../components/Message.tsx";
import { useMessageStore } from "../stores/common";

const Setting: React.FC = () => {
  const [activeTab, setActiveTab] = useState<"general" | "terminal" | "about">(
    "general",
  );
  const [config, setConfig] = useState<Config | null>(null);
  const [loading, setLoading] = useState(true);
  const { errorMessage } = useMessageStore();

  useEffect(() => {
    loadConfig();
  }, []);

  const loadConfig = async () => {
    try {
      const configData = await ReadConfig();
      setConfig(configData || new Config());
      setLoading(false);
    } catch (error) {
      console.error("Failed to load config:", error);
      setLoading(false);
    }
  };

  const handleSave = async () => {
    if (config) {
      try {
        await SaveConfig(config);
        ConfigService.CloseWindow();
      } catch (error) {
        console.error("Failed to save config:", error);
        errorMessage("配置保存失败");
      }
    }
  };

  const handleCancel = () => {
    // 重新加载配置以取消更改
    loadConfig().then(() => {});
    ConfigService.CloseWindow().then(() => {});
  };

  const handleGeneralChange = (field: keyof Config["General"], value: any) => {
    if (config) {
      const updatedConfig = { ...config };
      updatedConfig.General = { ...updatedConfig.General, [field]: value };
      setConfig(updatedConfig);
    }
  };

  const handleTerminalChange = (
    field: keyof Config["Terminal"],
    value: any,
  ) => {
    if (config) {
      const updatedConfig = { ...config };
      updatedConfig.Terminal = { ...updatedConfig.Terminal, [field]: value };
      setConfig(updatedConfig);
    }
  };

  const handleSelectDirectory = async () => {
    try {
      const selectedDir = await SelectDirectory();
      if (selectedDir) {
        handleGeneralChange("UserDataDir", selectedDir);
      }
    } catch (error) {
      console.error("Failed to select directory:", error);
      errorMessage("选择目录失败");
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
          width: 100,
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

  if (loading) {
    return (
      <Container maxWidth="lg" sx={{ py: 3 }}>
        <Typography>加载中...</Typography>
      </Container>
    );
  }

  return (
    <>
      <Box sx={{ height: "100%", display: "flex", flexDirection: "column" }}>
        <Box sx={{ display: "flex", flex: 1, overflow: "hidden" }}>
          {/* 左侧导航 */}
          <Paper
            sx={{
              width: 150,
              flexShrink: 0,
              borderRight: 1,
              borderColor: "divider",
              display: "flex",
              flexDirection: "column",
            }}
            elevation={0}
            square
          >
            <Box sx={{ p: 2, pb: 1.5 }}>
              <Typography
                variant="h6"
                sx={{ fontWeight: 600, fontSize: "1.2rem" }}
              >
                设置
              </Typography>
            </Box>
            <Divider />
            <List sx={{ py: 1 }}>
              <ListItem disablePadding>
                <ListItemButton
                  selected={activeTab === "general"}
                  onClick={() => setActiveTab("general")}
                  sx={{
                    py: 1,
                    px: 2,
                    "&.Mui-selected": {
                      backgroundColor: "primary.main",
                      color: "primary.contrastText",
                      "&:hover": {
                        backgroundColor: "primary.dark",
                      },
                    },
                  }}
                >
                  <ListItemText primary="通用" />
                </ListItemButton>
              </ListItem>
              <ListItem disablePadding>
                <ListItemButton
                  selected={activeTab === "terminal"}
                  onClick={() => setActiveTab("terminal")}
                  sx={{
                    py: 1,
                    px: 2,
                    "&.Mui-selected": {
                      backgroundColor: "primary.main",
                      color: "primary.contrastText",
                      "&:hover": {
                        backgroundColor: "primary.dark",
                      },
                    },
                  }}
                >
                  <ListItemText primary="终端" />
                </ListItemButton>
              </ListItem>
              <ListItem disablePadding>
                <ListItemButton
                  selected={activeTab === "about"}
                  onClick={() => setActiveTab("about")}
                  sx={{
                    py: 1,
                    px: 2,
                    "&.Mui-selected": {
                      backgroundColor: "primary.main",
                      color: "primary.contrastText",
                      "&:hover": {
                        backgroundColor: "primary.dark",
                      },
                    },
                  }}
                >
                  <ListItemText primary="关于" />
                </ListItemButton>
              </ListItem>
            </List>
          </Paper>

          {/* 右侧内容区 */}
          <Box
            sx={{
              flex: 1,
              display: "flex",
              flexDirection: "column",
              overflow: "hidden",
            }}
          >
            <Box
              sx={{
                flex: 1,
                overflowY: "auto",
                p: 3,
              }}
            >
              {activeTab === "general" && config && (
                <Box>
                  <Typography
                    variant="h5"
                    gutterBottom
                    sx={{ mb: 3, fontWeight: 600 }}
                  >
                    通用设置
                  </Typography>
                  <Paper sx={{ p: 3 }} elevation={1}>
                    <FormRow label="用户数据目录">
                      <Box sx={{ display: "flex", gap: 1 }}>
                        <TextField
                          fullWidth
                          size="small"
                          value={config.General.UserDataDir || ""}
                          onChange={(e) =>
                            handleGeneralChange("UserDataDir", e.target.value)
                          }
                          placeholder="请输入用户数据目录路径"
                          slotProps={{
                            input: {
                              endAdornment: (
                                <IconButton
                                  size="small"
                                  onClick={handleSelectDirectory}
                                >
                                  <MoreHoriz />
                                </IconButton>
                              ),
                            },
                          }}
                        />
                      </Box>
                    </FormRow>
                  </Paper>
                </Box>
              )}

              {activeTab === "terminal" && config && (
                <Box>
                  <Typography
                    variant="h5"
                    gutterBottom
                    sx={{ mb: 3, fontWeight: 600 }}
                  >
                    终端设置
                  </Typography>
                  <Paper sx={{ p: 3 }} elevation={1}>
                    <Stack spacing={2}>
                      <FormRow label="字体">
                        <TextField
                          required
                          fullWidth
                          size="small"
                          value={config.Terminal.fontFamily || ""}
                          onChange={(e) =>
                            handleTerminalChange("fontFamily", e.target.value)
                          }
                          placeholder="例如: Consolas, Monaco"
                        />
                      </FormRow>
                      <FormRow label="字体大小">
                        <TextField
                          fullWidth
                          size="small"
                          type="number"
                          slotProps={{
                            htmlInput: {
                              min: 1,
                              step: 1,
                            },
                          }}
                          value={config.Terminal.fontSize || ""}
                          onChange={(e) => {
                            const value = e.target.value;
                            // 允许空值
                            if (value === "" || value === "-") {
                              handleTerminalChange("fontSize", 14);
                              return;
                            }
                            const num = parseInt(value);
                            // 允许所有数字输入，包括中间状态
                            if (!isNaN(num)) {
                              // 确保最小值为1，但允许用户输入过程
                              handleTerminalChange(
                                "fontSize",
                                num >= 1 ? Math.floor(num) : 1,
                              );
                            }
                          }}
                          onBlur={(e) => {
                            // 失去焦点时验证最小值
                            const value = e.target.value;
                            if (value === "" || parseInt(value) < 1) {
                              handleTerminalChange("fontSize", 14);
                            }
                          }}
                          placeholder="例如: 14"
                        />
                      </FormRow>
                      <FormRow label="行高">
                        <TextField
                          fullWidth
                          size="small"
                          type="number"
                          slotProps={{
                            htmlInput: {
                              min: 0.1,
                              step: 0.1,
                            },
                          }}
                          value={config.Terminal.lineHeight || ""}
                          onChange={(e) => {
                            const value = e.target.value;
                            // 允许空值和小数点
                            if (
                              value === "" ||
                              value === "." ||
                              value === "-"
                            ) {
                              handleTerminalChange("lineHeight", 0);
                              return;
                            }
                            const num = parseFloat(value);
                            // 允许所有有效数字输入
                            if (!isNaN(num)) {
                              handleTerminalChange(
                                "lineHeight",
                                num > 0 ? num : 0.1,
                              );
                            }
                          }}
                          onBlur={(e) => {
                            // 失去焦点时验证最小值
                            const value = e.target.value;
                            if (value === "" || parseFloat(value) <= 0) {
                              handleTerminalChange("lineHeight", 1.2);
                            }
                          }}
                          placeholder="例如: 1.2"
                        />
                      </FormRow>
                    </Stack>
                  </Paper>
                </Box>
              )}

              {activeTab === "about" && (
                <Box>
                  <Typography
                    variant="h5"
                    gutterBottom
                    sx={{ mb: 3, fontWeight: 600 }}
                  >
                    关于
                  </Typography>
                  <Paper sx={{ p: 3 }} elevation={1}>
                    <Stack spacing={2} sx={{ textAlign: "center" }}>
                      <Typography variant="h4" sx={{ fontWeight: 800 }}>
                        Vexo
                      </Typography>
                      <Typography variant="body1" sx={{ fontWeight: 500 }}>
                        v0.0.1
                      </Typography>
                      <Box
                        sx={{
                          display: "flex",
                          alignItems: "center",
                          justifyContent: "center",
                        }}
                      >
                        <svg
                          xmlns="http://www.w3.org/2000/svg"
                          width="20"
                          height="20"
                          viewBox="0 0 24 24"
                          fill="currentColor"
                        >
                          <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z" />
                        </svg>
                        <Typography
                          variant="body1"
                          component="a"
                          href="https://github.com/ilaziness/vexo"
                          target="_blank"
                          rel="noopener noreferrer"
                          sx={{
                            ml: 0.8,
                            fontWeight: 500,
                            color: "primary.main",
                            textDecoration: "none",
                          }}
                        >
                          https://github.com/ilaziness/vexo
                        </Typography>
                      </Box>
                    </Stack>
                  </Paper>
                </Box>
              )}
            </Box>
          </Box>
        </Box>

        {/* 底部按钮 */}
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
          <Button variant="outlined" onClick={handleCancel}>
            取消
          </Button>
          <Button variant="contained" onClick={handleSave}>
            保存
          </Button>
        </Box>
      </Box>
      <Message />
    </>
  );
};

export default Setting;
