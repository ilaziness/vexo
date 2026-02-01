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
  Paper,
  Stack,
  IconButton,
} from "@mui/material";
import { MoreHoriz } from "@mui/icons-material";
import {
  SelectDirectory,
  GetAppInfo,
} from "../../bindings/github.com/ilaziness/vexo/services/appservice";
import About from "../components/About";
import { AppInfo } from "../../bindings/github.com/ilaziness/vexo/services";
import Message from "../components/Message.tsx";
import Loading from "../components/Loading.tsx";
import { useMessageStore } from "../stores/message.ts";
import FormRow from "../components/FormRow";
import OpBar from "../components/OpBar.tsx";

const Setting: React.FC = () => {
  const [activeTab, setActiveTab] = useState<"general" | "terminal" | "about">(
    "general",
  );
  const [config, setConfig] = useState<Config | null>(null);
  const [loading, setLoading] = useState(true);
  const { errorMessage } = useMessageStore();
  const [appinfo, setAppInfo] = useState<AppInfo>({} as AppInfo);

  useEffect(() => {
    loadConfig();
    GetAppInfo().then((info) => {
      console.log("App Info:", info);
      setAppInfo(info);
    });
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
    if (!config) return;

    // 验证必填字段
    const errors: string[] = [];

    // 通用设置验证
    if (
      !config.General.UserDataDir ||
      config.General.UserDataDir.trim() === ""
    ) {
      errors.push("用户数据目录不能为空");
    }

    // 终端设置验证
    if (
      !config.Terminal.fontFamily ||
      config.Terminal.fontFamily.trim() === ""
    ) {
      errors.push("字体不能为空");
    }
    if (!config.Terminal.fontSize || config.Terminal.fontSize < 1) {
      errors.push("字体大小必须大于等于1");
    }
    if (!config.Terminal.lineHeight || config.Terminal.lineHeight <= 0) {
      errors.push("行高必须大于0");
    }

    if (errors.length > 0) {
      errorMessage(errors.join("-"));
      return;
    }

    try {
      await SaveConfig(config);
      ConfigService.CloseWindow();
    } catch (error) {
      console.error("Failed to save config:", error);
      errorMessage("配置保存失败");
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

  if (loading) {
    return (
      <Box
        sx={{
          height: "100%",
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
        }}
      >
        <Loading />
      </Box>
    );
  }

  return (
    <>
      <Box sx={{ height: 40 }}>
        <OpBar />
      </Box>
      <Box
        sx={{
          height: "calc(100% - 40px)",
          display: "flex",
          flexDirection: "column",
        }}
      >
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
                p: 2,
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
                  <Paper sx={{ p: 2 }} elevation={1}>
                    <Stack spacing={1.5}>
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
                          onChange={(e) =>
                            handleTerminalChange(
                              "fontSize",
                              e.target.value === ""
                                ? null
                                : parseInt(e.target.value) || null,
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
                          slotProps={{
                            htmlInput: {
                              min: 0.1,
                              step: 0.1,
                            },
                          }}
                          value={config.Terminal.lineHeight || ""}
                          onChange={(e) =>
                            handleTerminalChange(
                              "lineHeight",
                              e.target.value === ""
                                ? null
                                : parseFloat(e.target.value) || null,
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
                <Paper sx={{ p: 3 }} elevation={1}>
                  <About appinfo={appinfo} />
                </Paper>
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
