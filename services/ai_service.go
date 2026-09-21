package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ilaziness/vexo/internal/ai"
	"github.com/ilaziness/vexo/internal/database"
	"github.com/ilaziness/vexo/internal/ssh"
	"github.com/wailsapp/wails/v3/pkg/application"
	"go.uber.org/zap"
)

const EventAIStreamChunk = "eventAIStreamChunk"

type AIStreamChunkData struct {
	SessionID string `json:"sessionId"`
	Type      string `json:"type"`
	Chunk     string `json:"chunk"`
}

func init() {
	application.RegisterEvent[AIStreamChunkData](EventAIStreamChunk)
}

var ErrAINotEnabled = errors.New("AI 助手尚未启用，请前往「设置 → AI」完成配置并启用后再试")

type AIMessage struct {
	ID        string `json:"id"`
	SessionID string `json:"session_id"`
	Role      string `json:"role"`
	Content   string `json:"content"`
	Parts     string `json:"parts"`
	Timestamp int64  `json:"timestamp"`
}

type AISession struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

type SSHContext struct {
	LinkID string `json:"link_id,omitempty"`
	Host   string `json:"host"`
	Port   int    `json:"port"`
	User   string `json:"user"`
}

type ChatRequest struct {
	SessionID  string      `json:"session_id"`
	NewMessage string      `json:"new_message"`
	SSHContext *SSHContext `json:"ssh_context,omitempty"`
}

type ChatResponse struct {
	Message AIMessage `json:"message"`
}

type AIService struct {
	app           *application.App
	logger        *zap.Logger
	configService *ConfigService
	sshService    *SSHService
	engine        *ai.AIEngine
	sessionRepo   database.AISessionRepository
}

func NewAIService(app *application.App, logger *zap.Logger, configService *ConfigService, sshService *SSHService, db *database.Database) *AIService {
	s := &AIService{app: app, logger: logger, configService: configService, sshService: sshService, sessionRepo: db.AISessionRepo}
	if cfg, err := s.GetConfig(); err == nil && cfg.Enabled {
		engine := ai.NewAIEngine(logger)
		engineCfg, err := s.engineConfig(cfg)
		if err != nil {
			logger.Warn("failed to init AI engine", zap.Error(err))
		} else if err := engine.SetConfig(engineCfg); err != nil {
			logger.Warn("failed to init AI engine", zap.Error(err))
		} else {
			s.engine = engine
		}
	}
	return s
}

func (s *AIService) GetConfig() (*AIConfig, error) {
	if s.configService == nil {
		return nil, fmt.Errorf("config service not initialized")
	}
	cfg := s.configService.aiConfig()
	return &cfg, nil
}

func (s *AIService) engineConfig(cfg *AIConfig) (*ai.Config, error) {
	apiKey := cfg.APIKey
	if apiKey == "" && cfg.APIKeyEncrypted != "" {
		decrypted, err := s.configService.decryptAPIKey(cfg.APIKeyEncrypted)
		if err != nil {
			return nil, err
		}
		apiKey = decrypted
	}
	out := *cfg
	out.APIKey = apiKey
	return &out, nil
}

func (s *AIService) SaveConfig(cfg AIConfig) error {
	runtimeAPIKey := cfg.APIKey
	if cfg.APIKey != "" && cfg.APIKeyEncrypted == "" {
		encrypted, err := s.configService.encryptAPIKey(cfg.APIKey)
		if err != nil {
			return fmt.Errorf("encrypt api key failed: %w", err)
		}
		cfg.APIKeyEncrypted = encrypted
		cfg.APIKey = ""
	}
	engineConfig, err := s.engineConfig(&cfg)
	if err != nil {
		return err
	}
	if runtimeAPIKey != "" {
		engineConfig.APIKey = runtimeAPIKey
	}
	if s.engine == nil {
		s.engine = ai.NewAIEngine(s.logger)
	}
	if err := s.engine.SetConfig(engineConfig); err != nil {
		return fmt.Errorf("set ai engine config failed: %w", err)
	}
	return s.configService.saveAI(cfg)
}

func (s *AIService) ResetConfig() error {
	return s.SaveConfig(AIConfig{
		Enabled: false, Provider: ai.ProviderOllama, Model: "llama3.2",
		Endpoint: "http://localhost:11434", Temperature: 0.7, MaxTokens: 2048,
	})
}

func (s *AIService) GetProviders() []ai.ProviderInfo { return ai.GetAllProviders() }

