## Why

用户需要在客户端历史版本列表中删除不需要的历史记录，同时服务端版本存储上限需要增加以支持更多历史版本保留。

## What Changes

- **客户端**: 在历史版本列表每项后添加删除按钮，点击后弹框确认再执行删除操作
- **客户端**: 同步设置页面历史记录列表只展示前10个版本
- **客户端**: 恢复数据弹框内历史记录支持滚动加载，每次加载50条
- **服务端**: 版本记录上限从5个增加到500个

## Capabilities

### New Capabilities
<!-- Capabilities being introduced. Replace <name> with kebab-case identifier (e.g., user-auth, data-export, api-rate-limiting). Each creates specs/<name>/spec.md -->

### Modified Capabilities
- `sync-client`: 添加历史版本删除功能，包括删除按钮UI和确认弹框交互；历史记录列表分页展示（设置页显示前10个，恢复弹框滚动加载每次50条）
- `sync-server`: 版本存储上限从5个增加到500个

## Impact

- 前端UI组件: 历史版本列表页面
- 客户端服务: 添加删除历史版本API调用
- 服务端存储: 版本数量限制配置变更
