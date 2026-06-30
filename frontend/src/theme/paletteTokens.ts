/**
 * 三套主题的锚点色 — MUI palette 与 XTerm 终端共用，避免双源漂移。
 */
export const lightTokens = {
  bgDefault: "#F1F5F9",
  bgPaper: "#FFFFFF",
  textPrimary: "#0F172A",
  textSecondary: "#475569",
  primary: "#2563EB",
  primaryLight: "#60A5FA",
  primaryDark: "#1D4ED8",
  secondary: "#10B981",
  secondaryLight: "#34D399",
  secondaryDark: "#059669",
  aiBubble: "#FFFFFF",
  success: "#059669",
  divider: "rgba(15, 23, 42, 0.10)",
  terminalCyan: "#0891B2",
} as const;

export const darkTokens = {
  bgDefault: "#0D1117",
  bgPaper: "#21262D",
  textPrimary: "#E6EDF3",
  textSecondary: "#8B949E",
  primary: "#58A6FF",
  primaryLight: "#79C0FF",
  primaryDark: "#388BFD",
  secondary: "#22C55E",
  secondaryLight: "#4ADE80",
  secondaryDark: "#16A34A",
  aiBubble: "#30363D",
  success: "#3FB950",
  divider: "rgba(240, 246, 252, 0.10)",
  terminalCyan: "#0891B2",
  ansiGreen: "#3FB950",
  ansiBrightGreen: "#56D364",
} as const;

export const eyeCareTokens = {
  bgDefault: "#FFFBEB",
  bgPaper: "#FEF3C7",
  textPrimary: "#451A03",
  textSecondary: "#92400E",
  primary: "#B45309",
  primaryLight: "#D97706",
  primaryDark: "#78350F",
  secondary: "#0D9488",
  secondaryLight: "#14B8A6",
  secondaryDark: "#0F766E",
  aiBubble: "#FFFDF5",
  success: "#65A30D",
  divider: "rgba(69, 26, 3, 0.12)",
  assistantBubbleBorder: "rgba(69, 26, 3, 0.12)",
  terminalCyan: "#0D9488",
} as const;
