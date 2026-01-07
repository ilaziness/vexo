export interface SSHLinkInfo {
  host: string;
  port: number;
  user: string;
  password?: string;
  key?: string;
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
