package ai

import (
	"context"
	"fmt"
	"sync"

	genai "github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/anthropic"
	"github.com/firebase/genkit/go/plugins/compat_oai"
	"github.com/firebase/genkit/go/plugins/compat_oai/openai"
	"github.com/firebase/genkit/go/plugins/googlegenai"
	"github.com/firebase/genkit/go/plugins/ollama"
	azureaifoundry "github.com/xavidop/genkit-azure-foundry-go"
	"go.uber.org/zap"
)

// logger 用于 AI 包的日志
var logger *zap.Logger

// SetLogger 设置日志器
func SetLogger(l *zap.Logger) {
	logger = l
}

func logError(msg string, err error) {
	if logger != nil {
		logger.Error(msg, zap.Error(err))
	}
}

func logWarn(msg string, fields ...zap.Field) {
	if logger != nil {
		logger.Warn(msg, fields...)
	}
}

// AIEngine AI 引擎
type AIEngine struct {
	genkit *genkit.Genkit
	config *Config
	mu     sync.RWMutex
}

// Config AI 配置
type Config struct {
	Enabled          bool
	Provider         Provider
	Model            string
	APIKey           string
	APIKeyEncrypted  string
	Endpoint         string // For Ollama or OpenAI-compatible
	Temperature      float64
	MaxTokens        int
	SafetyCheckLevel SafetyLevel
}

// CommandGenerateRequest 命令生成请求
type CommandGenerateRequest struct {
	Input            string   `json:"input"`
	CurrentDirectory string   `json:"current_directory"`
	OSInfo           string   `json:"os_info"`
	UserLevel        string   `json:"user_level"`
	RecentCommands   []string `json:"recent_commands"`
}

// CommandGenerateResponse 命令生成响应
type CommandGenerateResponse struct {
	Command       string             `json:"command"`
	Explanation   string             `json:"explanation"`
	SafetyWarning *SafetyCheckResult `json:"safety_warning"`
	Alternatives  []string           `json:"alternatives"`
	Confidence    float64            `json:"confidence"`
}

// CommandExplainRequest 命令解释请求
type CommandExplainRequest struct {
	Command          string `json:"command"`
	CurrentDirectory string `json:"current_directory"`
	OSInfo           string `json:"os_info"`
	UserLevel        string `json:"user_level"`
}

// CommandExplainResponse 命令解释响应
type CommandExplainResponse struct {
	Explanation   string             `json:"explanation"`
	Parts         []CommandPart      `json:"parts"`
	SafetyWarning *SafetyCheckResult `json:"safety_warning"`
}

// CommandPart 命令部分解释
type CommandPart struct {
	Part    string `json:"part"`
	Meaning string `json:"meaning"`
}

// NewAIEngine 创建 AI 引擎
func NewAIEngine() (*AIEngine, error) {
	engine := &AIEngine{
		config: &Config{
			Enabled:          false,
			Provider:         ProviderOllama,
			Model:            "llama3.2",
			Temperature:      0.7,
			MaxTokens:        2048,
			SafetyCheckLevel: SafetyMedium,
		},
	}

	if engine.config.Enabled {
		if err := engine.reinit(); err != nil {
			logWarn("failed to initialize AI provider plugin", zap.Error(err))
		}
	}

	return engine, nil
}

// SetConfig 设置配置
func (e *AIEngine) SetConfig(config *Config) {
	if config == nil {
		return
	}

	newConfig := *config

	e.mu.RLock()
	oldConfig := e.config
	e.mu.RUnlock()

	var (
		reinitializedGenkit *genkit.Genkit
		err                 error
	)

	// 仅在启用 AI 且相关配置变化时重新初始化 genkit
	needReinit := newConfig.Enabled && (oldConfig == nil ||
		!oldConfig.Enabled ||
		oldConfig.Provider != newConfig.Provider ||
		oldConfig.APIKey != newConfig.APIKey ||
		oldConfig.APIKeyEncrypted != newConfig.APIKeyEncrypted ||
		oldConfig.Endpoint != newConfig.Endpoint ||
		oldConfig.Model != newConfig.Model)

	if needReinit {
		reinitializedGenkit, err = initGenkitWithProvider(context.Background(), &newConfig)
		if err != nil {
			logError("failed to reinitialize AI provider plugin", err)
			return
		}
	}

	e.mu.Lock()
	e.config = &newConfig
	if !newConfig.Enabled {
		e.genkit = nil
	} else if reinitializedGenkit != nil {
		e.genkit = reinitializedGenkit
	}
	e.mu.Unlock()
}

