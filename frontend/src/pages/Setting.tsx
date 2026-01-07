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
  Grid,
  Card,
  CardContent,
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

  if (loading) {
    return (
      <Container maxWidth="lg" sx={{ py: 3 }}>
        <Typography>加载中...</Typography>
      </Container>
    );
  }

  return (
    <>
      <Container
        maxWidth={false}
        sx={{ height: "100vh", p: 0, display: "flex", flexDirection: "column" }}
      >
        <Grid container sx={{ height: "100%", flexGrow: 1 }}>
          {/* 左侧导航 */}
          <Grid size={3}>
            <Paper sx={{ height: "100%", borderRight: "1px solid #e0e0e0" }}>
              <Box sx={{ p: 2 }}>
                <Typography variant="h6" sx={{ fontWeight: "bold", mb: 2 }}>
                  设置
                </Typography>
              </Box>
              <Divider />
              <List>
                <ListItem disablePadding>
                  <ListItemButton
                    selected={activeTab === "general"}
                    onClick={() => setActiveTab("general")}
                  >
                    <ListItemText primary="通用" />
                  </ListItemButton>
                </ListItem>
                <ListItem disablePadding>
                  <ListItemButton
                    selected={activeTab === "terminal"}
                    onClick={() => setActiveTab("terminal")}
                  >
                    <ListItemText primary="终端" />
                  </ListItemButton>
                </ListItem>
                <ListItem disablePadding>
                  <ListItemButton
                    selected={activeTab === "about"}
                    onClick={() => setActiveTab("about")}
                  >
                    <ListItemText primary="关于" />
                  </ListItemButton>
                </ListItem>
              </List>
            </Paper>
          </Grid>

          {/* 右侧面板 */}
          <Grid size={9} direction="column">
            <Grid sx={{ p: 3, overflowY: "auto", flexGrow: 1 }}>
              {activeTab === "general" && config && (
                <Card>
                  <CardContent>
                    <Typography variant="h5" gutterBottom>
                      通用设置
                    </Typography>
                    <TextField
                      fullWidth
                      label="用户数据目录"
                      value={config.General.UserDataDir || ""}
                      onChange={(e) =>
                        handleGeneralChange("UserDataDir", e.target.value)
                      }
                      margin="normal"
                    />
                  </CardContent>
                </Card>
              )}

              {activeTab === "terminal" && config && (
                <Card>
                  <CardContent>
                    <Typography variant="h5" gutterBottom>
                      终端设置
                    </Typography>
                    <TextField
                      fullWidth
                      label="字体"
                      value={config.Terminal.Font || ""}
                      onChange={(e) =>
                        handleTerminalChange("Font", e.target.value)
                      }
                      margin="normal"
                    />
                    <TextField
                      fullWidth
                      label="字体大小"
                      type="number"
                      value={config.Terminal.FontSize || 0}
                      onChange={(e) =>
                        handleTerminalChange(
                          "FontSize",
                          parseInt(e.target.value),
                        )
                      }
                      margin="normal"
                    />
                    <TextField
                      fullWidth
                      label="行高"
                      type="number"
                      inputProps={{ step: 0.1 }} // 使用 inputProps 而不是直接使用 step
                      value={config.Terminal.LineHeight || 0}
                      onChange={(e) =>
                        handleTerminalChange(
                          "LineHeight",
                          parseFloat(e.target.value),
                        )
                      }
                      margin="normal"
                    />
                  </CardContent>
                </Card>
              )}

              {activeTab === "about" && (
                <Card>
                  <CardContent>
                    <Typography variant="h5" gutterBottom>
                      关于
                    </Typography>
                    <Typography>
                      <strong>应用名称:</strong> Vexo
                    </Typography>
                    <Typography>
                      <strong>版本:</strong> v1.0.0
                    </Typography>
                    <Typography>
                      <strong>描述:</strong> 一个现代化的SSH和SFTP客户端工具
                    </Typography>
                  </CardContent>
                </Card>
              )}
            </Grid>
          </Grid>
        </Grid>

        {/* 底部按钮 - 位于整个布局的最底部 */}
        <Grid
          sx={{
            borderTop: "1px solid #e0e0e0",
            p: 2,
            display: "flex",
            justifyContent: "flex-end",
            mt: "auto",
          }}
        >
          <Button variant="outlined" onClick={handleCancel} sx={{ mr: 1 }}>
            取消
          </Button>
          <Button variant="contained" onClick={handleSave}>
            保存
          </Button>
        </Grid>
      </Container>
      <Message />
    </>
  );
};

export default Setting;
