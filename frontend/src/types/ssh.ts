import type { ITheme } from "@xterm/xterm";

export enum ConnectionStatus {
  Connected = "connected",
  Disconnected = "disconnected",
  Connecting = "connecting",
}

export interface SSHLinkInfo {
  linkID?: string;
  bookmarkID?: string;
  host: string;
  port: number;
  user: string;
  password?: string;
  key?: string;
  keyPassword?: string;
  proxyJumpID?: string;
  useAgent?: boolean;
  certificate?: string;
}

/** Matches services.ConnectRequest JSON fields (camelCase). */
export interface ConnectRequest {
  host: string;
  port: number;
  user: string;
  password: string;
  key: string;
  keyPassword: string;
  proxyJumpID: string;
  useAgent: boolean;
  certificate: string;
}

/** Build ConnectRequest with stable defaults for optional SSHLinkInfo fields. */
export function toConnectRequest(li: SSHLinkInfo): ConnectRequest {
  const key = li.key || "";
  return {
    host: li.host,
    port: li.port,
    user: li.user,
    password: li.password || "",
    key,
    keyPassword: key ? li.keyPassword || "" : "",
    proxyJumpID: li.proxyJumpID || "",
    useAgent: !!li.useAgent,
    certificate: li.certificate || "",
  };
}

export interface SSHTab {
  index: string;
  name: string;
  sshInfo?: SSHLinkInfo;
  connectionStatus?: ConnectionStatus;
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
  eyeCare: ITheme;
};

// 应用主题类型定义
export type AppTheme = "light" | "dark" | "eyeCare";
