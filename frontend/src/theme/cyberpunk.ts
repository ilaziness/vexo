import type { ColorSystemOptions } from "@mui/material/styles";

/**
 * 赛博朋克主题 - Cyberpunk Theme
 *
 * 设计理念：
 * - 融合 Tech Innovation 的电光蓝与 Midnight Galaxy 的深紫色调
 * - 高对比度霓虹配色，营造未来科技感
 * - 深色背景配合荧光色前景，视觉冲击力与可读性兼具
 * - 适合追求个性和 modern 风格的用户
 *
 * 配色方案参考：
 * - Primary: Electric Blue (#0066ff) 变体 - 科技感主色调
 * - Secondary: Neon Cyan (#00ffff) 与紫色融合 - 霓虹强调色
 * - Background: Dark Gray (#1e1e1e) 偏紫 - 深夜都市背景
 * - Accent: 亮紫粉色调 - 赛博朋克标志性荧光色
 */
const cyberpunk: ColorSystemOptions = {
  palette: {
    mode: "dark",
    // 主色调：电光蓝，科技感强烈
    primary: {
      main: "#00d4ff", // 霓虹青色/蓝色
      light: "#5ce1ff",
      dark: "#0099cc",
      contrastText: "#0a0a0f",
    },
    // 次要色：霓虹紫粉色，与主色形成对比
    secondary: {
      main: "#ff00a0", // 霓虹品红/粉色
      light: "#ff5cc8",
      dark: "#cc007d",
      contrastText: "#ffffff",
    },
    // 背景色：深紫黑色，营造深夜都市感
    background: {
      default: "#0f0f1a", // 深蓝紫黑色主背景
      paper: "#1a1a2e", // 稍浅的卡片背景，带紫色调
    },
    // 文本颜色：高对比度荧光色调
    text: {
      primary: "#e8f4f8", // 荧光蓝白色
      secondary: "#a0a8b8", // 淡蓝灰色
      disabled: "#5a6070", // 暗蓝灰色
    },
    // 分割线：半透明的霓虹色
    divider: "rgba(0, 212, 255, 0.15)",
    // 动作按钮
    action: {
      active: "#00d4ff",
      hover: "rgba(0, 212, 255, 0.12)",
      selected: "rgba(0, 212, 255, 0.2)",
      disabled: "rgba(168, 178, 198, 0.3)",
      disabledBackground: "rgba(168, 178, 198, 0.12)",
    },
    // 语义化颜色：保持霓虹风格但协调
    success: {
      main: "#00ff88", // 霓虹绿色
      light: "#5cFFB5",
      dark: "#00cc6a",
    },
    error: {
      main: "#ff3366", // 霓虹红色
      light: "#ff6699",
      dark: "#cc0044",
    },
    warning: {
      main: "#ffcc00", // 霓虹黄色
      light: "#ffe066",
      dark: "#cc9900",
    },
    info: {
      main: "#00d4ff", // 与主色一致
      light: "#5ce1ff",
      dark: "#0099cc",
    },
  },
};

export default cyberpunk;
