package ai

import (
	"fmt"
	"strings"

	"github.com/ilaziness/vexo/internal/system"
)

// SSHPromptContext SSH 连接上下文（用于系统提示词）
type SSHPromptContext struct {
	Host string
	Port int
	User string
}

const promptBase = `你是 Vexo 内置 AI 助手。Vexo 是一款 SSH 桌面客户端；用户在侧边栏与你对话，在另一个标签页的终端中操作远程机器。

## 能力边界
- 你可以：解释命令、分析报错、编写脚本/配置示例、给出 Linux/Unix 运维建议
- 你不能：查看终端实时输出、执行命令、访问远程文件系统
- 若系统提示词中已给出远程 OS 信息，优先依据这些信息回答
- 未提供或未采集的信息（具体目录、已安装软件、完整报错）不要编造；先说明假设，或请用户补充

## 回答规范
- 使用与用户相同的语言（用户用中文则中文回答）
- 命令与脚本必须用 Markdown 代码块包裹，便于复制
- 多步骤操作用编号列表；先给简短结论，再给详细步骤
- 涉及 rm、dd、mkfs、fdisk、kill -9、chmod 777、> /dev/ 等破坏性/高风险操作时，必须标注风险，并优先给出更安全的替代（如 ls 预览、--dry-run、先备份）

## 安全
- 禁止请求、输出、猜测密码、私钥、API Key 或 token
- 禁止建议明文存储或传输凭证`

const promptNoSSH = `

## 当前 SSH 上下文
用户当前没有打开的 SSH 连接。请提供通用的 Linux/Unix、SSH 使用与故障排查帮助，不要假设特定远程主机或环境。`

// BuildSystemPrompt 构建系统提示词
func BuildSystemPrompt(ctx *SSHPromptContext, remote *system.RemoteSystemInfo) string {
	if !isValidSSHContext(ctx) {
		return promptBase + promptNoSSH
	}

	var b strings.Builder
	b.WriteString(promptBase)
	b.WriteString("\n\n## 当前 SSH 上下文\n")
	fmt.Fprintf(&b, "- 连接：%s@%s:%d\n\n", ctx.User, ctx.Host, ctx.Port)
	b.WriteString("默认将用户问题理解为与上述远程环境相关；命令与路径按 Linux/Unix 环境给出。")

	if remote != nil && remote.Ready && remote.HasContent() {
		b.WriteString("\n\n## 远程系统信息\n")
		if remote.Hostname != "" {
			fmt.Fprintf(&b, "- 主机名：%s\n", remote.Hostname)
		}
		if remote.OSPretty != "" || remote.OSID != "" || remote.OSVersion != "" {
			fmt.Fprintf(&b, "- 操作系统：%s（%s %s）\n", remote.OSPretty, remote.OSID, remote.OSVersion)
		}
		if remote.Kernel != "" {
			fmt.Fprintf(&b, "- 内核：%s\n", remote.Kernel)
		}
		if remote.Arch != "" {
			fmt.Fprintf(&b, "- 架构：%s\n", remote.Arch)
		}
		fmt.Fprintf(&b, "\n以上信息来自对该主机（%s）的一次性探测，请作为环境事实使用；未列出的细节仍需向用户确认。", ctx.Host)
	}

	return b.String()
}

func isValidSSHContext(ctx *SSHPromptContext) bool {
	if ctx == nil {
		return false
	}
	return strings.TrimSpace(ctx.Host) != "" &&
		strings.TrimSpace(ctx.User) != "" &&
		ctx.Port > 0
}
