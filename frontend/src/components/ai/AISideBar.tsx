import React, { useState, useEffect, useMemo } from "react";
import { Slide, Box, Avatar, Button } from "@mui/material";
import AddIcon from "@mui/icons-material/Add";
import { ChatBox, ChatMessageAvatar, ChatConversationList } from "@mui/x-chat";
import type { ChatConversation, ChatMessage } from "@mui/x-chat/headless";
import { useAIAssistantStore } from "../../stores/aiAssistant";
import { LogService, AIService } from "../../../bindings/github.com/ilaziness/vexo/services/index";
import { GenkitAdapter, aiUser, currentUser } from "./GenkitAdapter";
import { useMessageStore } from "../../stores/message";
import { parseCallServiceError } from "../../func/service";

// 字母头像组件：ownerState 通过 slotProps 自动注入
const LetterAvatar = React.forwardRef<HTMLDivElement, any>(
  function LetterAvatar(props, ref) {
    const { ownerState, ...other } = props;
    const isUser = ownerState?.role === "user";
    return (
      <Avatar
        ref={ref}
        {...other}
        sx={{
          width: 32,
          height: 32,
          fontSize: isUser ? 14 : 11,
          fontWeight: 600,
          bgcolor: isUser ? "primary.main" : "secondary.main",
        }}
      >
        {isUser ? "Y" : "AI"}
      </Avatar>
    );
  }
);

// 包装 ChatMessageAvatar，只替换 avatar slot
const CustomMessageAvatar = React.forwardRef<HTMLDivElement, any>(
  function CustomMessageAvatar(props, ref) {
    return (
      <ChatMessageAvatar
        ref={ref}
        {...props}
        slots={{ ...props.slots, avatar: LetterAvatar }}
      />
    );
  }
);

// 自定义会话列表组件，顶部添加新建会话按钮
const CustomConversationList = React.forwardRef<HTMLDivElement, any>(
  function CustomConversationList(props, ref) {
    const { onCreateConversation, ...rest } = props;
    return (
      <Box sx={{ display: "flex", flexDirection: "column", height: "100%" }}>
        <Box sx={{ p: 1.5, borderBottom: 1, borderColor: "divider" }}>
          <Button
            fullWidth
            variant="outlined"
            size="small"
            startIcon={<AddIcon />}
            onClick={onCreateConversation}
          >
            新建会话
          </Button>
        </Box>
        <Box sx={{ flex: 1, overflow: "auto" }}>
          <ChatConversationList ref={ref} {...rest} />
        </Box>
      </Box>
    );
  }
);

