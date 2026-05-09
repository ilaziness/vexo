# AI 助手侧边栏 (AI Assistant Sidebar)

## 概述

AI 助手侧边栏是 Vexo 主界面的核心功能模块，以右侧可收起侧边栏形式提供通用 AI 对话能力。用户通过左侧栏按钮唤起，侧边栏占据主界面 30% 宽度，支持多轮会话管理、终端内容联动、流畅动画过渡。

## 功能需求

### 1. 侧边栏布局与动画

#### 1.1 布局规格
- 侧边栏位于主界面右侧，与终端区域同级
- 展开宽度：主界面宽度的 30%
- 收起宽度：0px
- 终端区域宽度随侧边栏状态动态变化：展开时 70%，收起时 100%
- 侧边栏内部采用纵向 flex 布局：顶部工具栏 + 中间内容区 + 底部输入区

#### 1.2 动画规格
- 展开动画：
  - 侧边栏 width: 0 → 30%
  - 终端区域 width: 100% → 70%
  - 侧边栏内容 opacity: 0 → 1
  - 时长：300ms
  - 缓动：cubic-bezier(0.4, 0, 0.2, 1)
- 收起动画：
  - 侧边栏 width: 30% → 0
  - 终端区域 width: 70% → 100%
  - 侧边栏内容 opacity: 1 → 0
  - 时长：300ms
  - 缓动：cubic-bezier(0.4, 0, 0.2, 1)
- 内容区域在宽度动画完成后（或同步）淡入淡出，避免内容挤压

#### 1.3 响应式行为
- 侧边栏设置最小有效宽度 320px，若 30% 小于 320px 则按 320px 计算
- 收起状态下，侧边栏完全不可见，不占用布局空间

### 2. 顶部工具栏

#### 2.1 布局
- 高度：48px
- 背景：与主题背景色一致（`background.paper`）
- 底部分隔线：1px，`divider` 颜色
- 内边距：左右 12px
- 内容横向排列，两端对齐

#### 2.2 元素
- 左侧：AI 助手标题 + 图标
  - 图标：`AutoAwesomeIcon`，颜色 `primary.main`
  - 标题："AI 助手"，Typography variant="subtitle1"
- 右侧：三个图标按钮（间距 4px）
  1. **新建会话按钮**
     - 图标：`AddCircleOutlineIcon` 或 `EditNoteIcon`
     - Tooltip："新建会话"
     - 点击：创建新会话，切换到聊天视图
  2. **历史会话按钮**
     - 图标：`HistoryIcon`
     - Tooltip："历史会话"
     - 点击：切换到历史会话列表视图；若已在历史视图，则返回聊天视图
  3. **关闭按钮**
     - 图标：`CloseIcon`
     - Tooltip："关闭侧边栏"
     - 点击：收起侧边栏

### 3. 中间内容区

#### 3.1 视图切换
内容区支持两种视图切换：
- **聊天视图**（默认）：展示当前会话的消息列表
- **历史视图**：展示历史会话列表

切换方式：通过顶部工具栏的"历史会话"按钮切换，无页面跳转。

#### 3.2 聊天视图

##### 3.2.1 消息列表
- 占据内容区全部空间，可纵向滚动
- 消息按时间正序排列（从上到下）
- 自动滚动到底部（新消息到达时）
- 空状态：当当前会话无消息时，显示空状态提示
  - 居中显示 `AutoAwesomeIcon`（大图标，透明度 0.5）
  - 提示文字："有什么可以帮你的？"
  - 下方可展示示例问题（如"如何查看端口占用？"、"解释这条命令的含义"）

##### 3.2.2 消息气泡样式
- **用户消息**：
  - 布局：横向 flex，row-reverse（右对齐）
  - 头像：32px 圆形，背景 `primary.main`，显示"我"文字或 `PersonIcon`
  - 气泡：最大宽度 80%，背景 `primary.light`（深色主题下需调整），文字白色，圆角 12px（左下、右下、左上圆角，右上小圆角），内边距 12px
  - 时间：hover 时显示在气泡下方，Typography variant="caption"，颜色 `text.secondary`
- **AI 消息**：
  - 布局：横向 flex，row（左对齐）
  - 头像：32px 圆形，背景 `secondary.main`，显示 `AutoAwesomeIcon`
  - 气泡：最大宽度 80%，背景 `grey.100`（深色主题下使用 `grey.800` 或调整），文字 `text.primary`，圆角 12px（左下、右下、右上圆角，左上小圆角），内边距 12px
  - 时间：hover 时显示
