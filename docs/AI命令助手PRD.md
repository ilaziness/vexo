# AI 命令助手 - 产品需求文档 (PRD)

## 1. 文档信息

| 项目     | 内容             |
| -------- | ---------------- |
| 产品名称 | Vexo AI 命令助手 |
| 文档版本 | v1.0             |
| 创建日期 | 2026-05-06       |
| 状态     | 待评审           |

---

## 2. 产品概述

### 2.1 背景

Vexo 是一款跨平台 SSH/SFTP 桌面客户端，用户需要频繁在终端中执行各种 Linux/Unix 命令。许多命令语法复杂、参数众多，用户（尤其是初学者）经常需要查阅文档或搜索引擎来获取正确的命令格式，降低了工作效率。

### 2.2 产品目标

在 Vexo 中集成 AI 能力，让用户能够通过自然语言描述需求，AI 自动生成对应的 SSH 命令，并提供命令解释和安全提示，从而：

- 降低复杂命令的学习门槛
- 提高工作效率，减少查阅文档的时间
- 帮助用户理解命令含义，促进学习
- 提供安全警告，避免危险操作

### 2.3 目标用户

- 系统管理员、DevOps 工程师
- 后端开发人员
- 对 Linux 命令不熟悉的初学者
- 需要快速执行复杂任务的高级用户

---

## 3. 功能需求

### 3.1 核心功能

#### 3.1.1 自然语言生成命令

**功能描述**：用户输入自然语言描述，AI 生成对应的 SSH 命令。

**用户流程**：

1. 用户在终端界面或命令面板中点击"AI 命令助手"按钮
2. 弹出 AI 命令输入框
3. 用户输入自然语言描述（如"查找最近7天修改的大于100MB的文件"）
4. AI 生成命令并显示结果
5. AI 生成命令并展示结果（命令、解释、安全警告、替代方案）
6. 用户可以选择：
   - **运行**：调用 `CommandService.SendCommand` 发送到指定 SSH 会话执行
   - **复制到剪贴板**：复制命令文本
   - **编辑**：在输入框中修改命令后重新生成或运行

**示例**：

| 用户输入                         | AI 生成命令                                       |
| -------------------------------- | ------------------------------------------------- |
| 查找最近7天修改的大于100MB的文件 | `find / -size +100M -mtime -7 2>/dev/null`        |
| 查看占用磁盘空间最多的前10个目录 | `du -ah / \| sort -rh \| head -n 10`              |
| 将所有 .log 文件压缩并删除原文件 | `find . -name "*.log" -exec gzip {} \;`           |
| 查看8080端口被哪个进程占用       | `lsof -i :8080` 或 `netstat -tulpn \| grep :8080` |

#### 3.1.2 命令解释

**功能描述**：AI 对生成的命令进行逐部分解释，帮助用户理解。

**展示内容**：

- 命令各部分的含义
- 参数说明
- 预期输出描述

**示例输出**：

```
命令：find / -size +100M -mtime -7 2>/dev/null

解释：
- find /          : 从根目录开始搜索
- -size +100M     : 查找大于100MB的文件
- -mtime -7       : 修改时间在7天以内
- 2>/dev/null     : 将错误输出重定向到空设备（忽略权限错误）

预期：列出符合条件的文件路径列表
```

#### 3.1.3 安全警告

**功能描述**：AI 识别潜在的危险操作并提供警告。

**危险命令类型**：

- 删除操作（rm -rf, 批量删除）
- 覆盖文件操作
- 权限修改（chmod, chown）
- 系统配置修改
- 需要 root 权限的操作

**警告展示**：

- 醒目的警告标识
- 危险等级（高/中/低）
- 具体风险提示
- 建议的安全替代方案

#### 3.1.4 命令优化建议

**功能描述**：AI 提供命令的优化版本或替代方案。

**优化类型**：

- 性能优化（更快的执行方式）
- 安全增强（添加安全参数）
- 兼容性建议（跨平台替代命令）
- 更简洁的写法

### 3.2 辅助功能

#### 3.2.1 命令历史学习

**功能描述**：AI 分析用户的历史命令使用习惯，提供个性化建议。

**实现方式**：

- 记录用户常用的命令模式
- 基于历史使用推荐相似命令
- 学习用户偏好的命令参数

#### 3.2.2 多命令组合

**功能描述**：用户描述复杂任务，AI 生成包含管道、重定向的复合命令。

**示例**：

- 用户输入："统计日志文件中每个IP的访问次数，按降序排列"
- AI 生成：`awk '{print $1}' access.log | sort | uniq -c | sort -rn`

#### 3.2.3 上下文感知

**功能描述**：AI 根据当前终端的工作目录、已执行的命令等上下文信息，提供更准确的建议。

**上下文信息**：

- 当前工作目录（通过解析 SSH 会话输出获取）
- 当前用户权限（通过解析 SSH 会话输出获取）
- 操作系统类型（通过解析 `uname` 输出获取）
- 最近执行的命令（通过解析 SSH 会话输出获取）

**注意**：未连接的 SSH 标签页无法使用 AI 功能，AI 面板在该标签页下置灰或隐藏

---

## 4. 技术架构

### 4.1 整体架构

```
┌─────────────────────────────────────────────────────────────┐
│                      前端 (React + MUI)                      │
│  ┌─────────────────┐  ┌─────────────────┐  ┌──────────────┐ │
│  │  AI命令输入框    │  │  命令结果展示    │  │  安全警告组件 │ │
│  └────────┬────────┘  └────────┬────────┘  └──────┬───────┘ │
│           └────────────────────┼───────────────────┘         │
│                              │                              │
└──────────────────────────────┼──────────────────────────────┘
                               │ Wails IPC
┌──────────────────────────────┼──────────────────────────────┐
│                      后端 (Go)                               │
│                              │                              │
│  ┌───────────────────────────▼───────────────────────────┐  │
│  │              AIService (新建)                          │  │
│  │  ┌─────────────┐  ┌──────────────┐  ┌──────────────┐ │  │
│  │  │ Genkit Flow │  │ 命令安全分析  │  │ 上下文收集器  │ │  │
│  │  │ 命令生成     │  │ 危险命令检测  │  │ 终端状态获取  │ │  │
│  │  └─────────────┘  └──────────────┘  └──────────────┘ │  │
│  └───────────────────────────┬───────────────────────────┘  │
│                              │                              │
│  ┌───────────────────────────▼───────────────────────────┐  │
│  │              SSHService (现有)                         │  │
│  │         命令执行 / 会话管理                            │  │
│  └───────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                               │
                    ┌──────────┴──────────┐
                    │    AI 模型提供商     │
                    │  ┌───────────────┐  │
                    │  │ Ollama (本地) │  │
                    │  │ OpenAI        │  │
                    │  │ Anthropic     │  │
                    │  │ Google AI     │  │
                    │  └───────────────┘  │
                    └─────────────────────┘
```

### 4.2 后端服务设计

#### 4.2.1 AIService 结构

