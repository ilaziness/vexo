// Genkit Adapter - 将MUI X Chat连接到后端AIService
import type { ChatAdapter, ChatMessageChunk, ChatStreamEnvelope, ChatUser } from '@mui/x-chat/headless';
import { Events } from '@wailsio/runtime';
import { AIService } from '../../../bindings/github.com/ilaziness/vexo/services';
import { ChatRequest, AIMessage } from '../../../bindings/github.com/ilaziness/vexo/services/models';
import { parseCallServiceError } from '../../func/service';
import { getCurrentSSHContext } from '../../func/aiContext';
import { useMessageStore } from '../../stores/message';
import { useAIAssistantStore } from '../../stores/aiAssistant';

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
type ListMessagesInput = Parameters<NonNullable<ChatAdapter['listMessages']>>[0];
type ListMessagesResult = Awaited<ReturnType<NonNullable<ChatAdapter['listMessages']>>>;

export class GenkitAdapter implements ChatAdapter {
  async sendMessage(input: SendMessageInput): Promise<ReadableStream<ChatMessageChunk | ChatStreamEnvelope>> {
    const { signal } = input;
    const textPart = input.message.parts.find((part) => part.type === 'text');
    const newMessage = (textPart?.text || '').trim();
    if (!newMessage) {
      const message = '消息内容不能为空';
      useMessageStore.getState().errorMessage(message);
      throw new Error(message);
    }

    const sessionId = input.conversationId || '';
    if (!sessionId) {
      const message = '会话未就绪，请稍后重试';
      useMessageStore.getState().errorMessage(message);
      throw new Error(message);
    }

    let chatPromise: ReturnType<typeof AIService.Chat>;
    let messageId: string;
    const capturedSessionId = sessionId;

    try {
      const sshContext = getCurrentSSHContext();
      const request = new ChatRequest({
        session_id: sessionId,
        new_message: newMessage,
        ...(sshContext ? { ssh_context: sshContext } : {}),
      });

      messageId = `msg-${Date.now()}-${Math.random().toString(36).slice(2, 9)}`;
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

        const endStreaming = () => {
          useAIAssistantStore.getState().setStreaming(false);
        };

        const ensureStarted = () => {
          if (!hasStarted) {
            hasStarted = true;
            useAIAssistantStore.getState().setStreaming(true);
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
          endStreaming();
          controller.close();
        };
        signal.addEventListener('abort', onAbort);

        const unsubscribe = Events.On(EVENT_AI_STREAM_CHUNK, (event: any) => {
          if (aborted) return;
          if (useAIAssistantStore.getState().activeSessionId !== capturedSessionId) return;
          const data = event.data;
          if (data?.sessionId !== capturedSessionId) return;

          ensureStarted();
          const chunk: string = data.chunk || '';

          if (data.type === 'reasoning') {
            if (!hasReasoningStarted) {
              hasReasoningStarted = true;
              controller.enqueue({ type: 'reasoning-start', id: reasoningId } as ChatMessageChunk);
            }
            controller.enqueue({ type: 'reasoning-delta', id: reasoningId, delta: chunk } as ChatMessageChunk);
          } else {
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
            if (useAIAssistantStore.getState().activeSessionId !== capturedSessionId) {
              endStreaming();
              controller.close();
              return;
            }
            ensureStarted();
            if (hasReasoningStarted) {
              controller.enqueue({ type: 'reasoning-end', id: reasoningId } as ChatMessageChunk);
            }
            if (!hasTextStarted) {
              controller.enqueue({ type: 'text-start', id: textId } as ChatMessageChunk);
            }
            controller.enqueue({ type: 'text-end', id: textId } as ChatMessageChunk);
            controller.enqueue({ type: 'finish', messageId } as ChatMessageChunk);
            endStreaming();
            controller.close();
          })
          .catch((err: any) => {
            if (aborted) return;
            unsubscribe();
            signal.removeEventListener('abort', onAbort);
            if (useAIAssistantStore.getState().activeSessionId !== capturedSessionId) {
              endStreaming();
              controller.close();
              return;
            }
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
            endStreaming();
            controller.close();
          });
      },
    });
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
