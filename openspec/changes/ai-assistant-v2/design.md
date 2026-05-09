## Context

Vexo 是一个基于 Go + Wails v3 + React 的 SSH 桌面 GUI 应用。当前已有一个基于终端内嵌的 AI 命令助手（`TerminalAICommandPanel`），功能包括命令生成、命令解释、发送到终端执行等。

当前 AI 相关技术栈：

- 前端状态管理：Zustand（`aiCommand.ts` store）
- UI 组件库：Material UI v7.3
- AI 框架：Go 后端使用 Genkit 框架集成 AI 能力
- 后端服务：`AIService` 暴露 `GenerateCommand`、`ExplainCommand` 等方法到前端
- 现有类型定义：`frontend/src/types/ai.ts`

当前主界面布局：

- `App.tsx`：横向布局，`Header`（左侧栏，宽度 42px）+ `SSHTabs`（终端区域，flex: 1）
- `SSHTabBody`：每个标签页包含 Terminal/SFTP 切换、终端内容区、StatusBar
- `Terminal` 组件内嵌 `TerminalAICommandPanel`，通过 Slide 动画从右侧滑出，宽度 40%

## Goals / Non-Goals

**Goals:**

- 在左侧栏（`Header`）新增 AI 助手按钮，点击切换右侧 AI 助手侧边栏
- AI 助手侧边栏宽度占主界面 30%，终端区域相应收缩为 70%
- 侧边栏展开/收起均有平滑动画（宽度变化 + 内容淡入淡出）
- 侧边栏顶部工具栏：新建会话、历史会话、关闭按钮
- 对话区域采用 ChatGPT 风格的聊天气泡布局
- 底部输入框支持多行（Shift+Enter 换行，Enter 发送）
- 支持多轮会话管理（新建、切换历史会话、删除会话）
- AI 回复时显示加载动画，支持停止生成
- 支持选中终端内容并发送到 AI 提问
- 主题完全跟随主应用深色/浅色模式

**Non-Goals:**

- 不实现 AI 模型切换 UI（在设置页面已完成配置）
- 不实现历史会话的云端同步
- 不实现 AI 助手的独立窗口模式
- 不实现语音输入/输出
- 不实现文件附件上传

## Decisions

### 1. 布局架构调整

当前 `App.tsx` 布局：

```
<Box flexDirection="row" height="100%" width="100%">
  <Header />           <!-- 左侧栏，固定宽度 -->
  <SSHTabs />          <!-- 终端区域，flex: 1 -->
</Box>
```

新版布局：

```
<Box flexDirection="row" height="100%" width="100%">
  <Header />           <!-- 左侧栏，固定宽度 -->
  <Box flex={1} display="flex" flexDirection="row" overflow="hidden">
    <SSHTabs />        <!-- 终端区域，宽度动态：100% 或 70% -->
    <AIAssistantSidebar />  <!-- AI 助手侧边栏，宽度 30%，可收起 -->
  </Box>
</Box>
```

**Rationale**: 将 AI 侧边栏提升到 `App.tsx` 层级，与 `SSHTabs` 同级，便于统一控制宽度分配和动画。

### 2. 侧边栏动画实现

使用 MUI 的 `Collapse` 或 CSS `transition` 实现宽度动画：

- 展开：width 从 0 过渡到 30%，同时内容区域 opacity 从 0 到 1
- 收起：width 从 30% 过渡到 0，同时内容区域 opacity 从 1 到 0
- 动画时长：300ms，缓动函数：`cubic-bezier(0.4, 0, 0.2, 1)`
- 终端区域宽度同步过渡，避免布局跳变

**Rationale**: 宽度动画比 Slide 更符合侧边栏的认知，且能让终端区域同步收缩，视觉更连贯。

### 3. 会话数据模型

```typescript
interface AIChatSession {
  id: string; // 会话唯一标识
  title: string; // 会话标题（首条用户消息前 20 字，或默认"新会话"）
  messages: AIMessage[]; // 消息列表
  createdAt: number; // 创建时间戳
  updatedAt: number; // 最后更新时间戳
  model?: string; // 使用的模型（预留）
}

interface AIMessage {
  id: string;
  role: "user" | "assistant" | "system";
  content: string;
  timestamp: number;
  isError?: boolean; // 是否为错误消息
  isStreaming?: boolean; // 是否正在流式输出
}
```

**Rationale**: 与现有 `AIMessage` 类型兼容，扩展会话层级管理。

### 4. 会话管理策略

- **新建会话**：创建新的 `AIChatSession`，当前会话自动保存到历史列表，切换到新会话
- **历史会话列表**：按 `updatedAt` 倒序排列，展示标题 + 最后消息时间
- **切换会话**：点击历史会话项，加载该会话的消息列表到当前对话区
- **删除会话**：点击删除按钮，弹出确认对话框，确认后从列表移除
- **当前会话状态**：独立存储当前激活的会话 ID，与历史列表分离

**Rationale**: 与常见 AI 聊天应用（ChatGPT、Claude）的交互模式保持一致。

### 5. 状态管理重构

