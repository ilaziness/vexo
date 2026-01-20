import { create } from "zustand";
import { ITheme, TerminalThemeMode, AppTheme } from "../types/ssh";
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
  getCurrentTheme: () => ITheme;
  syncWithGlobalTheme: (globalTheme: AppTheme) => void;
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
  getCurrentTheme: (): ITheme => {
    const { theme } = get();
    return terminalThemes[theme];
  },
  // 根据总体主题模式自动切换终端主题
  syncWithGlobalTheme: (globalTheme: AppTheme) => {
    set(() => ({ theme: globalTheme }));
  },
}));

export default useTerminalStore;
