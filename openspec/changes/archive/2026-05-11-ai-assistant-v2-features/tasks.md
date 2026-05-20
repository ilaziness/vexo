## 1. GenkitAdapter流式输出

- [x] 1.1 研究MUI X Chat的ChatAdapter接口streamMessage方法签名
- [x] 1.2 在GenkitAdapter中实现streamMessage方法
- [x] 1.3 实现流式响应解析（如后端支持）
- [ ] 1.4 测试流式输出效果

## 2. 停止生成功能

- [x] 2.1 在GenkitAdapter中添加stopGeneration方法
- [x] 2.2 调用后端AIService.StopGeneration接口
- [x] 2.3 在AIAssistantSidebar中添加停止生成按钮
- [x] 2.4 实现按钮点击事件调用stopGeneration
- [x] 2.5 测试停止生成功能

## 3. 错误处理和重试机制

- [x] 3.1 在GenkitAdapter的sendMessage中添加try-catch
- [x] 3.2 使用parseCallServiceError解析错误信息
- [x] 3.3 使用useMessageStore显示错误提示
- [x] 3.4 在错误提示中添加重试按钮
- [x] 3.5 实现重试逻辑（重新调用sendMessage）
- [ ] 3.6 测试网络错误场景
- [ ] 3.7 测试AI服务异常场景

## 4. Zustand Store数据同步

- [x] 4.1 实现loadSessions方法（调用AIService.ListSessions）
- [x] 4.2 实现createSession方法（调用AIService.CreateSession）
- [x] 4.3 实现deleteSession方法（调用AIService.DeleteSession）
- [x] 4.4 实现loadMessages方法（调用AIService.ListMessages）
- [x] 4.5 添加错误处理（try-catch、parseCallServiceError、useMessageStore）
- [x] 4.6 添加加载状态（loadingSessions、loadingMessages）
- [x] 4.7 测试数据同步功能

## 5. 会话标题自动生成

- [x] 5.1 在createSession方法中实现标题生成逻辑
- [x] 5.2 提取首条用户消息前20字作为标题
- [x] 5.3 处理空消息或超短消息情况（使用默认标题"新会话"）
- [x] 5.4 在后端创建会话时传递title参数
- [ ] 5.5 测试标题自动生成功能

## 6. 会话列表组件

- [x] 6.1 创建frontend/src/components/ai/SessionList.tsx
- [x] 6.2 使用MUI X Chat的ChatConversationList组件
- [x] 6.3 从aiAssistant store获取sessions数据
- [x] 6.4 实现会话列表渲染（标题、更新时间）
- [x] 6.5 实现点击切换会话（setCurrentSessionId）
- [x] 6.6 实现删除会话功能（调用deleteSession）
- [x] 6.7 集成到AIAssistantSidebar
- [ ] 6.8 测试会话列表功能

## 7. 删除确认对话框

- [x] 7.1 在SessionList中添加删除确认对话框（MUI Dialog）
- [x] 7.2 实现确认逻辑（点击确认后调用deleteSession）
- [x] 7.3 实现取消逻辑（关闭对话框）
- [ ] 7.4 测试删除确认功能

## 8. AIAssistantSidebar集成

- [x] 8.1 集成流式输出到ChatBox
- [x] 8.2 集成停止生成按钮
- [x] 8.3 集成错误提示和重试
- [x] 8.4 集成SessionList组件
- [x] 8.5 实现侧边栏打开时自动加载会话列表
- [x] 8.6 实现切换会话时自动加载消息列表
- [ ] 8.7 测试完整功能流程

## 9. 编译测试

- [x] 9.1 在frontend目录运行：`npm run build:dev`
- [x] 9.2 验证编译无错误
- [ ] 9.3 运行应用测试流式输出
- [ ] 9.4 测试停止生成功能
- [ ] 9.5 测试错误处理和重试
- [ ] 9.6 测试会话管理（新建、切换、删除）
- [ ] 9.7 测试会话标题自动生成
