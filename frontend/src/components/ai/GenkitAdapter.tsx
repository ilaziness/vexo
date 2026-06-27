// Genkit Adapter - 将MUI X Chat连接到后端AIService
import type { ChatAdapter, ChatMessageChunk, ChatStreamEnvelope, ChatUser } from '@mui/x-chat/headless';
import { Events } from '@wailsio/runtime';
import { AIService } from '../../../bindings/github.com/ilaziness/vexo/services';
import { ChatRequest, AIMessage, AISession } from '../../../bindings/github.com/ilaziness/vexo/services/models';
import { parseCallServiceError } from '../../func/service';
import { useMessageStore } from '../../stores/message';

const EVENT_AI_STREAM_CHUNK = 'eventAIStreamChunk';

export const currentUser: ChatUser = {
  id: 'user',
  displayName: 'You',
};

export const aiUser: ChatUser = {
  id: 'assistant',
  displayName: 'AI',
};

type SendMessageInput = Parameters<ChatAdapter['sendMessage']>[0];
type ListConversationsResult = Awaited<ReturnType<NonNullable<ChatAdapter['listConversations']>>>;
type ListMessagesInput = Parameters<NonNullable<ChatAdapter['listMessages']>>[0];
type ListMessagesResult = Awaited<ReturnType<NonNullable<ChatAdapter['listMessages']>>>;

export class GenkitAdapter implements ChatAdapter {
  lastCreatedSessionId: string | null = null;

  async sendMessage(input: SendMessageInput): Promise<ReadableStream<ChatMessageChunk | ChatStreamEnvelope>> {
    const { signal } = input;
    const textPart = input.message.parts.find((part) => part.type === 'text');
    const newMessage = (textPart?.text || '').trim();
    if (!newMessage) {
      const message = '消息内容不能为空';
      useMessageStore.getState().errorMessage(message);
      throw new Error(message);
    }

    let sessionId: string;
    let chatPromise: ReturnType<typeof AIService.Chat>;
    let messageId: string;
    let capturedSessionId: string;

    try {
      // 若没有 conversationId，先创建会话
      sessionId = input.conversationId || '';
      if (!sessionId) {
        const session = await AIService.CreateSession();
        if (!session) throw new Error('create session failed');
        sessionId = session.id;
        this.lastCreatedSessionId = sessionId;
      }

      const request = new ChatRequest({
        session_id: sessionId,
        new_message: newMessage,
      });

      messageId = `msg-${Date.now()}-${Math.random().toString(36).slice(2, 9)}`;
      capturedSessionId = sessionId;
      chatPromise = AIService.Chat(request);
    } catch (err) {
      const message = parseCallServiceError(err);
      useMessageStore.getState().errorMessage(message);
      throw new Error(message, { cause: err });
    }

    return new ReadableStream({
      start(controller) {
        const reasoningId = `reasoning-${messageId}`;
        const textId = `text-${messageId}`;
        let hasReasoningStarted = false;
        let hasTextStarted = false;
        let hasStarted = false;
        let aborted = false;

        const ensureStarted = () => {
          if (!hasStarted) {
            hasStarted = true;
            controller.enqueue({ type: 'start', messageId, author: aiUser } as ChatMessageChunk);
          }
        };

        const onAbort = () => {
          aborted = true;
          unsubscribe();
          ensureStarted();
          if (hasReasoningStarted) {
            controller.enqueue({ type: 'reasoning-end', id: reasoningId } as ChatMessageChunk);
          }
          if (hasTextStarted) {
            controller.enqueue({ type: 'text-end', id: textId } as ChatMessageChunk);
          }
          controller.enqueue({ type: 'abort', messageId } as ChatMessageChunk);
          controller.close();
        };
        signal.addEventListener('abort', onAbort);

        const unsubscribe = Events.On(EVENT_AI_STREAM_CHUNK, (event: any) => {
          if (aborted) return;
          const data = event.data;
          if (data?.sessionId !== capturedSessionId) return;

          ensureStarted();
          const chunk: string = data.chunk || '';

          // 思考过程 chunk（以 <think> 开头或标记为 reasoning）
          if (data.type === 'reasoning') {
            if (!hasReasoningStarted) {
              hasReasoningStarted = true;
              controller.enqueue({ type: 'reasoning-start', id: reasoningId } as ChatMessageChunk);
            }
            controller.enqueue({ type: 'reasoning-delta', id: reasoningId, delta: chunk } as ChatMessageChunk);
          } else {
            // 普通文本 chunk
            if (!hasTextStarted) {
              hasTextStarted = true;
              controller.enqueue({ type: 'text-start', id: textId } as ChatMessageChunk);
            }
            controller.enqueue({ type: 'text-delta', id: textId, delta: chunk } as ChatMessageChunk);
          }
        });

        chatPromise
          .then(() => {
            if (aborted) return;
            unsubscribe();
            signal.removeEventListener('abort', onAbort);
            ensureStarted();
            if (hasReasoningStarted) {
              controller.enqueue({ type: 'reasoning-end', id: reasoningId } as ChatMessageChunk);
            }
            if (!hasTextStarted) {
              // 至少要有一个 text part
              controller.enqueue({ type: 'text-start', id: textId } as ChatMessageChunk);
            }
            controller.enqueue({ type: 'text-end', id: textId } as ChatMessageChunk);
            controller.enqueue({ type: 'finish', messageId } as ChatMessageChunk);
            controller.close();
          })
          .catch((err: any) => {
            if (aborted) return;
            unsubscribe();
            signal.removeEventListener('abort', onAbort);
            ensureStarted();
            const message = parseCallServiceError(err);
            useMessageStore.getState().errorMessage(message);
            if (hasReasoningStarted) {
              controller.enqueue({ type: 'reasoning-end', id: reasoningId } as ChatMessageChunk);
            }
            if (!hasTextStarted) {
              controller.enqueue({ type: 'text-start', id: textId } as ChatMessageChunk);
              controller.enqueue({ type: 'text-delta', id: textId, delta: `Error: ${message}` } as ChatMessageChunk);
            }
            controller.enqueue({ type: 'text-end', id: textId } as ChatMessageChunk);
            controller.enqueue({ type: 'finish', messageId } as ChatMessageChunk);
            controller.close();
          });
      },
    });
  }

