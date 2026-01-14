import { create } from "zustand";
import { ProgressData } from "../../bindings/github.com/ilaziness/vexo/services/models";

export interface TransferStore {
  transfers: Map<string, ProgressData[]>; // sessionID -> ProgressData[]
  addProgress: (progress: ProgressData) => void;
  removeProgress: (sessionID: string, id: string) => void;
  getTransfersBySession: (sessionID: string) => ProgressData[];
  clearSession: (sessionID: string) => void;
}

export const useTransferStore = create<TransferStore>((set, get) => ({
  transfers: new Map(),
  addProgress: (progress: ProgressData) => {
    set((state) => {
      const newTransfers = new Map(state.transfers);
      const sessionID = progress.SessionID;
      if (!newTransfers.has(sessionID)) {
        newTransfers.set(sessionID, []);
      }
      const sessionTransfers = newTransfers.get(sessionID)!;
      // 检查是否已有相同ID的progress，如果有则更新，否则添加
      const existingIndex = sessionTransfers.findIndex((p) => p.ID === progress.ID);
      if (existingIndex >= 0) {
        sessionTransfers[existingIndex] = progress;
      } else {
        sessionTransfers.push(progress);
      }
      newTransfers.set(sessionID, sessionTransfers);
      return { transfers: newTransfers };
    });
  },
  removeProgress: (sessionID: string, id: string) => {
    set((state) => {
      const newTransfers = new Map(state.transfers);
      const sessionTransfers = newTransfers.get(sessionID) || [];
      const filtered = sessionTransfers.filter((p) => p.ID !== id);
      newTransfers.set(sessionID, filtered);
      return { transfers: newTransfers };
    });
  },
  getTransfersBySession: (sessionID: string) => {
    return get().transfers.get(sessionID) || [];
  },
  clearSession: (sessionID: string) => {
    set((state) => {
      const newTransfers = new Map(state.transfers);
      newTransfers.delete(sessionID);
      return { transfers: newTransfers };
    });
  },
}));
