import { create } from "zustand";

export interface TerminalSettings {
  fontFamily: string;
  fontSize: number;
  lineHeight: number;
}

export const useTerminalStore = create<TerminalSettings>((set) => ({
  fontFamily:
    '"Cascadia Code", "Fira Code", Menlo, Monaco, Consolas, "DejaVu Sans Mono", "Ubuntu Mono", "Courier New", monospace',
  fontSize: 15,
  lineHeight: 1.2,
  setFontFamily: (f: string) => set(() => ({ fontFamily: f })),
  setFontSize: (s: number) => set(() => ({ fontSize: s })),
  setLineHeight: (l: number) => set(() => ({ lineHeight: l })),
}));

export default useTerminalStore;