  async listConversations(): Promise<ListConversationsResult> {
    try {
      const sessions = await AIService.ListSessions(50);
      const conversations: ListConversationsResult['conversations'] = (sessions || [])
        .filter((s): s is AISession => s !== null)
        .map((s) => ({
          id: s.id,
          title: s.title || '新会话',
          createdAt: new Date(s.created_at * 1000).toISOString(),
          lastMessageAt: new Date(s.updated_at * 1000).toISOString(),
          participants: [currentUser, aiUser],
        }));
      return { conversations };
    } catch (err) {
      useMessageStore.getState().errorMessage(parseCallServiceError(err));
      return { conversations: [] };
    }
  }

  async listMessages(input: ListMessagesInput): Promise<ListMessagesResult> {
    try {
      const msgs = await AIService.ListMessages(input.conversationId);
      const messages: ListMessagesResult['messages'] = (msgs || [])
        .filter((m): m is AIMessage => m !== null)
        .map((m) => {
          let parts: ListMessagesResult['messages'][number]['parts'];
          try {
            const parsed: unknown = m.parts ? JSON.parse(m.parts) : null;
            parts = Array.isArray(parsed)
              ? (parsed as ListMessagesResult['messages'][number]['parts'])
              : [{ type: 'text', text: m.content }];
          } catch {
            parts = [{ type: 'text', text: m.content }];
          }
          return {
            id: m.id,
            role: m.role === 'user' ? 'user' : 'assistant',
            status: 'sent',
            parts,
            createdAt: new Date(m.timestamp * 1000).toISOString(),
            conversationId: m.session_id,
            author: m.role === 'user' ? currentUser : aiUser,
          };
        });
      return { messages };
    } catch (err) {
      useMessageStore.getState().errorMessage(parseCallServiceError(err));
      return { messages: [] };
    }
  }

  stop(): void {
  }
}
