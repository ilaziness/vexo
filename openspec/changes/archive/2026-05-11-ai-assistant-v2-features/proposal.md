## Why

阶段二已完成前端基础架构搭建，但GenkitAdapter仅实现了基础版本，缺少流式输出、错误处理、会话管理等核心功能。阶段三负责实现这些关键功能，使AI助手具备完整的对话交互能力。

## What Changes

- 实现GenkitAdapter的流式输出（streamMessage方法）
- 实现停止生成功能（调用后端StopGeneration）
- 实现错误处理和重试机制（网络错误、AI服务异常）
- 实现Zustand Store与后端数据同步（会话列表、消息列表）
- 实现ChatConversationList数据绑定（会话列表显示、切换、删除）
- 实现会话标题自动生成（首条消息前20字）
- 实现删除会话确认对话框

## Capabilities

### New Capabilities
- `ai-streaming-output`: AI回复流式输出，实时显示生成内容
- `ai-stop-generation`: 停止AI生成功能，中断正在进行的回复
- `ai-error-handling`: 错误处理和重试机制，处理网络错误和AI服务异常
- `ai-session-title`: 会话标题自动生成（基于首条用户消息）
- `ai-session-management`: 会话列表数据绑定，支持切换和删除确认

### Modified Capabilities
- `genkit-adapter`: 扩展GenkitAdapter实现流式输出和错误处理
- `ai-assistant-store`: 完善状态管理，实现后端数据同步和错误状态

## Impact

**前端代码**：
- 修改frontend/src/components/ai/GenkitAdapter.tsx（实现流式输出、错误处理）
- 修改frontend/src/components/ai/AIAssistantSidebar.tsx（集成停止生成、错误提示）
- 修改frontend/src/stores/aiAssistant.ts（添加数据同步、错误状态）
- 新增frontend/src/components/ai/SessionList.tsx（会话列表组件）

**后端代码**：
- 无需修改（后端已在阶段一完成StopGeneration接口）

**用户体验**：
- AI回复实时流式显示，提升交互体验
- 支持停止生成，避免长时间等待
- 错误时友好提示，支持重试
- 会话标题自动生成，无需手动输入
- 删除会话需确认，防止误操作
