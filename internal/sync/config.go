package sync

// SyncConfig 同步配置
type SyncConfig struct {
	ServerURL string `toml:"server_url" json:"serverUrl"` // 同步服务器地址
	SyncID    string `toml:"sync_id" json:"syncId"`       // 同步ID
	UserKey   string `toml:"user_key" json:"userKey"`     // 用户密钥（用于派生加密密钥）
}

// IsConfigured 检查是否已配置
func (c *SyncConfig) IsConfigured() bool {
	return c.ServerURL != "" && c.SyncID != "" && c.UserKey != ""
}
