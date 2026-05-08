package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/ilaziness/vexo/internal/ai"
	"github.com/ilaziness/vexo/internal/secret"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
)

// AIConfig AI 配置结构体 - 用于配置文件持久化
type AIConfig struct {
	Enabled          bool    `json:"enabled" toml:"enabled"`
	Provider         string  `json:"provider" toml:"provider"`
	Model            string  `json:"model" toml:"model"`
	APIKey           string  `json:"api_key,omitempty" toml:"api_key,omitempty"`
	APIKeyEncrypted  string  `json:"api_key_encrypted,omitempty" toml:"api_key_encrypted,omitempty"`
	Endpoint         string  `json:"endpoint" toml:"endpoint"` // 自定义端点（Ollama/OpenAI-compatible等）
	Temperature      float64 `json:"temperature" toml:"temperature"`
	MaxTokens        int     `json:"max_tokens" toml:"max_tokens"`
	SafetyCheckLevel string  `json:"safety_check_level" toml:"safety_check_level"`
}

// ActiveSession 活跃会话信息
type ActiveSession struct {
	ID       string `json:"id"`
	IsActive bool   `json:"is_active"`
}

// AIService AI 服务 - 前端接口层
type AIService struct {
	configService *ConfigService
	sshService    *SSHService
	engine        *ai.AIEngine
}

// NewAIService 创建 AI 服务实例
func NewAIService(configService *ConfigService, sshService *SSHService) *AIService {
	// 先设置日志器，确保 AI 包初始化错误能被记录
	ai.SetLogger(Logger)

	s := &AIService{
		configService: configService,
		sshService:    sshService,
	}

	// 初始化时加载配置，仅在启用 AI 时创建引擎
	if cfg, err := s.GetConfig(); err == nil && cfg.Enabled {
		engine, err := ai.NewAIEngine()
		if err != nil {
			Logger.Warn("failed to create AI engine", zap.Error(err))
		} else {
			engine.SetConfig(s.getEngineConfig(cfg))
			s.engine = engine
		}
	}

	return s
}

// GetConfig 获取当前 AI 配置
func (s *AIService) GetConfig() (*AIConfig, error) {
	if s.configService == nil || s.configService.Config == nil {
		return nil, fmt.Errorf("config service not initialized")
	}
	return &s.configService.Config.AI, nil
}

// getEngineConfig 转换为引擎配置
func (s *AIService) getEngineConfig(cfg *AIConfig) *ai.Config {
	apiKey := cfg.APIKey
	if apiKey == "" && cfg.APIKeyEncrypted != "" {
		decrypted, err := s.decryptAPIKey(cfg.APIKeyEncrypted)
		if err == nil {
			apiKey = decrypted
		}
	}

	return &ai.Config{
		Enabled:          cfg.Enabled,
		Provider:         ai.Provider(cfg.Provider),
		Model:            cfg.Model,
		APIKey:           apiKey,
		APIKeyEncrypted:  cfg.APIKeyEncrypted,
		Endpoint:         cfg.Endpoint,
		Temperature:      cfg.Temperature,
		MaxTokens:        cfg.MaxTokens,
		SafetyCheckLevel: ai.SafetyLevel(cfg.SafetyCheckLevel),
	}
}

// SaveConfig 保存 AI 配置
func (s *AIService) SaveConfig(config AIConfig) error {
	if s.configService == nil || s.configService.Config == nil {
		return fmt.Errorf("config service not initialized")
	}

	runtimeAPIKey := config.APIKey

	// 如果有明文的 API Key，加密存储
	if config.APIKey != "" && config.APIKeyEncrypted == "" {
		encrypted, err := s.encryptAPIKey(config.APIKey)
		if err != nil {
			return fmt.Errorf("encrypt api key failed: %w", err)
		}
		config.APIKeyEncrypted = encrypted
		config.APIKey = "" // 清空明文
	}

	// 更新引擎配置
	engineConfig := s.getEngineConfig(&config)
	if runtimeAPIKey != "" {
		engineConfig.APIKey = runtimeAPIKey
	}
	if s.engine != nil {
		s.engine.SetConfig(engineConfig)
	}

	// 保存到配置文件
	s.configService.Config.AI = config

	if err := s.configService.SaveConfig(*s.configService.Config); err != nil {
		return fmt.Errorf("save config failed: %w", err)
	}

	return nil
}

// ResetConfig 重置 AI 配置为默认值
func (s *AIService) ResetConfig() error {
	defaultConfig := AIConfig{
		Enabled:          false,
		Provider:         string(ai.ProviderOllama),
		Model:            "llama3.2",
		Temperature:      0.7,
		MaxTokens:        2048,
		SafetyCheckLevel: string(ai.SafetyMedium),
	}
	return s.SaveConfig(defaultConfig)
}

// GenerateCommand 生成命令 - 前端调用接口
func (s *AIService) GenerateCommand(req ai.CommandGenerateRequest) (*ai.CommandGenerateResponse, error) {
	if s.engine == nil {
		return nil, fmt.Errorf("AI engine not initialized")
	}
	return s.engine.GenerateCommand(context.Background(), &req)
}

