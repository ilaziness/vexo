## Why

当前 SSH 标签页按照固定顺序排列，用户无法根据工作习惯调整标签页顺序。当打开多个连接时，无法将相关的会话放在一起，影响工作效率。添加拖拽排序功能可以让用户自由组织标签页，提升多会话管理体验。

## What Changes

- 为 SSH 标签页添加鼠标拖拽排序功能
- 用户可以通过拖拽标签页来调整其顺序
- 拖拽时显示视觉反馈（如半透明效果、占位符）
- 标签顺序调整后保持状态直到应用关闭

## Capabilities

### New Capabilities
- `tab-drag-reorder`: 标签页拖拽排序功能，支持鼠标拖拽调整 SSH 标签页顺序

### Modified Capabilities
- 无

## Impact

- 前端组件 `SSHTabs.tsx` 需要集成拖拽功能
- Zustand store `ssh.ts` 需要添加标签页重排序方法
- 需要引入拖拽库 `@hello-pangea/dnd`（react-beautiful-dnd 的维护分支）
- 不影响后端服务或数据持久化
