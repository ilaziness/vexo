import type { ITheme } from "@xterm/xterm";

export type { ITheme };
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
  infoMessage: (message: string) => void;
}

// 终端主题配置映射
export type TerminalThemes = {
  light: ITheme;
  dark: ITheme;
  blueDark: ITheme;
  atom: ITheme;
  deep: ITheme;
  eyeCare: ITheme;
};

// 应用主题类型定义
export type AppTheme =
  | "light"
  | "dark"
  | "blueDark"
  | "atom"
  | "deep"
  | "eyeCare";

export type TerminalThemeMode = AppTheme;
