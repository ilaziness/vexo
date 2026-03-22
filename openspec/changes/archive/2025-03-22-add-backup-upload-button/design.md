## Context

当前应用的数据备份功能位于设置页面的"同步"标签页中，用户需要打开设置窗口才能执行备份操作。为了提供更便捷的备份方式，需要在主界面的 Header 组件中添加一个快捷备份按钮。

Header.tsx 位于 `frontend/src/components/Header.tsx`，是一个垂直侧边栏组件，包含 Logo、新连接、书签、主题切换、命令面板和设置按钮。

SyncSettings.tsx 中已实现 `handleUpload` 函数，调用 `UploadSync()` 绑定方法执行备份上传，包含进度轮询和状态管理。

Loading.tsx 是一个可复用的加载组件，支持自定义消息和大小。

## Goals / Non-Goals

**Goals:**
- 在 Header 的设置按钮前添加备份上传按钮
- 点击后弹出确认对话框
- 确认后执行与 SyncSettings 相同的备份上传逻辑
- 上传期间显示全屏 Loading 覆盖整个界面
- 完成后显示成功/失败消息提示

**Non-Goals:**
- 不修改 SyncSettings 组件的现有功能
- 不添加新的后端 API
- 不修改备份上传的核心逻辑

## Decisions

**1. 使用 MUI Dialog 组件实现确认对话框**
- 使用 `@mui/material` 的 Dialog、DialogTitle、DialogContent、DialogActions
- 与现有代码风格保持一致（SyncSettings 中已使用 Dialog）

**2. 复用 SyncSettings 的 UploadSync 调用逻辑**
- 直接调用 `UploadSync()` 绑定方法
- 使用相同的进度轮询机制（`GetSyncProgress`）
- 保持错误处理方式一致（`parseCallServiceError`）

**3. 全屏 Loading 实现方式**
- 使用 MUI 的 Backdrop 组件实现全屏遮罩
- 在 Backdrop 中嵌入 Loading 组件
- Backdrop 设置 `sx={{ zIndex: (theme) => theme.zIndex.modal + 1 }}` 确保在最上层

**4. 按钮位置**
- 放置在设置按钮之前（ThemeSwitcher 和命令面板之后）
- 使用 CloudUpload 图标，与 SyncSettings 中的备份按钮保持一致

## Risks / Trade-offs

**[Risk] 用户可能误触备份按钮**
→ Mitigation: 添加确认对话框，明确提示用户操作后果

**[Risk] 备份期间用户无法操作界面**
→ Mitigation: 全屏 Loading 是预期行为，防止用户在备份过程中进行其他操作导致数据不一致

**[Risk] 同步配置未设置时点击备份按钮**
→ Mitigation: 在点击时检查配置，如未配置则提示用户先前往设置页面配置同步参数
