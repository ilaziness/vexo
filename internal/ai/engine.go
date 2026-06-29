package ai

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/compat_oai"
	"github.com/firebase/genkit/go/plugins/compat_oai/openai"
	"github.com/firebase/genkit/go/plugins/googlegenai"
	"github.com/firebase/genkit/go/plugins/ollama"
	"go.uber.org/zap"
)

var logger *zap.Logger

// SetLogger 设置AI包的日志器
func SetLogger(l *zap.Logger) {
	logger = l
}

// ChatMessage 聊天消息
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest 聊天请求
type ChatRequest struct {
	SessionID    string        `json:"session_id"`
	Messages     []ChatMessage `json:"messages"`
	NewMessage   string        `json:"new_message"`
	SystemPrompt string        `json:"system_prompt"`
}

// ChatResponse 聊天响应
type ChatResponse struct {
	Message ChatMessage `json:"message"`
}

// StreamChunk 流式数据块
type StreamChunk struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// AIEngine AI引擎
type AIEngine struct {
	genkit    *genkit.Genkit
	config    *Config
	model     ai.Model
	genConfig any
	mu        sync.RWMutex
}

// NewAIEngine 创建AI引擎
func NewAIEngine() *AIEngine {
	return &AIEngine{
		config: &Config{
			Enabled:     false,
			Provider:    ProviderOllama,
			Model:       "llama3.2",
			Endpoint:    "http://localhost:11434",
			Temperature: 0.7,
			MaxTokens:   2048,
		},
	}
}

// Init 初始化插件、定义模型、构建配置
func (e *AIEngine) Init(ctx context.Context, cfg *Config) error {
	var ollamaPlugin *ollama.Ollama
	var compatPlugin *compat_oai.OpenAICompatible
	var openaiPlugin *openai.OpenAI
	var g *genkit.Genkit

	switch cfg.Provider {
	case ProviderOllama:
		endpoint := cfg.Endpoint
		if endpoint == "" {
			endpoint = "http://localhost:11434"
		}
		ollamaPlugin = &ollama.Ollama{
			ServerAddress: endpoint,
			Timeout:       300,
		}
		g = genkit.Init(ctx, genkit.WithPlugins(ollamaPlugin))

	case ProviderGoogle:
		g = genkit.Init(ctx, genkit.WithPlugins(&googlegenai.GoogleAI{APIKey: cfg.APIKey}))

	case ProviderOpenAI:
		openaiPlugin = &openai.OpenAI{
			APIKey: cfg.APIKey,
		}
		g = genkit.Init(ctx, genkit.WithPlugins(openaiPlugin))

	case ProviderOpenAICompatible:
		if cfg.Endpoint == "" {
			return fmt.Errorf("endpoint is required for openai_compatible")
		}
		compatPlugin = &compat_oai.OpenAICompatible{
			Provider: "custom",
			APIKey:   cfg.APIKey,
			BaseURL:  cfg.Endpoint,
		}
		g = genkit.Init(ctx, genkit.WithPlugins(compatPlugin))

	default:
		return fmt.Errorf("unsupported provider: %s", cfg.Provider)
	}

	var model ai.Model
	switch cfg.Provider {
	case ProviderOllama:
		if ollamaPlugin == nil {
			return fmt.Errorf("ollama plugin not initialized")
		}
		ollamaPlugin.DefineModel(g, ollama.ModelDefinition{
			Name: cfg.Model,
			Type: "chat",
		}, nil)
		model = ollama.Model(g, cfg.Model)

	case ProviderOpenAI:
		openaiPlugin.DefineModel(cfg.Model, ai.ModelOptions{})

	case ProviderOpenAICompatible:
		if compatPlugin == nil {
			return fmt.Errorf("compat plugin not initialized")
		}
		compatPlugin.DefineModel(string(cfg.Provider), cfg.Model, ai.ModelOptions{})
		model = compatPlugin.Model(g, cfg.Model)

	case ProviderGoogle:
		model = genkit.LookupModel(g, cfg.Model)
	}

	var genConfig any
	if cfg.Provider != ProviderGoogle {
		genConfig = cfg.buildGenConfig()
	}

	e.mu.Lock()
	e.genkit = g
	e.config = cfg
	e.model = model
	e.genConfig = genConfig
	e.mu.Unlock()

	return nil
}

// ChatStream 流式对话
func (e *AIEngine) ChatStream(ctx context.Context, req *ChatRequest, onChunk func(StreamChunk) error) (*ChatResponse, error) {
	e.mu.RLock()
	model := e.model
	genConfig := e.genConfig
	g := e.genkit
	cfg := e.config
	e.mu.RUnlock()

	if g == nil || !cfg.Enabled {
		return nil, fmt.Errorf("ai engine not initialized or not enabled")
	}

	messages := buildMessages(req.Messages, req.NewMessage, req.SystemPrompt)

	var opts []ai.GenerateOption
	if model != nil {
		opts = append(opts, ai.WithModel(model))
	} else {
		opts = append(opts, ai.WithModelName("googleai/"+cfg.Model))
	}
	opts = append(opts, ai.WithMessages(messages...))
	if genConfig != nil {
		opts = append(opts, ai.WithConfig(genConfig))
	}

	var fullText strings.Builder
	stream := genkit.GenerateStream(ctx, g, opts...)

	for result, err := range stream {
		if err != nil {
			return nil, fmt.Errorf("generate error: %w", err)
		}
		if result.Done {
			break
		}

		// 同一 chunk 可能同时带 Reasoning 与 Text；优先走 reasoning，避免思考过程重复写入正文
		if reasoning := result.Chunk.Reasoning(); reasoning != "" {
			chunk := StreamChunk{Type: "reasoning", Text: reasoning}
			if err := onChunk(chunk); err != nil {
				return nil, err
			}
			continue
		}

		if text := result.Chunk.Text(); text != "" {
			chunk := StreamChunk{Type: "text", Text: text}
			if err := onChunk(chunk); err != nil {
				return nil, err
			}
			fullText.WriteString(text)
		}
	}

	return &ChatResponse{
		Message: ChatMessage{
			Role:    "assistant",
			Content: fullText.String(),
		},
	}, nil
}

// buildMessages 构建消息列表
func buildMessages(history []ChatMessage, newMsg, systemPrompt string) []*ai.Message {
	var msgs []*ai.Message
	if strings.TrimSpace(systemPrompt) != "" {
		msgs = append(msgs, ai.NewTextMessage(ai.RoleSystem, systemPrompt))
	}
	for _, m := range history {
		role := ai.RoleUser
		if m.Role == "assistant" {
			role = ai.RoleModel
		}
		msgs = append(msgs, ai.NewTextMessage(role, m.Content))
	}
	msgs = append(msgs, ai.NewUserTextMessage(newMsg))
	return msgs
}

// IsEnabled 是否启用
func (e *AIEngine) IsEnabled() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.config.Enabled && e.genkit != nil
}

// SetConfig 更新配置
func (e *AIEngine) SetConfig(cfg *Config) error {
	if !cfg.Enabled {
		e.mu.Lock()
		e.config = cfg
		e.genkit = nil
		e.model = nil
		e.genConfig = nil
		e.mu.Unlock()
		return nil
	}
	return e.Init(context.Background(), cfg)
}

// GetConfig 获取当前配置
func (e *AIEngine) GetConfig() Config {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return *e.config
}
