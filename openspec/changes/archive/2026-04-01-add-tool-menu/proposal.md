## Why

运维工作中经常需要使用各种工具（如端口检测、编码解码、正则测试等），目前这些工具分散在各处或需要额外打开其他应用。在 Vexo 中集成常用运维工具，可以让用户在同一个界面内快速访问和使用这些工具，提高工作效率，减少窗口切换。

## What Changes

- 在 Header 组件中添加"工具"菜单按钮
- 点击工具按钮打开新窗口展示工具面板
- 首屏以卡片式布局展示常用运维工具列表
- 点击卡片进入具体工具使用界面
- 工具界面采用左侧列表导航、右侧内容区的布局
- 支持工具间的快速切换
- 增加富文本编辑器工具，支持可视化编辑和HTML源代码查看

## Capabilities

### New Capabilities
- `tool-menu`: 工具菜单入口和窗口管理
- `tool-dashboard`: 工具首屏卡片式展示页面
- `tool-layout`: 工具界面布局（左侧列表+右侧内容）
- `tool-list`: 可用工具列表由后端 ToolService 提供
- `tool-port-check`: 端口连通性检测工具（后端实现检测逻辑）
- `tool-encoder`: 编码解码工具（后端实现编解码逻辑）
- `tool-regex`: 正则表达式测试工具（后端实现匹配逻辑）
- `tool-rich-text`: 富文本编辑器工具（前端使用 react-quill 实现，支持源代码查看）

### Modified Capabilities
- （无现有 spec 需要修改）

## Impact

- 前端：新增工具相关页面组件和路由
- Go 服务层：新增 ToolService 用于管理工具窗口和工具逻辑
- 路由：新增 `/tools` 路由
- 工具列表：由后端 `ToolService.GetTools()` 动态提供
- 工具逻辑：端口检测、编码解码、正则匹配等逻辑在 Go 后端实现
