# AGENTS.md

本项目是一个基于 Go 语言与 Wails v3 框架构建的 SSH 桌面 GUI 应用，旨在为用户提供跨平台、安全、高效的终端连接体验。本文档描述了项目所采用的技术栈、架构规范、开发约定以及各模块职责。

---

## 1. 技术栈概览

| 类别      | 技术/库              | 版本要求    |
|---------|-------------------|---------|
| 后端语言    | Go                | ≥ 1.25  |
| 桌面框架    | Wails             | v3.x    |
| 前端 UI 库 | Material UI (MUI) | v7.3.x  |
| 前端路由    | React Router      | v7.11.x |
| 终端模拟器   | XTerm.js          | v5.5.x  |
| 构建工具    | Wails CLI + Vite  | 默认集成    |

---

## 2. 项目结构规范

```
├── service/                  # 核心逻辑
├── internal                  # 应用内部功能包
│   ├── secret                # 加解密相关
├── frontend/                 # 前端react项目目录
│   ├── bindings/github.com/ilaziness/vexo/services    # go暴露到前端的函数、类型等
│   ├── public                # 公共静态资源
│   ├── src/                  # 前端源文件
│   │   ├── components/       # React UI 组件（使用 MUI）
│   │   ├── func              # 功能函数
│   │   ├── language          # 多语言翻译文件
│   │   ├── pages/            # 页面级组件（配合 React Router）
│   │   ├── stores/           # 状态管理
│   │   ├── styles            # 组件css
│   │   ├── theme             # 主题配色
│   │   ├── types             # ts类型定义
│   │   ├── main.tsx          # 前端react主入口
│   │   ├── routes.ts         # 路由
│   ├── public/
├── main.go                   # main入口   
└── go.mod
```

---

## 3. 编码规范

模块化，可复用，减少重复模块代码

---

## 4. 前端集成规范

React开发时启用了严格模式StrictMode。

### 4.1 XTerm.js 集成

- 使用 `@xterm/xterm@5.5` 及 `@xterm/addon-fit@0.10`、`@xterm/addon-web-links@0.12` 等官方插件。
- 每个终端实例封装为 React 组件 `<Terminal sessionID={id} />`。
- 输入通过 `term.onData()` 捕获并转发至 Go 后端。
- 输出通过监听 `session.output` 事件写入 `term.write()`。

### 4.2 路由设计（React Router v7.11）


### 4.3 UI 规范（Material UI v7.3）

- 主题：使用 `createTheme` 定制深色/浅色模式。
- 组件：优先使用 MUI 官方组件（如 `Tabs`, `Drawer`, `TextField`, `Button`）。

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

---

## 7. 工具使用

如果没有相应知识，Context7工具可用的情况下使用Context7查询最新知识

---