// ExplainCommand 解释命令 - 前端调用接口
func (s *AIService) ExplainCommand(req ai.CommandExplainRequest) (*ai.CommandExplainResponse, error) {
	if s.engine == nil {
		return nil, fmt.Errorf("AI engine not initialized")
	}
	return s.engine.ExplainCommand(context.Background(), &req)
}

// GetProviders 获取所有支持的提供商列表 - 前端下拉选项
func (s *AIService) GetProviders() []ai.ProviderInfo {
	return ai.GetAllProviders()
}

// ListModels 获取指定提供商的推荐模型列表 - 前端下拉选项
func (s *AIService) ListModels(provider string) []string {
	endpoint := ""
	if cfg, err := s.GetConfig(); err == nil && cfg != nil {
		endpoint = cfg.Endpoint
	}
	return ai.GetProviderModels(ai.Provider(provider), endpoint)
}

// GetActiveSessions 获取活跃 SSH 会话列表
func (s *AIService) GetActiveSessions() []ActiveSession {
	if s.sshService == nil {
		return []ActiveSession{}
	}

	sessions := make([]ActiveSession, 0)
	for _, sess := range s.sshService.GetActiveSessions() {
		id, ok := sess["id"].(string)
		if !ok || id == "" {
			continue
		}
		sessions = append(sessions, ActiveSession{
			ID:       id,
			IsActive: true,
		})
	}
	return sessions
}

// encryptAPIKey 加密 API Key 使用用户密码
func (s *AIService) encryptAPIKey(apiKey string) (string, error) {
	if s.configService == nil {
		return "", fmt.Errorf("config service not initialized")
	}
	password, err := s.configService.getPasswordWithPrompt("需要密码来加密 API Key")
	if err != nil {
		return "", err
	}
	return secret.Encrypt(password, apiKey)
}

// decryptAPIKey 解密 API Key 使用用户密码
func (s *AIService) decryptAPIKey(encrypted string) (string, error) {
	if s.configService == nil {
		return "", fmt.Errorf("config service not initialized")
	}
	password, err := s.configService.getPasswordWithPrompt("需要密码来解密 API Key")
	if err != nil {
		return "", err
	}
	return secret.Decrypt(password, encrypted)
}

// SafetyLevelInfo 安全等级信息
type SafetyLevelInfo struct {
	Value       string `json:"value"`
	Label       string `json:"label"`
	Description string `json:"description"`
}

// GetSafetyLevels 获取所有安全检查等级选项
func (s *AIService) GetSafetyLevels() []SafetyLevelInfo {
	return []SafetyLevelInfo{
		{
			Value:       "low",
			Label:       "低",
			Description: "仅警告高风险命令",
		},
		{
			Value:       "medium",
			Label:       "中",
			Description: "警告中高风险命令 (推荐)",
		},
		{
			Value:       "high",
			Label:       "高",
			Description: "警告所有潜在风险",
		},
	}
}

// SessionContext SSH会话上下文信息
type SessionContext struct {
	CurrentDirectory string   `json:"current_directory"` // 当前工作目录
	OSInfo           string   `json:"os_info"`           // 操作系统信息
	UserLevel        string   `json:"user_level"`        // 用户级别
	RecentCommands   []string `json:"recent_commands"`   // 最近执行的命令
}

// GetSessionContext 获取指定SSH会话的上下文信息
func (s *AIService) GetSessionContext(sessionID string) (*SessionContext, error) {
	if s.sshService == nil {
		return nil, fmt.Errorf("ssh service not initialized")
	}

	// 获取会话连接
	connAny, ok := s.sshService.SSHConnects.Load(sessionID)
	if !ok {
		return nil, fmt.Errorf("session %s not found", sessionID)
	}
	conn := connAny.(*SSHConnect)

	ctx := &SessionContext{
		UserLevel:        "intermediate",
		RecentCommands:   []string{},
		CurrentDirectory: "~",
		OSInfo:           "Linux",
	}

	// 尝试获取当前目录
	if conn.session != nil {
		// 发送 pwd 命令获取当前目录
		output, err := s.executeCommand(conn.session, "pwd")
		if err == nil && len(output) > 0 {
			ctx.CurrentDirectory = strings.TrimSpace(output)
		}

		// 发送 uname 命令获取操作系统信息
		output, err = s.executeCommand(conn.session, "uname -a")
		if err == nil && len(output) > 0 {
			ctx.OSInfo = strings.TrimSpace(output)
		}

		// 发送 whoami 命令获取当前用户
		output, err = s.executeCommand(conn.session, "whoami")
		if err == nil && len(output) > 0 {
			user := strings.TrimSpace(output)
			if user == "root" {
				ctx.UserLevel = "advanced"
			}
		}
	}

	return ctx, nil
}

// executeCommand 在SSH会话中执行命令并返回输出
func (s *AIService) executeCommand(session *ssh.Session, command string) (string, error) {
	output, err := session.CombinedOutput(command + "\n")
	if err != nil {
		return "", err
	}
	return string(output), nil
}
