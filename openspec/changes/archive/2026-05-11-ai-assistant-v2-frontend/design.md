## Context

AI命令助手V2需要全新的前端UI实现，采用全局侧边栏架构替代原有的终端内嵌方案。阶段二负责搭建前端基础架构。

**当前状态**：
- 后端已在阶段一完成（数据库、AIService扩展、Genkit Flow）
- 前端尚未实现MUI X Chat集成
- 需要从零搭建侧边栏UI组件

**约束条件**：
- 必须使用MUI X Chat组件库（当前版本9.0.0-alpha.5）
- 必须使用Zustand进行状态管理
- 侧边栏宽度30%，支持动画过渡
- 需要适配深色/浅色主题
- 前端调用Go方法需要catch错误并使用parseCallServiceError解析

## Goals / Non-Goals

**Goals:**
- 搭建基于MUI X Chat的AI助手侧边栏UI
- 实现Zustand状态管理，管理侧边栏开关、会话数据
- 实现GenkitAdapter连接后端AIService
- 集成侧边栏到主界面布局（SSHTabs）
- 在Header中添加AI助手入口按钮

**Non-Goals:**
- 不实现具体的功能逻辑（流式输出、错误处理、会话管理等）- 这些在阶段三完成
- 不实现终端联动功能 - 这些在阶段四完成
- 不清理旧功能代码 - 这些在阶段五完成

## Decisions

**1. 使用MUI X Chat而非自研组件**
- **理由**：MUI X Chat提供完整的聊天UI组件（消息列表、输入框、代码块渲染、markdown支持），大幅减少开发工作量
- **替代方案**：完全自研聊天组件 - 工作量大，难以达到MUI的UI质量

**2. Zustand作为状态管理方案**
- **理由**：轻量级、简单易用，适合管理侧边栏状态和会话数据
- **替代方案**：React Context + useReducer - 代码冗余，性能较差

**3. GenkitAdapter模式连接后端**
- **理由**：MUI X Chat使用Adapter模式抽象后端连接，实现GenkitAdapter可复用现有AIService
- **替代方案**：直接在组件中调用AIService - 违反单一职责原则，难以测试

**4. 侧边栏宽度使用CSS transition动画**
- **理由**：CSS transition性能好，实现简单
- **替代方案**：JavaScript动画 - 性能较差，代码复杂

**5. 侧边栏集成在SSHTabs组件中**
- **理由**：侧边栏与终端区域是平级布局关系，SSHTabs控制整个SSH Tabs区域
- **替代方案**：在App组件中集成 - 布局层级混乱，难以管理

## Risks / Trade-offs

**[风险] MUI X Chat为alpha版本，稳定性未知**
- **缓解**：密切关注官方更新，遇到问题及时降级或寻找替代方案

**[风险] Genkit Adapter实现复杂度未知**
- **缓解**：参考MUI X Chat官方文档和示例，逐步实现必需方法

**[风险] 宽度动画可能影响性能**
- **缓解**：使用CSS transform而非width属性（如果性能问题），限制动画时长<300ms

**[权衡] 侧边栏固定30%宽度 vs 可拖拽调整**
- **决策**：固定30%宽度，简化实现，满足基本需求
- **权衡**：牺牲灵活性换取开发速度

## Migration Plan

1. 安装依赖：npm install @mui/x-chat zustand
2. 创建类型定义文件
3. 创建Zustand store
4. 创建AIAssistantSidebar组件（初始版本，集成MUI X Chat）
5. 实现GenkitAdapter（基础版本，仅实现必需方法）
6. 修改SSHTabs集成侧边栏布局
7. 修改Header添加AI助手按钮
8. 生成前端绑定：wails3 generate bindings -ts
9. 前端编译测试：npm run build:dev

**回滚策略**：如遇问题可回退到阶段一完成状态，前端代码可随时撤销

## Open Questions

- MUI X Chat的ChatAdapter接口具体需要实现哪些方法？需要查阅官方文档确认
- Genkit Flow的流式输出是否支持？如不支持，阶段三需要调整实现
- MUI X Chat的ChatConversationList组件是否需要单独集成还是ChatBox内置？
