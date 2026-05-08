import { Menu, MenuItem, ListItemIcon, ListItemText } from "@mui/material";
import ContentCopyIcon from "@mui/icons-material/ContentCopy";
import ContentPasteIcon from "@mui/icons-material/ContentPaste";
import ClearAllIcon from "@mui/icons-material/ClearAll";
import AutoAwesomeIcon from '@mui/icons-material/AutoAwesome';
import { terminalInstances } from "../stores/terminalInstances";
import { useAICommandStore } from "../stores/aiCommand";

interface TerminalContextMenuProps {
  contextMenu: {
    mouseX: number;
    mouseY: number;
  } | null;
  onClose: () => void;
  linkID: string;
}

export default function TerminalContextMenu({
  contextMenu,
  onClose,
  linkID,
}: TerminalContextMenuProps) {
  const openAIPanel = useAICommandStore((state) => state.openPanel);
  const explainCommand = useAICommandStore((state) => state.explainCommand);

  const handleCopy = () => {
    const term = terminalInstances.get(linkID);
    const selection = term?.getSelection();
    if (selection) {
      navigator.clipboard.writeText(selection);
    }
    onClose();
  };

  const handleAIExplain = async () => {
    const term = terminalInstances.get(linkID);
    const selection = term?.getSelection();
    if (selection) {
      explainCommand(selection, linkID);
      await openAIPanel();
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
      <MenuItem onClick={handleClear}>
        <ListItemIcon>
          <ClearAllIcon fontSize="small" />
        </ListItemIcon>
        <ListItemText>清屏</ListItemText>
      </MenuItem>
      <MenuItem onClick={handleAIExplain}>
        <ListItemIcon>
          <AutoAwesomeIcon fontSize="small" />
        </ListItemIcon>
        <ListItemText>AI 解释此命令</ListItemText>
      </MenuItem>
    </Menu>
  );
}
