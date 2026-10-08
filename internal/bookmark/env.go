package bookmark

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ParseEnvVars decodes a JSON object stored in bookmarks.env_vars.
// Empty input yields an empty map. TERM keys are stripped (PTY term owns TERM).
func ParseEnvVars(raw string) (map[string]string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return map[string]string{}, nil
	}
	var m map[string]string
	if err := json.Unmarshal([]byte(raw), &m); err != nil {
		return nil, fmt.Errorf("环境变量 JSON 无效: %w", err)
	}
	if m == nil {
		return map[string]string{}, nil
	}
	out := make(map[string]string, len(m))
	for k, v := range m {
		k = strings.TrimSpace(k)
		if k == "" || strings.EqualFold(k, "TERM") {
			continue
		}
		out[k] = v
	}
	return out, nil
}

// FormatEnvVars encodes env map to JSON for storage. Empty map yields "".
func FormatEnvVars(m map[string]string) (string, error) {
	if len(m) == 0 {
		return "", nil
	}
	clean := make(map[string]string, len(m))
	for k, v := range m {
		k = strings.TrimSpace(k)
		if k == "" || strings.EqualFold(k, "TERM") {
			continue
		}
		clean[k] = v
	}
	if len(clean) == 0 {
		return "", nil
	}
	b, err := json.Marshal(clean)
	if err != nil {
		return "", fmt.Errorf("环境变量编码失败: %w", err)
	}
	return string(b), nil
}
