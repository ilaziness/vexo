import { Terminal } from "@xterm/xterm";

class TerminalInstances {
  private static instance: TerminalInstances;
  private instances: Map<string, Terminal>;

  private constructor() {
    this.instances = new Map();
  }

  public static getInstance(): TerminalInstances {
    if (!TerminalInstances.instance) {
      TerminalInstances.instance = new TerminalInstances();
    }
    return TerminalInstances.instance;
  }

  public get(id: string): Terminal | undefined {
    return this.instances.get(id);
  }

  public set(id: string, term: Terminal): void {
    this.instances.set(id, term);
  }

  public remove(id: string): void {
    this.instances.delete(id);
  }
}

export const terminalInstances = TerminalInstances.getInstance();
