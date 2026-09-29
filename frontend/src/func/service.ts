export const parseCallServiceError = (err: any): string => {
  if (err && typeof err === "object") {
    if ("message" in err && typeof err.message === "string") {
      try {
        const parsed = JSON.parse(err.message);
        if (
          parsed &&
          typeof parsed === "object" &&
          "message" in parsed &&
          typeof parsed.message === "string"
        ) {
          return parsed.message;
        }
      } catch {
        // 如果解析失败，继续使用原始消息
      }
      return err.message;
    }
    if ("error" in err && typeof err.error === "string") {
      return err.error;
    }
  }
  return String(err);
};

// 把 SSH 握手/认证失败转成可直接展示的短句。
export const formatSSHConnectError = (err: unknown): string => {
  const raw = parseCallServiceError(err).trim();
  if (raw.includes("已取消")) {
    return "已取消登录";
  }
  if (raw.includes("等待超时")) {
    return "验证超时，请重试";
  }
  if (
    raw.includes("unable to authenticate") ||
    raw.includes("no supported methods remain") ||
    raw.includes("没有可用认证")
  ) {
    return "登录失败，请检查密码或验证码";
  }
  return raw || "连接失败";
};

// genTabIndex generate tab index ID
export const genTabIndex = (): string => {
  return `${Date.now()}`;
};

export const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return "0 B";
  if (bytes < 1024) return bytes + " B";
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(2) + " KB";
  if (bytes < 1024 * 1024 * 1024)
    return (bytes / (1024 * 1024)).toFixed(2) + " MB";
  return (bytes / (1024 * 1024 * 1024)).toFixed(2) + " GB";
};

export const sleep = (ms: number): Promise<void> => {
  return new Promise((resolve) => setTimeout(resolve, ms));
};
