import { createColorScheme, createTheme } from "@mui/material/styles";
import type {} from "@mui/x-chat/themeAugmentation";
import light from "./light";
import dark from "./dark";
import eyeCare from "./eyeCare";
import { AppTheme } from "../types/ssh";

declare module "@mui/material/styles" {
  interface ColorSchemeOverrides {
    eyeCare: true;
  }
}

/**
 * 主题选项配置
 */
export interface ThemeOption {
  value: AppTheme;
  label: string;
  mode: "light" | "dark";
  color: string;
}

export const ThemeOptions: ThemeOption[] = [
  {
    value: "light",
    label: "日间",
    mode: "light",
    color: "#2563eb",
  },
  {
    value: "dark",
    label: "夜间",
    mode: "dark",
    color: "#60a5fa",
  },
  {
    value: "eyeCare",
    label: "护眼",
    mode: "light",
    color: "#92400e",
  },
];

/**
 * SSH GUI 客户端主题配置
 * 设计理念：清爽简洁，适合终端应用场景
 */
const theme = createTheme({
  // 配置颜色方案
  colorSchemes: {
    // Daylight - 亮色专业模式
    light,
    // Midnight - 暗色夜空模式
    dark: createColorScheme(dark),
    // Amber - 护眼复古模式
    eyeCare: createColorScheme(eyeCare),
  },

  // 字体排版
  typography: {
    fontFamily: [
      "-apple-system",
      "BlinkMacSystemFont",
      '"Segoe UI"',
      "Roboto",
      '"Helvetica Neue"',
      "Arial",
      "sans-serif",
      '"Apple Color Emoji"',
      '"Segoe UI Emoji"',
      '"Segoe UI Symbol"',
    ].join(","),
  },
});

export default theme;
