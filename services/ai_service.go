package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/ilaziness/vexo/internal/ai"
	"github.com/ilaziness/vexo/internal/database"
	"github.com/ilaziness/vexo/internal/secret"
	"go.uber.org/zap"
)

// AI 流式事件
const EventAIStreamChunk = "eventAIStreamChunk"

// AIConfig AI 配置结构体 - 用于配置文件持久化
type AIConfig struct {
	Enabled         bool    `json:"enabled" toml:"enabled"`
	Provider        string  `json:"provider" toml:"provider"`
	Model           string  `json:"model" toml:"model"`
	APIKey          string  `json:"api_key,omitempty" toml:"api_key,omitempty"`
	APIKeyEncrypted string  `json:"api_key_encrypted,omitempty" toml:"api_key_encrypted,omitempty"`
	Endpoint        string  `json:"endpoint" toml:"endpoint"`
	Temperature     float64 `json:"temperature" toml:"temperature"`
	MaxTokens       int     `json:"max_tokens" toml:"max_tokens"`
}

// ActiveSession 活跃会话信息
type ActiveSession struct {
	ID       string `json:"id"`
	IsActive bool   `json:"is_active"`
}

// AIMessage AI 消息结构（用于服务层）
type AIMessage struct {
	ID        string `json:"id"`
	SessionID string `json:"session_id"`
	Role      string `json:"role"` // 'user' or 'assistant'
	Content   string `json:"content"`
	Parts     string `json:"parts"` // JSON array of parts
	Timestamp int64  `json:"timestamp"`
}

// AISession AI 会话结构（用于服务层）
type AISession struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

// ChatRequest 多轮对话请求
type ChatRequest struct {
	SessionID  string `json:"session_id"`
	NewMessage string `json:"new_message"`
}

// ChatResponse 多轮对话响应
type ChatResponse struct {
	Message AIMessage `json:"message"`
}

// AIService AI 服务 - 前端接口层
type AIService struct {
	configService *ConfigService
	sshService    *SSHService
	engine        *ai.AIEngine
	sessionRepo   database.AISessionRepository
}

