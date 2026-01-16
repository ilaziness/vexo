export interface SSHLinkInfo {
  host: string;
  port: number;
  user: string;
  password?: string;
  key?: string;
  keyPassword?: string;
}

export interface SSHTab {
  index: string;
  name: string;
  sshInfo?: SSHLinkInfo;
}

export interface Message {
  open: boolean;
  text: string;
  type: "error" | "success" | "info";
}

export interface MessageStore {
  message: Message;
  setClose: () => void;
  errorMessage: (message: string) => void;
  successMessage: (message: string) => void;
}

// 终端主题类型定义
export interface TerminalTheme {
  background: string;
  foreground: string;
  cursor: string;
  cursorAccent: string;
  selection: string;
  black: string;
  red: string;
  green: string;
  yellow: string;
  blue: string;
  magenta: string;
  cyan: string;
  white: string;
  brightBlack: string;
  brightRed: string;
  brightGreen: string;
  brightYellow: string;
  brightBlue: string;
  brightMagenta: string;
  brightCyan: string;
  brightWhite: string;
}

// 终端主题配置映射
export type TerminalThemes = {
  light: TerminalTheme;
  dark: TerminalTheme;
  blueDark: TerminalTheme;
  atom: TerminalTheme;
  deep: TerminalTheme;
  eyeCare: TerminalTheme;
};

export type TerminalThemeMode = keyof TerminalThemes;
