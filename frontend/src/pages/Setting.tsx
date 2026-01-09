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
} from "@mui/material";
import Message from "../components/Message.tsx";

const Setting: React.FC = () => {
  const [activeTab, setActiveTab] = useState<"general" | "terminal" | "about">(
    "general",
  );
  const [config, setConfig] = useState<Config | null>(null);
  const [loading, setLoading] = useState(true);

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
        alert("配置保存成功");
      } catch (error) {
        console.error("Failed to save config:", error);
        alert("配置保存失败");
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
          width: 140,
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
              width: 200,
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
              <Typography variant="h6" sx={{ fontWeight: 600, fontSize: "1.1rem" }}>
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
                  <ListItemText 
                    primary="通用" 
                    primaryTypographyProps={{ fontSize: "0.95rem" }}
                  />
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
                  <ListItemText 
                    primary="终端" 
                    primaryTypographyProps={{ fontSize: "0.95rem" }}
                  />
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
                  <ListItemText 
                    primary="关于" 
                    primaryTypographyProps={{ fontSize: "0.95rem" }}
                  />
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
                  <Typography variant="h5" gutterBottom sx={{ mb: 3, fontWeight: 600 }}>
                    通用设置
                  </Typography>
                  <Paper sx={{ p: 3 }} elevation={1}>
                    <FormRow label="用户数据目录">
                      <TextField
                        fullWidth
                        size="small"
                        value={config.General.UserDataDir || ""}
                        onChange={(e) =>
                          handleGeneralChange("UserDataDir", e.target.value)
                        }
                        placeholder="请输入用户数据目录路径"
                      />
                    </FormRow>
                  </Paper>
                </Box>
              )}

              {activeTab === "terminal" && config && (
                <Box>
                  <Typography variant="h5" gutterBottom sx={{ mb: 3, fontWeight: 600 }}>
                    终端设置
                  </Typography>
                  <Paper sx={{ p: 3 }} elevation={1}>
                    <Stack spacing={2}>
                      <FormRow label="字体">
                        <TextField
                          fullWidth
                          size="small"
                          value={config.Terminal.Font || ""}
                          onChange={(e) =>
                            handleTerminalChange("Font", e.target.value)
                          }
                          placeholder="例如: Consolas, Monaco"
                        />
                      </FormRow>
                      <FormRow label="字体大小">
                        <TextField
                          fullWidth
                          size="small"
                          type="number"
                          value={config.Terminal.FontSize || 0}
                          onChange={(e) =>
                            handleTerminalChange(
                              "FontSize",
                              parseInt(e.target.value),
                            )
                          }
                          placeholder="例如: 14"
                        />
                      </FormRow>
                      <FormRow label="行高">
                        <TextField
                          fullWidth
                          size="small"
                          type="number"
                          inputProps={{ step: 0.1 }}
                          value={config.Terminal.LineHeight || 0}
                          onChange={(e) =>
                            handleTerminalChange(
                              "LineHeight",
                              parseFloat(e.target.value),
                            )
                          }
                          placeholder="例如: 1.2"
                        />
                      </FormRow>
                    </Stack>
                  </Paper>
                </Box>
              )}

              {activeTab === "about" && (
                <Box>
                  <Typography variant="h5" gutterBottom sx={{ mb: 3, fontWeight: 600 }}>
                    关于
                  </Typography>
                  <Paper sx={{ p: 3 }} elevation={1}>
                    <Stack spacing={2}>
                      <Box>
                        <Typography color="text.secondary" sx={{ mb: 0.5, fontSize: "0.9rem" }}>
                          应用名称
                        </Typography>
                        <Typography variant="body1" sx={{ fontWeight: 500 }}>
                          Vexo
                        </Typography>
                      </Box>
                      <Divider />
                      <Box>
                        <Typography color="text.secondary" sx={{ mb: 0.5, fontSize: "0.9rem" }}>
                          版本
                        </Typography>
                        <Typography variant="body1" sx={{ fontWeight: 500 }}>
                          v1.0.0
                        </Typography>
                      </Box>
                      <Divider />
                      <Box>
                        <Typography color="text.secondary" sx={{ mb: 0.5, fontSize: "0.9rem" }}>
                          描述
                        </Typography>
                        <Typography variant="body1" sx={{ fontWeight: 500 }}>
                          一个现代化的SSH和SFTP客户端工具
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
