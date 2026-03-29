## Context

当前 `SSHTabs.tsx` 使用 MUI 的 Tabs 组件渲染 SSH 标签页，标签页顺序由 `sshTabs` 数组决定。用户需要通过拖拽来调整标签页顺序，以便更好地组织多个 SSH 会话。

技术栈：React + TypeScript + Material UI + Zustand

## Goals / Non-Goals

**Goals:**
- 实现标签页鼠标拖拽排序功能
- 提供流畅的拖拽视觉反馈
- 保持拖拽后的会话状态
- 支持取消拖拽操作

**Non-Goals:**
- 标签页顺序持久化到本地存储
- 触摸设备手势支持
- 跨窗口拖拽

## Decisions

### 使用 @hello-pangea/dnd 作为拖拽库
选择 `@hello-pangea/dnd`（react-beautiful-dnd 的活跃维护分支）：
- 完全兼容 react-beautiful-dnd API，社区成熟
- 支持 React 18 和严格模式
- 对 MUI 组件有良好的兼容性
- 提供流畅的拖拽动画和视觉反馈

### 保持现有 Tabs 组件结构
在现有 MUI Tabs 基础上集成拖拽，而非替换为自定义组件：
- 保持 MUI 的滚动按钮和样式一致性
- 最小化 UI 改动，降低回归风险
- 使用 `DragDropContext` 和 `Droppable` 包裹 Tabs，每个 Tab 使用 `Draggable`

### Store 层添加 reorder 方法
在 `useSSHTabsStore` 中添加 `reorderTabs` 方法：
```typescript
reorderTabs: (oldIndex: number, newIndex: number) => void
```
- 使用数组操作交换位置
- 保持其他状态（currentTab、sshInfo）不变

## Risks / Trade-offs

**[Risk] 拖拽与 Tabs 滚动按钮冲突** → 确保拖拽仅在标签区域内生效，滚动按钮保持独立

**[Risk] 拖拽过程中意外关闭标签** → 拖拽状态下禁用标签关闭按钮

**[Risk] 大量标签页时性能下降** → 使用 `react-window` 或虚拟化列表（如有需要）

**[Trade-off] 不支持触摸设备** → 当前主要面向桌面端 SSH 客户端，触摸支持为低优先级
