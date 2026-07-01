import React, { useState } from 'react';
import {
  Box,
  IconButton,
  Typography,
  Drawer,
  List,
  ListItem,
  ListItemButton,
  ListItemText,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogContentText,
  DialogActions,
  Button,
  Tooltip,
} from '@mui/material';
import { Menu, Add, Delete, Close } from '@mui/icons-material';
import { useAIAssistantStore } from '../../stores/aiAssistant';
import { getSidebarContentWidth } from '../../func/aiSidebar';

interface SessionPanelProps {
  drawerContainer: HTMLDivElement | null;
}

function formatRelativeTime(unixSeconds: number): string {
  const diff = Date.now() - unixSeconds * 1000;
  const minutes = Math.floor(diff / 60000);
  if (minutes < 1) return '刚刚';
  if (minutes < 60) return `${minutes} 分钟前`;
  const hours = Math.floor(minutes / 60);
  if (hours < 24) return `${hours} 小时前`;
  const days = Math.floor(hours / 24);
  return `${days} 天前`;
}

const SessionPanel: React.FC<SessionPanelProps> = ({ drawerContainer }) => {
  const sessions = useAIAssistantStore((s) => s.sessions);
  const activeSessionId = useAIAssistantStore((s) => s.activeSessionId);
  const sidebarWidth = useAIAssistantStore((s) => s.sidebarWidth);
  const historyDrawerOpen = useAIAssistantStore((s) => s.historyDrawerOpen);
  const isStreaming = useAIAssistantStore((s) => s.isStreaming);
  const setHistoryDrawerOpen = useAIAssistantStore((s) => s.setHistoryDrawerOpen);
  const setSidebarOpen = useAIAssistantStore((s) => s.setSidebarOpen);
  const selectSession = useAIAssistantStore((s) => s.selectSession);
  const createSession = useAIAssistantStore((s) => s.createSession);
  const deleteSession = useAIAssistantStore((s) => s.deleteSession);

  const [deleteTargetId, setDeleteTargetId] = useState<string | null>(null);

  const activeSession = sessions.find((s) => s.id === activeSessionId);
  const deleteTarget = sessions.find((s) => s.id === deleteTargetId);
  const drawerPaperWidth = getSidebarContentWidth(sidebarWidth);

  const handleConfirmDelete = async () => {
    if (!deleteTargetId) return;
    await deleteSession(deleteTargetId);
    setDeleteTargetId(null);
  };

  const toggleHistoryDrawer = () => setHistoryDrawerOpen(!historyDrawerOpen);

  return (
    <>
      <Box
        sx={{
          height: 40,
          minHeight: 40,
          display: 'flex',
          alignItems: 'center',
          px: 1,
          gap: 0.5,
          borderBottom: 1,
          borderColor: 'divider',
          flexShrink: 0,
        }}
      >
        <Tooltip title={historyDrawerOpen ? '关闭会话历史' : '会话历史'}>
          <span>
            <IconButton
              size="small"
              disabled={isStreaming || !drawerContainer}
              onClick={toggleHistoryDrawer}
              aria-label={historyDrawerOpen ? '关闭会话历史' : '打开会话历史'}
              aria-expanded={historyDrawerOpen}
            >
              <Menu fontSize="small" />
            </IconButton>
          </span>
        </Tooltip>
        <Typography
          variant="body2"
          noWrap
          sx={{ flex: 1, fontWeight: 500, minWidth: 0 }}
        >
          {activeSession?.title || '新会话'}
        </Typography>
        <Tooltip title="新建会话">
          <span>
            <IconButton
              size="small"
              disabled={isStreaming}
              onClick={() => void createSession()}
              aria-label="新建会话"
            >
              <Add fontSize="small" />
            </IconButton>
          </span>
        </Tooltip>
        <Tooltip title="关闭 AI 助手">
          <IconButton
            size="small"
            onClick={() => setSidebarOpen(false)}
            aria-label="关闭 AI 助手"
          >
            <Close fontSize="small" />
          </IconButton>
        </Tooltip>
      </Box>

      {drawerContainer && (
        <Drawer
          anchor="left"
          open={historyDrawerOpen}
          onClose={() => setHistoryDrawerOpen(false)}
          variant="temporary"
          container={drawerContainer}
          slotProps={{
            root: { container: drawerContainer },
            paper: {
              sx: {
                position: 'absolute',
                width: drawerPaperWidth,
                height: '100%',
                boxSizing: 'border-box',
              },
            },
          }}
          sx={{
            position: 'absolute',
            zIndex: 1,
            '& .MuiBackdrop-root': { position: 'absolute' },
          }}
        >
          <Box sx={{ display: 'flex', flexDirection: 'column', height: '100%' }}>
            <Box
              sx={{
                px: 2,
                py: 1.5,
                borderBottom: 1,
                borderColor: 'divider',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'space-between',
              }}
            >
              <Typography variant="subtitle2">会话历史</Typography>
              <Button size="small" onClick={() => setHistoryDrawerOpen(false)}>
                关闭
              </Button>
            </Box>
            <List sx={{ flex: 1, overflow: 'auto', py: 0 }}>
              {sessions.map((session) => (
                <ListItem
                  key={session.id}
                  disablePadding
                  secondaryAction={
                    <IconButton
                      edge="end"
                      size="small"
                      disabled={isStreaming}
                      aria-label="删除会话"
                      onClick={() => setDeleteTargetId(session.id)}
                    >
                      <Delete fontSize="small" />
                    </IconButton>
                  }
                >
                  <ListItemButton
                    selected={session.id === activeSessionId}
                    disabled={isStreaming}
                    onClick={() => selectSession(session.id)}
                  >
                    <ListItemText
                      primary={session.title}
                      secondary={formatRelativeTime(session.updatedAt)}
                      slotProps={{
                        primary: { noWrap: true },
                        secondary: { variant: 'caption' },
                      }}
                    />
                  </ListItemButton>
                </ListItem>
              ))}
            </List>
          </Box>
        </Drawer>
      )}

      <Dialog open={deleteTargetId !== null} onClose={() => setDeleteTargetId(null)}>
        <DialogTitle>删除会话</DialogTitle>
        <DialogContent>
          <DialogContentText>
            确定删除会话「{deleteTarget?.title || '新会话'}」？此操作不可撤销。
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setDeleteTargetId(null)}>取消</Button>
          <Button color="error" onClick={() => void handleConfirmDelete()}>
            删除
          </Button>
        </DialogActions>
      </Dialog>
    </>
  );
};

export default SessionPanel;
