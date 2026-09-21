package ai

import (
	"strings"

	"github.com/firebase/genkit/go/plugins/ollama"
	"github.com/openai/openai-go"
	"google.golang.org/genai"
)

var (
	defaultTemperature = 0.5
)

// Config AI配置（配置文件与引擎共用）
type Config struct {
	Enabled         bool     `json:"enabled" toml:"enabled"`
	Provider        Provider `json:"provider" toml:"provider"`
	Model           string   `json:"model" toml:"model"`
	APIKey          string   `json:"api_key,omitempty" toml:"api_key,omitempty"`
	APIKeyEncrypted string   `json:"api_key_encrypted,omitempty" toml:"api_key_encrypted,omitempty"`
	Endpoint        string   `json:"endpoint" toml:"endpoint"`
	Temperature     float64  `json:"temperature" toml:"temperature"`
	MaxTokens       int      `json:"max_tokens" toml:"max_tokens"`
}

// buildGenConfig 构建各供应商的原生配置结构
func (c *Config) buildGenConfig() any {
	switch c.Provider {
	case ProviderOllama:
		cfg := &ollama.GenerateContentConfig{
			Think:       ollama.ThinkEnabled(true),
			Temperature: &defaultTemperature,
		}
		if strings.Contains(c.Model, "GPT") {
			cfg.Think = ollama.ThinkEffort("medium")
		}
		if c.Temperature > 0 {
			cfg.Temperature = new(c.Temperature)
		}
		if c.MaxTokens > 0 {
			cfg.NumPredict = new(c.MaxTokens)
		}
		return cfg

	case ProviderOpenAI, ProviderOpenAICompatible:
		cfg := &openai.ChatCompletionNewParams{
			Temperature:     openai.Float(defaultTemperature),
			ReasoningEffort: openai.ReasoningEffortMedium,
		}
		if c.Temperature > 0 {
			cfg.Temperature = openai.Float(c.Temperature)
		}
		if c.MaxTokens > 0 {
			cfg.MaxTokens = openai.Int(int64(c.MaxTokens))
		}
		return cfg
	case ProviderGoogle:
		cfg := &genai.GenerateContentConfig{
			Temperature: genai.Ptr(float32(0.5)),
			ThinkingConfig: &genai.ThinkingConfig{
				IncludeThoughts: true,
			},
		}
		if c.Temperature > 0 {
			cfg.Temperature = genai.Ptr(float32(c.Temperature))
		}
		if c.MaxTokens > 0 {
			cfg.MaxOutputTokens = int32(c.MaxTokens)
		}
		return cfg
	}
	return nil
}
