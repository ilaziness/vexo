# AGENTS.md

本项目是一个基于 Go 语言与 Wails v3 框架构建的 SSH 桌面 GUI 应用，旨在为用户提供跨平台、安全、高效的终端连接体验。本文档描述了项目所采用的技术栈、架构规范、开发约定以及各模块职责。

---

## 1. 技术栈概览

| 类别       | 技术/库           | 版本要求  |
| ---------- | ----------------- | --------- |
| 后端语言   | Go                | ≥ 1.27.1  |
| 桌面框架   | Wails             | v3.x      |
| 前端框架   | React             | v19.x     |
| 前端 UI 库 | Material UI (MUI) | v9.1.x    |
| 前端路由   | React Router      | v8.x      |
| 终端模拟器 | XTerm.js          | v6.x      |
| 构建工具   | Wails CLI + Vite  | Vite v8.x |
| AI 框架    | Genkit            | v1.9.x    |

---

## 2. 项目结构规范

```
├── services/                 # Wails 适配层（唯一允许 import Wails 的业务代码）
│   ├── register.go           # 组合根：注入依赖、注册 Service、Shutdown
│   ├── windows.go            # 多窗口生命周期
│   ├── app_service.go        # 主窗口/对话框/更新/WS 地址
│   ├── ssh_service.go        # SSH RPC 转发与关会话编排
│   ├── ssh_tunnel.go         # 隧道 RPC
│   ├── sftp_service.go       # SFTP RPC 与系统对话框
│   ├── bookmark_service.go   # 书签 RPC
│   ├── command_service.go    # 命令 RPC
│   ├── config_service.go     # 配置 RPC 与口令提示
│   ├── sync_service.go       # 同步 RPC
│   ├── ai_service.go         # AI RPC 与流式事件
│   ├── tool_service.go       # 工具 RPC
│   └── log_service.go        # 日志（持有 *zap.Logger）
├── internal/                 # 引擎与基础设施（禁止 import services / Wails）
│   ├── ssh/                  # 连接、会话、known_hosts、远程信息
│   ├── sftp/                 # SFTP 协议
│   ├── tunnel/               # 端口转发
│   ├── termws/               # 本机终端 WebSocket
│   ├── transfer/             # 传输进度
│   ├── bookmark/             # 书签业务与跳板链解析
│   ├── command/              # 内置/用户命令
│   ├── config/               # TOML 配置
│   ├── secret/               # 加解密与进程内 Vault
│   ├── ai/                   # Genkit 引擎
│   ├── sync/                 # 同步客户端
│   ├── tools/                # 编解码/哈希/端口
│   ├── database/             # SQLite
│   ├── buildinfo/            # ldflags 构建信息
│   ├── updater/
│   ├── system/               # 可执行目录、SafeGo
│   └── utils/
├── sync-backend/             # 同步服务端（独立 go.mod，勿与桌面进程合并）
│   ├── main.go               # 服务端入口
│   ├── server.go             # HTTP 服务
│   ├── config.go             # 服务端配置
│   ├── handler.go            # HTTP 请求处理器
│   ├── service.go            # 业务逻辑（版本管理/用户管理）
│   ├── storage.go            # 文件存储管理
│   ├── model.go              # 数据模型
│   ├── db.go                 # 数据库连接
│   └── rate.go               # 限流
├── build/                    # 各平台构建配置（android/darwin/docker/ios/linux/windows）
├── openspec/                 # OpenSpec 规范驱动开发配置
│   ├── specs/                # 功能规格定义
│   └── changes/              # 变更提案（含 archive/ 归档）
├── docs/                     # 项目文档
├── data/                     # 内置默认数据（如 commands.json）
├── frontend/                 # 前端 React 项目
│   ├── bindings/             # Go 绑定（wails3 generate bindings -ts 生成）
│   │   └── github.com/ilaziness/vexo/services/  # 暴露给前端的函数、类型等
│   ├── public/               # 公共静态资源
│   ├── src/
│   │   ├── components/       # React UI 组件（MUI）
│   │   │   ├── ai/           # AI 助手（AISideBar、GenkitAdapter 等）
│   │   │   ├── settings/     # 设置页子组件
│   │   │   └── subwindow/    # 子窗口组件
│   │   ├── func/             # 功能函数（decode、service、validation 等）
│   │   ├── hooks/            # 自定义 React Hooks
│   │   ├── language/         # 多语言翻译文件
│   │   ├── pages/            # 页面级组件（配合 React Router）
│   │   ├── stores/           # 状态管理（zustand）
│   │   ├── styles/           # 组件 CSS
│   │   ├── theme/            # 主题配色
│   │   ├── types/            # TypeScript 类型定义
│   │   ├── main.tsx          # 前端入口
│   │   └── routes.ts         # 路由定义
│   ├── index.html
│   ├── vite.config.ts
│   └── package.json
├── main.go                   # 应用入口
├── Taskfile.yml              # 构建任务
├── Makefile
└── go.mod
```

---

## 3. 编码规范

模块化，可复用，减少重复模块代码。

分层硬规则：

