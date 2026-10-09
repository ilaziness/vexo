import { create } from "zustand";

interface UIState {
  bookmarkManageOpen: boolean;
  setBookmarkManageOpen: (open: boolean) => void;
  toggleBookmarkManage: () => void;
}

export const useUIStore = create<UIState>((set) => ({
  bookmarkManageOpen: false,
  setBookmarkManageOpen: (open) => set({ bookmarkManageOpen: open }),
  toggleBookmarkManage: () =>
    set((state) => ({ bookmarkManageOpen: !state.bookmarkManageOpen })),
}));