// GetConfig 获取配置
func (e *AIEngine) GetConfig() *Config {
	e.mu.RLock()
	defer e.mu.RUnlock()
	if e.config == nil {
		return nil
	}
	c := *e.config
	return &c
}

// IsEnabled 检查是否启用
func (e *AIEngine) IsEnabled() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.config != nil && e.config.Enabled
}

// reinit 重新初始化 Genkit 插件实例
func (e *AIEngine) reinit() error {
	e.mu.RLock()
	if e.config == nil {
		e.mu.RUnlock()
		return fmt.Errorf("config is nil")
	}
	if !e.config.Enabled {
		e.mu.RUnlock()
		return nil
	}
	cfg := *e.config
	e.mu.RUnlock()

	g, err := initGenkitWithProvider(context.Background(), &cfg)
	if err != nil {
		return err
	}

	e.mu.Lock()
	e.genkit = g
	e.mu.Unlock()
	return nil
}

func initGenkitWithProvider(ctx context.Context, cfg *Config) (*genkit.Genkit, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is nil")
	}

	modelName, err := GetGenkitModelName(cfg.Provider, cfg.Model)
	if err != nil {
		return nil, err
	}

	opts := []genkit.GenkitOption{
		genkit.WithDefaultModel(modelName),
	}

	switch cfg.Provider {
	case ProviderOllama:
		endpoint := cfg.Endpoint
		if endpoint == "" {
			endpoint = "http://localhost:11434"
		}
		opts = append(opts, genkit.WithPlugins(&ollama.Ollama{
			ServerAddress: endpoint,
			Timeout:       120,
		}))
	case ProviderGoogleGenAI:
		opts = append(opts, genkit.WithPlugins(&googlegenai.GoogleAI{APIKey: cfg.APIKey}))
	case ProviderGoogleVertex:
		opts = append(opts, genkit.WithPlugins(&googlegenai.VertexAI{}))
	case ProviderAnthropic:
		opts = append(opts, genkit.WithPlugins(&anthropic.Anthropic{APIKey: cfg.APIKey}))
	case ProviderAzureAIFoundry:
		if cfg.Endpoint == "" {
			return nil, fmt.Errorf("endpoint is required for provider: %s", cfg.Provider)
		}
		opts = append(opts, genkit.WithPlugins(&azureaifoundry.AzureAIFoundry{
			Endpoint: cfg.Endpoint,
			APIKey:   cfg.APIKey,
		}))
	case ProviderOpenAI:
		opts = append(opts, genkit.WithPlugins(&openai.OpenAI{APIKey: cfg.APIKey}))
	case ProviderOpenAICompatible, ProviderDeepSeek, ProviderXAI:
		providerName, err := getCompatProviderName(cfg.Provider)
		if err != nil {
			return nil, err
		}
		if cfg.Endpoint == "" {
			return nil, fmt.Errorf("endpoint is required for provider: %s", cfg.Provider)
		}
		opts = append(opts, genkit.WithPlugins(&compat_oai.OpenAICompatible{
			Provider: providerName,
			APIKey:   cfg.APIKey,
			BaseURL:  cfg.Endpoint,
		}))
	default:
		return nil, fmt.Errorf("unsupported provider: %s", cfg.Provider)
	}

	return genkit.Init(ctx, opts...), nil
}

func getCompatProviderName(provider Provider) (string, error) {
	switch provider {
	case ProviderOpenAICompatible:
		return "openai_compatible", nil
	case ProviderDeepSeek:
		return "deepseek", nil
	case ProviderXAI:
		return "xai", nil
	default:
		return "", fmt.Errorf("provider %s is not compatible_oai provider", provider)
	}
}

