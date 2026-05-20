import { create } from 'zustand';

interface AIAssistantState {
  sidebarOpen: boolean;
  toggleSidebarOpen: () => void;
}

export const useAIAssistantStore = create<AIAssistantState>((set) => ({
  sidebarOpen: false,
  toggleSidebarOpen: () => set((state) => ({ sidebarOpen: !state.sidebarOpen })),
}));