```go
// services/ai_service.go
type AIService struct {
    genkit       *genkit.Genkit        // Genkit 实例
    sshService   *SSHService           // SSH 服务引用
    configService *ConfigService       // 配置服务
    modelProvider string               // 当前使用的模型提供商
    model        string                // 当前使用的模型名称
}

// NewAIService 创建 AI 服务实例
func NewAIService(sshService *SSHService, configService *ConfigService) *AIService {
    return &AIService{
        sshService:    sshService,
        configService: configService,
    }
}

// AICommandRequest AI 命令生成请求
type AICommandRequest struct {
    NaturalLanguage string   `json:"natural_language"` // 自然语言描述
    SessionID       string   `json:"session_id"`       // 可选：当前会话ID
    Context         AIContext `json:"context"`          // 可选：上下文信息
}

// AIContext 上下文信息
type AIContext struct {
    CurrentDirectory string   `json:"current_directory"`
    OSInfo           string   `json:"os_info"`
    RecentCommands   []string `json:"recent_commands"`
    UserLevel        string   `json:"user_level"` // beginner/intermediate/advanced
}

// AICommandResponse AI 命令生成响应
type AICommandResponse struct {
    Command         string            `json:"command"`          // 生成的命令
    Explanation     string            `json:"explanation"`      // 命令解释
    SafetyWarning   *AISafetyWarning  `json:"safety_warning"`   // 安全警告（可选）
    Alternatives    []string          `json:"alternatives"`     // 替代命令
    Confidence      float64           `json:"confidence"`       // 置信度 0-1
}

// AISafetyWarning 安全警告
type AISafetyWarning struct {
    Level       string `json:"level"`        // high/medium/low
    Message     string `json:"message"`      // 警告信息
    RiskDetail  string `json:"risk_detail"`  // 风险详情
    Suggestion  string `json:"suggestion"`   // 安全建议
}
```

#### 4.2.2 Genkit Flow 定义

```go
// 命令生成 Flow
commandGenFlow := genkit.DefineFlow(g, "commandGenerator",
    func(ctx context.Context, input *AICommandRequest) (*AICommandResponse, error) {
        // 构建 prompt
        prompt := buildCommandPrompt(input)

        // 调用 AI 生成命令
        resp, err := genkit.GenerateData[AICommandResponse](ctx, g,
            ai.WithPrompt(prompt),
            ai.WithModelName(modelName),
        )

        return resp, err
    })
```

#### 4.2.3 Prompt 模板

```
你是一个专业的 Linux 命令助手。你的任务是根据用户的自然语言描述，生成准确、安全、高效的 SSH 命令。

## 规则：
1. 只生成可以在 SSH 终端中执行的命令
2. 优先使用常见、标准的命令格式
3. 如果命令有潜在风险，必须在 safety_warning 中标明
4. 提供详细的命令解释
5. 如果有多重实现方式，提供最优解并在 alternatives 中列出其他方案

## 用户当前环境：
- 工作目录: {{.CurrentDirectory}}
- 操作系统: {{.OSInfo}}
- 用户水平: {{.UserLevel}}

## 用户最近执行的命令：
{{range .RecentCommands}}- {{.}}
{{end}}

## 用户请求：
{{.NaturalLanguage}}

请以 JSON 格式返回结果，包含以下字段：
- command: 生成的命令字符串
- explanation: 命令的详细解释
- safety_warning: 如果有风险，包含 level, message, risk_detail, suggestion
- alternatives: 其他可行的命令选项
- confidence: 你对这个命令准确性的置信度 (0-1)
```

### 4.3 前端组件设计

#### 4.3.1 两种 AI 使用模式

Vexo 支持两种 AI 使用模式，用户可根据场景选择：

**模式 A：全局 AI 助手（左侧菜单入口）**

- 入口：左侧菜单栏 AI 图标
- 范围：全局适用，所有标签页共享
- 面板位置：右侧弹出，覆盖标签页区域
- 适用场景：通用问答、命令生成、学习咨询、不针对特定终端的操作

**模式 B：终端内 AI 助手（终端右键/工具栏入口）**

- 入口：终端右键菜单或终端工具栏
- 范围：仅针对当前 SSH 连接
- 面板位置：终端右侧，与终端并排
- 适用场景：针对当前终端上下文的命令生成、命令解释、终端操作辅助

两种模式共享同一个 AI 服务后端，但交互方式和上下文不同。

#### 4.3.2 AI 助手面板组件（AIPanel）

两种模式共用同一个底层组件 `AIPanel`，通过 `mode` 属性区分行为：

```tsx
// frontend/src/components/AIPanel.tsx
interface AIPanelProps {
  mode: "global" | "terminal"; // 全局模式或终端内模式
  onClose: () => void;
  sessionID?: string; // 终端模式必填，全局模式可选
  initialCommand?: string; // 可选：从右键菜单传入的选中命令
  initialMode?: "generate" | "explain"; // 模式：生成或解释
}

// 组件结构：
// ┌─────────────────────────────────────────────────┐
// │ AI 助手                                    [×]  │
// ├─────────────────────────────────────────────────┤
// │ 目标会话：[全部会话 ▼]（仅全局模式显示）        │
// ├─────────────────────────────────────────────────┤
// │                                                 │
// │  [AI 头像] 你好！我是 Vexo AI 助手，            │
// │           有什么可以帮你的？                     │
// │                                                 │
// │           [用户] 如何查找大文件？                │
// │  [AI 头像] 你可以使用以下命令：                  │
// │           find / -size +100M                    │
// │           [运行] [复制]                         │
// │                                                 │
// │           [用户] 解释一下这个命令                │
// │  [AI 头像] 这个命令的含义是...                   │
// │                                                 │
// ├─────────────────────────────────────────────────┤
// │ [输入框：输入你的问题...]                 [发送] │
// └─────────────────────────────────────────────────┘
```

**全局模式与终端内模式的区别**：

| 特性           | 全局模式 (`mode="global"`)   | 终端内模式 (`mode="terminal"`) |
| -------------- | ---------------------------- | ------------------------------ |
| 入口           | 左侧菜单栏 AI 图标           | 终端右键菜单/工具栏            |
| 面板位置       | 右侧覆盖标签页               | 终端右侧并排                   |
| 目标会话选择器 | 显示（可选择全部或指定会话） | 隐藏（自动使用当前会话）       |
| 适用范围       | 所有标签页                   | 仅当前 SSH 连接                |
| 对话历史       | 持久化保存                   | 会话级别                       |

#### 4.3.3 命令展示组件

命令结果以纯文本形式展示，不需要语法高亮：

```
生成结果：
┌─────────────────────────────────────┐
│ $ find / -size +100M -mtime -7     │
│ [复制] [运行] [编辑]                │
└─────────────────────────────────────┘

命令解释：
- find /          : 从根目录开始搜索
- -size +100M     : 查找大于100MB的文件
- -mtime -7       : 修改时间在7天以内
- 2>/dev/null     : 忽略权限错误

⚠️ 安全警告
此命令将搜索整个文件系统，可能需要较长时间
建议指定搜索路径以提高效率
```

#### 4.3.4 集成位置

**全局 AI 助手入口**：

1. **左侧菜单栏**：在 Header 组件中添加 AI 助手图标（SmartToyIcon）
2. 点击后在 SSHTabs 区域右侧滑出 AIChatPanel

**终端内 AI 助手入口**：

1. **终端右键菜单**：在 TerminalContextMenu 中添加"AI 解释此命令"和"AI 生成命令"选项
2. **终端工具栏**（可选）：在终端标签页的工具栏添加 AI 助手图标按钮

#### 4.3.5 两种模式的对比

| 特性       | 全局 AI 助手   | 终端内 AI 助手      |
| ---------- | -------------- | ------------------- |
| 入口位置   | 左侧菜单栏     | 终端右键菜单/工具栏 |
| 面板位置   | 右侧覆盖标签页 | 终端右侧并排        |
| 适用范围   | 所有标签页     | 仅当前 SSH 连接     |
| 对话历史   | 持久化保存     | 会话级别            |
| 上下文感知 | 可选择目标会话 | 自动获取当前会话    |
| 主要用途   | 通用问答、学习 | 命令生成、命令解释  |
| 切换标签页 | 面板保持打开   | 面板随终端关闭      |
| 运行命令   | 需选择目标会话 | 直接运行到当前终端  |

