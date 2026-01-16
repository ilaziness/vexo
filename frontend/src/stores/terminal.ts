import { create } from "zustand";
import { TerminalTheme, TerminalThemeMode } from "../types/ssh";
import terminalThemes from "../theme/terminalThemes";

export interface TerminalSettings {
  fontFamily: string;
  fontSize: number;
  lineHeight: number;
  theme: TerminalThemeMode;
  setFontFamily: (f: string) => void;
  setFontSize: (s: number) => void;
  setLineHeight: (l: number) => void;
  setTheme: (theme: TerminalThemeMode) => void;
  getCurrentTheme: () => TerminalTheme;
  syncWithGlobalTheme: (globalTheme: string) => void;
}

export const useTerminalStore = create<TerminalSettings>((set, get) => ({
  fontFamily:
    '"Noto Sans Mono", "Cascadia Code", "Fira Code", Menlo, Monaco, Consolas, "DejaVu Sans Mono", "Ubuntu Mono", "Courier New", monospace',
  fontSize: 14,
  lineHeight: 1,
  theme: "dark", // 默认使用暗色主题
  setFontFamily: (f: string) => set(() => ({ fontFamily: f })),
  setFontSize: (s: number) => set(() => ({ fontSize: s })),
  setLineHeight: (l: number) => set(() => ({ lineHeight: l })),
  setTheme: (theme: TerminalThemeMode) => set(() => ({ theme })),
  getCurrentTheme: (): TerminalTheme => {
    const { theme } = get();
    return terminalThemes[theme];
  },
  // 根据总体主题模式自动切换终端主题
  syncWithGlobalTheme: (globalTheme: string) => {
    let terminalTheme: TerminalThemeMode = "dark"; // 默认

    // 根据全局主题名称映射到终端主题
    switch (globalTheme) {
      case "light":
        terminalTheme = "light";
        break;
      case "blueDark":
        terminalTheme = "blueDark";
        break;
      case "atom":
        terminalTheme = "atom";
        break;
      case "deep":
        terminalTheme = "deep";
        break;
      case "eyeCare":
        terminalTheme = "eyeCare";
        break;
      default:
        terminalTheme = "dark";
        break;
    }

    set(() => ({ theme: terminalTheme }));
  },
}));

export default useTerminalStore;