// NewAIService 创建 AI 服务实例
func NewAIService(configService *ConfigService, sshService *SSHService, db *database.Database) *AIService {
	// 先设置日志器，确保 AI 包初始化错误能被记录
	ai.SetLogger(Logger)

	s := &AIService{
		configService: configService,
		sshService:    sshService,
		sessionRepo:   db.AISessionRepo,
	}

	// 初始化时加载配置，仅在启用 AI 时创建引擎
	if cfg, err := s.GetConfig(); err == nil && cfg.Enabled {
		engine := ai.NewAIEngine()
		if err := engine.SetConfig(s.getEngineConfig(cfg)); err != nil {
			Logger.Warn("failed to init AI engine", zap.Error(err))
		} else {
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
		Enabled:     cfg.Enabled,
		Provider:    ai.Provider(cfg.Provider),
		Model:       cfg.Model,
		APIKey:      apiKey,
		Endpoint:    cfg.Endpoint,
		Temperature: cfg.Temperature,
		MaxTokens:   cfg.MaxTokens,
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
	if s.engine == nil {
		s.engine = ai.NewAIEngine()
	}
	if err := s.engine.SetConfig(engineConfig); err != nil {
		return fmt.Errorf("set ai engine config failed: %w", err)
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
		Enabled:     false,
		Provider:    string(ai.ProviderOllama),
		Model:       "llama3.2",
		Endpoint:    "http://localhost:11434",
		Temperature: 0.7,
		MaxTokens:   2048,
	}
	return s.SaveConfig(defaultConfig)
}

// GetProviders 获取所有支持的提供商列表 - 前端下拉选项
func (s *AIService) GetProviders() []ai.ProviderInfo {
	return ai.GetAllProviders()
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

// CreateSession 创建新会话
func (s *AIService) CreateSession() (*AISession, error) {
	session := &database.AISession{
		Title: "新会话",
	}
	if err := s.sessionRepo.CreateSession(context.Background(), session); err != nil {
		return nil, fmt.Errorf("create session failed: %w", err)
	}
	return &AISession{
		ID:        session.ID,
		Title:     session.Title,
		CreatedAt: session.CreatedAt.Unix(),
		UpdatedAt: session.UpdatedAt.Unix(),
	}, nil
}

// GetSession 获取会话
func (s *AIService) GetSession(id string) (*AISession, error) {
	session, err := s.sessionRepo.GetSession(context.Background(), id)
	if err != nil {
		return nil, fmt.Errorf("get session failed: %w", err)
	}
	return &AISession{
		ID:        session.ID,
		Title:     session.Title,
		CreatedAt: session.CreatedAt.Unix(),
		UpdatedAt: session.UpdatedAt.Unix(),
	}, nil
}

// ListSessions 列出会话
func (s *AIService) ListSessions(limit int) ([]*AISession, error) {
	sessions, err := s.sessionRepo.ListSessions(context.Background(), limit)
	if err != nil {
		return nil, fmt.Errorf("list sessions failed: %w", err)
	}
	result := make([]*AISession, len(sessions))
	for i, s := range sessions {
		result[i] = &AISession{
			ID:        s.ID,
			Title:     s.Title,
			CreatedAt: s.CreatedAt.Unix(),
			UpdatedAt: s.UpdatedAt.Unix(),
		}
	}
	return result, nil
}

// DeleteSession 删除会话
func (s *AIService) DeleteSession(id string) error {
	if err := s.sessionRepo.DeleteSession(context.Background(), id); err != nil {
		return fmt.Errorf("delete session failed: %w", err)
	}
	return nil
}

// ListMessages 列出会话的消息
func (s *AIService) ListMessages(sessionID string) ([]*AIMessage, error) {
	messages, err := s.sessionRepo.ListMessages(context.Background(), sessionID)
	if err != nil {
		return nil, fmt.Errorf("list messages failed: %w", err)
	}
	result := make([]*AIMessage, len(messages))
	for i, m := range messages {
		result[i] = &AIMessage{
			ID:        m.ID,
			SessionID: m.SessionID,
			Role:      m.Role,
			Content:   m.Content,
			Parts:     m.Parts,
			Timestamp: m.Timestamp.Unix(),
		}
	}
	return result, nil
}

// UpdateSession 更新会话
func (s *AIService) UpdateSession(session *AISession) error {
	dbSession := &database.AISession{
		ID:        session.ID,
		Title:     session.Title,
		CreatedAt: time.Unix(session.CreatedAt, 0),
		UpdatedAt: time.Unix(session.UpdatedAt, 0),
	}
	if err := s.sessionRepo.UpdateSession(context.Background(), dbSession); err != nil {
		return fmt.Errorf("update session failed: %w", err)
	}
	return nil
}

// Chat 多轮对话（流式），通过事件推送 chunk 到前端
func (s *AIService) Chat(req *ChatRequest) (*ChatResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("chat request is nil")
	}

	sessionID := strings.TrimSpace(req.SessionID)
	if sessionID == "" {
		return nil, fmt.Errorf("session id is required")
	}

	newMessage := strings.TrimSpace(req.NewMessage)
	if newMessage == "" {
		return nil, fmt.Errorf("new message is required")
	}

	if s.engine == nil {
		return nil, fmt.Errorf("AI engine not initialized")
	}

	history, err := s.sessionRepo.ListMessages(context.Background(), sessionID)
	if err != nil {
		return nil, fmt.Errorf("list history messages failed: %w", err)
	}

	messages := make([]ai.ChatMessage, len(history))
	for i, msg := range history {
		messages[i] = ai.ChatMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	// 保存用户消息
	userParts, err := json.Marshal([]map[string]string{{"type": "text", "text": newMessage}})
	if err != nil {
		return nil, fmt.Errorf("marshal user message parts failed: %w", err)
	}
	userMsg := &database.AIMessage{
		SessionID: sessionID,
		Role:      "user",
		Content:   newMessage,
		Parts:     string(userParts),
		Timestamp: time.Now(),
	}
	if err := s.sessionRepo.CreateMessage(context.Background(), userMsg); err != nil {
		return nil, fmt.Errorf("save user message failed: %w", err)
	}

	aiReq := &ai.ChatRequest{
		SessionID:  sessionID,
		Messages:   messages,
		NewMessage: newMessage,
	}

	// 会话无标题时，用本次用户输入截取作为标题
	if session, err := s.sessionRepo.GetSession(context.Background(), sessionID); err == nil && session.Title == "" {
		title := []rune(newMessage)
		if len(title) > 20 {
			title = title[:20]
		}
		if len(title) == 0 {
			session.Title = "新会话"
		} else {
			session.Title = string(title)
		}
		if err := s.sessionRepo.UpdateSession(context.Background(), session); err != nil {
			Logger.Warn("failed to update session title", zap.Error(err))
		}
	}

	// 流式调用 AI 引擎
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	aiResp, err := s.engine.ChatStream(ctx, aiReq, func(chunk ai.StreamChunk) error {
		if app != nil {
			app.Event.Emit(EventAIStreamChunk, map[string]string{
				"sessionId": sessionID,
				"type":      chunk.Type,
				"chunk":     chunk.Text,
			})
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("AI chat failed: %w", err)
	}

	// 保存 AI 回复
	assistantParts, err := json.Marshal([]map[string]string{{"type": "text", "text": aiResp.Message.Content}})
	if err != nil {
		return nil, fmt.Errorf("marshal assistant message parts failed: %w", err)
	}
	assistantMsg := &database.AIMessage{
		SessionID: sessionID,
		Role:      aiResp.Message.Role,
		Content:   aiResp.Message.Content,
		Parts:     string(assistantParts),
		Timestamp: time.Now(),
	}
	if err := s.sessionRepo.CreateMessage(context.Background(), assistantMsg); err != nil {
		return nil, fmt.Errorf("save assistant message failed: %w", err)
	}

	// 更新会话时间戳
	session, err := s.sessionRepo.GetSession(context.Background(), sessionID)
	if err == nil {
		session.UpdatedAt = time.Now()
		if err := s.sessionRepo.UpdateSession(context.Background(), session); err != nil {
			Logger.Warn("failed to update session timestamp", zap.Error(err))
		}
	}

	return &ChatResponse{
		Message: AIMessage{
			ID:        assistantMsg.ID,
			SessionID: assistantMsg.SessionID,
			Role:      assistantMsg.Role,
			Content:   assistantMsg.Content,
			Parts:     assistantMsg.Parts,
			Timestamp: assistantMsg.Timestamp.Unix(),
		},
	}, nil
}
