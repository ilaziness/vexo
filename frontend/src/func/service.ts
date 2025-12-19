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
  return `${new Date().getMilliseconds()}`;
};

export const decodeBase64 = (base64: string): Uint8Array => {
  const binStr = atob(base64);
  return Uint8Array.from(binStr, (c) => c.charCodeAt(0));
};

export const encodeBase64 = (data: string | Uint8Array): string => {
  let uint8Array: Uint8Array;

  if (typeof data === "string") {
    // xterm.js 的 onData 回调传入的是 string，但它代表原始字节序列（如 "\x1b[D"）
    // 必须按 Latin1（即每个字符 charCode = byte）方式转为 Uint8Array
    uint8Array = Uint8Array.from(data, (c) => c.charCodeAt(0));
  } else {
    uint8Array = data;
  }
  // 转为二进制字符串（每个字节 → 一个字符）
  let binary = "";
  for (let i = 0; i < uint8Array.length; i++) {
    binary += String.fromCharCode(uint8Array[i]);
  }
  // Base64 编码
  return btoa(binary);
};
