import { create } from 'zustand';
import { AIConfig, AIService } from '../../bindings/github.com/ilaziness/vexo/services';
import { parseCallServiceError } from '../func/service';

interface AIConfigState {
  config: AIConfig | null;
  isLoading: boolean;
  error: string | null;
  providers: { name: string; label: string; needs_api_key: boolean; needs_endpoint: boolean }[];
  isLoadingProviders: boolean;

  loadConfig: () => Promise<void>;
  saveConfig: (config: AIConfig) => Promise<boolean>;
  resetConfig: () => Promise<void>;
  loadProviders: () => Promise<void>;
  updatePartialConfig: (partial: Partial<AIConfig>) => void;
  clearError: () => void;
}

const defaultConfig: AIConfig = {
  enabled: false,
  provider: 'ollama',
  model: 'llama3.2',
  endpoint: 'http://localhost:11434',
  temperature: 0.7,
  max_tokens: 2048,
};

export const useAIConfigStore = create<AIConfigState>((set) => ({
  config: null,
  isLoading: false,
  error: null,
  providers: [],
  isLoadingProviders: false,

  loadConfig: async () => {
    set({ isLoading: true, error: null });
    try {
      const config = await AIService.GetConfig();
      set({ config: config || defaultConfig, isLoading: false });
    } catch (err) {
      set({ error: parseCallServiceError(err), isLoading: false });
    }
  },

  loadProviders: async () => {
    set({ isLoadingProviders: true });
    try {
      const providers = await AIService.GetProviders();
      set({ providers: providers || [], isLoadingProviders: false });
    } catch (err) {
      set({ providers: [], isLoadingProviders: false });
    }
  },

  saveConfig: async (config: AIConfig) => {
    set({ isLoading: true, error: null });
    try {
      await AIService.SaveConfig(config);
      set({ config, isLoading: false });
      return true;
    } catch (err) {
      set({ error: parseCallServiceError(err), isLoading: false });
      return false;
    }
  },

  resetConfig: async () => {
    set({ isLoading: true, error: null });
    try {
      await AIService.ResetConfig();
      const config = await AIService.GetConfig();
      set({ config: config || defaultConfig, isLoading: false });
    } catch (err) {
      set({ error: parseCallServiceError(err), isLoading: false });
    }
  },

  updatePartialConfig: (partial: Partial<AIConfig>) => {
    set((state) => ({
      config: state.config ? { ...state.config, ...partial } : { ...defaultConfig, ...partial },
    }));
  },

  clearError: () => set({ error: null }),
}));