#### 4.3.6 全局 AI 助手 - 目标会话选择

全局 AI 助手支持选择命令发送的目标：

```
目标会话选择器：
┌─────────────────────────────────────────────┐
│ 目标会话：[ 全部会话 ▼ ]                    │
│                                             │
│ 下拉选项：                                  │
│ • 全部会话（默认）                           │
│ • ──────────────                            │
│ • user1@192.168.1.100:22 (当前活跃)         │
│ • user2@10.0.0.50:22                        │
│ • root@172.16.0.1:22                        │
│ • ...                                       │
└─────────────────────────────────────────────┘

行为说明：
- 选择"全部会话"：命令通过 `CommandService.SendCommand` 发送到所有活跃 SSH 会话执行
- 选择特定会话：命令仅发送到选中会话执行
- 切换标签页不影响目标会话选择
- 如果选中的会话被关闭，自动切换回"全部会话"
- 如果没有任何活跃会话，"运行"按钮置灰不可用
- 命令发送后直接执行（通过 SSH stdin 写入并追加 `\n`），无需用户额外确认

**会话列表数据源**：
- 会话列表从后端 `SSHService` 获取，通过 Wails 绑定方法暴露给前端
- 后端维护活跃 SSH 会话列表，包含会话 ID、连接信息（user@host:port）
- 前端调用 `AIService.GetActiveSessions()` 获取列表
```

#### 4.3.7 布局集成方案

全局 AI 助手集成到主应用布局：

```
──────────────────────────────────────────────────────────────┐
│                        App.tsx 布局                          │
│                                                              │
│  ┌──────┐  ┌──────────────────────────────────────────────┐ │
│  │      │  │                  SSHTabs                     │ │
│  │Header│  │  ┌────────────────────────────────────────  │ │
│  │(42px)│  │  │  Tab Bar                               │  │ │
│  │      │  │  ├────────────────────────────────────────┤  │ │
│  │    │  │  │                                        │  │ │
│  │ AI   │  │  │         终端内容区域                    │  │ │
│  │      │  │  │                                        │  │ │
│  │      │  │  └────────────────────────────────────────┘  │ │
│  ──────┘  └──────────────────────────────────────────────┘ │
│                                                              │
│  点击 AI 图标后：                                            │
│  ┌──────┐  ┌──────────────────────────┬───────────────────┐ │
│  │Header│  │       SSHTabs            │   AIChatPanel     │ │
│  │      │  │  (宽度缩减)              │   (450px)         │ │
│  │      │  │                          │                   │ │
│  │      │  │    终端内容区域           │   AI 对话界面     │ │
│  │      │  │                          │                   │ │
│  │      │  │                          │   [×] 关闭        │ │
│  └──────┘  └──────────────────────────┴───────────────────┘ │
└──────────────────────────────────────────────────────────────┘
```

### 4.4 数据流

```
用户输入自然语言
    ↓
前端调用 AIService.GenerateCommand()
    ↓
后端收集上下文信息（当前目录、OS、历史命令）
    ↓
Genkit Flow 调用 AI 模型
    ↓
AI 返回结构化响应（命令、解释、警告）
    ↓
后端进行安全分析（二次校验）
    ↓
前端展示结果
    ↓
用户选择：发送/复制/编辑
```

### 4.5 前端完整交互流程

#### 4.5.1 场景一：终端内使用 AI 命令助手

**入口**：用户在 SSH 终端界面，点击工具栏的 AI 图标按钮

**完整交互流程**：

```
步骤 1: 打开 AI 面板
┌─────────────────────────────────────────────────┐
│ 用户操作：点击终端工具栏的 AI 图标                │
│                                                  │
│ 系统响应：                                        │
│ 1. 在终端右侧滑出 AICommandPanel 面板             │
│ 2. 面板宽度占终端区域的 40%（可拖动调整）          │
│ 3. 显示加载状态："正在初始化 AI 助手..."          │
│ 4. 后端检查 AI 服务状态和配置                     │
│                                                  │
│ 前端状态变化：                                    │
│ - aiPanelOpen: true                              │
│ - aiPanelLoading: true                           │
│ - aiPanelError: null                             │
└─────────────────────────────────────────────────┘

步骤 2: AI 面板就绪
┌─────────────────────────────────────────────────┐
│ 系统响应：                                        │
│ 1. AI 服务初始化完成                             │
│ 2. 显示输入框和提示文本                           │
│ 3. 提示文本："描述你想要的操作，AI 将生成命令"    │
│ 4. 输入框自动获得焦点                             │
│                                                  │
│ 界面元素：                                        │
│ ┌─────────────────────────────────────────────┐ │
│ │ 🔍 [输入框]                          [发送] │ │
│ │                                             │ │
│ │ 💡 示例：                                    │ │
│ │ • 查找大于100MB的文件                        │ │
│ │ • 查看端口占用                               │ │
│ │ • 压缩日志文件                               │ │
│ └─────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────┘

步骤 3: 用户输入并发送
┌─────────────────────────────────────────────────┐
│ 用户操作：输入"查找最近7天修改的大于100MB的文件"   │
│ 用户操作：点击"发送"按钮或按 Enter 键             │
│                                                  │
│ 前端状态变化：                                    │
│ - aiPanelLoading: true                           │
│ - aiPanelInput: "查找最近7天..."                 │
│ - aiPanelResult: null                            │
│                                                  │
│ 系统响应：                                        │
│ 1. 显示加载动画（骨架屏或旋转图标）               │
│ 2. 禁用发送按钮，防止重复提交                     │
│ 3. 调用 AIService.GenerateCommand()               │
│ 4. 传递参数：                                     │
│    - natural_language: 用户输入                   │
│    - session_id: 当前终端会话 ID                  │
│    - context: {                                   │
│        current_directory: 自动获取                │
│        os_info: 自动获取                          │
│        recent_commands: 最近5条命令               │
│      }                                            │
└─────────────────────────────────────────────────┘

步骤 4: 显示 AI 生成结果
┌─────────────────────────────────────────────────┐
│ 系统响应（成功）：                                │
│ 1. 隐藏加载动画                                   │
│ 2. 显示生成结果区域                               │
│ 3. 高亮显示命令（代码块样式）                     │
│ 4. 显示命令解释                                   │
│ 5. 如有安全警告，显示警告卡片                     │
│ 6. 显示操作按钮                                   │
│                                                  │
│ 界面展示：                                        │
│ ┌─────────────────────────────────────────────┐ │
│ │ ✅ 生成成功                                  │ │
│ │                                             │ │
│ │ 📋 命令：                                    │ │
│ │ ┌─────────────────────────────────────────┐ │ │
│ │ │ $ find / -size +100M -mtime -7 2>/dev   │ │ │
│ │ │ [📋复制] [▶️运行] [✏️编辑]              │ │ │
│ │ └─────────────────────────────────────────┘ │ │
│ │                                             │ │
│ │ 📖 命令解释：                                │ │
│ │ • find /          : 从根目录开始搜索         │ │
│ │ • -size +100M     : 查找大于100MB的文件      │ │
│ │ • -mtime -7       : 修改时间在7天以内        │ │
│ │ • 2>/dev/null     : 忽略权限错误             │ │
│ │                                             │ │
│ │ ⚠️ 安全提示：                                │ │
│ │ 此命令将搜索整个文件系统，可能需要较长时间    │ │
│ │ 建议指定搜索路径，如 find /home ...          │ │
│ │                                             │ │
│ │ 💡 替代方案：                                │ │
│ │ • find /home -size +100M -mtime -7          │ │
│ │ • find . -maxdepth 3 -size +100M -mtime -7  │ │
│ └─────────────────────────────────────────────┘ │
│                                                  │
│ 前端状态：                                        │
│ - aiPanelLoading: false                          │
│ - aiPanelResult: { command, explanation, ... }   │
│ - aiPanelError: null                             │
└─────────────────────────────────────────────────┘

步骤 5: 用户操作结果
┌─────────────────────────────────────────────────┐
│ 选项 A - 复制到剪贴板：                           │
│ 1. 用户点击"复制"按钮                             │
│ 2. 命令文本复制到系统剪贴板                       │
│ 3. 显示成功提示："已复制到剪贴板"                 │
│ 4. 按钮短暂显示"✓ 已复制"反馈                     │
│                                                  │
│ 选项 B - 运行命令：                               │
│ 1. 用户点击"运行"按钮                             │
│ 2. 调用 `CommandService.SendCommand` 发送到目标会话执行 │
│ 3. 显示成功提示："命令已发送到终端执行"           │
│ 4. AI 面板可选择保持打开或自动关闭                │
│ 5. 终端获得焦点                                   │
│                                                  │
│ 选项 C - 编辑后运行：                             │
│ 1. 用户点击"编辑"按钮                             │
│ 2. 命令变为可编辑状态（文本框）                   │
│ 3. 用户修改命令                                   │
│ 4. 显示"保存并运行"和"取消"按钮                   │
│ 5. 保存后，调用 `CommandService.SendCommand` 执行 │
│                                                  │
│ 选项 D - 继续提问：                               │
│ 1. 用户在输入框输入新的问题                       │
│ 2. 重复步骤 3-4                                   │
│ 3. 历史记录保留在面板中（可滚动查看）             │
└─────────────────────────────────────────────────┘

步骤 6: 关闭 AI 面板
┌─────────────────────────────────────────────────┐
│ 用户操作：点击 AI 面板右上角的 × 按钮             │
│                                                  │
│ 系统响应：                                        │
│ 1. AI 面板向右滑出关闭                            │
│ 2. 终端区域恢复全宽                               │
│ 3. 终端获得焦点                                   │
│ 4. 前端状态重置：                                 │
│    - aiPanelOpen: false                          │
│    - aiPanelResult: null                         │
│    - aiPanelInput: ""                            │
└─────────────────────────────────────────────────┘
```

