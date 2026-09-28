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

export interface FormatResult {
  success: boolean;
  result: string;
  error?: string;
}

export interface CronResult {
  success: boolean;
  nextTimes?: string[];
  error?: string;
}

export interface CIDRResult {
  success: boolean;
  network?: string;
  broadcast?: string;
  netmask?: string;
  prefixLen?: number;
  firstHost?: string;
  lastHost?: string;
  hostCount?: string;
  totalAddrs?: string;
  contains?: boolean | null;
  containsIP?: string;
  error?: string;
}

export interface JWTResult {
  success: boolean;
  header?: string;
  payload?: string;
  signature?: string;
  error?: string;
}

export interface RandomResult {
  success: boolean;
  result: string;
  error?: string;
}

export interface BaseConvertResult {
  success: boolean;
  result: string;
  error?: string;
}

export enum DiffLineType {
  Equal = "equal",
  Add = "add",
  Remove = "remove",
}

export interface DiffLine {
  type: DiffLineType | string;
  content: string;
  oldLine?: number;
  newLine?: number;
}

export interface DiffResult {
  success: boolean;
  lines: DiffLine[];
  error?: string;
}

export interface ChmodResult {
  success: boolean;
  octal?: string;
  symbolic?: string;
  error?: string;
}

export enum EncodingType {
  Base64 = "base64",
  URL = "url",
  HTML = "html",
}

export enum HashAlgorithm {
  MD5 = "md5",
  SHA1 = "sha1",
  SHA256 = "sha256",
  SHA512 = "sha512",
}

export enum JSONYAMLAction {
  Format = "format",
  Minify = "minify",
  Validate = "validate",
  Json2Yaml = "json2yaml",
  Yaml2Json = "yaml2json",
}

export enum RandomKind {
  UuidV4 = "uuid-v4",
  UuidV7 = "uuid-v7",
  String = "string",
}

export enum RandomCharset {
  Alphanum = "alphanum",
  Alpha = "alpha",
  Numeric = "numeric",
  Hex = "hex",
  Base64 = "base64",
  Custom = "custom",
}

export enum ChmodDirection {
  ToSymbolic = "toSymbolic",
  ToOctal = "toOctal",
}
