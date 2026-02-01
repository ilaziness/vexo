import React, { useState } from "react";
import { Browser } from "@wailsio/runtime";
import {
  Button,
  Typography,
  Box,
  Stack,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  CircularProgress,
} from "@mui/material";
import { GitHub as GitHubIcon } from "@mui/icons-material";
import { CheckUpdate } from "../../bindings/github.com/ilaziness/vexo/services/appservice";
import { AppInfo } from "../../bindings/github.com/ilaziness/vexo/services";
import { useMessageStore } from "../stores/message";

interface Props {
  appinfo: AppInfo;
}

const About: React.FC<Props> = ({ appinfo }) => {
  const { errorMessage, infoMessage } = useMessageStore();
  const [checkingUpdate, setCheckingUpdate] = useState(false);
  const [newVersion, setNewVersion] = useState<any | null>(null);
  const [updateDialogOpen, setUpdateDialogOpen] = useState(false);

  const handleCheckUpdate = async () => {
    try {
      setCheckingUpdate(true);
      const [ok, ver] = await CheckUpdate();
      setCheckingUpdate(false);
      if (ok) {
        setNewVersion(ver);
        setUpdateDialogOpen(true);
      } else {
        infoMessage("当前已是最新版本");
      }
    } catch (err) {
      setCheckingUpdate(false);
      console.error("CheckUpdate failed", err);
      errorMessage("检查更新失败");
    }
  };

  return (
    <Box>
      <Typography variant="h5" gutterBottom sx={{ mb: 3, fontWeight: 600 }}>
        关于
      </Typography>
      <Box>
        <Stack spacing={2} sx={{ textAlign: "center", alignItems: "center" }}>
          <img
            src="/appicon.png"
            alt="Logo"
            style={{ width: "64px", height: "64px" }}
          />
          <Typography variant="h4" sx={{ fontWeight: 800 }}>
            Vexo
          </Typography>

          <Box
            sx={{
              display: "flex",
              alignItems: "center",
              justifyContent: "center",
            }}
          >
            <GitHubIcon
              sx={{ "&:hover": { cursor: "pointer" } }}
              onClick={() => {
                Browser.OpenURL(appinfo.HomeURL);
              }}
            />
            <Typography
              component="span"
              variant="body1"
              sx={{ fontWeight: 500, ml: 1 }}
            >
              {appinfo.Version}
            </Typography>
            <Button
              size="small"
              sx={{ ml: 1 }}
              onClick={handleCheckUpdate}
              disabled={checkingUpdate}
              variant="outlined"
            >
              {checkingUpdate ? <CircularProgress size={16} /> : "检查更新"}
            </Button>
          </Box>

          <Dialog
            open={updateDialogOpen}
            onClose={() => setUpdateDialogOpen(false)}
          >
            <DialogTitle>发现新版本</DialogTitle>
            <DialogContent>
              <Typography sx={{ fontWeight: 600 }}>
                {newVersion?.Version}
              </Typography>
              <Typography sx={{ whiteSpace: "pre-wrap", mt: 1 }}>
                {newVersion?.Notes}
              </Typography>
            </DialogContent>
            <DialogActions>
              <Button onClick={() => setUpdateDialogOpen(false)}>关闭</Button>
              <Button
                onClick={() => {
                  if (newVersion?.URL) {
                    Browser.OpenURL(newVersion.URL);
                  }
                }}
              >
                打开下载页
              </Button>
            </DialogActions>
          </Dialog>

          <Typography
            color="text.secondary"
            variant="body2"
            sx={{ fontWeight: 500, ml: 1 }}
          >
            SSH桌面客户端，支持Window、Linux和MacOS
          </Typography>
        </Stack>
      </Box>
    </Box>
  );
};

export default About;