- `services/`：RPC 导出、Wails 事件、系统对话框、多窗口；把 UI 结果转成引擎入参。不写 SSH handshake、SFTP 拷贝循环、SQL。
- `internal/`：禁止 import `services` 和 Wails。SSH 引擎只接收已解析的 hop 链（`[]ssh.Endpoint`），不认识书签 ID。
- 禁止包级可变业务状态。允许：`embed`、ldflags（`internal/buildinfo`）、错误 sentinel、常量。
- 日志：组合根创建 `*zap.Logger` 并注入 `NewXxx(logger)`；前端仍用 `LogService`。
- 可单独封装的包放在 `internal` 下；不为小函数膨胀 `utils`。
- 和前端交互的导出方法放在 `services`，方法名保持稳定。

适当位置添加注释，不要太多。

### 1. 日志记录

- 前端记录日志使用 go 绑定的 `LogService` 方法。
- Go 内部记录日志使用注入的 `*zap.Logger`，不要新增包级 `Logger`。

### 2. AI集成

AI 相关功能使用 `Genkit` 框架实现，不要直接调用 LLM 商 API。

---

## 4. 前端集成规范

安装依赖使用`npm install`命令，卸载依赖使用`npm uninstall`命令，不要之间便更加`package.json`文件。

React开发时启用了严格模式StrictMode。

Go绑定到前端的方法和数据类型在`frontend\bindings\github.com\ilaziness\vexo\services`。

搜索代码应该排除`dist`目录。

前端调用go生成的绑定方法需要catch错误，例如：

```typescript
try {
  const bookmark = await BookmarkService.GetBookmarkByID(bookmarkID);
  if (!bookmark) {
    errorMessage("bookmark not found");
    return;
  }
  setCurrentTab(newIndex);
} catch (err) {
  errorMessage(parseCallServiceError(err));
}
```

使用`parseCallServiceError`解析错误提示，消息提醒使用`useMessageStore`立马的方法。

### 4.1 XTerm.js 集成

- 使用 `@xterm/xterm` 及 `@xterm/addon-fit`、`@xterm/addon-web-links` 等官方插件。
- 每个终端实例封装为 React 组件 `<Terminal sessionID={id} />`。
- 输入通过 `term.onData()` 捕获并转发至 Go 后端。
- 输出通过监听 `session.output` 事件写入 `term.write()`。

### 4.2 路由设计（React Router v8）

路由渲染页面组件在`frontend\src\pages`，定义于 `frontend\src\routes.ts`（Hash Router）。

- 主路由 `/` 渲染 `<App />`，主界面。
- 子路由：
  - `/setting` 对应 `<Setting />` 组件，应用设置。
  - `/command` 对应 `<Command />` 组件，命令管理。
  - `/tools` 对应 `<Tools />` 组件，工具箱列表。
  - `/tools/:toolId` 对应 `<ToolDetail />` 组件，工具详情。
  - `/submainwindow` 对应 `<SubMainWindow />` 组件，子窗口。

### 4.3 UI 规范（Material UI v9）

- 主题：使用 `createTheme` 定制深色/浅色模式。
- 组件：优先使用 MUI 官方组件（如 `Tabs`, `Drawer`, `TextField`, `Button`）。

### 4.4 TypeScript 类型定义规范

- **枚举类型**：对于有限的选项集合（如视图模式、消息角色、状态等），必须使用 `enum` 定义，不要使用联合类型字符串穷举。
  - ❌ 错误：`view: 'chat' | 'history'`
  - ✅ 正确：
    ```typescript
    enum AIAssistantView {
      Chat = "chat",
      History = "history",
    }
    view: AIAssistantView;
    ```
- **类型定义位置**：共享的 TypeScript 类型定义放在 `frontend/src/types` 目录下，组件内部专用的类型可定义在组件文件中。

---

## 5. 安全与最佳实践

- **密钥存储**：敏感信息（如私钥、密码）不持久化明文，使用系统 Keychain（macOS）或 Credential Vault（Windows）。
- **输入过滤**：后端对所有来自前端的 SSH 配置参数进行校验与转义。
- **超时控制**：空闲连接自动断开（可配置）。
- **日志脱敏**：日志中禁止记录密码、私钥内容。

---

## 6. 开发与调试

**重要**：使用这里列出来的命令编译测试，不要用其他的。

### 1. 前端编译调试

`frontend`执行下面命令编译调试：

```bash
npm run build:dev
```

### 2. Go编译调试

项目根目录执行下面命令编译调试：

生成前端绑定：`wails3 generate bindings -ts`

```bash
# 编译go，项目根目录执行
go build -o {name} .
# 或
go build .
```

`{name}` 为平台对应的可执行文件名，如 `vexo.exe`（Windows）、`vexo`（macOS/Linux）。

# dev运行，使用ctrl+c结束

## wails3 dev

## 7. 工具使用

如果没有相应知识，Context7工具可用的情况下使用Context7查询最新知识

---

## 8. 相关文档

- genkit doc：[guid doc](https://genkit.dev/docs/go/get-started)
- genkit api参考：[api ref](https://pkg.go.dev/github.com/firebase/genkit/go)
- wails V3 api参考: [application api ref](https://pkg.go.dev/github.com/wailsapp/wails/v3/pkg/application)
- wails v3 event: [event ref](https://pkg.go.dev/github.com/wailsapp/wails/v3/pkg/events)
