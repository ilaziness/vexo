// Package ai provides AI functionality using Google Genkit
package ai

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var defaultOllamaModels = []string{"llama3.2", "qwen2.5", "mistral"}

var ollamaHTTPClient = &http.Client{Timeout: 5 * time.Second}

// Provider 模型提供商类型
type Provider string

const (
	ProviderGoogleGenAI      Provider = "googleai"          // Google Generative AI
	ProviderGoogleVertex     Provider = "googlevertex"      // Google Vertex AI
	ProviderOpenAI           Provider = "openai"            // OpenAI
	ProviderOpenAICompatible Provider = "openai_compatible" // OpenAI-Compatible APIs
	ProviderAnthropic        Provider = "anthropic"         // Anthropic (Claude)
	ProviderAWSBedrock       Provider = "aws_bedrock"       // AWS Bedrock
	ProviderAzureAIFoundry   Provider = "azure_foundry"     // Azure AI Foundry
	ProviderXAI              Provider = "xai"               // xAI (Grok)
	ProviderDeepSeek         Provider = "deepseek"          // DeepSeek
	ProviderOllama           Provider = "ollama"            // Ollama
)

// ProviderInfo 提供商信息
type ProviderInfo struct {
	Name          string `json:"name"`
	Label         string `json:"label"`
	Description   string `json:"description"`
	NeedsAPIKey   bool   `json:"needs_api_key"`
	NeedsEndpoint bool   `json:"needs_endpoint"`
}

// GetAllProviders 获取所有支持的提供商列表
func GetAllProviders() []ProviderInfo {
	return []ProviderInfo{
		{
			Name:        string(ProviderGoogleGenAI),
			Label:       "Google Generative AI",
			Description: "Google Gemini models",
			NeedsAPIKey: true,
		},
		{
			Name:        string(ProviderGoogleVertex),
			Label:       "Google Vertex AI",
			Description: "Enterprise Google AI",
			NeedsAPIKey: true,
		},
		{
			Name:        string(ProviderOpenAI),
			Label:       "OpenAI",
			Description: "GPT-4, GPT-3.5 models",
			NeedsAPIKey: true,
		},
		{
			Name:          string(ProviderOpenAICompatible),
			Label:         "OpenAI-Compatible APIs",
			Description:   "Compatible with OpenAI API format",
			NeedsAPIKey:   true,
			NeedsEndpoint: true,
		},
		{
			Name:        string(ProviderAnthropic),
			Label:       "Anthropic (Claude)",
			Description: "Claude models",
			NeedsAPIKey: true,
		},
		{
			Name:        string(ProviderAWSBedrock),
			Label:       "AWS Bedrock",
			Description: "AWS managed models",
			NeedsAPIKey: true,
		},
		{
			Name:          string(ProviderAzureAIFoundry),
			Label:         "Azure AI Foundry",
			Description:   "Azure AI models",
			NeedsAPIKey:   true,
			NeedsEndpoint: true,
		},
		{
			Name:          string(ProviderXAI),
			Label:         "xAI (Grok)",
			Description:   "Grok models",
			NeedsAPIKey:   true,
			NeedsEndpoint: true,
		},
		{
			Name:          string(ProviderDeepSeek),
			Label:         "DeepSeek",
			Description:   "DeepSeek models",
			NeedsAPIKey:   true,
			NeedsEndpoint: true,
		},
		{
			Name:          string(ProviderOllama),
			Label:         "Ollama",
			Description:   "Local models",
			NeedsAPIKey:   false,
			NeedsEndpoint: true,
		},
	}
}

// GetProviderModels 获取指定提供商的推荐模型列表
// endpoint 仅对需要自定义端点的提供商（如 Ollama）有效
func GetProviderModels(provider Provider, endpoint string) []string {
	switch provider {
	case ProviderGoogleGenAI:
		return []string{"gemini-2.5-flash", "gemini-2.5-pro", "gemini-1.5-flash"}
	case ProviderGoogleVertex:
		return []string{"gemini-2.5-flash", "gemini-2.5-pro"}
	case ProviderOpenAI:
		return []string{"gpt-4o", "gpt-4o-mini", "gpt-3.5-turbo"}
	case ProviderAnthropic:
		return []string{"claude-sonnet-4-20250514", "claude-haiku-4-20250514"}
	case ProviderXAI:
		return []string{"grok-2", "grok-1.5"}
	case ProviderDeepSeek:
		return []string{"deepseek-chat", "deepseek-coder"}
	case ProviderOllama:
		return fetchOllamaModels(endpoint)
	default:
		return nil
	}
}

// ollamaModelsResponse Ollama /api/tags 响应结构
type ollamaModelsResponse struct {
	Models []struct {
		Name string `json:"name"`
	} `json:"models"`
}

// fetchOllamaModels 从 Ollama 端点动态获取模型列表
func fetchOllamaModels(endpoint string) []string {
	if endpoint == "" {
		endpoint = "http://localhost:11434"
	}

	resp, err := ollamaHTTPClient.Get(endpoint + "/api/tags")
	if err != nil {
		return defaultOllamaModels
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return defaultOllamaModels
	}

	var result ollamaModelsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return defaultOllamaModels
	}

	if len(result.Models) == 0 {
		return defaultOllamaModels
	}

	models := make([]string, 0, len(result.Models))
	for _, m := range result.Models {
		if m.Name != "" {
			models = append(models, m.Name)
		}
	}
	return models
}

// IsGenkitSupportedProvider returns whether a provider is currently supported by Genkit integration.
func IsGenkitSupportedProvider(provider Provider) bool {
	switch provider {
	case ProviderGoogleGenAI,
		ProviderGoogleVertex,
		ProviderOpenAI,
		ProviderOpenAICompatible,
		ProviderAnthropic,
		ProviderAzureAIFoundry,
		ProviderXAI,
		ProviderDeepSeek,
		ProviderOllama:
		return true
	default:
		return false
	}
}

// GetGenkitModelName resolves the model name used by Genkit.Generate.
func GetGenkitModelName(provider Provider, model string) (string, error) {
	if model == "" {
		return "", fmt.Errorf("model is required")
	}

	switch provider {
	case ProviderOllama:
		return fmt.Sprintf("ollama/%s", model), nil
	case ProviderGoogleGenAI:
		return fmt.Sprintf("googleai/%s", model), nil
	case ProviderGoogleVertex:
		return fmt.Sprintf("vertexai/%s", model), nil
	case ProviderAzureAIFoundry:
		return fmt.Sprintf("azureaifoundry/%s", model), nil
	case ProviderOpenAI:
		return fmt.Sprintf("openai/%s", model), nil
	case ProviderOpenAICompatible:
		return fmt.Sprintf("openai_compatible/%s", model), nil
	case ProviderAnthropic:
		return fmt.Sprintf("anthropic/%s", model), nil
	case ProviderDeepSeek:
		return fmt.Sprintf("deepseek/%s", model), nil
	case ProviderXAI:
		return fmt.Sprintf("xai/%s", model), nil
	default:
		return "", fmt.Errorf("unsupported provider: %s", provider)
	}
}
