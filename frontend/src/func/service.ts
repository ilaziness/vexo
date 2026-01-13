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

export const getTabIndex = (): string => {
  return `${new Date().getTime()}`;
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