#### 4.5.2 场景二：终端右键菜单 - AI 解释选中命令

**入口**：用户在终端中选中一段命令文本，右键点击

**完整交互流程**：

```
步骤 1: 选中命令并右键
┌─────────────────────────────────────────────────┐
│ 用户操作：                                       │
│ 1. 在终端中用鼠标选中一段命令文本                 │
│    例如：`find / -name "*.log" -delete`          │
│ 2. 右键点击选中的文本                             │
│                                                  │
│ 系统响应：                                        │
│ 1. 显示终端右键上下文菜单                         │
│ 2. 菜单项包含：                                   │
│    - 复制                                         │
│    - 粘贴                                         │
│    - 清屏                                         │
│    - ─────────────                                │
│    - 🤖 AI 解释此命令（新增）                     │
│    - 🤖 AI 优化此命令（新增）                     │
└─────────────────────────────────────────────────┘

步骤 2: 点击"AI 解释此命令"
┌─────────────────────────────────────────────────┐
│ 用户操作：点击"AI 解释此命令"菜单项               │
│                                                  │
│ 系统响应：                                        │
│ 1. 获取终端中选中的文本                           │
│ 2. 在终端下方展开解释面板（高度约 200px）         │
│ 3. 显示加载状态："正在分析命令..."               │
│ 4. 调用 AIService.ExplainCommand()                │
│ 5. 传递参数：                                     │
│    - command: 选中的命令文本                      │
│    - session_id: 当前终端会话 ID                  │
└─────────────────────────────────────────────────┘

步骤 3: 显示解释结果
┌─────────────────────────────────────────────────┐
│ 界面展示：                                        │
│ ┌─────────────────────────────────────────────┐ │
│ │ 🤖 AI 命令解释                               │ │
│ │                                             │ │
│ │ 原命令：                                     │ │
│ │ $ find / -name "*.log" -delete              │ │
│ │                                             │ │
│ │ 📖 解释：                                    │ │
│ │ • find /          : 从根目录开始查找         │ │
│ │ • -name "*.log"   : 匹配所有 .log 结尾的文件 │ │
│ │ • -delete         : 直接删除匹配的文件       │ │
│ │                                             │ │
│ │ ⚠️ 安全警告（高风险）：                       │ │
│ │ 此命令将删除系统中所有 .log 文件，且不可恢复  │ │
│ │ 建议：                                       │ │
│ │ 1. 先用 -print 预览：find / -name "*.log"    │ │
│ │ 2. 确认无误后再执行删除                       │ │
│ │ 3. 或先备份重要日志                           │ │
│ │                                             │ │
│ │ [关闭] [运行] [复制优化版本]                  │ │
│ └─────────────────────────────────────────────┘ │
│                                                  │
│ 交互选项：                                        │
│ • 点击"关闭"：收起解释面板                        │
│ • 点击"运行"：调用 `CommandService.SendCommand` 执行命令 │
│ • 点击"复制优化版本"：复制 AI 建议的安全版本      │
└─────────────────────────────────────────────────┘
```

#### 4.5.3 场景三：AI 未配置时提示

**入口**：用户首次点击 AI 图标（全局 AI 助手或终端内 AI 助手），且 AI 功能未启用

**完整交互流程**：

```
步骤 1: 检测到未配置
┌─────────────────────────────────────────────────┐
│ 系统响应：                                        │
│ 1. 弹出提示对话框（Modal Dialog）                 │
│ 2. 对话框内容：                                   │
│                                                  │
│ ┌─────────────────────────────────────────────┐ │
│ │ ⚙️ AI 功能未配置                             │ │
│ │                                             │ │
│ │ 请先在设置页面配置 AI 模型服务后再使用。       │ │
│ │                                             │ │
│ │ 支持的模型提供商：                            │ │
│ │ • Ollama（推荐，本地运行，隐私保护）          │ │
│ │ • OpenAI（需要 API Key）                     │ │
│ │ • Anthropic（需要 API Key）                  │ │
│ │ • Google AI（需要 API Key）                  │ │
│ │                                             │ │
│ │ [稍后设置] [前往设置]                         │ │
│ └─────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────┘

步骤 2: 用户选择
┌─────────────────────────────────────────────────┐
│ 选项 A - 点击"前往设置"：                         │
│ 1. 关闭提示对话框                                 │
│ 2. 打开设置窗口，自动切换到"AI"标签页             │
│ 3. 用户在设置页面完成配置后，关闭设置窗口         │
│ 4. 再次点击 AI 图标即可正常使用                   │
│                                                  │
│ 选项 B - 点击"稍后设置"：                         │
│ 1. 关闭提示对话框                                 │
│ 2. 不执行任何操作                                 │
│ 3. 下次点击 AI 图标时仍会弹出提示                 │
└─────────────────────────────────────────────────┘
```

#### 4.5.4 错误处理流程

