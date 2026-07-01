const AI_NOT_ENABLED =
  "AI 助手尚未启用，请前往「设置 → AI」完成配置并启用后再试";

const LEGACY_ERROR_PATTERNS = [
  /ai engine not initialized/i,
  /ai 助手尚未启用/i,
  /ai 助手初始化失败/i,
  /ai 配置加载失败/i,
];

/** 将 AI 对话相关错误转换为用户可读提示 */
export const formatAIChatError = (raw: string): string => {
  const message = raw.trim().replace(/^AI chat failed:\s*/i, "").trim();
  if (LEGACY_ERROR_PATTERNS.some((pattern) => pattern.test(message))) {
    return AI_NOT_ENABLED;
  }
  if (message.startsWith("AI 助手")) {
    return AI_NOT_ENABLED;
  }
  return message;
};
