declare module '@mui/x-chat-headless/types' {
  interface ChatToolDefinitionMap {
    run_ssh_command: {
      input: { command: string };
      output: { output: string };
    };
    upsert_plan: {
      input: {
        title: string;
        tasks: Array<{ id: string; title: string; status: string }>;
      };
      output: {
        title: string;
        tasks: Array<{ id: string; title: string; status: string }>;
      };
    };
  }
}

export {};