```
错误场景 1: AI 服务不可用
┌─────────────────────────────────────────────────┐
│ 触发条件：                                        │
│ • AI 服务未启动                                   │
│ • 模型未安装                                      │
│ • API Key 无效                                    │
│                                                  │
│ 界面展示：                                        │
│ ┌─────────────────────────────────────────────┐ │
│ │ ❌ AI 服务不可用                              │ │
│ │                                             │ │
│ │ 错误信息：模型 codellama 未安装               │ │
│ │                                             │ │
│ │ 建议操作：                                    │ │
│ │ 1. 运行 ollama pull codellama 安装模型        │ │
│ │ 2. 或在设置中更换其他模型                     │ │
│ │                                             │ │
│ │ [打开设置] [重试]                             │ │
│ └─────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────┘

错误场景 2: 网络请求超时
┌─────────────────────────────────────────────────┐
│ 触发条件：                                        │
│ • 云端模型 API 响应超时（> 30 秒）               │
│ • 网络连接不稳定                                  │
│                                                  │
│ 界面展示：                                        │
│ ┌─────────────────────────────────────────────┐ │
│ │ ⏱️ 请求超时                                  │ │
│ │                                             │ │
│ │ AI 响应超时，请检查网络连接或稍后重试         │ │
│ │                                             │ │
│ │ [重试] [切换到本地模型]                       │ │
│ └─────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────┘

错误场景 3: AI 生成结果解析失败
┌─────────────────────────────────────────────────┐
│ 触发条件：                                        │
│ • AI 返回格式不符合预期                           │
│ • JSON 解析失败                                   │
│                                                  │
│ 界面展示：                                        │
│ ┌─────────────────────────────────────────────┐ │
│ │ ⚠️ 解析失败                                  │ │
│ │                                             │ │
│ │ AI 返回结果格式异常，已尝试重新生成           │ │
│ │                                             │ │
│ │ [重试] [查看原始响应]                         │ │
│ └─────────────────────────────────────────────┘ │
│                                                  │
│ 系统行为：                                        │
│ 1. 自动重试 1 次                                  │
│ 2. 如果仍然失败，显示错误提示                     │
│ 3. 提供"查看原始响应"用于调试                     │
└─────────────────────────────────────────────────┘
```

#### 4.5.5 状态管理设计

```typescript
// frontend/src/stores/aiCommand.ts
import { create } from "zustand";

interface AICommandState {
  // 面板状态
  aiPanelOpen: boolean;
  aiPanelLoading: boolean;
  aiPanelError: string | null;

  // 输入状态
  aiPanelInput: string;
  aiPanelInputHistory: string[]; // 输入历史

  // 结果状态
  aiPanelResult: AICommandResponse | null;
  aiPanelResultHistory: AICommandResponse[]; // 对话历史

  // 配置状态
  aiConfig: AIConfig | null;
  aiConfigLoading: boolean;

  // 操作方法
  openPanel: () => void;
  closePanel: () => void;
  setInput: (input: string) => void;
  generateCommand: () => Promise<void>;
  sendToTerminal: (command: string) => void;
  copyToClipboard: (text: string) => Promise<void>;
  loadConfig: () => Promise<void>;
  saveConfig: (config: AIConfig) => Promise<void>;
  resetState: () => void;
}

const useAICommandStore = create<AICommandState>((set, get) => ({
  aiPanelOpen: false,
  aiPanelLoading: false,
  aiPanelError: null,
  aiPanelInput: "",
  aiPanelInputHistory: [],
  aiPanelResult: null,
  aiPanelResultHistory: [],
  aiConfig: null,
  aiConfigLoading: false,

  openPanel: () => set({ aiPanelOpen: true }),
  closePanel: () => set({ aiPanelOpen: false, aiPanelResult: null }),
  setInput: (input) => set({ aiPanelInput: input }),

  generateCommand: async () => {
    const { aiPanelInput, aiConfig } = get();
    if (!aiPanelInput.trim()) return;

    set({ aiPanelLoading: true, aiPanelError: null });

    try {
      const result = await AIService.GenerateCommand({
        natural_language: aiPanelInput,
        context: await getTerminalContext(),
      });

      set({
        aiPanelLoading: false,
        aiPanelResult: result,
        aiPanelResultHistory: [...get().aiPanelResultHistory, result],
        aiPanelInputHistory: [...get().aiPanelInputHistory, aiPanelInput],
      });
    } catch (error) {
      set({
        aiPanelLoading: false,
        aiPanelError: parseCallServiceError(error),
      });
    }
  },

  runCommand: (command: string, sessionIDs: string[]) => {
    // 调用 CommandService.SendCommand 发送到 SSH 会话执行
    CommandService.SendCommand({
      command: command + "\n",
      session_ids: sessionIDs,
    });
  },

  copyToClipboard: async (text: string) => {
    await navigator.clipboard.writeText(text);
  },

  loadConfig: async () => {
    set({ aiConfigLoading: true });
    try {
      const config = await AIService.GetConfig();
      set({ aiConfig: config, aiConfigLoading: false });
    } catch (error) {
      set({ aiConfigLoading: false });
    }
  },

  saveConfig: async (config: AIConfig) => {
    set({ aiConfigLoading: true });
    try {
      await AIService.SaveConfig(config);
      set({ aiConfig: config, aiConfigLoading: false });
    } catch (error) {
      set({ aiConfigLoading: false });
    }
  },

  resetState: () =>
    set({
      aiPanelInput: "",
      aiPanelResult: null,
      aiPanelError: null,
    }),
}));
```

#### 4.5.6 组件集成方式

```tsx
// frontend/src/components/Terminal.tsx 修改示例
import AICommandPanel from './AICommandPanel';
import { useAICommandStore } from '../stores/aiCommand';

function Terminal(props: { readonly linkID: string }) {
  const aiPanelOpen = useAICommandStore((state) => state.aiPanelOpen);

  return (
    <Box sx={{ display: 'flex', width: '100%', height: '100%' }}>
      {/* 终端区域 */}
      <Box
        ref={termRef}
        sx={{
          flex: aiPanelOpen ? '0 0 60%' : 1, // AI 面板打开时终端占 60%
          transition: 'flex 0.3s ease',
        }}
      />

      {/* AI 命令面板 */}
      {aiPanelOpen && (
        <AICommandPanel
          sessionID={props.linkID}
          onCommandSend={(command) => {
            // 发送命令到终端
            const term = terminalInstances.get(props.linkID);
            term?.paste(command);
            term?.focus();
          }}
          onClose={() => useAICommandStore.getState().closePanel()}
        />
      )}
    </Box>
  );
}

// frontend/src/components/TerminalContextMenu.tsx 修改示例
import SmartToyIcon from '@mui/icons-material/SmartToy';

function TerminalContextMenu({ contextMenu, onClose, linkID }) {
  const handleAIExplain = async () => {
    const term = terminalInstances.get(linkID);
    const selectedText = term?.getSelection();
    if (selectedText) {
      // 打开 AI 解释面板
      useAICommandStore.getState().openPanel();
      useAICommandStore.getState().explainCommand(selectedText);
    }
    onClose();
  };

  return (
    <Menu ...>
      <MenuItem onClick={handleCopy}>...</MenuItem>
      <MenuItem onClick={handlePaste}>...</MenuItem>
      <MenuItem onClick={handleClear}>...</MenuItem>
      <Divider />
      <MenuItem onClick={handleAIExplain}>
        <ListItemIcon>
          <SmartToyIcon fontSize="small" />
        </ListItemIcon>
        <ListItemText>AI 解释此命令</ListItemText>
      </MenuItem>
    </Menu>
  );
}
```

#### 4.5.7 关键设计决策与澄清

**1. 对话历史存储策略**

