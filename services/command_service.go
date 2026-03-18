package services

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ilaziness/vexo/internal/database"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
	"go.uber.org/zap"
)

// CommandsJSON embed 的 commands.json 数据
var CommandsJSON []byte

// BuiltinCommand 内置命令结构
type BuiltinCommand struct {
	Name        string `json:"name"`        // 命令名称
	Command     string `json:"command"`     // 完整命令
	Description string `json:"description"` // 描述和用法说明
}

// CategoryCommands 分类及其命令（用于 JSON 解析）
type CategoryCommands struct {
	Category string           `json:"category"` // 分类名
	Commands []BuiltinCommand `json:"commands"` // 该分类下的命令列表
}

// UserCommand 用户自定义命令结构
type UserCommand struct {
	Category    string `json:"category"`
	Name        string `json:"name"`
	Command     string `json:"command"`
	Description string `json:"description"`
	CreatedAt   int64  `json:"created_at"`
}

// CommandHistory 命令历史记录
type CommandHistory struct {
	Timestamp int64  `json:"timestamp"`
	Command   string `json:"command"`
}

// CommandInfo 返回给前端的命令信息（合并内置和用户）
type CommandInfo struct {
	Category    string `json:"category"`
	Name        string `json:"name"`
	Command     string `json:"command"`
	Description string `json:"description"`
	IsCustom    bool   `json:"is_custom"` // 是否为用户自定义
}

// SendCommandRequest 发送命令请求
type SendCommandRequest struct {
	Command    string   `json:"command"`     // 要发送的命令内容
	SessionIDs []string `json:"session_ids"` // 目标 SSH 会话 ID 列表
}

// CommandService 命令服务
type CommandService struct {
	builtinCommands map[string][]BuiltinCommand // 内置命令（按分类组织）
	sshService      *SSHService                 // SSH 服务引用
	db              *database.Database          // 数据库实例
}

// NewCommandService 创建命令服务实例
func NewCommandService(sshService *SSHService, db *database.Database) *CommandService {
	cs := &CommandService{
		sshService: sshService,
		db:         db,
	}

	// 加载内置命令
	if err := cs.loadBuiltinCommands(); err != nil {
		Logger.Error("load builtin commands failed", zap.Error(err))
	}

	return cs
}

// loadBuiltinCommands 从 embed JSON 文件加载内置命令
func (cs *CommandService) loadBuiltinCommands() error {
	var categories []CategoryCommands
	if err := json.Unmarshal(CommandsJSON, &categories); err != nil {
		return fmt.Errorf("unmarshal builtin commands: %w", err)
	}

	// 转换为 map[string][]BuiltinCommand
	builtinCmds := make(map[string][]BuiltinCommand)
	for _, cat := range categories {
		builtinCmds[cat.Category] = cat.Commands
	}

	cs.builtinCommands = builtinCmds
	Logger.Debug("loaded builtin commands", zap.Int("categories", len(categories)))
	return nil
}

// GetUserCommands 获取用户自定义命令
func (cs *CommandService) GetUserCommands() []UserCommand {
	dbCommands, err := cs.db.UserCommandRepo.GetAllCommands()
	if err != nil {
		Logger.Error("get user commands failed", zap.Error(err))
		return []UserCommand{}
	}

	commands := make([]UserCommand, 0, len(dbCommands))
	for _, cmd := range dbCommands {
		commands = append(commands, UserCommand{
			Category:    cmd.Category,
			Name:        cmd.Name,
			Command:     cmd.Command,
			Description: cmd.Description,
			CreatedAt:   cmd.CreatedAt.UnixNano() / 1000000,
		})
	}
	return commands
}

// SaveUserCommand 保存用户命令到数据库
func (cs *CommandService) SaveUserCommand(cmd UserCommand) error {
	existing, err := cs.db.UserCommandRepo.GetCommandByCategoryAndName(cmd.Category, cmd.Name)

	cmd.CreatedAt = time.Now().UnixNano() / 1000000

	if err == nil && existing != nil {
		// 更新
		dbCmd := &database.UserCommandDB{
			ID:          existing.ID,
			Category:    cmd.Category,
			Name:        cmd.Name,
			Command:     cmd.Command,
			Description: cmd.Description,
			CreatedAt:   time.Unix(cmd.CreatedAt/1000, 0),
		}
		return cs.db.UserCommandRepo.UpdateCommand(dbCmd)
	} else {
		// 新增
		dbCmd := &database.UserCommandDB{
			Category:    cmd.Category,
			Name:        cmd.Name,
			Command:     cmd.Command,
			Description: cmd.Description,
			CreatedAt:   time.Unix(cmd.CreatedAt/1000, 0),
		}
		_, err := cs.db.UserCommandRepo.InsertCommand(dbCmd)
		return err
	}
}

