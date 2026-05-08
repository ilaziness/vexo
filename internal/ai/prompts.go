package ai

import (
	"encoding/json"
	"fmt"
)

// BuildCommandGeneratorPrompt 构建命令生成 Prompt
func BuildCommandGeneratorPrompt(input, currentDir, osInfo, userLevel string, recentCommands []string) string {
	recentCmdsStr := "[]"
	if len(recentCommands) > 0 {
		b, _ := json.Marshal(recentCommands)
		recentCmdsStr = string(b)
	}

	return fmt.Sprintf(`You are an expert Linux command assistant.

Context:
- Current Directory: %s
- Operating System: %s
- User Level: %s
- Recent Commands: %s

User Request: %s

Generate a command with fields: command, explanation, safety_warning (level, message, risk_detail, suggestion), alternatives (array), confidence (0-1).

Safety levels: none (safe), low (review), medium (understand), high (confirmation required).`,
		currentDir, osInfo, userLevel, recentCmdsStr, input)
}

// BuildCommandExplainerPrompt 构建命令解释 Prompt
func BuildCommandExplainerPrompt(command, currentDir, osInfo, userLevel string) string {
	return fmt.Sprintf(`You are an expert Linux command explainer.

Context:
- Current Directory: %s
- Operating System: %s
- User Level: %s

Command: %s

Explain this command with fields: explanation (string), parts (array of {part, meaning}), safety_warning (level, message, risk_detail, suggestion).

Safety levels: none (safe), low (minor risk), medium (moderate risk), high (dangerous).`,
		currentDir, osInfo, userLevel, command)
}
