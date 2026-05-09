# AI 助手 V2 开发任务

## 阶段一：基础架构与 Store

### 1.1 创建 AI 助手 Store
- [ ] 创建 `frontend/src/stores/aiAssistant.ts`
- [ ] 定义 `AIChatSession`、`AIMessage` 类型
- [ ] 实现侧边栏开关状态管理
- [ ] 实现会话管理（创建、切换、删除）
- [ ] 实现消息发送和加载状态管理
- [ ] 实现视图切换（chat / history）
- [ ] 集成现有 `AIService` 调用

### 1.2 创建基础组件结构
- [ ] 创建 `frontend/src/components/ai/AIAssistantSidebar.tsx`（容器组件）
- [ ] 创建 `frontend/src/components/ai/AIToolbar.tsx`（顶部工具栏）
- [ ] 创建 `frontend/src/components/ai/AIChatPanel.tsx`（聊天面板）
- [ ] 创建 `frontend/src/components/ai/AIMessageList.tsx`（消息列表）
- [ ] 创建 `frontend/src/components/ai/AIMessageItem.tsx`（消息项）
- [ ] 创建 `frontend/src/components/ai/AIMessageInput.tsx`（输入框）
- [ ] 创建 `frontend/src/components/ai/AIHistoryList.tsx`（历史列表）
- [ ] 创建 `frontend/src/components/ai/AIHistoryItem.tsx`（历史项）

## 阶段二：布局与动画

### 2.1 修改主界面布局
- [ ] 修改 `frontend/src/pages/App.tsx`
- [ ] 集成 `AIAssistantSidebar` 到主布局
- [ ] 实现终端区域宽度动态变化（100% ↔ 70%）
- [ ] 实现侧边栏宽度动画（0 ↔ 30%）
- [ ] 实现内容淡入淡出动画
- [ ] 设置最小宽度 320px 保护

### 2.2 左侧栏按钮
- [ ] 修改 `frontend/src/components/Header.tsx`
- [ ] 添加 AI 助手按钮（图标 + Tooltip）
- [ ] 按钮状态与侧边栏开关同步（高亮/普通）
- [ ] 点击切换侧边栏展开/收起

## 阶段三：聊天功能

### 3.1 消息列表
- [ ] 实现消息气泡布局（用户右对齐、AI 左对齐）
- [ ] 实现头像展示（用户/AI）
- [ ] 实现时间戳（hover 显示）
- [ ] 实现空状态提示
- [ ] 实现自动滚动到底部
- [ ] 实现代码块渲染和复制按钮

### 3.2 输入框
- [ ] 实现多行输入框（TextField multiline）
- [ ] 实现 Enter 发送、Shift+Enter 换行
- [ ] 实现发送按钮状态（启用/禁用）
- [ ] 实现停止生成按钮（加载时替换发送按钮）
- [ ] 实现输入框自动聚焦

### 3.3 AI 交互
- [ ] 集成 `AIService` 调用
- [ ] 实现加载状态展示
- [ ] 实现停止生成功能
- [ ] 实现错误消息展示
- [ ] 实现重试机制

## 阶段四：会话管理

### 4.1 历史会话列表
- [ ] 实现历史会话视图切换
- [ ] 实现列表项渲染（标题 + 相对时间）
- [ ] 实现按更新时间倒序排列
- [ ] 实现空状态
- [ ] 实现删除按钮（hover 显示）

### 4.2 会话操作
- [ ] 实现新建会话功能
- [ ] 实现切换会话功能
- [ ] 实现删除会话确认对话框
- [ ] 实现会话标题自动生成（首条消息前 20 字）

## 阶段五：终端联动

### 5.1 右键菜单扩展
- [ ] 修改 `frontend/src/components/TerminalContextMenu.tsx`
- [ ] 添加"发送到 AI 提问"菜单项
- [ ] 实现选中文本获取
- [ ] 实现点击后打开侧边栏并填充输入框

### 5.2 边界处理
- [ ] 处理选中文本过长截断
- [ ] 处理侧边栏已打开时的填充逻辑

## 阶段六：旧功能移除

### 6.1 移除终端内嵌 AI 面板
- [ ] 删除 `frontend/src/components/ai/TerminalAICommandPanel.tsx`
- [ ] 修改 `frontend/src/components/Terminal.tsx`，移除 `TerminalAICommandPanel` 引用
- [ ] 修改 `frontend/src/components/StatusBar.tsx`，移除 AI 切换逻辑
- [ ] 清理 `frontend/src/stores/aiCommand.ts`（或标记废弃）

## 阶段七：主题与优化

### 7.1 主题适配
- [ ] 验证深色模式下气泡颜色正确
- [ ] 验证浅色模式下气泡颜色正确
- [ ] 验证输入框、按钮、分隔线主题变量使用正确

### 7.2 性能优化
- [ ] 消息列表大数据量时虚拟滚动（如需要）
- [ ] 动画性能优化（使用 transform 替代 width 如可能）

## 阶段八：测试验证

### 8.1 功能测试
- [ ] 侧边栏展开/收起动画流畅
- [ ] 按钮切换状态正确
- [ ] 新建/切换/删除会话正常
- [ ] 消息发送/接收/停止正常
- [ ] 终端右键发送到 AI 正常
- [ ] 主题切换适配正确

### 8.2 边界测试
- [ ] 空输入禁用发送
- [ ] 网络错误提示正确
- [ ] 长文本输入处理
- [ ] 快速切换会话稳定性