将原有的 `aiCommand.ts` store 重构为 `aiAssistant.ts` store：

```typescript
interface AIAssistantState {
  // 侧边栏状态
  isOpen: boolean;
  isLoading: boolean;
  error: string | null;

  // 会话状态
  sessions: AIChatSession[];
  currentSessionId: string | null;

  // 视图状态
  view: "chat" | "history"; // 当前显示聊天界面还是历史列表

  // 输入状态
  input: string;
  abortController: AbortController | null;

  // 操作方法
  toggleSidebar: () => void;
  openSidebar: () => void;
  closeSidebar: () => void;

  // 会话操作
  createSession: () => void;
  switchSession: (sessionId: string) => void;
  deleteSession: (sessionId: string) => void;
  updateSessionTitle: (sessionId: string, title: string) => void;

  // 消息操作
  sendMessage: (content: string, context?: string) => Promise<void>;
  stopGeneration: () => void;
  clearCurrentSession: () => void;

  // 视图切换
  showChatView: () => void;
  showHistoryView: () => void;
}
```

**Rationale**: 统一管理与 AI 助手相关的所有状态，避免多个 store 分散逻辑。

### 6. 终端内容联动

在 `TerminalContextMenu` 中新增选项：

- 当用户选中终端文本并右键时，菜单显示"发送到 AI 提问"
- 点击后，自动打开 AI 侧边栏（如未打开），将选中文本填充到输入框
- 用户可继续补充问题后发送

**Rationale**: 利用现有右键菜单机制，最小化改动，提供自然的交互路径。

### 7. AI 服务调用方式

复用现有 `AIService`，但扩展为通用对话接口：

- 当前已有 `GenerateCommand` 和 `ExplainCommand`
- 新增或复用通用 `Chat` 接口，接收消息历史数组，返回流式响应
- 如后端暂无可用的通用对话接口，可先用现有接口过渡，PRD 中预留扩展

**Rationale**: 后端 AI 框架（Genkit）已集成，前端只需调用暴露的服务方法。

### 8. 消息展示样式

采用 ChatGPT 风格聊天气泡：

- 用户消息：右对齐，主色调（primary）背景，白色文字，圆角气泡
- AI 消息：左对齐，灰色背景（grey.100 或根据主题调整），主文字色，圆角气泡
- 每条消息上方显示发送时间（可选，hover 显示）
- 代码块使用等宽字体，带复制按钮
- AI 消息头像使用 `AutoAwesomeIcon` 或机器人图标
- 用户消息头像使用"我"文字或用户图标

**Rationale**: 用户已确认使用 ChatGPT 风格，且 MUI 组件易于实现。

### 9. 历史会话列表样式

- 侧边栏内嵌历史列表视图，通过顶部"历史会话"按钮切换
- 列表项展示：会话标题（单行截断）、最后消息时间（相对时间，如"2 小时前"）
- 每项右侧有删除按钮（hover 显示）
- 空状态：显示"暂无历史会话"提示

**Rationale**: 内嵌切换避免弹窗/下拉菜单的层级复杂度。

### 10. 主题适配

- AI 助手侧边栏完全跟随主应用主题（深色/浅色）
- 聊天气泡在深色模式下使用调整后的色值，确保对比度
- 输入框、按钮、分隔线等均使用 MUI 主题变量

**Rationale**: 用户已确认跟随主主题，保持视觉一致性。

## Risks / Trade-offs

| Risk                                     | Mitigation                                       |
| ---------------------------------------- | ------------------------------------------------ |
| 移除旧的终端内嵌 AI 面板可能影响习惯用户 | 在更新日志中说明，新入口更易访问                 |
| 会话历史存储在内存，刷新后丢失           | 后续迭代增加本地持久化（localStorage 或 SQLite） |
| 宽度 30% 在小屏幕上可能过宽              | 设置最小宽度 320px，或后续支持拖拽调整           |
| 流式响应实现复杂度                       | 首版可用轮询或完整响应过渡，后续优化为流式       |

## Migration Plan

1. 创建 `aiAssistant.ts` store，实现会话管理和消息发送逻辑
2. 创建 `AIAssistantSidebar` 组件（侧边栏容器 + 动画）
3. 创建 `AIChatPanel` 组件（消息列表 + 输入框）
4. 创建 `AIHistoryList` 组件（历史会话列表）
5. 修改 `App.tsx`，集成侧边栏到主布局
6. 修改 `Header.tsx`，添加 AI 助手按钮
7. 修改 `TerminalContextMenu.tsx`，添加"发送到 AI"选项
8. 移除 `TerminalAICommandPanel.tsx` 及 `Terminal.tsx`、`StatusBar.tsx` 中的相关引用
9. 清理旧的 `aiCommand.ts` store（或保留兼容）
10. 测试验证所有交互路径

## Open Questions

1. 是否需要将会话历史持久化到本地数据库？（首版可存内存，后续迭代）
2. AI 通用对话接口是否已就绪？如未就绪，是否先用 `GenerateCommand` 过渡？
3. 是否需要支持导出会话为文本/Markdown？
