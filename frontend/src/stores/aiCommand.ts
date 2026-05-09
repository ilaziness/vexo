import { create } from "zustand";
import {
  AIService,
  CommandService,
} from "../../bindings/github.com/ilaziness/vexo/services";
import {
  CommandGenerateRequest,
  CommandGenerateResponse,
  CommandExplainRequest,
} from "../../bindings/github.com/ilaziness/vexo/internal/ai/models";
import { parseCallServiceError } from "../func/service";
import { useAIConfigStore } from "./aiConfig";

// 生成唯一 ID
const generateId = () => `${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;

// AI 消息类型
export interface AIMessage {
  id: string;
  role: "user" | "assistant";
  content: string;
  timestamp: number;
  customData?: {
    command?: CommandGenerateResponse;
  };
}

interface AICommandState {
  // 面板状态
  isOpen: boolean;
  isLoading: boolean;
  error: string | null;

  // 配置弹框状态
  showConfigDialog: boolean;

  // 输入状态
  input: string;
  inputHistory: string[];

  // 对话历史
  messages: AIMessage[];

  // 当前结果
  currentResult: unknown;

  // 操作方法
  openPanel: () => Promise<boolean>;
  closePanel: () => void;
  togglePanel: () => Promise<boolean>;
  setInput: (input: string) => void;
  generateCommand: (sessionID?: string) => Promise<boolean>;
  explainCommand: (command: string, sessionID?: string) => Promise<boolean>;
  sendToTerminal: (command: string, sessionIDs?: string[]) => Promise<void>;
  copyToClipboard: (text: string) => Promise<void>;
  addMessage: (message: AIMessage) => void;
  clearHistory: () => void;
  resetError: () => void;

  // 配置弹框控制
  openConfigDialog: () => void;
  closeConfigDialog: () => void;
}

// 检查 AI 是否启用的辅助函数
const checkAIEnabled = async (): Promise<boolean> => {
  const state = useAIConfigStore.getState();
  // 如果配置为 null，先加载配置
  if (state.config === null && !state.isLoading) {
    await state.loadConfig();
  }
  const config = useAIConfigStore.getState().config;
  return config?.enabled || false;
};

export const useAICommandStore = create<AICommandState>((set, get) => ({
  isOpen: false,
  isLoading: false,
  error: null,
  showConfigDialog: false,
  input: "",
  inputHistory: [],
  messages: [],
  currentResult: null,

  openPanel: async () => {
    if (!(await checkAIEnabled())) {
      set({ showConfigDialog: true });
      return false;
    }
    set({ isOpen: true });
    return true;
  },

  closePanel: () => set({ isOpen: false, currentResult: null }),

  togglePanel: async () => {
    if (!(await checkAIEnabled())) {
      set({ showConfigDialog: true });
      return false;
    }
    const newState = !get().isOpen;
    set({ isOpen: newState, currentResult: null });
    return newState;
  },

  openConfigDialog: () => set({ showConfigDialog: true }),
  closeConfigDialog: () => set({ showConfigDialog: false }),

  setInput: (input: string) => set({ input }),

  generateCommand: async (sessionID?: string) => {
    if (!(await checkAIEnabled())) {
      set({ showConfigDialog: true });
      return false;
    }

    const { input, inputHistory } = get();
    if (!input.trim()) return false;

    set({ isLoading: true, error: null, currentResult: null });

    // 添加用户消息
    const userMessage: AIMessage = {
      id: generateId(),
      role: "user",
      content: input,
      timestamp: Date.now(),
    };
    set((state) => ({ messages: [...state.messages, userMessage] }));

    try {
      const request: CommandGenerateRequest = {
        input: input,
        session_id: sessionID || "",
      };

      const result = await AIService.GenerateCommand(request);
      if (!result) {
        throw new Error("生成结果为空");
      }

      // 添加 AI 回复
      const assistantMessage: AIMessage = {
        id: generateId(),
        role: "assistant",
        content: result.explanation,
        timestamp: Date.now(),
        customData: {
          command: result,
        },
      };

      set((state) => ({
        isLoading: false,
        currentResult: result,
        messages: [...state.messages, assistantMessage],
        inputHistory: [...inputHistory, input],
        input: "",
      }));
      return true;
    } catch (err) {
      const errorMsg = parseCallServiceError(err);
      set({
        isLoading: false,
        error: errorMsg,
      });

      // 添加错误消息
      const errorMessage: AIMessage = {
        id: generateId(),
        role: "assistant",
        content: `生成命令失败: ${errorMsg}`,
        timestamp: Date.now(),
      };
      set((state) => ({
        messages: [...state.messages, errorMessage],
      }));
      return false;
    }
  },

  explainCommand: async (command: string, sessionID?: string) => {
    if (!(await checkAIEnabled())) {
      set({ showConfigDialog: true });
      return false;
    }

    set({ isLoading: true, error: null });

    try {
      const request: CommandExplainRequest = {
        command: command,
        session_id: sessionID || "",
      };

      const result = await AIService.ExplainCommand(request);
      if (!result) {
        throw new Error("解释结果为空");
      }

      // 构造解释内容
      let explanation = result.explanation + "\n\n";
      if (result.parts && result.parts.length > 0) {
        explanation += "命令解析:\n";
        result.parts.forEach((part) => {
          explanation += `- ${part.part}: ${part.meaning}\n`;
        });
      }

      // 添加 AI 回复
      const assistantMessage: AIMessage = {
        id: generateId(),
        role: "assistant",
        content: explanation,
        timestamp: Date.now(),
      };

      set((state) => ({
        isLoading: false,
        messages: [...state.messages, assistantMessage],
        isOpen: true,
      }));
      return true;
    } catch (err) {
      const errorMsg = parseCallServiceError(err);
      set({
        isLoading: false,
        error: errorMsg,
      });
      return false;
    }
  },

  sendToTerminal: async (command: string, sessionIDs?: string[]) => {
    try {
      await CommandService.SendCommand({
        command: command + "\n",
        session_ids: sessionIDs || [],
      });
    } catch (err) {
      set({ error: parseCallServiceError(err) });
    }
  },

  copyToClipboard: async (text: string) => {
    try {
      await navigator.clipboard.writeText(text);
    } catch (err) {
      set({ error: "复制到剪贴板失败" });
    }
  },

  addMessage: (message: AIMessage) =>
    set((state) => ({ messages: [...state.messages, message] })),

  clearHistory: () => set({ messages: [], inputHistory: [], currentResult: null }),

  resetError: () => set({ error: null }),
}));