- **代码块**：
  - 使用等宽字体（`fontFamily: monospace`）
  - 背景色与气泡略有区分（如 `grey.200` / `grey.700`）
  - 圆角 4px，内边距 8px
  - 右上角显示复制按钮（小图标按钮，`ContentCopyIcon`）
- **错误消息**：
  - AI 角色，气泡边框或背景使用 `error.light` 色调
  - 显示错误图标

##### 3.2.3 加载状态
- AI 回复中显示加载指示器：
  - 位置：消息列表底部，AI 消息位置
  - 样式：`CircularProgress` size=24，配合闪烁的三个点动画或纯进度圈
  - 旁边显示文字："正在思考..."
- 加载期间，输入框禁用发送，显示停止按钮

##### 3.2.4 停止生成
- 加载期间，输入框右侧的发送按钮变为停止按钮
- 停止按钮图标：`StopIcon`，颜色 `error.main`
- 点击后中断当前 AI 响应，保留已生成的部分内容（如有）

#### 3.3 历史会话视图

##### 3.3.1 列表布局
- 占据内容区全部空间，可纵向滚动
- 列表项按 `updatedAt` 倒序排列（最新的在最上方）
- 每项高度：约 64px
- 每项内边距：12px 16px
- 项之间分隔线：1px，`divider` 颜色

##### 3.3.2 列表项内容
- 左侧：会话标题
  - Typography variant="body2"，单行显示，超出截断（`textOverflow: ellipsis`）
  - 标题规则：取会话第一条用户消息的前 20 个字符，若无则显示"新会话"
- 右侧：最后更新时间
  - Typography variant="caption"，颜色 `text.secondary`
  - 格式：相对时间（"刚刚"、"5 分钟前"、"2 小时前"、"昨天"、日期）
- 最右侧：删除按钮（hover 时显示或始终显示）
  - 图标：`DeleteOutlineIcon`，颜色 `error.main`（hover 时）
  - 点击：弹出确认对话框

##### 3.3.3 空状态
- 无历史会话时，居中显示
- 图标：`HistoryIcon`，大图标，透明度 0.5
- 文字："暂无历史会话"

##### 3.3.4 确认删除对话框
- 使用 MUI `Dialog` 组件
- 标题："删除会话"
- 内容："确定要删除会话\"{title}\"吗？此操作不可恢复。"
- 按钮："取消"（次要）、"删除"（主要，error 颜色）

### 4. 底部输入区

#### 4.1 布局
- 高度：自适应（最小约 56px，最大约 120px）
- 背景：`background.paper`
- 顶部分隔线：1px，`divider` 颜色
- 内边距：12px 16px

#### 4.2 输入框
- 使用 MUI `TextField` 组件，multiline
- 最大行数：5 行
- 最小行数：1 行
- placeholder："输入你的问题..."（按主题语言）
- 禁用状态：AI 加载中或初始化中
- 自动聚焦：打开侧边栏或切换到新会话时聚焦

#### 4.3 发送按钮
- 位置：输入框右侧（使用 `InputAdornment`）
- 图标：`SendIcon`
- 颜色：`primary.main`（可发送时），`text.disabled`（不可发送时）
- 禁用条件：输入为空或仅空白字符、AI 加载中

#### 4.4 停止按钮
- 位置：替换发送按钮
- 图标：`StopIcon`
- 颜色：`error.main`
- 显示条件：AI 正在生成回复时

#### 4.5 键盘交互
- Enter：发送消息
- Shift + Enter：换行
- 输入框为空时 Enter 不发送

### 5. 终端内容联动

#### 5.1 右键菜单扩展
- 在 `TerminalContextMenu` 中新增菜单项："发送到 AI 提问"
- 显示条件：用户选中了终端文本（选区非空）
- 位置：在现有菜单项中合适位置插入（如复制/粘贴之后）

#### 5.2 发送流程
1. 用户选中终端文本，右键点击
2. 点击"发送到 AI 提问"
3. 自动打开 AI 侧边栏（如未打开）
4. 将选中文本填充到输入框，格式：`"""{选中文本}"""`
5. 输入框获得焦点，用户可补充问题后发送

#### 5.3 边界情况
- 选中文本过长（> 2000 字符）：截断并提示
- 侧边栏已打开：直接填充输入框，不重复打开
- 当前正在生成回复：填充输入框，等待用户手动发送

### 6. 主题适配

