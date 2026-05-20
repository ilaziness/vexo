package ai

// Provider AI供应商类型
type Provider string

const (
	ProviderGoogle           Provider = "google"
	ProviderOpenAI           Provider = "openai"
	ProviderOpenAICompatible Provider = "openai_compatible"
	ProviderOllama           Provider = "ollama"
)

// ProviderInfo 供应商信息
type ProviderInfo struct {
	Name          string `json:"name"`
	Label         string `json:"label"`
	NeedsAPIKey   bool   `json:"needs_api_key"`
	NeedsEndpoint bool   `json:"needs_endpoint"`
}

// GetAllProviders 获取所有支持的供应商列表
func GetAllProviders() []ProviderInfo {
	return []ProviderInfo{
		{
			Name:        string(ProviderGoogle),
			Label:       "Google Gemini",
			NeedsAPIKey: true,
		},
		{
			Name:        string(ProviderOpenAI),
			Label:       "OpenAI",
			NeedsAPIKey: true,
		},
		{
			Name:          string(ProviderOpenAICompatible),
			Label:         "OpenAI Compatible",
			NeedsAPIKey:   true,
			NeedsEndpoint: true,
		},
		{
			Name:          string(ProviderOllama),
			Label:         "Ollama",
			NeedsEndpoint: true,
		},
	}
}
