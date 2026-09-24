import { Menu, MenuItem, ListItemIcon, ListItemText } from "@mui/material";
import ContentCopyIcon from "@mui/icons-material/ContentCopy";
import ContentPasteIcon from "@mui/icons-material/ContentPaste";
import ClearAllIcon from "@mui/icons-material/ClearAll";
import ChatIcon from "@mui/icons-material/Chat";
import { terminalInstances } from "../stores/terminalInstances";
import { useAIAssistantStore } from "../stores/aiAssistant";

interface TerminalContextMenuProps {
  contextMenu: {
    mouseX: number;
    mouseY: number;
  } | null;
  onClose: () => void;
  linkID: string;
}

function toShellCodeBlock(text: string): string {
  return "```shell\n" + text + "\n```";
}

export default function TerminalContextMenu({
  contextMenu,
  onClose,
  linkID,
}: TerminalContextMenuProps) {
  const selection = terminalInstances.get(linkID)?.getSelection() ?? "";
  const hasSelection = selection.length > 0;

  const handleCopy = () => {
    if (selection) {
      navigator.clipboard.writeText(selection);
    }
    onClose();
  };

  const handlePaste = async () => {
    const term = terminalInstances.get(linkID);
    try {
      const text = await navigator.clipboard.readText();
      if (text) {
        term?.paste(text);
        term?.focus();
      }
    } catch (err) {
      console.error("Failed to read clipboard contents: ", err);
    }
    onClose();
  };

  const handleClear = () => {
    const term = terminalInstances.get(linkID);
    term?.clear();
    onClose();
  };

  const handleAddToChat = () => {
    if (!selection) {
      onClose();
      return;
    }
    useAIAssistantStore.getState().appendToComposer(toShellCodeBlock(selection));
    onClose();
  };

  return (
    <Menu
      open={contextMenu !== null}
      onClose={onClose}
      anchorReference="anchorPosition"
      anchorPosition={
        contextMenu !== null
          ? { top: contextMenu.mouseY, left: contextMenu.mouseX }
          : undefined
      }
    >
      <MenuItem onClick={handleCopy}>
        <ListItemIcon>
          <ContentCopyIcon fontSize="small" />
        </ListItemIcon>
        <ListItemText>复制</ListItemText>
      </MenuItem>
      <MenuItem onClick={handlePaste}>
        <ListItemIcon>
          <ContentPasteIcon fontSize="small" />
        </ListItemIcon>
        <ListItemText>粘贴</ListItemText>
      </MenuItem>
      <MenuItem onClick={handleAddToChat} disabled={!hasSelection}>
        <ListItemIcon>
          <ChatIcon fontSize="small" />
        </ListItemIcon>
        <ListItemText>添加到聊天</ListItemText>
      </MenuItem>
      <MenuItem onClick={handleClear}>
        <ListItemIcon>
          <ClearAllIcon fontSize="small" />
        </ListItemIcon>
        <ListItemText>清屏</ListItemText>
      </MenuItem>
    </Menu>
  );
}
