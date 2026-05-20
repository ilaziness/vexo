# Cleanup V1 AI Features - Tasks

## 1. 后端 Go 代码清理

- [x] 1.1 移除 `services/ai_service.go` 中的 `GenerateCommand` 方法
- [x] 1.2 移除 `services/ai_service.go` 中的 `ExplainCommand` 方法
- [x] 1.3 移除 `services/ai_service.go` 中仅被 V1 使用的 `SessionContext` 类型及 `getSessionContext`、`clearSessionContext`、`executeCommand` 私有方法
- [x] 1.4 移除 `internal/ai/genkit.go` 中的 V1 专属类型：`CommandGenerateRequest`、`CommandGenerateResponse`、`CommandExplainRequest`、`CommandExplainResponse`、`CommandPart`
- [x] 1.5 移除 `internal/ai/genkit.go` 中的 V1 方法：`GenerateCommand`、`explainCommand`、`ExplainCommand`、`explainCommand`
- [x] 1.6 删除 `internal/ai/prompts.go`：仅包含 `BuildCommandGeneratorPrompt` 和 `BuildCommandExplainerPrompt`，均被 V1 独占使用
- [x] 1.7 检查 `internal/ai/` 目录，确认 `safety.go`、`providers.go` 无仅被 V1 使用的残留代码，予以保留
- [x] 1.8 执行 `go build` 验证后端编译通过

## 2. 前端绑定同步

- [x] 2.1 执行 `wails3 generate bindings -ts` 重新生成 TypeScript 绑定
- [x] 2.2 检查生成的绑定目录，确认 `GenerateCommand`、`ExplainCommand` 相关类型和方法已移除
- [x] 2.3 提交更新后的绑定文件

## 3. 前端叶子组件与状态管理清理

- [x] 3.1 删除 `frontend/src/components/ai/TerminalAICommandPanel.tsx`
- [x] 3.2 删除 `frontend/src/components/AIConfigDialog.tsx`
- [x] 3.3 删除 `frontend/src/stores/aiCommand.ts`
- [x] 3.4 执行 `npm run build:dev` 验证编译通过（此时 Terminal.tsx 等会报错，属于预期）

## 4. 前端类型定义清理

- [x] 4.1 清理 `frontend/src/types/ai.ts` 中的 V1 专属类型：`AICommandRequest`、`AICommandResponse`、`AISafetyWarning`、`ExplainCommandRequest`、`ExplainCommandResponse`、`CommandPart`、`AIMessage`、`AIPanelMode`
- [x] 4.2 确认保留的类型（`AIProvider`、`AIConfig`、`AIContext`）仍被 V2 代码引用（`AISettings.tsx`、`aiConfig.ts`）；`ActiveSession` 为 V1 遗留死代码，一并清理
- [x] 4.3 保留 `ai.ts`（仍有 V2 使用的类型未清空）

## 5. 前端父组件引用清理

- [x] 5.1 移除 `frontend/src/components/Terminal.tsx` 中对 `TerminalAICommandPanel`、`AIConfigDialog` 的 import、渲染及 `useAICommandStore` 状态引用
- [x] 5.2 调整 `Terminal.tsx` 布局：将 `flex: aiPanelOpen ? "0 0 60%" : 1` 恢复为 `flex: 1`，移除 `aiPanelOpen` 相关逻辑，确保终端始终占据 100% 宽度
- [x] 5.3 移除 `frontend/src/components/StatusBar.tsx` 中的 `AutoAwesomeIcon` 导入、AI 按钮及 `useAICommandStore` 引用
- [x] 5.4 移除 `frontend/src/components/TerminalContextMenu.tsx` 中与 V1 AI 相关的右键菜单项（"AI 解释此命令"）及对应的 `useAICommandStore`/`explainCommand` 引用
- [x] 5.5 检查其他文件（`frontend/src/components/ai/` 下仅存 `GenkitAdapter.tsx`，无 V1 残留）
- [x] 5.6 执行 `npm run build:dev` 验证前端编译通过

## 6. 最终验证

- [x] 6.1 运行完整应用，确认状态栏无 AI 切换按钮
- [x] 6.2 确认终端内无 V1 AI 面板，终端右键菜单无 "发送到 AI 提问" / "AI 解释此命令" 菜单项
- [x] 6.3 确认 V2 全局 AI 侧边栏功能正常：新建会话、发送消息、停止生成、终端文本联动
- [x] 6.4 执行 `go build` 验证完整后端编译通过
