// Tool 工具定义
export interface Tool {
  id: string;
  name: string;
  description: string;
  icon: string;
  category: string;
}

// PortCheckResult 端口检测结果
export interface PortCheckResult {
  success: boolean;
  host: string;
  port: number;
  responseTime: number; // 毫秒
  error?: string;
}

// EncodeRequest 编码请求
export interface EncodeRequest {
  toolType: string; // base64, url, html
  input: string;
}

// EncodeResponse 编码响应
export interface EncodeResponse {
  result: string;
  error?: string;
}

// Match 正则匹配结果
export interface Match {
  text: string;
  index: number;
  groups: string[];
}

// RegexMatchResult 正则匹配结果
export interface RegexMatchResult {
  matches: Match[];
  count: number;
  error?: string;
}

// HashResult 哈希计算结果
export interface HashResult {
  success: boolean;
  input: string;
  algorithm: string;
  result: string;
  error?: string;
}

// TimestampResult 时间戳转换结果
export interface TimestampResult {
  success: boolean;
  timestamp?: number;
  datetime?: string;
  error?: string;
}
