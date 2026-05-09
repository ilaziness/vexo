## Why

当前 Vexo 的 AI 助手以终端内嵌面板形式存在（`TerminalAICommandPanel`），功能聚焦于命令生成与解释，且与单个 SSH 会话强耦合。随着使用场景扩展，用户需要一个独立的、全局可访问的 AI 助手，能够：

1. 不依赖当前终端会话，随时进行通用技术问答、运维知识查询
2. 支持多轮会话管理，保留历史对话上下文
3. 与终端内容联动（选中终端文本发送到 AI 提问）
4. 提供更完善的对话交互体验（聊天气泡、加载状态、停止生成等）

因此需要将 AI 助手从终端内嵌面板升级为全局侧边栏形态的独立功能模块。

## What Changes

- 在主界面左侧栏（`Header` 组件区域）新增 AI 助手入口按钮
- 点击按钮在主界面右侧弹出 AI 助手侧边栏，按钮可切换侧边栏的展开/收起状态
- 侧边栏展开时占据主界面宽度的 30%，终端区域收缩为 70%；收起后终端恢复 100% 宽度
- 侧边栏弹出和收起均有平滑动画效果
- 侧边栏顶部工具栏包含：新建会话、历史会话、关闭按钮
- 中间区域为对话消息列表，采用类似 ChatGPT 的聊天气泡样式
- 底部为输入框（支持多行）和发送按钮
- 支持多轮会话管理：新建会话、切换历史会话、删除历史会话
- 支持选中终端内容发送到 AI 提问
- AI 回复时显示加载动画，支持停止生成
- 主题跟随主应用深色/浅色模式

## Capabilities

### New Capabilities
- `ai-assistant-sidebar`: AI 助手侧边栏主组件
- `ai-chat-panel`: AI 对话面板（消息列表 + 输入区域）
- `ai-session-manager`: 会话管理（新建、切换、删除、历史列表）
- `ai-message-list`: 消息列表组件（聊天气泡样式）
- `ai-message-input`: 多行输入框组件
- `ai-terminal-integration`: 终端内容选中并发送到 AI

### Modified Capabilities
- `header`: 左侧栏新增 AI 助手按钮
- `app-layout`: 主界面布局适配侧边栏展开/收起
- `terminal-context-menu`: 右键菜单增加"发送到 AI 提问"选项

### Removed Capabilities
- `terminal-ai-panel`: 移除原有的终端内嵌 AI 命令面板（`TerminalAICommandPanel`）
- `status-bar-ai-toggle`: 移除状态栏的 AI 助手切换入口

## Impact

- 前端：
  - 新增 AI 助手侧边栏相关组件
  - 修改 `App.tsx` 布局以支持侧边栏
  - 修改 `Header.tsx` 增加 AI 助手按钮
  - 移除 `TerminalAICommandPanel` 及相关引用
  - 移除状态栏 AI 切换逻辑
  - 新增/修改 Zustand store 管理 AI 会话状态
- Go 服务层：
  - 复用现有 `AIService`，无需新增后端接口
  - 可能需要扩展会话持久化存储（如需要保存历史会话到本地数据库）
- 用户体验：
  - AI 助手从终端附属功能升级为独立全局功能
  - 操作路径更直观，左侧栏一键唤起

## Migration Plan

1. 创建新的 AI 助手侧边栏组件和 store
2. 修改 `App.tsx` 布局，集成侧边栏
3. 修改 `Header.tsx`，添加 AI 助手按钮，移除旧入口
4. 修改 `Terminal.tsx` 和 `TerminalContextMenu.tsx`，添加选中内容发送到 AI 功能
5. 移除 `TerminalAICommandPanel.tsx` 和状态栏相关 AI 切换代码
6. 测试验证：侧边栏展开/收起动画、会话管理、终端联动、主题适配