// DeleteUserCommand 删除用户命令
func (cs *CommandService) DeleteUserCommand(category, name string) error {
	cmd, err := cs.db.UserCommandRepo.GetCommandByCategoryAndName(category, name)
	if err != nil {
		return fmt.Errorf("user command not found: category=%s, name=%s", category, name)
	}
	return cs.db.UserCommandRepo.DeleteCommand(cmd.ID)
}

// GetCommandHistory 获取命令历史
func (cs *CommandService) GetCommandHistory() []CommandHistory {
	dbHistory, err := cs.db.CommandHistoryRepo.GetHistory()
	if err != nil {
		Logger.Error("get command history failed", zap.Error(err))
		return []CommandHistory{}
	}

	history := make([]CommandHistory, 0, len(dbHistory))
	for _, h := range dbHistory {
		history = append(history, CommandHistory{
			Command:   h.Command,
			Timestamp: h.Timestamp.UnixNano() / 1000000,
		})
	}
	return history
}

// AddCommandHistory 添加命令到历史
func (cs *CommandService) AddCommandHistory(command string) error {
	_, err := cs.db.CommandHistoryRepo.InsertHistory(&database.CommandHistoryDB{
		Command:   command,
		Timestamp: time.Now(),
	})
	return err
}

// ClearCommandHistory 清空命令历史
func (cs *CommandService) ClearCommandHistory() error {
	return cs.db.CommandHistoryRepo.ClearHistory()
}

// GetBuiltinCommands 获取内置命令列表（按分类）
func (cs *CommandService) GetBuiltinCommands() map[string][]BuiltinCommand {
	return cs.builtinCommands
}

// GetAllCommands 按分类返回所有命令（合并内置和用户，用户覆盖内置）
func (cs *CommandService) GetAllCommands() map[string][]CommandInfo {
	result := make(map[string][]CommandInfo)

	// 先添加内置命令（已经是按分类组织的）
	for category, cmds := range cs.builtinCommands {
		for _, cmd := range cmds {
			result[category] = append(result[category], CommandInfo{
				Category:    category,
				Name:        cmd.Name,
				Command:     cmd.Command,
				Description: cmd.Description,
				IsCustom:    false,
			})
		}
	}

	// 用户命令覆盖
	userCmds := cs.GetUserCommands()
	for _, cmd := range userCmds {
		key := cmd.Category
		// 查找并替换同名的内置命令
		replaced := false
		for i, existing := range result[key] {
			if existing.Name == cmd.Name {
				result[key][i] = CommandInfo{
					Category:    cmd.Category,
					Name:        cmd.Name,
					Command:     cmd.Command,
					Description: cmd.Description,
					IsCustom:    true,
				}
				replaced = true
				break
			}
		}
		// 如果是新的则添加
		if !replaced {
			result[key] = append(result[key], CommandInfo{
				Category:    cmd.Category,
				Name:        cmd.Name,
				Command:     cmd.Command,
				Description: cmd.Description,
				IsCustom:    true,
			})
		}
	}

	return result
}

// SendCommand 发送命令到指定的 SSH 会话
func (cs *CommandService) SendCommand(req SendCommandRequest) error {
	if cs.sshService == nil {
		return fmt.Errorf("SSH service not initialized")
	}

	// 发送到各个会话
	var lastErr error
	for _, sessionID := range req.SessionIDs {
		err := cs.sshService.SendToSession(sessionID, req.Command)
		if err != nil {
			Logger.Error("send command failed", zap.String("session", sessionID), zap.Error(err))
			lastErr = err
		}
	}

	// 添加到历史
	if lastErr == nil {
		cs.AddCommandHistory(req.Command)
	}

	return lastErr
}

// ShowWindow 显示命令面板窗口
func (cs *CommandService) ShowWindow() {
	if AppInstance.CommandWindow == nil {
		AppInstance.CommandWindow = newCommandWindow()
	}
	AppInstance.CommandWindow.OnWindowEvent(events.Common.WindowClosing, func(event *application.WindowEvent) {
		AppInstance.CommandWindow = nil
	})
	AppInstance.CommandWindow.Show()
	AppInstance.CommandWindow.Focus()
}

// CloseWindow 关闭命令面板窗口
func (cs *CommandService) CloseWindow() {
	if AppInstance.CommandWindow != nil {
		AppInstance.CommandWindow.Close()
		AppInstance.CommandWindow = nil
	}
}
