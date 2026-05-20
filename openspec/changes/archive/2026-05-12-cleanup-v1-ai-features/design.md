# Cleanup V1 AI Features - Design

## Context

AI 命令助手 V2 采用全局侧边栏架构，已在前期阶段完成开发（阶段一至阶段四）。V1 版本基于终端内嵌面板（`TerminalAICommandPanel`），功能已被 V2 完全覆盖。V1 与 V2 数据不兼容，V1 的 `GenerateCommand`/`ExplainCommand` API 和前端状态管理均与 V2 的新架构（基于会话的多轮对话）无关。当前代码库中同时存在 V1 和 V2 两套 AI 代码，造成维护冗余和潜在的编译/绑定警告。

## Goals / Non-Goals

**Goals：**
- 彻底移除所有 V1 AI 命令助手的前端组件（`TerminalAICommandPanel`、`AIConfigDialog`）、状态管理（`aiCommand` store）和类型定义。
- 彻底移除后端 `AIService` 中与 V1 相关的绑定方法（`GenerateCommand`、`ExplainCommand`）、仅被 V1 使用的私有辅助方法（`SessionContext`、`getSessionContext`、`executeCommand`）及内部实现。
- 删除 `internal/ai/prompts.go`（仅被 V1 调用）。
- 更新所有引用 V1 AI 功能的父组件（`Terminal` 布局收缩逻辑、`StatusBar` AI 按钮、`TerminalContextMenu` 菜单项），确保应用编译和运行正常。
- 重新生成 Wails 前端绑定，确保绑定代码与后端同步，无残留 V1 类型。

**Non-Goals：**
- 不修改 V2 AI 助手（侧边栏、Genkit 适配器、会话持久化等）的任何代码。
- 不修改非 AI 相关的功能（如 SSH 终端、书签、SFTP、同步等）。
- 不引入新的功能或 UI 行为。

## Decisions

### 1. 一次性删除而非渐进式废弃

**选择**：直接删除 V1 相关文件和代码，不保留废弃标记或兼容层。

**理由**：V1 和 V2 完全不兼容，不存在逐步迁移的用户数据；V2 已完全替代 V1 功能，保留废弃代码只会增加维护成本。

**替代方案**：标记 `@deprecated` 并保留一版本后删除。被否决，因为无数据兼容需求，且会增加一次无意义的发布周期。

### 2. 按依赖顺序从叶子到根清理

**清理顺序**：
1. 后端 Go 代码：先移除 `internal/ai/prompts.go`（仅被 V1 调用），再移除 `internal/ai/genkit.go` 中的 V1 类型与方法，最后移除 `services/ai_service.go` 中的 V1 公开方法及仅被 V1 使用的私有辅助方法（`SessionContext` 等）。
2. 前端绑定：重新生成 `wails3 generate bindings -ts`，让绑定层自动同步。
3. 前端叶子组件：删除 `TerminalAICommandPanel.tsx`、`AIConfigDialog.tsx`。
4. 前端状态管理：删除 `aiCommand.ts` store。
5. 前端类型定义：清理 `types/ai.ts` 中仅 V1 使用的类型。
6. 前端父组件：移除 `Terminal.tsx` 中的引用与布局收缩逻辑、`StatusBar.tsx` 中的 AI 按钮、`TerminalContextMenu.tsx` 中的菜单项。

**理由**：从依赖树的叶子开始删除，确保每一步编译都通过，避免中间状态出现未定义引用。

### 3. `types/ai.ts` 部分保留

`AIConfig`、`AIContext`、`AIProvider` 等配置相关类型在 V2 中继续使用，因此不删除整个文件，仅移除 V1 专属类型（`AICommandRequest`、`AICommandResponse`、`ExplainCommandRequest`、`ExplainCommandResponse`、`CommandPart`、`AIMessage` 等）。

## Risks / Trade-offs

| 风险 | 缓解措施 |
|------|----------|
| 误删 V2 共享类型（如 `AIConfig`、`SafetyLevel`） | 清理前逐一确认类型的使用位置；优先删除仅被 `aiCommand.ts` 或 `TerminalAICommandPanel.tsx` 引用的类型。`safety.go`、`providers.go` 被 V2 配置功能使用，保留不删。 |
| 重新生成绑定后，V2 调用的 Go 方法签名发生变化 | V2 使用 `Chat`、`CreateSession` 等方法，与 V1 的 `GenerateCommand` 方法无交集，风险极低。生成后执行编译验证。 |
| 父组件（如 `Terminal`）因移除子组件和布局收缩逻辑导致显示异常 | 移除 `TerminalAICommandPanel` 和 `AIConfigDialog` 后，将 `Terminal.tsx` 的 flex 从条件收缩恢复为 `flex: 1`，确保终端始终占据 100% 宽度。 |
| 删除 `internal/ai/prompts.go` 时误删 V2 使用的 prompt 构建函数 | `prompts.go` 中仅包含 `BuildCommandGeneratorPrompt` 和 `BuildCommandExplainerPrompt`，均被 V1 独占使用；V2 的聊天 prompt 构建在 `genkit.go` 的 `buildChatPrompt` 中，不会受影响。 |
| 状态栏移除 AI 按钮后视觉不平衡 | 状态栏右侧原有三个图标（AI、隧道、传输列表），移除 AI 后仅剩两个，属于正常变化，无需额外处理。 |

## Migration Plan

1. **开发阶段**：按上述顺序逐个删除/清理文件，每步执行 `go build` 和 `npm run build:dev` 验证。
2. **验证阶段**：运行完整应用，确认：
   - 状态栏无 AI 按钮。
   - 终端右键菜单无 V1 AI 相关菜单项。
   - V2 全局 AI 侧边栏正常工作（新建会话、聊天、终端联动）。
3. **绑定同步**：执行 `wails3 generate bindings -ts` 并提交更新后的绑定文件。
4. **回滚策略**：本变更纯为删除操作，回滚即恢复被删除文件（通过 Git 撤销）。

## Open Questions

- `frontend/src/stores/aiConfig.ts`（V1/V2 共享的 AI 配置 store）是否需要调整？→ 不调整，V2 继续使用。
- `services/ai_service.go` 中 V2 的 `Chat`、`CreateSession` 等方法是否会受 V1 方法删除影响？→ 不会，二者内部实现独立。
- `frontend/src/components/AIConfigDialog.tsx` 是否被 V2 使用？→ 否，仅由 `aiCommandStore`（V1）的 `showConfigDialog` 状态控制打开，V2 使用独立的 `AISettings.tsx` 页面进行配置，因此 `AIConfigDialog` 一并删除。
- `internal/ai/prompts.go` 是否被 V2 使用？→ 否，其中仅包含 `BuildCommandGeneratorPrompt` 和 `BuildCommandExplainerPrompt`，均被 V1 独占调用，V2 的聊天使用 `genkit.go` 中的 `buildChatPrompt`，因此 `prompts.go` 整文件删除。
- `services/ai_service.go` 中的 `SessionContext`、`getSessionContext`、`executeCommand` 是否被 V2 使用？→ 否，仅被 V1 的 `GenerateCommand`/`ExplainCommand` 调用，V2 的 `Chat` 不依赖 SSH 会话上下文，因此一并移除。
