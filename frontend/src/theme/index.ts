import { createColorScheme, createTheme } from "@mui/material/styles";
import light from "./light";
import blueDark from "./blueDark";
import atom from "./atom";
import deep from "./deep";
import eyeCare from "./eyeCare";
import cyberpunk from "./cyberpunk";
import { AppTheme } from "../types/ssh";

declare module "@mui/material/styles" {
  interface ColorSchemeOverrides {
    blueDark: true;
    atom: true;
    deep: true;
    eyeCare: true;
    cyberpunk: true;
  }
}

/**
 * 主题选项配置
 */
export const ThemeOptions: {
  value: AppTheme;
  label: string;
  mode: "light" | "dark";
}[] = [
  { value: "light", label: "亮白", mode: "light" },
  { value: "dark", label: "暗色", mode: "dark" },
  { value: "blueDark", label: "蓝夜", mode: "dark" },
  { value: "atom", label: "极黑", mode: "dark" },
  { value: "deep", label: "深邃", mode: "dark" },
  { value: "eyeCare", label: "暖棕", mode: "light" },
  { value: "cyberpunk", label: "赛博朋克", mode: "dark" },
];

/**
 * SSH GUI 客户端主题配置
 * 设计理念：清爽简洁，适合终端应用场景
 */
const theme = createTheme({
  // 配置颜色方案
  colorSchemes: {
    // 亮白主题（自定义偏灰色调）
    light,
    // 暗色（MUI内置）
    dark: true,
    // 蓝夜主题
    blueDark: createColorScheme(blueDark),
    // 极黑
    atom: createColorScheme(atom),
    // 深邃
    deep: createColorScheme(deep),
    // 暖棕
    eyeCare: createColorScheme(eyeCare),
    // 赛博朋克
    cyberpunk: createColorScheme(cyberpunk),
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
