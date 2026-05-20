## Why

Vexo AI 命令助手 V2 需要全新的后端基础设施来支持全局侧边栏架构和多轮对话管理。现有的 V1 版本使用终端内嵌方案，无法满足 V2 的需求（独立会话管理、历史对话保留、终端联动等）。Phase 1 后端开发将建立数据持久化层和扩展 AI 服务，为前端功能提供必要的基础支撑。

## What Changes

- **数据库迁移**：在现有 SQLite 数据库基础上新增表结构（ai_sessions、ai_messages）用于持久化 AI 会话和消息数据
- **数据访问层**：实现 AISessionRepository 接口，提供会话和消息的 CRUD 操作
- **服务层扩展**：扩展 AIService 新增会话管理方法（CreateSession、GetSession、ListSessions、DeleteSession）和多轮对话方法（Chat、StopGeneration）
- **AI 流程扩展**：扩展 Genkit Flow 以支持多轮对话上下文传递
- **前端绑定生成**：运行 `wails3 generate bindings -ts` 生成最新的 TypeScript 绑定

## Capabilities

### New Capabilities
- `ai-session-persistence`: AI 会话和消息的本地数据库持久化能力
- `multi-turn-chat`: 支持历史上下文的多轮对话能力
- `session-management`: 会话创建、查询、列表、删除管理能力

### Modified Capabilities
- None（V2 为全新开发，不修改现有 spec 行为）

## Impact

- **修改文件**：
  - `internal/database/` - 新增 AI 会话相关数据库访问代码和迁移脚本
  - `services/ai_service.go` - 扩展会话管理和多轮对话方法
  - `internal/ai/genkit.go` - 扩展 Genkit Flow 支持多轮对话
- **生成文件**：
  - `frontend/bindings/` - 通过 wails3 生成的 TypeScript 绑定文件
- **数据库迁移**：在现有 SQLite 数据库基础上执行 schema 迁移创建新表
