## 1. 依赖安装

- [x] 1.1 在frontend目录安装MUI X Chat依赖：`npm install @mui/x-chat`
- [x] 1.2 在frontend目录安装zustand依赖：`npm install zustand`
- [x] 1.3 验证package.json中依赖已正确添加

## 2. 类型定义

- [x] 2.1 创建frontend/src/types/aiAssistant.ts文件
- [x] 2.2 定义MessageRole枚举（User、Assistant）
- [x] 2.3 定义AIChatSession接口（id、title、messages、createdAt、updatedAt）
- [x] 2.4 定义AIMessage接口（id、role、content、timestamp）

## 3. 状态管理

- [x] 3.1 创建frontend/src/stores/aiAssistant.ts文件
- [x] 3.2 定义AIAssistantState接口（sidebarOpen、sessions、currentSessionId等）
- [x] 3.3 实现setSidebarOpen方法
- [x] 3.4 实现loadSessions方法（调用后端AIService）
- [x] 3.5 实现createSession方法（调用后端AIService）
- [x] 3.6 实现deleteSession方法（调用后端AIService）
- [x] 3.7 实现setCurrentSessionId方法
- [x] 3.8 实现loadMessages方法（调用后端AIService）

## 4. Genkit Adapter

- [x] 4.1 创建frontend/src/components/ai/GenkitAdapter.tsx文件
- [x] 4.2 实现GenkitAdapter类，实现ChatAdapter接口
- [x] 4.3 实现sendMessage方法（调用AIService.Chat）
- [x] 4.4 实现streamMessage方法（流式输出，初始版本可为空）
- [x] 4.5 实现其他必需的Adapter方法（根据MUI X Chat文档）

## 5. AI助手侧边栏组件

- [x] 5.1 创建frontend/src/components/ai/AIAssistantSidebar.tsx文件
- [x] 5.2 导入MUI X Chat的ChatBox组件
- [x] 5.3 创建GenkitAdapter实例
- [x] 5.4 实现ChatBox组件集成，传入adapter和features配置
- [x] 5.5 添加关闭按钮和工具栏（新建会话、历史会话）
- [x] 5.6 应用主题样式（深色/浅色模式适配）

## 6. 主界面布局集成

- [x] 6.1 修改frontend/src/components/SSHTabs.tsx
- [x] 6.2 从aiAssistant store导入sidebarOpen状态
- [x] 6.3 实现侧边栏容器div，条件渲染（sidebarOpen时显示）
- [x] 6.4 设置SSH Tabs内容宽度（sidebarOpen时70%，否则100%）
- [x] 6.5 添加CSS transition动画（width 0.3s ease）
- [x] 6.6 集成AIAssistantSidebar组件，设置宽度30%
- [x] 6.7 测试侧边栏展开/收起动画效果

## 7. Header按钮

- [x] 7.1 修改frontend/src/components/Header.tsx
- [x] 7.2 从aiAssistant store导入setSidebarOpen方法
- [x] 7.3 在左侧栏添加AI助手按钮图标
- [x] 7.4 实现按钮点击事件：setSidebarOpen(true)
- [x] 7.5 测试按钮点击打开侧边栏功能

## 8. 前端绑定生成

- [x] 8.1 在项目根目录运行：`wails3 generate bindings -ts`
- [x] 8.2 验证frontend/bindings目录下AIService相关绑定已生成
- [x] 8.3 检查TypeScript类型定义是否正确

## 9. 编译测试

- [x] 9.1 在frontend目录运行：`npm run build:dev`
- [x] 9.2 验证编译无错误
- [x] 9.3 运行应用测试侧边栏UI是否正常显示
- [x] 9.4 测试侧边栏展开/收起动画流畅度
