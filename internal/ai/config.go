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

// Config AI配置
type Config struct {
	Enabled  bool
	Provider Provider
	Model    string
	APIKey   string
	Endpoint string

	Temperature float64
	MaxTokens   int
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
