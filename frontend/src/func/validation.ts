/**
 * 验证工具函数
 */

// 验证 IP:Port 格式
export const validateIpPort = (addr: string): { valid: boolean; error: string } => {
  if (!addr.trim()) {
    return { valid: false, error: "请输入地址" };
  }

  const parts = addr.split(":");
  if (parts.length !== 2) {
    return { valid: false, error: "格式应为 ip:port" };
  }

  const [ip, portStr] = parts;

  // 验证 IP 地址
  const ipRegex = /^(\d{1,3}\.){3}\d{1,3}$/;
  if (!ipRegex.test(ip)) {
    return { valid: false, error: "IP 地址格式不正确" };
  }

  // 验证 IP 各部分范围
  const ipParts = ip.split(".");
  for (const part of ipParts) {
    const num = Number.parseInt(part, 10);
    if (num < 0 || num > 255) {
      return { valid: false, error: "IP 地址各部分应在 0-255 之间" };
    }
  }

  // 验证端口
  const port = Number.parseInt(portStr, 10);
  if (Number.isNaN(port) || port < 1 || port > 65535) {
    return { valid: false, error: "端口必须是 1-65535 之间的数字" };
  }

  return { valid: true, error: "" };
};

// 验证远程地址格式 (host:port)
export const validateRemoteAddr = (
  addr: string,
): { valid: boolean; error: string } => {
  if (!addr.trim()) {
    return { valid: false, error: "请输入远程地址" };
  }

  const parts = addr.split(":");
  if (parts.length !== 2) {
    return { valid: false, error: "格式应为 host:port" };
  }

  const [host, portStr] = parts;
  if (!host.trim()) {
    return { valid: false, error: "主机名不能为空" };
  }

  const port = Number.parseInt(portStr, 10);
  if (Number.isNaN(port) || port < 1 || port > 65535) {
    return { valid: false, error: "端口必须是 1-65535 之间的数字" };
  }

  return { valid: true, error: "" };
};

// 验证端口号
export const validatePort = (portStr: string): { valid: boolean; error: string } => {
  if (!portStr.trim()) {
    return { valid: false, error: "请输入端口号" };
  }

  const port = Number.parseInt(portStr, 10);
  if (Number.isNaN(port) || port < 1 || port > 65535) {
    return { valid: false, error: "端口必须是 1-65535 之间的数字" };
  }

  return { valid: true, error: "" };
};

// 验证 SSH 隧道表单数据
export interface TunnelFormData {
  localAddr: string;
  remoteAddr: string;
  remotePort: string;
}

export interface TunnelFormErrors {
  localAddr: string;
  remoteAddr: string;
  remotePort: string;
}

export type TunnelType = "local" | "remote" | "dynamic";

/**
 * 计算隧道表单的错误状态
 */
export const calculateTunnelFormErrors = (
  tunnelType: TunnelType,
  formData: TunnelFormData,
): TunnelFormErrors => {
  const newErrors: TunnelFormErrors = {
    localAddr: "",
    remoteAddr: "",
    remotePort: "",
  };
  
  // 本地地址验证
  const localResult = validateIpPort(formData.localAddr);
  newErrors.localAddr = localResult.error;
  
  // 远程地址验证（仅本地转发）
  if (tunnelType === "local") {
    const remoteResult = validateRemoteAddr(formData.remoteAddr);
    newErrors.remoteAddr = remoteResult.error;
  }
  
  // 远程端口验证（仅远程转发）
  if (tunnelType === "remote") {
    const portResult = validatePort(formData.remotePort);
    newErrors.remotePort = portResult.error;
  }
  
  return newErrors;
};

/**
 * 验证隧道表单是否有效
 */
export const isTunnelFormValid = (
  tunnelType: TunnelType,
  formData: TunnelFormData,
): boolean => {
  const errors = calculateTunnelFormErrors(tunnelType, formData);
  return !errors.localAddr && !errors.remoteAddr && !errors.remotePort;
};

/**
 * 获取隧道类型显示名称
 */
export const getTunnelTypeName = (tunnelType: TunnelType): string => {
  switch (tunnelType) {
    case "local":
      return "本地端口转发";
    case "remote":
      return "远程端口转发";
    case "dynamic":
      return "动态端口转发";
    default:
      return "隧道";
  }
};

/**
 * 获取本地地址的 placeholder 提示
 */
export const getLocalAddrPlaceholder = (tunnelType: TunnelType): string => {
  return tunnelType === "remote" 
    ? "例如: 127.0.0.1:3306" 
    : "例如: 127.0.0.1:8080";
};

/**
 * 获取本地地址的 helper text
 */
export const getLocalAddrHelperText = (
  tunnelType: TunnelType,
  error: string,
): string => {
  if (error) return error;
  return tunnelType === "remote"
    ? "本地目标地址，格式: ip:port"
    : "本地监听地址，格式: ip:port";
};

/**
 * 获取远程地址的 helper text
 */
export const getRemoteAddrHelperText = (error: string): string => {
  return error || "远程目标地址，格式: host:port";
};

/**
 * 获取远程端口的 helper text
 */
export const getRemotePortHelperText = (error: string): string => {
  return error || "远端将在 127.0.0.1:<端口> 上监听";
};