func (s *AIService) CreateSession() (*AISession, error) {
	session := &database.AISession{Title: "新会话"}
	if err := s.sessionRepo.CreateSession(context.Background(), session); err != nil {
		return nil, err
	}
	return toAISession(session), nil
}
func (s *AIService) GetSession(id string) (*AISession, error) {
	session, err := s.sessionRepo.GetSession(context.Background(), id)
	if err != nil {
		return nil, err
	}
	return toAISession(session), nil
}
func (s *AIService) ListSessions(limit int) ([]*AISession, error) {
	sessions, err := s.sessionRepo.ListSessions(context.Background(), limit)
	if err != nil {
		return nil, err
	}
	out := make([]*AISession, len(sessions))
	for i, sess := range sessions {
		out[i] = toAISession(sess)
	}
	return out, nil
}
func (s *AIService) DeleteSession(id string) error {
	return s.sessionRepo.DeleteSession(context.Background(), id)
}
func (s *AIService) ListMessages(sessionID string) ([]*AIMessage, error) {
	messages, err := s.sessionRepo.ListMessages(context.Background(), sessionID)
	if err != nil {
		return nil, err
	}
	out := make([]*AIMessage, len(messages))
	for i, m := range messages {
		out[i] = toAIMessage(m)
	}
	return out, nil
}
func (s *AIService) UpdateSession(session *AISession) error {
	return s.sessionRepo.UpdateSession(context.Background(), &database.AISession{
		ID: session.ID, Title: session.Title,
		CreatedAt: time.Unix(session.CreatedAt, 0), UpdatedAt: time.Unix(session.UpdatedAt, 0),
	})
}

func (s *AIService) Chat(req *ChatRequest) (*ChatResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("chat request is nil")
	}
	sessionID := strings.TrimSpace(req.SessionID)
	newMessage := strings.TrimSpace(req.NewMessage)
	if sessionID == "" || newMessage == "" {
		return nil, fmt.Errorf("session id and message are required")
	}
	cfg, err := s.GetConfig()
	if err != nil || !cfg.Enabled || s.engine == nil {
		return nil, ErrAINotEnabled
	}
	history, err := s.sessionRepo.ListMessages(context.Background(), sessionID)
	if err != nil {
		return nil, err
	}
	messages := make([]ai.ChatMessage, len(history))
	for i, msg := range history {
		messages[i] = ai.ChatMessage{Role: msg.Role, Content: msg.Content}
	}
	userParts, _ := json.Marshal([]map[string]string{{"type": "text", "text": newMessage}})
	userMsg := &database.AIMessage{SessionID: sessionID, Role: "user", Content: newMessage, Parts: string(userParts), Timestamp: time.Now()}
	if err := s.sessionRepo.CreateMessage(context.Background(), userMsg); err != nil {
		return nil, err
	}
	sshCtx, remoteInfo := s.resolveSSHContext(req.SSHContext)
	aiReq := &ai.ChatRequest{
		SessionID: sessionID, Messages: messages, NewMessage: newMessage,
		SystemPrompt: ai.BuildSystemPrompt(sshCtx, remoteInfo),
	}
	if session, err := s.sessionRepo.GetSession(context.Background(), sessionID); err == nil && len(history) == 0 && (session.Title == "" || session.Title == "新会话") {
		title := []rune(newMessage)
		if len(title) > 20 {
			title = title[:20]
		}
		session.Title = string(title)
		if session.Title == "" {
			session.Title = "新会话"
		}
		_ = s.sessionRepo.UpdateSession(context.Background(), session)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()
	aiResp, err := s.engine.ChatStream(ctx, aiReq, func(chunk ai.StreamChunk) error {
		s.app.Event.Emit(EventAIStreamChunk, AIStreamChunkData{SessionID: sessionID, Type: chunk.Type, Chunk: chunk.Text})
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("AI chat failed: %w", err)
	}
	assistantParts, _ := json.Marshal([]map[string]string{{"type": "text", "text": aiResp.Message.Content}})
	assistantMsg := &database.AIMessage{
		SessionID: sessionID, Role: aiResp.Message.Role, Content: aiResp.Message.Content,
		Parts: string(assistantParts), Timestamp: time.Now(),
	}
	if err := s.sessionRepo.CreateMessage(context.Background(), assistantMsg); err != nil {
		return nil, err
	}
	if session, err := s.sessionRepo.GetSession(context.Background(), sessionID); err == nil {
		session.UpdatedAt = time.Now()
		_ = s.sessionRepo.UpdateSession(context.Background(), session)
	}
	return &ChatResponse{Message: *toAIMessage(assistantMsg)}, nil
}

func (s *AIService) resolveSSHContext(ctx *SSHContext) (*ai.SSHPromptContext, *ssh.RemoteSystemInfo) {
	if ctx == nil || s.sshService == nil {
		return nil, nil
	}
	linkID := strings.TrimSpace(ctx.LinkID)
	host := strings.TrimSpace(ctx.Host)
	user := strings.TrimSpace(ctx.User)
	if linkID == "" || host == "" || user == "" || ctx.Port <= 0 {
		return nil, nil
	}
	if !s.sshService.hasSession(linkID) {
		return nil, nil
	}
	return &ai.SSHPromptContext{Host: host, Port: ctx.Port, User: user}, s.sshService.GetOrFetchRemoteSystemInfo(linkID, host)
}

func toAISession(s *database.AISession) *AISession {
	return &AISession{ID: s.ID, Title: s.Title, CreatedAt: s.CreatedAt.Unix(), UpdatedAt: s.UpdatedAt.Unix()}
}
func toAIMessage(m *database.AIMessage) *AIMessage {
	return &AIMessage{ID: m.ID, SessionID: m.SessionID, Role: m.Role, Content: m.Content, Parts: m.Parts, Timestamp: m.Timestamp.Unix()}
}
