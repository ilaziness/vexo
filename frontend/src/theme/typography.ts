/**
 * 应用默认等宽字体栈。
 * Noto Sans Mono 由 public/style.css @font-face 加载，与终端默认字体一致。
 */
export const CODE_FONT_FAMILY = [
  '"Noto Sans Mono"',
  '"Cascadia Code"',
  '"Fira Code"',
  "Menlo",
  "Monaco",
  "Consolas",
  '"DejaVu Sans Mono"',
  '"Ubuntu Mono"',
  '"Courier New"',
  "monospace",
].join(", ");

/** 行内 code / 代码块共用的排版 */
export const CODE_TYPOGRAPHY = {
  fontFamily: CODE_FONT_FAMILY,
  fontSize: "0.875em",
  lineHeight: 1.6,
} as const;
