# Cleanup V1 AI Features

## Why

AI 命令助手 V2 已全新开发完成，采用全局侧边栏架构，提供独立的多轮对话 AI 助手功能。V1 版本的终端内嵌 AI 面板（`TerminalAICommandPanel`）已被 V2 的侧边栏完全替代，不再具备维护价值。为避免代码冗余、减少维护负担，并确保用户始终使用 V2 新架构，需要彻底清理 V1 相关代码。

## What Changes

本次变更将彻底移除 AI 命令助手 V1 的全部前端组件、状态管理、类型定义以及后端绑定的 V1 API，具体包括：

**前端：**
- **删除** `frontend/src/components/ai/TerminalAICommandPanel.tsx`：V1 终端内嵌 AI 面板组件。
- **删除** `frontend/src/stores/aiCommand.ts`：V1 AI 命令相关的 Zustand 状态管理。
- **删除** `frontend/src/components/AIConfigDialog.tsx`：V1 专属的 AI 配置弹窗，仅由 `aiCommandStore` 控制打开。
- **移除** `frontend/src/components/StatusBar.tsx` 中的 AI 切换逻辑：包括 `AutoAwesomeIcon` 图标按钮及其 `togglePanel` 调用。
- **移除** `frontend/src/components/Terminal.tsx` 中对 `TerminalAICommandPanel`、`AIConfigDialog` 的引用与渲染，以及 `aiPanelOpen` 相关的布局收缩逻辑。
- **移除** `frontend/src/components/TerminalContextMenu.tsx` 中与 V1 AI 相关的右键菜单项（"AI 解释此命令"）。
- **清理** `frontend/src/types/ai.ts` 中仅被 V1 使用的类型定义（`AICommandRequest`、`AICommandResponse`、`AISafetyWarning`、`ExplainCommandRequest`、`ExplainCommandResponse`、`CommandPart`、`AIMessage`、`AIPanelMode`）。
- **清理**前端绑定目录中因 V1 API 移除而不再需要的生成代码文件。

**后端：**
- **移除** `services/ai_service.go` 中 V1 专用的接口方法 `GenerateCommand`、`ExplainCommand`。
- **移除** `services/ai_service.go` 中仅被 V1 使用的 `SessionContext` 类型及 `getSessionContext`、`clearSessionContext`、`executeCommand` 私有方法。
- **移除** `internal/ai/genkit.go` 中 V1 专属的类型（`CommandGenerateRequest`、`CommandGenerateResponse`、`CommandExplainRequest`、`CommandExplainResponse`、`CommandPart`）及 `GenerateCommand`、`explainCommand`、`ExplainCommand`、`explainCommand` 方法。
- **删除** `internal/ai/prompts.go`：仅被 V1 的 `generateCommand`/`explainCommand` 调用的 Prompt 构建函数。
- 保留 `internal/ai/safety.go`、`internal/ai/providers.go` 中 V2 继续使用的类型和函数。

## Capabilities

### New Capabilities

（本次变更不涉及新功能，仅做代码清理。）

### Modified Capabilities

（本次变更不涉及现有功能行为变更，仅移除已废弃的 V1 实现。）

## Impact

- **前端**：`TerminalAICommandPanel`、`AIConfigDialog`、`aiCommand` store、StatusBar AI 按钮、Terminal 内嵌引用与布局收缩逻辑、TerminalContextMenu 相关菜单项全部移除。
- **后端**：`AIService.GenerateCommand`、`AIService.ExplainCommand` 及仅被 V1 使用的 `SessionContext`/`getSessionContext`/`executeCommand` 全部移除；`internal/ai/prompts.go` 整文件删除；`internal/ai/genkit.go` 中的 V1 类型与方法移除。
- **绑定代码**：Wails 绑定生成的 TypeScript 类型和调用方法同步更新，V1 专属类型不再生成。
- **构建产物**：无新增产物，代码体积减小。
- **用户体验**：用户界面中不再出现 V1 的 AI 切换按钮、终端内嵌面板和配置弹窗，统一由 V2 全局侧边栏提供 AI 服务。
