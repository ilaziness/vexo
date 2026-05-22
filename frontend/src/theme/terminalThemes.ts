import { TerminalThemes } from "../types/ssh";

/**
 * SSH GUI 客户端终端主题配置
 *
 * 本文件定义了终端(XTerm.js)的配色方案，与总体主题系统保持一致的设计理念。
 * 每个终端主题都适配对应的总体主题配色，提供一致的视觉体验。
 *
 * 主题映射关系：
 * - light (Daylight): 对应总体亮色主题，白底深色文本
 * - dark (Midnight): 对应总体暗色主题，深蓝灰底霓虹蓝青
 * - eyeCare (Amber): 对应总体护眼主题，暖琥珀底深棕字
 *
 * ANSI颜色说明：
 * - black/red/green/yellow/blue/magenta/cyan/white: 基础ANSI颜色
 * - bright前缀: 对应的亮色版本
 * - background: 终端背景色
 * - foreground: 默认文本颜色
 * - cursor: 光标颜色
 * - cursorAccent: 光标背景色（反色）
 * - selection: 选中文本的背景色
 */
const terminalThemes: TerminalThemes = {
  // Daylight - 白底黑字，标准ANSI色，清爽专业
  light: {
    background: "#f8f9fa",
    foreground: "#1f2937",
    cursor: "#2563eb",
    cursorAccent: "#f8f9fa",
    selectionBackground: "#2563eb",
    selectionForeground: "#ffffff",
    black: "#1f2937",
    red: "#ef4444",
    green: "#10b981",
    yellow: "#f59e0b",
    blue: "#2563eb",
    magenta: "#8b5cf6",
    cyan: "#0891b2",
    white: "#6b7280",
    brightBlack: "#9ca3af",
    brightRed: "#f87171",
    brightGreen: "#34d399",
    brightYellow: "#fbbf24",
    brightBlue: "#60a5fa",
    brightMagenta: "#a78bfa",
    brightCyan: "#22d3ee",
    brightWhite: "#e5e7eb",
  },
  // Midnight - 深蓝灰底，霓虹蓝青，夜空感
  dark: {
    background: "#0f172a",
    foreground: "#e2e8f0",
    cursor: "#22d3ee",
    cursorAccent: "#0f172a",
    selectionBackground: "#22d3ee",
    selectionForeground: "#0f172a",
    black: "#0f172a",
    red: "#f87171",
    green: "#34d399",
    yellow: "#fbbf24",
    blue: "#60a5fa",
    magenta: "#c084fc",
    cyan: "#22d3ee",
    white: "#e2e8f0",
    brightBlack: "#475569",
    brightRed: "#fca5a5",
    brightGreen: "#6ee7b7",
    brightYellow: "#fcd34d",
    brightBlue: "#93c5fd",
    brightMagenta: "#d8b4fe",
    brightCyan: "#67e8f9",
    brightWhite: "#f8fafc",
  },
  // Amber - 暖琥珀底，深棕字，复古终端感
  eyeCare: {
    background: "#fef3c7",
    foreground: "#451a03",
    cursor: "#92400e",
    cursorAccent: "#fef3c7",
    selectionBackground: "#92400e",
    selectionForeground: "#fef3c7",
    black: "#451a03",
    red: "#c2410c",
    green: "#65a30d",
    yellow: "#d97706",
    blue: "#78350f",
    magenta: "#7c2d12",
    cyan: "#0d9488",
    white: "#92400e",
    brightBlack: "#a16207",
    brightRed: "#ea580c",
    brightGreen: "#84cc16",
    brightYellow: "#fbbf24",
    brightBlue: "#b45309",
    brightMagenta: "#9a3412",
    brightCyan: "#14b8a6",
    brightWhite: "#78350f",
  },
};

export default terminalThemes;
