## Context

阶段二已完成前端基础架构（MUI X Chat集成、Zustand Store、GenkitAdapter基础版本），但缺少核心功能实现。阶段三负责实现流式输出、错误处理、会话管理等关键功能。

**当前状态**：
- GenkitAdapter仅实现了基础sendMessage方法
- Zustand Store定义了方法但未实现完整的数据同步
- ChatConversationList未实现数据绑定
- 无错误处理和重试机制

**约束条件**：
- 必须使用MUI X Chat的流式输出API
- 前端调用Go方法必须catch错误并使用parseCallServiceError解析
- 错误提示使用useMessageStore的方法
- 会话标题自动生成（首条消息前20字）
- 删除会话需确认对话框

## Goals / Non-Goals

**Goals:**
- 实现GenkitAdapter流式输出（streamMessage）
- 实现停止生成功能
- 实现错误处理和重试机制
- 实现Zustand Store与后端数据同步
- 实现ChatConversationList数据绑定
- 实现会话标题自动生成
- 实现删除会话确认对话框

**Non-Goals:**
- 不实现终端联动功能 - 这些在阶段四完成
- 不清理旧功能代码 - 这些在阶段五完成
- 不修改后端代码 - 后端已在阶段一完成

## Decisions

**1. 使用MUI X Chat的流式输出API**
- **理由**：MUI X Chat的ChatAdapter接口支持streamMessage方法，使用官方API保证兼容性
- **替代方案**：手动实现WebSocket或SSE - 复杂度高，与MUI X Chat集成困难
- **实际实现**：会话列表使用自定义 MUI List 组件（非 ChatConversationList），因 ChatConversationList 与当前数据模型集成复杂度较高

**2. 错误处理使用try-catch + parseCallServiceError**
- **理由**：符合AGENTS.md规范，统一错误处理方式
- **替代方案**：直接显示原始错误 - 用户体验差，不统一

**3. 重试机制在前端实现**
- **理由**：用户可手动重试，无需复杂的后端重试逻辑
- **替代方案**：后端实现自动重试 - 增加后端复杂度，可能造成重复请求

**4. 会话标题在创建会话时生成**
- **理由**：用户发送首条消息后立即生成标题，无需用户手动输入
- **替代方案**：用户手动输入标题 - 增加操作步骤，用户体验差

**5. 删除确认使用MUI Dialog**
- **理由**：与现有UI风格一致，MUI Dialog提供标准确认对话框
- **替代方案**：浏览器原生confirm - UI风格不统一，可定制性差

**6. Zustand Store数据同步在组件加载时触发**
- **理由**：侧边栏打开时自动加载会话列表，切换会话时加载消息列表
- **替代方案**：实时轮询 - 性能开销大，不必要

## Risks / Trade-offs

**[风险] Genkit流式输出可能不支持**
- **缓解**：如不支持，退化为非流式输出（等待完整响应后显示）

**[风险] 后端StopGeneration可能不支持**
- **缓解**：如不支持，前端禁用停止生成按钮或显示"不支持"提示

**[风险] 流式输出可能导致消息闪烁**
- **缓解**：使用MUI X Chat的内置流式渲染，避免手动处理DOM

**[权衡] 自动重试 vs 手动重试**
- **决策**：手动重试（失败后显示重试按钮）
- **权衡**：牺牲便利性换取控制权，避免自动重试造成重复请求

**[权衡] 实时同步 vs 按需加载**
- **决策**：按需加载（侧边栏打开/切换会话时加载）
- **权衡**：牺牲实时性换取性能，减少不必要的网络请求

## Migration Plan

1. 实现GenkitAdapter流式输出（streamMessage方法）
2. 实现停止生成功能（调用后端StopGeneration）
3. 实现错误处理和重试机制（try-catch、parseCallServiceError、useMessageStore）
4. 实现Zustand Store数据同步方法（loadSessions、loadMessages等）
5. 实现ChatConversationList组件（会话列表、切换、删除）
6. 实现会话标题自动生成逻辑
7. 实现删除确认对话框（MUI Dialog）
8. 集成所有功能到AIAssistantSidebar
9. 测试流式输出、错误处理、会话管理功能

**回滚策略**：如遇问题可回退到阶段二完成状态，功能代码可随时撤销

## Open Questions

- MUI X Chat的ChatAdapter接口streamMessage方法的具体签名是什么？需要查阅官方文档
- 后端AIService.Chat是否支持流式响应？如不支持，需要调整实现
- 后端StopGeneration接口的具体行为是什么？是否真正中断生成？