const AISideBar = () => {
  const sidebarOpen = useAIAssistantStore((state) => state.sidebarOpen);

  const adapter = useMemo(() => new GenkitAdapter(), []);

  const [activeConversationId, setActiveConversationId] = useState<string>("");
  const [conversations, setConversations] = useState<ChatConversation[]>([]);
  const [threads, setThreads] = useState<Record<string, ChatMessage[]>>({});

  // 新建会话
  const handleCreateConversation = async () => {
    try {
      const session = await AIService.CreateSession();
      if (!session) {
        useMessageStore.getState().errorMessage("创建会话失败");
        return;
      }
      // 构建新会话对象
      const newConv: ChatConversation = {
        id: session.id,
        title: "新会话",
        lastMessageAt: new Date().toISOString(),
        participants: [currentUser, aiUser],
      };
      // 添加到会话列表并激活
      setConversations((prev) => [newConv, ...prev]);
      setActiveConversationId(session.id);
      adapter.lastCreatedSessionId = null;
    } catch (err) {
      useMessageStore.getState().errorMessage(parseCallServiceError(err));
    }
  };

  const messages = threads[activeConversationId ?? ""] ?? [];

  const [width, setWidth] = useState(550);
  const [isResizing, setIsResizing] = useState(false);
  const [isHovering, setIsHovering] = useState(false);

  // 初始化时加载历史会话列表
  useEffect(() => {
    adapter.listConversations().then(({ conversations: convList }) => {
      setConversations(convList as ChatConversation[]);
      if (convList.length > 0) {
        setActiveConversationId(convList[0].id);
      }
    }).catch(() => {});
  }, [adapter]);

  // 切换会话时加载该会话的历史消息
  useEffect(() => {
    if (!activeConversationId) return;
    if (threads[activeConversationId]) return; // 已加载过，不重复请求
    adapter.listMessages({ conversationId: activeConversationId }).then(({ messages: msgList }) => {
      setThreads((prev) => ({ ...prev, [activeConversationId]: msgList as ChatMessage[] }));
    }).catch(() => {});
  }, [activeConversationId, adapter]);

  const handleMouseDown = (e: React.MouseEvent) => {
    e.preventDefault();
    setIsResizing(true);
  };

  useEffect(() => {
    const handleMouseMove = (e: MouseEvent) => {
      if (!isResizing) return;
      const newWidth = window.innerWidth - e.clientX;
      if (newWidth >= 200) {
        setWidth(newWidth);
      }
    };

    const handleMouseUp = () => {
      setIsResizing(false);
    };

    if (isResizing) {
      window.addEventListener("mousemove", handleMouseMove);
      window.addEventListener("mouseup", handleMouseUp);
    }

    return () => {
      window.removeEventListener("mousemove", handleMouseMove);
      window.removeEventListener("mouseup", handleMouseUp);
    };
  }, [isResizing]);

  const handleActiveConversationChange = (nextId: string | undefined) => {
    if (!nextId) return;
    setActiveConversationId(nextId);
  };

  const handleMessagesChange = (nextMessages: ChatMessage[]) => {
    // 首次对话时 activeConversationId 尚未更新，从 adapter 取新建的 sessionId
    const convId = activeConversationId || adapter.lastCreatedSessionId;
    if (!convId) return;
    setThreads((prev) => ({ ...prev, [convId]: nextMessages }));
  };

  const handleFinish = () => {
    // 首次对话后将 adapter 记录的新 sessionId 同步为 activeConversationId，并新建 conversation 条目
    const newId = adapter.lastCreatedSessionId;
    if (newId) {
      adapter.lastCreatedSessionId = null;
      // 从后端重新获取该会话的最新标题（后端会在首次消息时自动生成）
      adapter.listConversations().then(({ conversations: convList }) => {
        const newConv = (convList as ChatConversation[]).find((c) => c.id === newId);
        setConversations((prev) => {
          const exists = prev.find((c) => c.id === newId);
          if (exists) return prev;
          return [newConv ?? { id: newId, title: "新会话" }, ...prev];
        });
        if (!activeConversationId) {
          setActiveConversationId(newId);
        }
      }).catch(() => {
        setConversations((prev) => {
          if (prev.find((c) => c.id === newId)) return prev;
          return [{ id: newId, title: "新会话" }, ...prev];
        });
        if (!activeConversationId) {
          setActiveConversationId(newId);
        }
      });
    }
  };

  return (
    <Slide in={sidebarOpen} direction="left" mountOnEnter unmountOnExit>
      <Box
        sx={{
          width: width,
          minWidth: 200,
          height: "calc(100% - 25px)",
          position: "absolute",
          top: 0,
          right: 0,
          zIndex: 99999,
          display: "flex",
        }}
      >
        <Box
          onMouseDown={handleMouseDown}
          onMouseEnter={() => setIsHovering(true)}
          onMouseLeave={() => setIsHovering(false)}
          sx={{
            width: 4,
            cursor: "col-resize",
            backgroundColor: isHovering ? "primary.main" : "transparent",
            transition: isResizing ? "none" : "background-color 0.2s",
            "&:hover": {
              backgroundColor: "primary.main",
            },
          }}
        />
        <ChatBox
          adapter={adapter}
          currentUser={currentUser}
          members={[currentUser, aiUser]}
          activeConversationId={activeConversationId}
          onActiveConversationChange={handleActiveConversationChange}
          conversations={conversations}
          onConversationsChange={setConversations}
          messages={messages}
          onMessagesChange={handleMessagesChange}
          onFinish={handleFinish}
          onError={(err: any) =>
            LogService.Error(`[AI Chat] ${err.source}: ${err.message}`)
          }
          features={{ attachments: false }}
          slots={{
            messageAvatar: CustomMessageAvatar,
            conversationList: (props: any) => (
              <CustomConversationList {...props} onCreateConversation={handleCreateConversation} />
            ),
          }}
          slotProps={{
            composerInput: { placeholder: "输入消息..." },
          }}
          sx={{
            width: "100%",
            height: "100%",
            borderLeft: isHovering ? 3 : 1,
            borderColor: "divider",
          }}
        />
      </Box>
    </Slide>
  );
};

export default AISideBar;
