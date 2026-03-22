## Why

用户在主界面需要快速备份数据到同步服务器，而不必打开设置页面。这提供了更便捷的数据保护方式，提升用户体验。

## What Changes

- 在 Header.tsx 组件的设置按钮前添加备份上传按钮
- 点击按钮弹出确认对话框，询问是否上传备份
- 确认后执行与 SyncSettings 组件中"备份数据"按钮相同的 `UploadSync` 功能
- 上传期间显示全屏 Loading 组件覆盖整个界面
- 上传完成后关闭 Loading 并显示成功/失败提示

## Capabilities

### New Capabilities
- `header-backup-upload`: 在 Header 组件中添加快捷备份上传功能，包含确认对话框和全屏 Loading 状态管理

### Modified Capabilities
- （无）

## Impact

- **Header.tsx**: 添加备份按钮、确认对话框状态、Loading 状态、上传逻辑
- **SyncSettings.tsx**: 参考现有的 `handleUpload` 函数实现
- **Loading.tsx**: 复用现有组件，可能需要支持全屏模式
