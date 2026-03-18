export interface BuiltinCommand {
    name: string;
    command: string;
    description: string;
}

export interface UserCommand {
    category: string;
    name: string;
    command: string;
    description: string;
    created_at: number;
}

export interface CommandInfo {
    category: string;
    name: string;
    command: string;
    description: string;
    is_custom: boolean;
}

export interface CommandHistory {
    timestamp: number;
    command: string;
}

export interface SSHTunnelSession {
    id: string;
    clientKey: string;
}

export interface SendCommandRequest {
    command: string;
    session_ids: string[];
}
