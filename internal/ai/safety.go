package ai

import (
	"regexp"
	"strings"
)

// SafetyLevel 安全等级
type SafetyLevel string

const (
	SafetyNone   SafetyLevel = "none"
	SafetyLow    SafetyLevel = "low"
	SafetyMedium SafetyLevel = "medium"
	SafetyHigh   SafetyLevel = "high"
)

// SafetyRule 安全规则
type SafetyRule struct {
	Name       string
	Pattern    *regexp.Regexp
	Level      SafetyLevel
	Message    string
	RiskDetail string
	Suggestion string
}

// SafetyCheckResult 安全检查结果
type SafetyCheckResult struct {
	Level      SafetyLevel `json:"level"`
	Message    string      `json:"message"`
	RiskDetail string      `json:"risk_detail"`
	Suggestion string      `json:"suggestion"`
}

var (
	// safetyRules 安全规则列表
	safetyRules = []SafetyRule{
		{
			Name:       "rm_rf",
			Pattern:    regexp.MustCompile(`(?i)rm\s+(-[a-zA-Z]*f|--force)`),
			Level:      SafetyHigh,
			Message:    "危险：强制删除命令",
			RiskDetail: "此命令会强制删除文件或目录，且无法恢复。使用 -f 标志会跳过确认提示。",
			Suggestion: "考虑使用 'rm -i' 交互模式，或先使用 'ls' 确认要删除的内容",
		},
		{
			Name:       "rm_recursive",
			Pattern:    regexp.MustCompile(`(?i)rm\s+.*-r\s+/(home|etc|var|usr|bin|lib|sbin|opt|tmp|root)`),
			Level:      SafetyHigh,
			Message:    "危险：递归删除系统目录",
			RiskDetail: "递归删除系统关键目录可能导致系统无法启动或数据永久丢失。",
			Suggestion: "绝对不要删除系统目录。如需清理空间，请使用系统工具或手动删除特定文件。",
		},
		{
			Name:       "dd_command",
			Pattern:    regexp.MustCompile(`(?i)\bdd\s+if=`),
			Level:      SafetyHigh,
			Message:    "警告：dd 磁盘写入命令",
			RiskDetail: "dd 命令直接写入磁盘设备，操作不当会覆盖重要数据或破坏文件系统。",
			Suggestion: "仔细检查 if= 和 of= 参数，确保源和目标正确。先使用 --dry-run 或在小规模测试。",
		},
		{
			Name:       "chmod_777",
			Pattern:    regexp.MustCompile(`(?i)chmod\s+.*777`),
			Level:      SafetyMedium,
			Message:    "警告：完全开放权限",
			RiskDetail: "777 权限允许所有用户读写执行，存在安全风险，可能导致未授权访问。",
			Suggestion: "使用最小权限原则，通常 755 (目录) 或 644 (文件) 已足够。",
		},
		{
			Name:       "chown_root",
			Pattern:    regexp.MustCompile(`(?i)sudo\s+chown\s+.*root`),
			Level:      SafetyMedium,
			Message:    "警告：修改文件所有者",
			RiskDetail: "修改文件所有者为 root 后，普通用户可能无法访问或修改这些文件。",
			Suggestion: "确认确实需要更改所有权，或考虑使用 ACL (setfacl) 进行更细粒度的权限控制。",
		},
		{
			Name:       "sudo_command",
			Pattern:    regexp.MustCompile(`(?i)^sudo\s+`),
			Level:      SafetyLow,
			Message:    "提示：使用 sudo 提升权限",
			RiskDetail: "sudo 命令以 root 权限执行，可能导致系统级别的变更。",
			Suggestion: "确认命令来源可信，理解命令的作用后再执行。",
		},
		{
			Name:       "wget_curl_pipe",
			Pattern:    regexp.MustCompile(`(?i)(wget|curl)\s+.*\|\s*sh`),
			Level:      SafetyHigh,
			Message:    "危险：管道执行远程脚本",
			RiskDetail: "直接下载并执行远程脚本可能引入恶意代码，造成严重的安全问题。",
			Suggestion: "先下载脚本到本地审查，确认安全后再执行。",
		},
		{
			Name:       "mkfs_command",
			Pattern:    regexp.MustCompile(`(?i)^mkfs\.`),
			Level:      SafetyHigh,
			Message:    "警告：格式化文件系统",
			RiskDetail: "mkfs 命令会格式化磁盘分区，导致该分区上所有数据丢失。",
			Suggestion: "双重确认目标设备，确保没有挂载重要数据。使用 -n (dry-run) 先检查。",
		},
		{
			Name:       "killall_pkill",
			Pattern:    regexp.MustCompile(`(?i)^(killall|pkill)\s+-9`),
			Level:      SafetyMedium,
			Message:    "警告：强制终止进程",
			RiskDetail: "强制终止进程可能导致数据丢失或服务中断。",
			Suggestion: "先尝试不带 -9 的 kill/killall，让进程优雅退出。",
		},
	}
)

// CheckCommandSafety 检查命令安全等级
func CheckCommandSafety(command string) *SafetyCheckResult {
	highestLevel := SafetyNone
	var messages []string
	var riskDetails []string
	var suggestions []string

	for _, rule := range safetyRules {
		if rule.Pattern.MatchString(command) {
			if GetSafetyLevelPriority(rule.Level) > GetSafetyLevelPriority(highestLevel) {
				highestLevel = rule.Level
			}
			messages = append(messages, rule.Message)
			riskDetails = append(riskDetails, rule.RiskDetail)
			suggestions = append(suggestions, rule.Suggestion)
		}
	}

	if highestLevel == SafetyNone {
		return &SafetyCheckResult{
			Level:      SafetyNone,
			Message:    "命令安全",
			RiskDetail: "未发现明显风险",
			Suggestion: "",
		}
	}

	return &SafetyCheckResult{
		Level:      highestLevel,
		Message:    strings.Join(messages, "; "),
		RiskDetail: strings.Join(riskDetails, "; "),
		Suggestion: strings.Join(suggestions, "; "),
	}
}

// GetSafetyLevelPriority 获取安全等级优先级
func GetSafetyLevelPriority(level SafetyLevel) int {
	switch level {
	case SafetyNone:
		return 0
	case SafetyLow:
		return 1
	case SafetyMedium:
		return 2
	case SafetyHigh:
		return 3
	default:
		return 0
	}
}
