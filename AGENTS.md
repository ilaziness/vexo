# AGENTS.md

本项目是一个基于 Go 语言与 Wails v3 框架构建的 SSH 桌面 GUI 应用，旨在为用户提供跨平台、安全、高效的终端连接体验。本文档描述了项目所采用的技术栈、架构规范、开发约定以及各模块职责。

---

## 1. 技术栈概览

| 类别       | 技术/库           | 版本要求   |
| ---------- | ----------------- | ---------- |
| 后端语言   | Go                | ≥ 1.26     |
| 桌面框架   | Wails             | v3.x       |
| 前端框架   | React             | v19.x      |
| 前端 UI 库 | Material UI (MUI) | v9.1.x     |
| 前端路由   | React Router      | v8.x       |
| 终端模拟器 | XTerm.js          | v6.x       |
| 构建工具   | Wails CLI + Vite  | Vite v8.x  |
| AI 框架    | Genkit            | v1.9.x     |

---

## 2. 项目结构规范

```
├── services/                 # Wails 服务层（暴露给前端调用）
│   ├── service.go            # 服务注册与初始化
│   ├── app.go                # 应用生命周期
│   ├── ssh_service.go        # SSH 连接与会话
│   ├── ssh_tunnel.go         # SSH 隧道
│   ├── sftp_service.go       # SFTP 文件传输
│   ├── bookmark_service.go   # 书签管理
│   ├── command_service.go    # 命令管理
│   ├── config_service.go     # 应用配置
│   ├── sync_service.go       # 数据同步
│   ├── ai_service.go         # AI 助手
│   ├── tool_service.go       # 工具箱
│   ├── log_service.go        # 日志
│   └── ...                   # hostkey、transfer、websocket 等
├── internal/                 # 应用内部功能包
│   ├── ai/                   # AI 引擎（Genkit 集成）
│   ├── secret/               # 加解密相关
│   ├── sync/                 # 数据同步客户端（上传/下载/版本管理）
│   ├── updater/              # 自动更新功能
│   ├── utils/                # 通用工具函数
│   ├── system/               # 系统相关功能
│   └── database/             # 数据库访问层（SQLite、迁移、Repository）
├── sync-backend/             # 同步服务端（独立可部署，独立 go.mod）
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

模块化，可复用，减少重复模块代码

可单独封装的包放在`internal`下面，合理取包名称。

和前端交互强相关的功能放在`services`下面对应文件，合理分文件编写，防止一个文件行数太多。

适当的位置添加注释，注释不要太多

### 1. 日志记录

- 前端记录日志使用go绑定到前端的`LogService`里面的方法。
- go里面记录日志使用`LogService`里面的方法。

## 2. AI集成

AI相关功能使用`Genkit`AI开发框架来实现，不是直接调用LLM商的API的方式来做。

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
      Chat = 'chat',
      History = 'history',
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