| 模式           | 存储位置                                | 保留策略                 | 存储上限                                  |
| -------------- | --------------------------------------- | ------------------------ | ----------------------------------------- |
| 全局 AI 助手   | 本地 SQLite 数据库（internal/database） | 永久保留，用户可手动清除 | 最多 100 条对话记录，超出后自动清理最早的 |
| 终端内 AI 助手 | 内存（Zustand store）                   | 终端关闭即清除           | 无上限（会话期间）                        |

**2. 上下文信息获取方式**

- **当前工作目录**：通过解析 SSH 会话输出（如 `pwd` 命令回显）获取
- **操作系统信息**：通过解析 `uname -a` 输出获取（会话建立时缓存，不重复请求）
- **最近命令**：通过解析 SSH 会话输出获取最近 5 条命令
- **获取上下文不会阻塞用户操作**：上下文信息在 AI 请求发起时异步收集，超时自动降级为无上下文模式
- **未连接会话**：未连接的 SSH 标签页无法使用 AI 功能

**3. 安全警告机制**

- **AI 模型判断**：由 AI 模型在生成命令时根据 prompt 规则判断，返回 `safety_warning` 字段
- **警告等级**：`high`（高风险）/`medium`（中风险）/`low`（低风险）
- **展示方式**：前端根据 `safety_warning.level` 显示对应颜色的警告卡片
- **Prompt 规则**：在 system prompt 中明确要求 AI 识别危险操作（如删除、覆盖、权限修改等）并返回警告信息

**4. 全局 AI 助手完整交互流程**

````
步骤 1: 打开全局 AI 面板
┌─────────────────────────────────────────────────┐
│ 用户操作：点击左侧菜单栏 AI 图标                  │
│                                                  │
│ 系统响应：                                        │
│ 1. SSHTabs 区域宽度缩减，右侧滑出 AIChatPanel     │
│ 2. 面板宽度固定 450px                             │
│ 3. 加载对话历史（从本地数据库）                   │
│ 4. 初始化目标会话选择器（默认"全部会话"）         │
└─────────────────────────────────────────────────┘

步骤 2: 用户输入问题
┌─────────────────────────────────────────────────┐
│ 用户操作：在输入框输入自然语言问题                │
│ 例如："如何查找大文件？"                          │
│                                                  │
│ 系统响应：                                        │
│ 1. 调用 AIService.Chat()（通用对话接口）          │
│ 2. 传递参数：                                     │
│    - message: 用户输入                            │
│    - target_session_id: 选中的目标会话（可选）    │
│    - history: 最近 10 条对话历史                  │
│ 3. 流式输出响应（打字机效果）                     │
└─────────────────────────────────────────────────┘

步骤 3: 显示 AI 回复
┌─────────────────────────────────────────────────┐
│ 界面展示：                                        │
│ ┌─────────────────────────────────────────────┐ │
│ │ 🤖 AI 助手                                  │ │
│ │                                             │ │
│ │ [用户] 如何查找大文件？                      │ │
│ │                                             │ │
│ │ [AI] 你可以使用以下命令：                    │ │
│ │                                             │ │
│ │ ```bash                                     │ │
│ │ find / -size +100M 2>/dev/null              │ │
│ │ ```                                         │ │
│ │                                             │ │
│ │ 这个命令会查找大于 100MB 的文件...           │ │
│ │                                             │ │
│ │ [📋 复制] [▶️ 运行]                          │ │
│ └─────────────────────────────────────────────┘ │
│                                                  │
│ 注意：全局模式的回复是通用对话格式，              │
│ 不一定包含结构化命令信息                          │
└─────────────────────────────────────────────────┘

步骤 4: 运行命令
┌─────────────────────────────────────────────────┐
│ 用户操作：点击"运行"                              │
│                                                  │
│ 系统响应：                                        │
│ 1. 从 AI 回复中提取代码块中的命令                 │
│ 2. 根据目标会话选择器的值：                       │
│    - 全部会话：调用 `CommandService.SendCommand` 发送到所有活跃会话执行 │
│    - 特定会话：调用 `CommandService.SendCommand` 发送到指定会话执行 │
│ 3. 命令通过 SSH stdin 直接执行（追加 `\n`）       │
│ 4. 目标终端获得焦点                               │
└─────────────────────────────────────────────────┘
````

**5. 流式输出实现方案**

- 使用 Genkit 的流式生成 API（`genkit.GenerateStream`）
- 后端通过 Wails Events 将流式数据片段推送到前端
- 前端使用打字机效果逐字显示，提升用户体验
- 流式输出仅适用于通用对话模式（`Chat` 接口）
- 结构化命令生成（`GenerateCommand` 接口）使用一次性返回，保证 JSON 完整性

**流式输出数据流**：

```
用户输入问题
    ↓
前端调用 AIService.ChatStream()（返回 EventStream）
    ↓
后端调用 genkit.GenerateStream()
    ↓
AI 模型逐 token 返回
    ↓
后端通过 Wails Events.Emit("aiStreamChunk", chunk) 推送
    ↓
前端监听 Events.On("aiStreamChunk") 逐字渲染
    ↓
流结束，后端发送 Events.Emit("aiStreamDone")
```

**6. API Key 存储**

- AI 配置保存在 `config.toml` 配置文件中
- API Key 使用 `internal/secret` 包加密（Argon2 + AES-GCM）后存储在配置文件
- 加密后的 API Key 作为 `api_key_encrypted` 字段保存
- 仅在需要时解密读取，不在内存中长期保留

### 4.7 响应式设计

| 屏幕宽度        | AI 面板宽度   | 终端宽度 | 备注           |
| --------------- | ------------- | -------- | -------------- |
| < 1200px        | 50%           | 50%      | 小屏幕平分     |
| 1200px - 1600px | 40%           | 60%      | 中等屏幕       |
| > 1600px        | 450px（固定） | 剩余     | 大屏幕固定宽度 |

---

## 5. 模型配置

### 5.1 支持的模型提供商

| 提供商       | 模型                      | 特点                 | 适用场景   |
| ------------ | ------------------------- | -------------------- | ---------- |
| Ollama       | llama3/codellama          | 本地运行，隐私保护   | 默认推荐   |
| OpenAI       | gpt-4o/gpt-4o-mini        | 高质量，需要 API Key | 追求准确性 |
| Anthropic    | claude-sonnet-4           | 高质量，需要 API Key | 追求准确性 |
| Google AI    | gemini-2.5-flash          | 性价比高             | 平衡方案   |
| Azure OpenAI | gpt-4o                    | 企业级部署           | 企业用户   |
| Mistral      | mistral-large             | 高性能开源           | 开源偏好   |
| Cohere       | command-r                 | 专为命令优化         | 特定场景   |

### 5.2 配置项

在现有配置系统中添加 AI 配置：

```go
type AIConfig struct {
    Enabled          bool   `json:"enabled"`            // 是否启用 AI 功能
    Provider         string `json:"provider"`           // 模型提供商
    Model            string `json:"model"`              // 模型名称
    APIKey           string `json:"api_key"`            // API 密钥（云端模型需要）
    OllamaEndpoint   string `json:"ollama_endpoint"`    // Ollama 端点地址
    Temperature      float64 `json:"temperature"`       // 温度参数
    MaxTokens        int    `json:"max_tokens"`         // 最大 token 数
    SafetyCheckLevel string `json:"safety_check_level"` // 安全检查等级
}
```

### 5.3 默认配置

```json
{
  "ai": {
    "enabled": false,
    "provider": "ollama",
    "model": "codellama:latest",
    "ollama_endpoint": "http://localhost:11434",
    "temperature": 0.1,
    "max_tokens": 1000,
    "safety_check_level": "medium"
  }
}
```

