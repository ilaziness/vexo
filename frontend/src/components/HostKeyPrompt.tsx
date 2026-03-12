import React, { useEffect, useState } from "react";
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  Typography,
} from "@mui/material";
import { Events } from "@wailsio/runtime";
import { SSHService } from "../../bindings/github.com/ilaziness/vexo/services";

interface Payload {
  host: string;
  address: string;
  fingerprint: string;
  key_type?: string;
  key_base64?: string;
}

const HostKeyPrompt: React.FC = () => {
  const [open, setOpen] = useState(false);
  const [payload, setPayload] = useState<Payload | null>(null);

  useEffect(() => {
    const unsubscribe = Events.On("eventHostKeyPrompt", (event: any) => {
      try {
        const data =
          typeof event.data === "string" ? JSON.parse(event.data) : event.data;
        setPayload(data as Payload);
        setOpen(true);
      } catch (e) {
        console.error("Invalid host key prompt payload", e);
      }
    });

    return () => {
      unsubscribe();
    };
  }, []);

  const handleClose = async (trust: boolean) => {
    if (payload) {
      try {
        await SSHService.SetHostKeyDecision(payload.host, trust);
      } catch (err) {
        console.error("Failed to send host key decision", err);
      }
    }
    setOpen(false);
    setPayload(null);
  };

  return (
    <Dialog
      open={open}
      onClose={() => handleClose(false)}
      maxWidth="sm"
      fullWidth
    >
      <DialogTitle>新的主机密钥</DialogTitle>
      <DialogContent>
        {payload && (
          <>
            <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
              检测到未知主机密钥。主机: {payload.host} ({payload.address})
            </Typography>
            <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
              指纹: {payload.fingerprint}
            </Typography>
            <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
              若信任此主机，选择“信任并继续”，将会把该主机密钥保存到文件中。
            </Typography>
          </>
        )}
      </DialogContent>
      <DialogActions>
        <Button onClick={() => handleClose(false)}>拒绝</Button>
        <Button onClick={() => handleClose(true)} variant="contained">
          信任并继续
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default HostKeyPrompt;
