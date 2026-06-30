import { create } from "zustand";
import type { ITheme } from "@xterm/xterm";
import { AppTheme } from "../types/ssh";
import terminalThemes from "../theme/terminalThemes";
import { CODE_FONT_FAMILY } from "../theme/typography";

const VALID_THEMES: AppTheme[] = ["light", "dark", "eyeCare"];

function normalizeTheme(theme: AppTheme): AppTheme {
  return VALID_THEMES.includes(theme) ? theme : "dark";
}

export interface TerminalSettings {
  fontFamily: string;
  fontSize: number;
  lineHeight: number;
  theme: AppTheme;
  setFontFamily: (f: string) => void;
  setFontSize: (s: number) => void;
  setLineHeight: (l: number) => void;
  setTheme: (theme: AppTheme) => void;
  getCurrentTheme: () => ITheme;
  syncWithGlobalTheme: (globalTheme: AppTheme) => void;
}

export const useTerminalStore = create<TerminalSettings>((set, get) => ({
  fontFamily: CODE_FONT_FAMILY,
  fontSize: 14,
  lineHeight: 1,
  theme: "dark", // 默认使用暗色主题
  setFontFamily: (f: string) => set(() => ({ fontFamily: f })),
  setFontSize: (s: number) => set(() => ({ fontSize: s })),
  setLineHeight: (l: number) => set(() => ({ lineHeight: l })),
  setTheme: (theme: AppTheme) => set(() => ({ theme: normalizeTheme(theme) })),
  getCurrentTheme: (): ITheme => {
    const { theme } = get();
    const validTheme = normalizeTheme(theme);
    return terminalThemes[validTheme];
  },
  // 根据总体主题模式自动切换终端主题
  syncWithGlobalTheme: (globalTheme: AppTheme) => {
    set(() => ({ theme: normalizeTheme(globalTheme) }));
  },
}));

export default useTerminalStore;