### 5.4 设置页面集成

在现有设置页面（`frontend/src/pages/Setting.tsx`）中新增"AI"标签页，与通用、终端、同步、关于并列。

#### 5.4.1 配置结构变更

Go 后端 `Config` 结构体新增 `AI` 字段：

```go
type Config struct {
    General  GeneralConfig   `toml:"general"`
    Terminal TerminalConfig  `toml:"terminal"`
    Sync     sync.SyncConfig `toml:"sync"`
    AI       AIConfig        `toml:"ai"`
}
```

#### 5.4.2 设置页面 UI 设计

```
设置页面左侧导航新增"AI"标签：
┌──────────────┐
│ 设置         │
├──────────────┤
│ 通用         │
│ 终端         │
│ 同步         │
│ AI  ← 新增   │
│ 关于         │
└──────────────┘
```

点击"AI"标签后，右侧显示 AI 配置表单：

```
┌──────────────────────────────────────────────────────────────┐
│ AI 设置                                                      │
│                                                              │
│ ┌──────────────────────────────────────────────────────────┐ │
│ │ 启用 AI 功能  [开关]                                      │ │
│ │                                                          │ │
│ │ 模型提供商  [ Ollama ▼ ]                                 │ │
│ │                                                          │ │
│ │ ┌─ Ollama 配置 ──────────────────────────────────────┐  │ │
│ │ │ 端点地址                                            │  │ │
│ │ │ [http://localhost:11434                    ] [测试] │  │ │
│ │ │                                                     │  │ │
│ │ │ 模型                                                │  │ │
│ │ │ [ codellama:latest ▼ ] [刷新模型列表]               │  │ │
│ │ │                                                     │  │ │
│ │ │ 已安装的模型：                                      │  │ │
│ │ │ • codellama:latest (3.8B)                           │  │ │
│ │ │ • llama3:latest (8B)                                │  │ │
│ │ └─────────────────────────────────────────────────────┘  │ │
│ │                                                          │ │
│ │ ┌─ OpenAI 配置 ──────────────────────────────────────┐  │ │
│ │ │ API Key                                             │  │ │
│ │ │ [sk-xxxxxxxxxxxxxxxxxxxx                   ] [显示] │  │ │
│ │ │                                                     │  │ │
│ │ │ 模型                                                │  │ │
│ │ │ [ gpt-4o-mini ▼ ]                                  │  │ │
│ │ └─────────────────────────────────────────────────────┘  │ │
│ │                                                          │ │
│ │ ┌─ Anthropic 配置 ───────────────────────────────────┐  │ │
│ │ │ API Key                                             │  │ │
│ │ │ [sk-ant-xxxxxxxxxxxxxxxxxxxx               ] [显示] │  │ │
│ │ │                                                     │  │ │
│ │ │ 模型                                                │  │ │
│ │ │ [ claude-sonnet-4-20250514 ▼ ]                     │  │ │
│ │ └─────────────────────────────────────────────────────┘  │ │
│ │                                                          │ │
│ │ ┌─ 高级选项 ─────────────────────────────────────────┐  │ │
│ │ │ Temperature  [0.1 ▼] (0.0-1.0，越低越确定)         │  │ │
│ │ │                                                     │  │ │
│ │ │ Max Tokens   [1000 ▼]                              │  │ │
│ │ │                                                     │  │ │
│ │ │ 安全检查等级  [ 中 ▼ ] (低/中/高)                   │  │ │
│ │ └─────────────────────────────────────────────────────┘  │ │
│ │                                                          │ │
│ │ [恢复默认设置]                                            │ │
│ └──────────────────────────────────────────────────────────┘ │
└──────────────────────────────────────────────────────────────┘
```

#### 5.4.3 交互说明

**1. 启用/禁用开关**

- 关闭时：灰化所有配置项，AI 功能不可用
- 开启时：激活配置项，自动初始化 AI 服务

**2. 模型提供商切换**

- 切换提供商时，自动显示对应提供商的配置区域
- 切换时不丢失已填写的配置（各提供商配置独立保存）

**3. Ollama 端点测试**

- 点击"测试"按钮，调用 `AIService.TestConnection()`
- 成功：显示"连接成功，发现 N 个模型"，绿色提示
- 失败：显示错误原因，红色提示

**4. 模型列表刷新**

- 点击"刷新模型列表"，重新从 Ollama 获取已安装模型
- 下拉框自动更新选项

**5. API Key 输入**

- 默认以掩码显示（`sk-****`）
- 点击"显示"按钮切换明文/掩码
- 保存时通过 `internal/secret` 包加密存储，不写明文到配置文件

**6. 恢复默认设置**

- 弹出确认对话框
- 确认后重置所有 AI 配置为默认值
- 不关闭"启用 AI 功能"开关

#### 5.4.4 前端组件结构

```
frontend/src/pages/Setting.tsx          # 设置页面主入口（新增 ai 标签）
frontend/src/components/AISettings.tsx  # AI 设置表单组件（新建）
frontend/src/stores/aiConfig.ts         # AI 配置状态管理（新建）
```

#### 5.4.5 后端接口

| 接口                               | 说明                       |
| ---------------------------------- | -------------------------- |
| `AIService.GetConfig()`            | 获取当前 AI 配置           |
| `AIService.SaveConfig(config)`     | 保存 AI 配置               |
| `AIService.TestConnection(config)` | 测试模型连接               |
| `AIService.ListModels()`           | 获取可用模型列表（Ollama） |
| `AIService.ResetConfig()`          | 恢复默认配置               |

---

## 6. 非功能需求

### 6.1 性能要求

| 指标             | 要求                                   |
| ---------------- | -------------------------------------- |
| 命令生成响应时间 | < 10 秒（本地模型），< 20 秒（云端模型） |
| 流式输出首字时间 | < 1 秒                                 |
| 内存占用增加     | < 100MB                                |
| 模型加载时间     | < 10 秒（首次）                        |

### 6.2 安全要求

- **API 密钥安全**：API 密钥使用 `internal/secret` 包加密（Argon2 + AES-GCM）后存储在配置文件中，不持久化明文
- **命令安全**：所有 AI 生成的命令必须经过本地安全规则校验
- **隐私保护**：使用本地模型时，数据不离开本机
- **权限控制**：AI 功能默认关闭，用户主动开启

### 6.3 可用性要求

- **离线可用**：使用 Ollama 本地模型时，无需网络连接
- **降级方案**：AI 服务不可用时，优雅降级并提示用户
- **错误处理**：清晰的错误提示和恢复建议

### 6.4 兼容性要求

| 平台    | 支持情况 |
| ------- | -------- |
| Windows | ✓        |
| macOS   | ✓        |
| Linux   | ✓        |

---

## 7. 开发计划

### 7.1 阶段划分

#### Phase 1: 基础架构与配置系统（第1-2周）

**后端**

- [x] 集成 Genkit Go 框架到项目
- [x] 创建 `services/ai_service.go` 基础结构（AIService、请求/响应类型定义）
- [x] 扩展 `Config` 结构体，新增 `AI` 字段
- [x] 实现 AI 配置的读取、保存、重置功能
- [x] 实现 Ollama 本地模型接入（端点连接、模型列表获取）- *模型调用待实现*
- [x] 实现 API Key 安全存储（集成 `internal/secret` 包）
- [ ] 实现 `TestConnection()` 连接测试接口

**前端**

