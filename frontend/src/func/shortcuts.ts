export const EVENT_SHORTCUT = "eventShortcut";

export enum ShortcutAction {
  TerminalFind = "terminal.find",
  TerminalCopy = "terminal.copy",
  TerminalPaste = "terminal.paste",
  TerminalAddToChat = "terminal.addToChat",
  TerminalToggleLog = "terminal.toggleLog",
  TerminalClear = "terminal.clear",
  PanelAI = "panel.ai",
  PanelBookmarks = "panel.bookmarks",
}

export function parseShortcutAction(data: unknown): string | null {
  try {
    const payload =
      typeof data === "string"
        ? (JSON.parse(data) as { action?: string })
        : (data as { action?: string } | undefined);
    return payload?.action ?? null;
  } catch {
    return null;
  }
}

export function toShellCodeBlock(text: string): string {
  return "```shell\n" + text + "\n```";
}

/** True when focus is in a real UI editor, not xterm's helper textarea. */
export function isEditableFocus(): boolean {
  const el = document.activeElement;
  if (!el || !(el instanceof HTMLElement)) {
    return false;
  }
  // xterm keeps focus on a hidden textarea inside .xterm.
  if (el.closest(".xterm")) {
    return false;
  }
  const tag = el.tagName;
  if (tag === "INPUT" || tag === "TEXTAREA" || tag === "SELECT") {
    return true;
  }
  if (el.isContentEditable) {
    return true;
  }
  return Boolean(el.closest("[contenteditable='true']"));
}