#### 6.1 深色模式
- 侧边栏背景：`background.paper`
- AI 消息气泡背景：`grey.800`
- 用户消息气泡背景：`primary.dark` 或调整后的 `primary.main`
- 输入框背景：`background.default`
- 文字颜色：使用 MUI 主题 `text.primary` / `text.secondary`

#### 6.2 浅色模式
- 侧边栏背景：`background.paper`
- AI 消息气泡背景：`grey.100`
- 用户消息气泡背景：`primary.light`
- 输入框背景：`background.default`

### 7. 错误处理

#### 7.1 AI 服务错误
- 错误以 AI 消息形式展示在对话列表中
- 消息样式使用错误色调
- 内容显示错误信息（如"服务暂不可用，请稍后重试"）
- 提供"重试"按钮，重新发送最后一条用户消息

#### 7.2 网络错误
- 检测网络异常，提示用户检查网络连接
- 支持重试机制

#### 7.3 输入校验
- 空输入或仅空白字符：禁用发送按钮
- 输入过长（> 4000 字符）：提示截断或分片发送

## 接口需求

### 前端 Store 接口

```typescript
// AI 助手状态管理（Zustand Store）
interface AIAssistantState {
  isOpen: boolean;
  isLoading: boolean;
  error: string | null;
  sessions: AIChatSession[];
  currentSessionId: string | null;
  view: "chat" | "history";
  input: string;

  toggleSidebar: () => void;
  openSidebar: () => void;
  closeSidebar: () => void;
  createSession: () => void;
  switchSession: (sessionId: string) => void;
  deleteSession: (sessionId: string) => void;
  sendMessage: (content: string, context?: string) => Promise<void>;
  stopGeneration: () => void;
  setInput: (input: string) => void;
  showChatView: () => void;
  showHistoryView: () => void;
}
```

### Go 后端接口（复用/扩展）

复用现有 `AIService` 方法：
- `GenerateCommand(request)`：命令生成（过渡期间可用）
- `ExplainCommand(request)`：命令解释（过渡期间可用）

建议新增通用对话接口：
- `Chat(messages: AIMessage[], context?: string)`：通用对话，返回流式或完整响应

## 组件清单

| 组件名 | 路径 | 职责 |
|--------|------|------|
| AIAssistantSidebar | `components/ai/AIAssistantSidebar.tsx` | 侧边栏容器，管理动画和布局 |
| AIChatPanel | `components/ai/AIChatPanel.tsx` | 聊天视图（消息列表 + 输入框） |
| AIMessageList | `components/ai/AIMessageList.tsx` | 消息列表渲染 |
| AIMessageItem | `components/ai/AIMessageItem.tsx` | 单条消息气泡 |
| AIMessageInput | `components/ai/AIMessageInput.tsx` | 底部输入框和发送按钮 |
| AIHistoryList | `components/ai/AIHistoryList.tsx` | 历史会话列表视图 |
| AIHistoryItem | `components/ai/AIHistoryItem.tsx` | 单条历史会话项 |
| AIToolbar | `components/ai/AIToolbar.tsx` | 顶部工具栏 |
| aiAssistant.ts | `stores/aiAssistant.ts` | Zustand 状态管理 |

## 验收标准

- [ ] 左侧栏点击 AI 助手按钮，右侧侧边栏平滑展开，宽度占 30%
- [ ] 再次点击按钮，侧边栏平滑收起，终端恢复 100% 宽度
- [ ] 侧边栏展开/收起动画时长 300ms，视觉流畅无卡顿
- [ ] 顶部工具栏包含新建会话、历史会话、关闭三个按钮，功能正常
- [ ] 聊天视图空状态显示正确，示例问题可点击填充输入框
- [ ] 用户消息右对齐，AI 消息左对齐，气泡样式符合设计
- [ ] 输入框支持多行，Shift+Enter 换行，Enter 发送
- [ ] AI 回复时显示加载动画，可点击停止按钮中断
- [ ] 历史会话列表按时间倒序排列，展示标题和相对时间
- [ ] 点击历史会话可切换到对应会话，消息加载正确
- [ ] 删除历史会话弹出确认对话框，确认后移除
- [ ] 新建会话创建空白会话，原会话保存到历史列表
- [ ] 终端右键菜单包含"发送到 AI 提问"，点击后填充输入框
- [ ] 深色/浅色模式切换时，AI 侧边栏样式正确适配
- [ ] 输入为空时发送按钮禁用
- [ ] AI 服务错误时，对话列表显示错误消息
