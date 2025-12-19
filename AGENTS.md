# AGENTS.md

本项目是一个基于 Go 语言与 Wails v3 框架构建的 SSH 桌面 GUI 应用，旨在为用户提供跨平台、安全、高效的终端连接体验。本文档描述了项目所采用的技术栈、架构规范、开发约定以及各模块职责。

---

## 1. 技术栈概览

| 类别           | 技术/库                          | 版本要求     |
|----------------|----------------------------------|--------------|
| 后端语言       | Go                               | ≥ 1.24       |
| 桌面框架       | Wails                            | v3.x         |
| 前端 UI 库     | Material UI (MUI)                | v7.3.x       |
| 前端路由       | React Router                     | v7.11.x      |
| 终端模拟器     | XTerm.js                         | v5.5.x       |
| 构建工具       | Wails CLI + Vite                 | 默认集成     |

---

## 2. 项目结构规范

```
.
├── cmd/
│   └── app/                  # Wails 主程序入口
├── internal/
│   ├── agents/               # SSH 连接代理逻辑（Go）
│   ├── models/               # 数据模型（如会话配置、密钥等）
│   └── services/             # 业务服务层（如会话管理、日志记录）
├── frontend/
│   ├── src/
│   │   ├── components/       # React UI 组件（使用 MUI）
│   │   ├── pages/            # 页面级组件（配合 React Router）
│   │   ├── hooks/            # 自定义 React Hooks
│   │   ├── lib/              # 第三方库封装（如 xterm.js 集成）
│   │   ├── routes/           # 路由配置（React Router v7.11）
│   │   └── App.tsx           # 根组件
│   └── public/
├── wails.json                # Wails 项目配置
└── go.mod
```

---

## 3. Agent 设计规范

### 3.1 SSH Agent 职责

- **连接管理**：建立、维持、关闭 SSH 会话。
- **认证支持**：支持密码、私钥（含 passphrase）、Agent Forwarding。
- **通道复用**：单个连接支持多终端标签页（multiplexing）。
- **安全隔离**：每个会话运行在独立 goroutine，避免状态污染。
- **错误处理**：提供结构化错误码，便于前端展示用户友好提示。

### 3.2 接口约定（Go ↔ Frontend）

所有前端调用通过 Wails 的 `runtime.Events` 和 `bridge methods` 实现：

```go
// 示例：暴露给前端的方法
func (a *App) Connect(sessionID string, config *SSHConfig) error
func (a *App) SendInput(sessionID, data string)
func (a *App) Disconnect(sessionID string)
func (a *App) ListSessions() []SessionInfo
```

前端通过 `window.go.app.*` 调用上述方法，并监听事件如：

- `session.output`：接收终端输出
- `session.status`：连接状态变更（connecting / connected / disconnected / error）

---

## 4. 前端集成规范

### 4.1 XTerm.js 集成

- 使用 `@xterm/xterm@5.5` 及 `@xterm/addon-fit@0.10`、`@xterm/addon-web-links@0.12` 等官方插件。
- 每个终端实例封装为 React 组件 `<Terminal sessionID={id} />`。
- 输入通过 `term.onData()` 捕获并转发至 Go 后端。
- 输出通过监听 `session.output` 事件写入 `term.write()`。

### 4.2 路由设计（React Router v7.11）

- `/`：会话列表页
- `/session/:id`：终端会话页（支持多标签）
- `/settings`：全局设置（主题、默认 SSH 配置等）
- `/connections/new`：新建连接表单

使用 `createBrowserRouter` + `RouterProvider` 实现声明式路由。

### 4.3 UI 规范（Material UI v7.3）

- 主题：使用 `createTheme` 定制深色/浅色模式。
- 组件：优先使用 MUI 官方组件（如 `Tabs`, `Drawer`, `TextField`, `Button`）。
- 响应式：适配桌面窗口缩放，最小宽度 800px。
- 键盘导航：支持 Tab 切换、Ctrl+T 新建标签等快捷键。

---

## 5. 安全与最佳实践

- **密钥存储**：敏感信息（如私钥、密码）不持久化明文，使用系统 Keychain（macOS）或 Credential Vault（Windows）。
- **输入过滤**：后端对所有来自前端的 SSH 配置参数进行校验与转义。
- **超时控制**：空闲连接自动断开（可配置）。
- **日志脱敏**：日志中禁止记录密码、私钥内容。

---

## 6. 开发与调试

- 启动开发模式：
  ```bash
  wails3 dev
  ```
- 前端热重载由 Vite 提供，Go 后端修改需重启。
- 使用 `wails inspect` 调试嵌入式 WebView（基于 Chromium）。

---

## 7. 工具使用

如果没有相应知识库，Context7工具可用的情况下使用Context7查询最新知识

---