// GenerateCommand 生成命令
func (e *AIEngine) GenerateCommand(ctx context.Context, req *CommandGenerateRequest) (*CommandGenerateResponse, error) {
	if !e.IsEnabled() {
		return nil, fmt.Errorf("AI is not enabled")
	}
	return e.generateCommand(ctx, req)
}

// generateCommand 内部实现
func (e *AIEngine) generateCommand(ctx context.Context, req *CommandGenerateRequest) (*CommandGenerateResponse, error) {
	prompt := BuildCommandGeneratorPrompt(
		req.Input,
		req.CurrentDirectory,
		req.OSInfo,
		req.UserLevel,
		req.RecentCommands,
	)

	e.mu.RLock()
	cfg := e.config
	g := e.genkit
	e.mu.RUnlock()

	result, err := callModel[CommandGenerateResponse](ctx, g, cfg, prompt)
	if err != nil {
		logError("failed to generate command", err)
		return nil, fmt.Errorf("model call failed: %w", err)
	}

	// 后端安全二次校验（仅当 AI 未返回安全警告时补充）
	safetyLevel := cfg.SafetyCheckLevel
	if result.SafetyWarning == nil && safetyLevel != SafetyNone {
		safetyResult := CheckCommandSafety(result.Command)
		if GetSafetyLevelPriority(safetyResult.Level) > GetSafetyLevelPriority(SafetyNone) {
			result.SafetyWarning = safetyResult
		}
	}

	return result, nil
}

// ExplainCommand 解释命令
func (e *AIEngine) ExplainCommand(ctx context.Context, req *CommandExplainRequest) (*CommandExplainResponse, error) {
	if !e.IsEnabled() {
		return nil, fmt.Errorf("AI is not enabled")
	}
	return e.explainCommand(ctx, req)
}

// explainCommand 内部实现
func (e *AIEngine) explainCommand(ctx context.Context, req *CommandExplainRequest) (*CommandExplainResponse, error) {
	prompt := BuildCommandExplainerPrompt(
		req.Command,
		req.CurrentDirectory,
		req.OSInfo,
		req.UserLevel,
	)

	e.mu.RLock()
	cfg := e.config
	g := e.genkit
	e.mu.RUnlock()

	result, err := callModel[CommandExplainResponse](ctx, g, cfg, prompt)
	if err != nil {
		logError("failed to explain command", err)
		return nil, fmt.Errorf("model call failed: %w", err)
	}

	// 后端安全校验
	safetyResult := CheckCommandSafety(req.Command)
	if GetSafetyLevelPriority(safetyResult.Level) > GetSafetyLevelPriority(SafetyNone) {
		result.SafetyWarning = safetyResult
	}

	return result, nil
}

// callModel 调用模型，使用泛型函数获取结构化输出。
func callModel[T any](ctx context.Context, g *genkit.Genkit, cfg *Config, prompt string) (*T, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is nil")
	}
	if g == nil {
		return nil, fmt.Errorf("genkit is not initialized")
	}

	modelName, err := GetGenkitModelName(cfg.Provider, cfg.Model)
	if err != nil {
		return nil, err
	}

	options := []genai.GenerateOption{
		genai.WithModelName(modelName),
		genai.WithPrompt(prompt),
	}
	providerConfig := buildProviderGenerateConfig(cfg.Provider, cfg.Temperature, cfg.MaxTokens)
	if providerConfig != nil {
		options = append(options, genai.WithConfig(providerConfig))
	}

	result, _, err := genkit.GenerateData[T](ctx, g, options...)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func buildProviderGenerateConfig(provider Provider, temperature float64, maxTokens int) interface{} {
	if temperature == 0 && maxTokens == 0 {
		return nil
	}

	switch provider {
	case ProviderOllama:
		cfg := map[string]interface{}{}
		if temperature != 0 {
			cfg["temperature"] = temperature
		}
		if maxTokens > 0 {
			cfg["num_predict"] = maxTokens
		}
		return cfg
	default:
		cfg := map[string]interface{}{}
		if temperature != 0 {
			cfg["temperature"] = temperature
		}
		if maxTokens > 0 {
			cfg["max_tokens"] = maxTokens
		}
		return cfg
	}
}