- [x] 设置页面新增"AI"标签页（`Setting.tsx` 修改）
- [x] 创建 `AISettings.tsx` 配置表单组件
- [x] 实现模型提供商切换、Ollama 端点测试、模型列表刷新 - *测试连接待实现*
- [x] 实现 API Key 输入（掩码/明文切换）
- [x] 实现启用/禁用开关、高级选项（Temperature、Max Tokens、安全检查等级）
- [x] 创建 `aiConfig.ts` 状态管理

**验收标准**

- 用户可在设置页面完成 AI 模型配置并保存
- Ollama 本地模型可连接并获取模型列表
- 配置数据正确持久化

#### Phase 2: 核心命令生成与解释（第3-4周）

**后端**

- [x] 实现命令生成 Prompt 模板
- [x] 实现 `commandGenerator` Genkit Flow（结构化输出）- *模型调用待实现*
- [x] 实现 `ExplainCommand()` 命令解释 Flow - *模型调用待实现*
- [ ] 实现上下文收集器（获取当前目录、OS 信息、最近命令）
- [x] 实现 `internal/ai/safety_rules.go` 安全规则库
- [x] 实现后端安全二次校验逻辑
- [x] 实现 `GenerateCommand()` 和 `ExplainCommand()` Wails 绑定方法

**前端**

- [ ] 创建 `aiCommand.ts` 状态管理（Zustand store）
- [ ] 创建 `TerminalAICommandPanel.tsx` 终端内 AI 面板组件
- [ ] 实现终端工具栏 AI 图标入口
- [ ] 实现终端右键菜单"AI 解释此命令"、"AI 优化此命令"选项
- [ ] 实现命令结果展示（代码块、解释、安全警告、替代方案）
- [ ] 实现复制、发送到终端、编辑后发送功能
- [ ] 实现未配置时的弹框提示（Modal Dialog）

**验收标准**

- 用户可通过终端工具栏或右键菜单触发 AI 命令助手
- AI 能生成准确的命令并提供解释和安全警告
- 命令可复制到剪贴板或发送到终端

#### Phase 3: 全局 AI 助手（第5-6周）

**后端**

- [ ] 实现 `Chat()` 通用对话 Flow（流式输出）
- [ ] 实现流式输出支持（`GenerateStream` + Wails Events 推送）
- [ ] 实现对话历史的数据库存储（`internal/database` 新增 AI 对话表）
- [ ] 实现对话历史的查询、删除接口
- [x] 实现目标会话选择相关接口（获取活跃会话列表）

**前端**

- [ ] 左侧菜单栏新增 AI 助手图标（`Header.tsx` 修改）
- [ ] 创建 `AIChatPanel.tsx` 全局 AI 面板组件
- [ ] 实现 SSHTabs 布局集成（面板滑入/滑出动画）
- [ ] 实现目标会话选择器（全部会话/指定会话）
- [ ] 实现流式输出的打字机效果
- [ ] 实现对话历史加载和展示
- [ ] 实现"清除对话历史"功能
- [ ] 实现 AI 回复中代码块的提取和"发送到终端"功能

**验收标准**

- 全局 AI 助手可从左侧菜单打开，支持多轮对话
- 对话历史持久化保存，切换标签页不丢失
- 可选择目标会话发送命令

#### Phase 4: 增强功能与优化（第7-8周）

**后端**

- [ ] 支持 OpenAI 模型接入
- [ ] 支持 Anthropic 模型接入
- [ ] 支持 Google AI 模型接入
- [ ] 实现命令历史学习功能（分析用户常用命令模式）
- [ ] 实现多命令组合生成（管道、重定向等复合命令）
- [ ] 实现相同问题缓存机制（避免重复调用 AI）
- [ ] 实现 Token 用量统计

**前端**

- [ ] AI 图标状态指示（加载中、错误、可用）
- [ ] 实现面板宽度拖动调整
- [ ] 实现响应式布局适配（不同屏幕宽度下的面板宽度）
- [ ] 实现多窗口 AI 助手隔离（每个窗口独立状态）
- [ ] 完善错误处理（服务不可用、超时、解析失败等场景）

**验收标准**

- 支持多种云端模型，用户可自由切换
- 面板宽度可拖动调整，响应式布局正常
- 错误处理完善，用户体验流畅

#### Phase 5: 测试与发布（第9-10周）

- [ ] 单元测试（AIService、安全规则、配置管理）
- [ ] 集成测试（端到端 AI 命令生成流程）
- [ ] 安全审计（API Key 存储、命令安全校验）
- [ ] 性能测试（响应时间、内存占用）
- [ ] 多平台测试（Windows、macOS、Linux）
- [ ] 用户测试与反馈收集
- [ ] 文档完善（用户手册、FAQ）

### 7.2 技术依赖

| 依赖                          | 版本     | 用途             |
| ----------------------------- | -------- | ---------------- |
| github.com/firebase/genkit/go | latest   | Genkit Go SDK    |
| Ollama                        | >= 0.4.0 | 本地 AI 模型运行 |
| @xterm/xterm                  | 现有     | 终端集成         |
| @mui/material                 | v7.3.x   | UI 组件          |
| zustand                       | 现有     | 状态管理         |

### 7.3 开发优先级

| 优先级 | 功能           | 说明                   |
| ------ | -------------- | ---------------------- |
| P0     | AI 配置管理    | 基础功能，必须首先完成 |
| P0     | 命令生成       | 核心功能               |
| P0     | 命令解释       | 核心功能               |
| P0     | 安全警告       | 安全必需               |
| P1     | 终端内 AI 助手 | 主要使用场景           |
| P1     | 全局 AI 助手   | 增强使用场景           |
| P1     | 流式输出       | 用户体验优化           |
| P2     | 对话历史持久化 | 增强功能               |
| P2     | 多模型支持     | 扩展性                 |
| P2     | 命令历史学习   | 个性化功能             |
| P3     | Token 用量统计 | 可选功能               |
| P3     | 面板宽度拖动   | 体验优化               |

---

## 8. 风险与应对

| 风险               | 影响 | 应对措施                        |
| ------------------ | ---- | ------------------------------- |
| AI 生成错误命令    | 高   | 本地安全规则校验 + 用户确认机制 |
| 本地模型资源占用大 | 中   | 提供云端模型选项 + 资源监控     |
| API 密钥泄露       | 高   | 使用系统 Keychain 存储          |
| 网络不稳定         | 中   | 支持离线本地模型                |
| 用户过度依赖 AI    | 低   | 提供命令解释，促进学习          |

---

## 9. 成功指标

| 指标           | 目标值           |
| -------------- | ---------------- |
| 命令生成准确率 | > 90%            |
| 用户满意度     | > 4.0/5.0        |
| 平均响应时间   | < 3 秒           |
| 安全警告准确率 | > 95%            |
| 功能使用率     | > 30% 的活跃用户 |

---

## 10. 未来扩展

- **脚本生成**：根据复杂需求生成完整的 Shell 脚本
- **命令调试**：AI 帮助调试失败的命令，分析错误原因
- **团队协作**：共享 AI 生成的命令模板
- **命令库集成**：将 AI 生成的优质命令保存到命令库

---

## 附录

### A. 术语表

| 术语   | 说明                          |
| ------ | ----------------------------- |
| Genkit | Google 开源的 AI 应用开发框架 |
| Flow   | Genkit 中的类型安全工作流单元 |
| Ollama | 本地运行大语言模型的工具      |
| Token  | AI 模型处理文本的基本单位     |

### B. 参考文档

- [Genkit Go 官方文档](https://genkit.dev/docs/go/overview/)
- [Genkit GitHub](https://github.com/firebase/genkit)
- [Ollama 官方文档](https://ollama.com/docs)
- [Vexo AGENTS.md](./AGENTS.md)
