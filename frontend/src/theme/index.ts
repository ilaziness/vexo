import { createTheme } from "@mui/material/styles";

/**
 * SSH GUI 客户端主题配置
 * 设计理念：清爽简洁，适合终端应用场景
 */
const theme = createTheme({
  // 配置颜色方案
  colorSchemes: {
    // 亮色模式主题
    light: {
      palette: {
        // 主色调：清爽的蓝色，代表技术感和专业性
        primary: {
          main: "#1976d2", // 蓝色主色
          light: "#42a5f5",
          dark: "#1565c0",
          contrastText: "#ffffff",
        },
        // 次要色：柔和的青色
        secondary: {
          main: "#0288d1",
          light: "#03a9f4",
          dark: "#01579b",
          contrastText: "#ffffff",
        },
        // 背景色：简洁的浅色系
        background: {
          default: "#f0f0f0", // 主背景 - 更柔和的灰色，减少眩光
          paper: "#fefefe", // 卡片、面板背景 - 略带灰调的白色
        },
        // 文本颜色
        text: {
          primary: "rgba(0, 0, 0, 0.87)",
          secondary: "rgba(0, 0, 0, 0.6)",
          disabled: "rgba(0, 0, 0, 0.38)",
        },
        // 分割线
        divider: "rgba(0, 0, 0, 0.12)",
        // 动作按钮
        action: {
          active: "rgba(0, 0, 0, 0.54)",
          hover: "rgba(0, 0, 0, 0.04)",
          selected: "rgba(0, 0, 0, 0.08)",
          disabled: "rgba(0, 0, 0, 0.26)",
          disabledBackground: "rgba(0, 0, 0, 0.12)",
        },
        // 成功、错误、警告、信息
        success: {
          main: "#2e7d32",
          light: "#4caf50",
          dark: "#1b5e20",
        },
        error: {
          main: "#d32f2f",
          light: "#ef5350",
          dark: "#c62828",
        },
        warning: {
          main: "#ed6c02",
          light: "#ff9800",
          dark: "#e65100",
        },
        info: {
          main: "#0288d1",
          light: "#03a9f4",
          dark: "#01579b",
        },
      },
    },
    // 暗黑模式主题
    dark: {
      palette: {
        // 主色调：柔和的蓝色，护眼且专业
        primary: {
          main: "#90caf9", // 偏浅的蓝色，在暗色背景下更清晰
          light: "#bbdefb",
          dark: "#42a5f5",
          contrastText: "rgba(0, 0, 0, 0.87)",
        },
        // 次要色：青色系
        secondary: {
          main: "#26c6da",
          light: "#4dd0e1",
          dark: "#0097a7",
          contrastText: "rgba(0, 0, 0, 0.87)",
        },
        // 背景色：深色系，避免纯黑，更护眼
        background: {
          default: "#0a0e27", // 深蓝黑色主背景
          paper: "#1a1f3a", // 稍浅的卡片背景
        },
        // 文本颜色：高对比度，易读
        text: {
          primary: "#e3f2fd", // 浅蓝白色，更柔和
          secondary: "rgba(227, 242, 253, 0.7)",
          disabled: "rgba(227, 242, 253, 0.5)",
        },
        // 分割线
        divider: "rgba(227, 242, 253, 0.12)",
        // 动作按钮
        action: {
          active: "#90caf9",
          hover: "rgba(144, 202, 249, 0.08)",
          selected: "rgba(144, 202, 249, 0.16)",
          disabled: "rgba(227, 242, 253, 0.3)",
          disabledBackground: "rgba(227, 242, 253, 0.12)",
        },
        // 成功、错误、警告、信息
        success: {
          main: "#66bb6a",
          light: "#81c784",
          dark: "#388e3c",
        },
        error: {
          main: "#f44336",
          light: "#e57373",
          dark: "#d32f2f",
        },
        warning: {
          main: "#ffa726",
          light: "#ffb74d",
          dark: "#f57c00",
        },
        info: {
          main: "#29b6f6",
          light: "#4fc3f7",
          dark: "#0288d1",
        },
      },
    },
  },

  // 组件默认样式配置
  components: {
    // 按钮样式
    MuiButton: {
      styleOverrides: {
        root: {
          borderRadius: 8,
          textTransform: "none", // 取消默认的全大写
          fontWeight: 500,
        },
      },
    },
    // 图标按钮样式
    MuiIconButton: {
      styleOverrides: {
        root: {
          borderRadius: 8,
        },
      },
    },
    // 输入框样式
    MuiTextField: {
      defaultProps: {
        variant: "outlined",
      },
      styleOverrides: {
        root: {
          "& .MuiOutlinedInput-root": {
            borderRadius: 8,
          },
        },
      },
    },
    // 卡片样式
    MuiCard: {
      styleOverrides: {
        root: {
          borderRadius: 12,
        },
      },
    },
    // 纸张样式
    MuiPaper: {
      styleOverrides: {
        root: {
          borderRadius: 8,
        },
        rounded: {
          borderRadius: 8,
        },
      },
    },
    // 对话框样式
    MuiDialog: {
      styleOverrides: {
        paper: {
          borderRadius: 12,
        },
      },
    },
    // 标签页样式
    MuiTab: {
      styleOverrides: {
        root: {
          textTransform: "none",
          fontWeight: 500,
        },
      },
    },
    // 工具提示
    MuiTooltip: {
      styleOverrides: {
        tooltip: {
          borderRadius: 6,
          fontSize: "0.75rem",
        },
      },
    },
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

  // 形状配置
  shape: {
    borderRadius: 8,
  },

  // 间距配置
  spacing: 8,
});

export default theme;
