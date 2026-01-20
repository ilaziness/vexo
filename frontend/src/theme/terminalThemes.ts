import { TerminalThemes } from "../types/ssh";

/**
 * SSH GUI 客户端终端主题配置
 *
 * 本文件定义了终端(XTerm.js)的配色方案，与总体主题系统保持一致的设计理念。
 * 每个终端主题都适配对应的总体主题配色，提供一致的视觉体验。
 *
 * 设计理念：
 * - 护眼舒适：避免高对比度，减少长时间使用时的视觉疲劳
 * - 一致性：终端主题与总体UI主题在色彩上保持协调
 * - 可读性：确保文本在背景上有足够的对比度
 * - 专业性：采用适合终端应用的配色方案
 *
 * 主题映射关系：
 * - light: 对应总体亮色主题，深色文本，浅色背景
 * - dark: 对应总体默认暗色主题，经典的终端配色
 * - blueDark: 对应总体蓝色深色主题，护眼蓝色系
 * - atom: 对应总体极客黑主题，Atom编辑器风格配色
 * - deep: 对应总体深邃夜主题，纯黑背景高对比度
 * - eyeCare: 对应总体护眼模式，暖色系护眼配色
 * - terminal: 对应总体终端CLI模式，经典绿色终端配色
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
  // 亮色模式 - 使用深色文本，浅色背景，适合白天使用
  light: {
    background: "#ffffff",
    foreground: "#000000",
    cursor: "#000000",
    cursorAccent: "#ffffff",
    selectionBackground: "rgba(100, 149, 237, 0.3)",
    selectionForeground: "#ffffff",
    black: "#000000",
    red: "#cd0000",
    green: "#00cd00",
    yellow: "#cdcd00",
    blue: "#0000cd",
    magenta: "#cd00cd",
    cyan: "#00cdcd",
    white: "#e5e5e5",
    brightBlack: "#7f7f7f",
    brightRed: "#ff0000",
    brightGreen: "#00ff00",
    brightYellow: "#ffff00",
    brightBlue: "#5c5cff",
    brightMagenta: "#ff00ff",
    brightCyan: "#00ffff",
    brightWhite: "#ffffff",
  },
  // 默认暗色模式 - 经典的终端配色方案
  dark: {
    background: "#1e1e1e",
    foreground: "#cccccc",
    cursor: "#ffffff",
    cursorAccent: "#000000",
    selectionBackground: "rgba(255, 255, 255, 0.3)",
    selectionForeground: "#000000",
    black: "#000000",
    red: "#cd3131",
    green: "#0dbc79",
    yellow: "#e5e510",
    blue: "#2472c8",
    magenta: "#bc3fbc",
    cyan: "#11a8cd",
    white: "#e5e5e5",
    brightBlack: "#666666",
    brightRed: "#f14c4c",
    brightGreen: "#23d18b",
    brightYellow: "#f5f543",
    brightBlue: "#3b8eea",
    brightMagenta: "#d670d6",
    brightCyan: "#29b8db",
    brightWhite: "#ffffff",
  },
  // 蓝色深色模式 - 护眼蓝色系，适合长时间工作
  blueDark: {
    background: "#0a0e27",
    foreground: "#e3f2fd",
    cursor: "#90caf9",
    cursorAccent: "#0a0e27",
    selectionBackground: "rgba(144, 202, 249, 0.3)",
    selectionForeground: "#000000",
    black: "#000000",
    red: "#f44336",
    green: "#66bb6a",
    yellow: "#ffa726",
    blue: "#29b6f6",
    magenta: "#ab47bc",
    cyan: "#26c6da",
    white: "#e3f2fd",
    brightBlack: "#546e7a",
    brightRed: "#ef5350",
    brightGreen: "#81c784",
    brightYellow: "#ffb74d",
    brightBlue: "#4fc3f7",
    brightMagenta: "#ba68c8",
    brightCyan: "#4dd0e1",
    brightWhite: "#ffffff",
  },
  // 极客黑模式 - Atom编辑器风格配色，程序员偏好
  atom: {
    background: "#282c34",
    foreground: "#abb2bf",
    cursor: "#61afef",
    cursorAccent: "#282c34",
    selectionBackground: "rgba(97, 175, 239, 0.3)",
    selectionForeground: "#ffffff",
    black: "#000000",
    red: "#e06c75",
    green: "#98c379",
    yellow: "#e5c07b",
    blue: "#61afef",
    magenta: "#c678dd",
    cyan: "#56b6c2",
    white: "#abb2bf",
    brightBlack: "#5c6370",
    brightRed: "#e06c75",
    brightGreen: "#98c379",
    brightYellow: "#e5c07b",
    brightBlue: "#61afef",
    brightMagenta: "#c678dd",
    brightCyan: "#56b6c2",
    brightWhite: "#ffffff",
  },
  // 深邃夜模式 - 纯黑背景，高对比度护眼设计
  deep: {
    background: "#000000",
    foreground: "rgba(255, 255, 255, 0.85)",
    cursor: "#ffffff",
    cursorAccent: "#000000",
    selectionBackground: "rgba(255, 255, 255, 0.2)",
    selectionForeground: "#000000",
    black: "#000000",
    red: "#ff4d4f",
    green: "#52c41a",
    yellow: "#faad14",
    blue: "#1890ff",
    magenta: "#722ed1",
    cyan: "#13c2c2",
    white: "rgba(255, 255, 255, 0.85)",
    brightBlack: "#434343",
    brightRed: "#ff7875",
    brightGreen: "#73d13d",
    brightYellow: "#ffc069",
    brightBlue: "#40a9ff",
    brightMagenta: "#b37feb",
    brightCyan: "#36cfc9",
    brightWhite: "#ffffff",
  },
  // 护眼模式 - 暖色系配色，专门为护眼设计
  eyeCare: {
    background: "#fdf6e3",
    foreground: "#423629",
    cursor: "#5d4037",
    cursorAccent: "#fdf6e3",
    selectionBackground: "rgba(93, 64, 55, 0.2)",
    selectionForeground: "#ffffff",
    black: "#000000",
    red: "#a0522d",
    green: "#556b2f",
    yellow: "#cd853f",
    blue: "#5d4037",
    magenta: "#8b4513",
    cyan: "#daa520",
    white: "#423629",
    brightBlack: "#9c8c74",
    brightRed: "#daa520",
    brightGreen: "#556b2f",
    brightYellow: "#cd853f",
    brightBlue: "#8b4513",
    brightMagenta: "#a0522d",
    brightCyan: "#daa520",
    brightWhite: "#000000",
  },
  // 终端CLI模式 - 经典绿色终端配色，模仿phosphor monitor
  terminal: {
    background: "#0a0a0a",
    foreground: "#33ff00",
    cursor: "#33ff00",
    cursorAccent: "#0a0a0a",
    selectionBackground: "rgba(51, 255, 0, 0.3)",
    selectionForeground: "#000000",
    black: "#000000",
    red: "#ff3333",
    green: "#33ff00",
    yellow: "#ffb000",
    blue: "#1f521f",
    magenta: "#ff33ff",
    cyan: "#33ffff",
    white: "#1f521f",
    brightBlack: "#0f290f",
    brightRed: "#ff6666",
    brightGreen: "#66ff33",
    brightYellow: "#ffcc33",
    brightBlue: "#33ff33",
    brightMagenta: "#ff66ff",
    brightCyan: "#66ffff",
    brightWhite: "#33ff00",
  },
};

export default terminalThemes;
