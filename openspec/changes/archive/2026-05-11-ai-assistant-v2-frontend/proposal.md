## Why

AI命令助手V2需要全新的前端UI实现，采用全局侧边栏架构替代原有的终端内嵌方案。阶段二负责搭建前端基础架构，包括安装MUI X Chat依赖、创建核心组件、实现状态管理和布局集成，为后续功能实现奠定基础。

## What Changes

- 安装前端依赖：@mui/x-chat、zustand
- 创建类型定义文件（frontend/src/types/aiAssistant.ts）
- 创建Zustand状态管理store（frontend/src/stores/aiAssistant.ts）
- 创建AIAssistantSidebar.tsx组件（集成MUI X Chat）
- 实现GenkitAdapter.tsx（连接后端AIService）
- 修改SSHTabs.tsx集成侧边栏布局（30%宽度，动画过渡）
- 修改Header.tsx添加AI助手入口按钮

## Capabilities

### New Capabilities
- `ai-assistant-sidebar`: 全局AI助手侧边栏UI组件，基于MUI X Chat实现聊天界面
- `ai-assistant-store`: 前端状态管理，管理侧边栏开关状态、会话数据、当前会话等
- `genkit-adapter`: 前端适配器，将MUI X Chat组件连接到后端Genkit服务

### Modified Capabilities
- `main-layout`: 修改主界面布局，支持侧边栏展开/收起的响应式布局
- `header`: 在Header组件中添加AI助手入口按钮

## Impact

**前端依赖**：
- 新增@mui/x-chat（当前版本9.0.0-alpha.5，alpha版本需注意稳定性）
- 新增zustand状态管理库
- MUI X Chat依赖@emotion/react、@emotion/styled

**前端代码**：
- 新增frontend/src/components/ai/目录及组件
- 新增frontend/src/types/aiAssistant.ts类型定义
- 新增frontend/src/stores/aiAssistant.ts状态管理
- 修改frontend/src/components/SSHTabs.tsx布局逻辑
- 修改frontend/src/components/Header.tsx添加按钮

**后端**：
- 无需修改后端代码（后端已在阶段一完成）
- 需重新生成前端绑定：wails3 generate bindings -ts
