## Context

Vexo 是一个基于 Go + Wails v3 + React 的 SSH 桌面 GUI 应用。目前已有设置窗口和命令面板窗口的实现模式可以参考。

当前窗口管理方式：
- `AppInstance` 统一管理应用窗口实例（MainWindow, SettingWindow, CommandWindow）
- 每个窗口通过独立的路由加载（如 `/#/setting`, `/#/command`）
- 窗口创建使用 `application.WebviewWindowOptions` 配置
- 窗口关闭时通过 `OnWindowEvent` 清理实例引用

前端技术栈：
- React Router v7.11 用于路由管理
- Material UI v7.3 用于 UI 组件
- 路由定义在 `frontend/src/routes.ts`

## Goals / Non-Goals

**Goals:**
- 在 Header 中添加工具菜单按钮，点击打开工具窗口
- 工具窗口首屏展示常用运维工具卡片列表
- 点击卡片进入具体工具使用界面
- 工具界面采用左侧列表导航 + 右侧内容区的布局
- 支持工具间的快速切换
- 初始实现 3 个基础工具：端口检测、编码解码、正则测试

**Non-Goals:**
- 不实现工具插件化架构（后续可扩展）
- 不实现用户自定义工具
- 不实现工具配置持久化（首版简化）

## Decisions

### 1. 窗口管理方式
参考 CommandWindow 的实现模式，新增 ToolWindow：
- 在 `AppInstance` 中添加 `ToolWindow` 字段
- 创建 `newToolWindow()` 函数，加载 `/#/tools` 路由
- 窗口尺寸：1200x800，无边框设计

**Rationale**: 与现有设置窗口、命令面板保持一致，降低维护成本。

### 2. 前端路由结构
```
/tools           -> 工具首屏（卡片列表）
/tools/:toolId   -> 具体工具界面
```

**Rationale**: 简单的路由结构，支持直接访问特定工具。

### 3. 工具界面布局
采用左侧固定宽度列表 + 右侧自适应内容区的经典布局：
- 左侧：工具列表导航，宽度 200px
- 右侧：工具具体内容区域
- 顶部：返回按钮 + 当前工具标题

**Rationale**: 与命令面板的设计模式保持一致，用户学习成本低。

### 4. 工具数据结构
Go 后端定义：
```go
type Tool struct {
    ID          string `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description"`
    Icon        string `json:"icon"`
    Category    string `json:"category"`
}
```

前端通过 `ToolService.GetTools()` 获取可用工具列表。

### 5. 初始工具集
- **端口检测** (`port-check`): TCP 端口连通性测试（后端实现）
- **编码解码** (`encoder`): Base64、URL、HTML 编解码（后端实现）
- **正则测试** (`regex`): 正则表达式匹配测试（后端实现）
- **富文本编辑器** (`rich-text`): 基于 react-quill 的富文本编辑器（前端实现，支持源代码查看）

### 6. 工具逻辑实现方式
大部分工具的核心逻辑在 Go 后端实现，前端负责界面展示和调用。富文本编辑器为纯前端实现：

- **端口检测**: `ToolService.CheckPort(host string, port int) PortCheckResult`
- **编码解码**: `ToolService.Encode(toolType string, input string) string` / `ToolService.Decode(toolType string, input string) string`
- **正则测试**: `ToolService.RegexMatch(pattern string, text string, flags string) RegexMatchResult`
- **富文本编辑器**: 前端使用 `react-quill` 实现，支持可视化编辑和 HTML 源代码查看

## Risks / Trade-offs

| Risk | Mitigation |
|------|------------|
| 工具数量增多后代码膨胀 | 每个工具作为独立组件，按需加载 |
| 工具界面风格不统一 | 提供统一的 ToolLayout 包装组件 |
| 窗口过多内存占用 | 工具窗口单例模式，关闭时释放 |

## Migration Plan

无需迁移，新功能增量添加。

## Open Questions

1. 是否需要支持工具搜索功能？（首版暂不实现）
2. 是否需要工具使用历史记录？（首版暂不实现）
