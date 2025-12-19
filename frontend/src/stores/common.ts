import { create } from "zustand";
import { Message, MessageStore } from "../types/ssh";

export const useMessageStore = create<MessageStore>((set) => ({
  message: {} as Message,
  setClose: () =>
    set((state) => ({
      message: { ...state.message, open: false },
    })),
  errorMessage: (message: string) =>
    set((state) => ({
      message: { open: true, text: message, type: "error" },
    })),
  successMessage: (message: string) =>
    set((state) => ({
      message: { open: true, text: message, type: "success" },
    })),
}));
