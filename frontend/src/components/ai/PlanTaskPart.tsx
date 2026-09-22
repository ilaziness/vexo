import {
  Box,
  Chip,
  Collapse,
  List,
  ListItem,
  ListItemIcon,
  ListItemText,
  Typography,
} from '@mui/material';
import CheckCircleOutlinedIcon from '@mui/icons-material/CheckCircleOutlined';
import RadioButtonUncheckedOutlinedIcon from '@mui/icons-material/RadioButtonUncheckedOutlined';
import ErrorOutlinedIcon from '@mui/icons-material/ErrorOutlined';
import PlayCircleOutlinedIcon from '@mui/icons-material/PlayCircleOutlined';

export type PlanTaskStatus = 'pending' | 'running' | 'done' | 'failed';

export interface PlanTask {
  id: string;
  title: string;
  status: PlanTaskStatus;
}

export interface PlanSnapshot {
  title: string;
  tasks: PlanTask[];
}

function asPlan(value: unknown): PlanSnapshot | null {
  if (!value || typeof value !== 'object') return null;
  const v = value as Record<string, unknown>;
  if (!Array.isArray(v.tasks)) return null;
  const tasks: PlanTask[] = v.tasks.map((t: any, i: number) => ({
    id: String(t?.id ?? i),
    title: String(t?.title ?? `任务 ${i + 1}`),
    status: normalizeStatus(t?.status),
  }));
  return {
    title: String(v.title ?? '计划'),
    tasks,
  };
}

function normalizeStatus(s: unknown): PlanTaskStatus {
  const v = String(s || 'pending').toLowerCase();
  if (v === 'running' || v === 'done' || v === 'failed' || v === 'pending') return v;
  return 'pending';
}

function StatusIcon({ status }: { status: PlanTaskStatus }) {
  switch (status) {
    case 'done':
      return <CheckCircleOutlinedIcon color="success" fontSize="small" />;
    case 'failed':
      return <ErrorOutlinedIcon color="error" fontSize="small" />;
    case 'running':
      return <PlayCircleOutlinedIcon color="primary" fontSize="small" />;
    default:
      return <RadioButtonUncheckedOutlinedIcon color="disabled" fontSize="small" />;
  }
}

function statusLabel(status: PlanTaskStatus): string {
  switch (status) {
    case 'done':
      return '完成';
    case 'failed':
      return '失败';
    case 'running':
      return '进行中';
    default:
      return '待办';
  }
}

export interface PlanTaskPartProps {
  plan: PlanSnapshot | null;
}

/** Renders upsert_plan tool output as a task list. */
export function PlanTaskPart({ plan }: PlanTaskPartProps) {
  if (!plan || plan.tasks.length === 0) {
    return (
      <Typography variant="body2" color="text.secondary">
        （空计划）
      </Typography>
    );
  }

  const allDone = plan.tasks.every((t) => t.status === 'done' || t.status === 'failed');

  return (
    <Collapse in appear>
      <Box
        sx={{
          border: 1,
          borderColor: 'divider',
          borderRadius: 1,
          px: 1.5,
          py: 1,
          bgcolor: 'background.paper',
          minWidth: 220,
        }}
      >
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 0.5 }}>
          <Typography variant="subtitle2" sx={{ fontWeight: 600, flex: 1 }}>
            {plan.title}
          </Typography>
          <Chip
            size="small"
            label={allDone ? '已结束' : '进行中'}
            color={allDone ? 'default' : 'primary'}
            variant="outlined"
          />
        </Box>
        <List dense disablePadding>
          {plan.tasks.map((task) => (
            <ListItem key={task.id} disableGutters sx={{ py: 0.25 }}>
              <ListItemIcon sx={{ minWidth: 32 }}>
                <StatusIcon status={task.status} />
              </ListItemIcon>
              <ListItemText
                primary={task.title}
                secondary={statusLabel(task.status)}
                slotProps={{
                  primary: { variant: 'body2' },
                  secondary: { variant: 'caption' },
                }}
              />
            </ListItem>
          ))}
        </List>
      </Box>
    </Collapse>
  );
}

export function planFromToolPart(part: {
  toolInvocation?: { input?: unknown; output?: unknown; toolName?: string };
}): PlanSnapshot | null {
  const inv = part.toolInvocation;
  if (!inv || inv.toolName !== 'upsert_plan') return null;
  return asPlan(inv.output) || asPlan(inv.input);
}
