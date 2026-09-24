import { Box, Typography } from '@mui/material';
import { ChatConfirmation } from '@mui/x-chat';
import { useChat, type ChatPartRendererMap } from '@mui/x-chat/headless';
import { PlanTaskPart, planFromToolPart } from './PlanTaskPart';
import { SSH_TARGET_INPUT_KEY } from '../../func/aiContext';

/** Custom renderers must cover every tool of that part type — returning null hides the built-in widget. */
export const aiPartRenderers: ChatPartRendererMap = {
  'dynamic-tool': ({ part }) => <AgentToolPart part={part as any} />,
  tool: ({ part }) => <AgentToolPart part={part as any} />,
};

function AgentToolPart({ part }: { part: any }) {
  const { addToolApprovalResponse } = useChat();
  const inv = part?.toolInvocation;
  if (!inv) return null;

  const plan = planFromToolPart(part);
  if (plan) {
    return <PlanTaskPart plan={plan} />;
  }

  const approvalId = String(inv.approvalId || inv.toolCallId || '');
  const command = sshCommandFromInput(inv.input);
  const sshTarget = sshTargetFromInput(inv.input);
  const title = String(inv.toolName || 'tool');

  if (inv.state === 'approval-requested') {
    const hostLine = sshTarget ? `在 ${sshTarget} 执行？` : '批准执行 SSH 命令？';
    return (
      <ChatConfirmation
        message={command ? `${hostLine}\n$ ${command}` : `批准调用 ${title}？`}
        confirmLabel="批准"
        cancelLabel="拒绝"
        onConfirm={() => {
          void addToolApprovalResponse({ id: approvalId, approved: true });
        }}
        onCancel={() => {
          void addToolApprovalResponse({
            id: approvalId,
            approved: false,
            reason: 'user denied',
          });
        }}
      />
    );
  }

  return (
    <Box
      sx={{
        border: 1,
        borderColor: 'divider',
        borderRadius: 1,
        px: 1.25,
        py: 0.75,
        bgcolor: 'background.paper',
      }}
    >
      <Typography variant="caption" color="text.secondary">
        {title} · {toolStateLabel(inv.state)}
        {sshTarget ? ` · ${sshTarget}` : ''}
      </Typography>
      {command ? (
        <Typography variant="body2" sx={{ fontFamily: 'monospace', whiteSpace: 'pre-wrap' }}>
          $ {command}
        </Typography>
      ) : null}
      {inv.output != null ? (
        <Typography variant="body2" sx={{ whiteSpace: 'pre-wrap', mt: 0.5 }}>
          {formatUnknown(inv.output)}
        </Typography>
      ) : null}
      {inv.errorText ? (
        <Typography variant="body2" color="error">
          {String(inv.errorText)}
        </Typography>
      ) : null}
      {inv.approval?.reason ? (
        <Typography variant="caption" color="text.secondary">
          {String(inv.approval.reason)}
        </Typography>
      ) : null}
    </Box>
  );
}

function sshCommandFromInput(input: unknown): string {
  const v = parseJSONValue(input);
  if (v && typeof v === 'object' && 'command' in (v as object)) {
    return String((v as { command?: unknown }).command ?? '');
  }
  return '';
}

function sshTargetFromInput(input: unknown): string {
  const v = parseJSONValue(input);
  if (v && typeof v === 'object' && SSH_TARGET_INPUT_KEY in (v as object)) {
    return String((v as Record<string, unknown>)[SSH_TARGET_INPUT_KEY] ?? '').trim();
  }
  return '';
}

function parseJSONValue(v: unknown): unknown {
  if (typeof v === 'string') {
    try {
      return JSON.parse(v);
    } catch {
      return v;
    }
  }
  return v;
}

function formatUnknown(v: unknown): string {
  const parsed = parseJSONValue(v);
  if (parsed && typeof parsed === 'object' && 'output' in (parsed as object)) {
    return String((parsed as { output?: unknown }).output ?? '');
  }
  if (typeof parsed === 'string') return parsed;
  try {
    return JSON.stringify(parsed, null, 2);
  } catch {
    return String(parsed);
  }
}

function toolStateLabel(state: string): string {
  switch (state) {
    case 'input-streaming':
      return '准备中';
    case 'input-available':
      return '待执行';
    case 'approval-responded':
      return '已批准';
    case 'output-available':
      return '已完成';
    case 'output-error':
      return '失败';
    case 'output-denied':
      return '已拒绝';
    default:
      return state || '';
  }
}
