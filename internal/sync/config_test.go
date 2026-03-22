package sync

import (
	"testing"
)

func TestSyncConfigIsConfigured(t *testing.T) {
	tests := []struct {
		name     string
		config   SyncConfig
		expected bool
	}{
		{
			name:     "empty config",
			config:   SyncConfig{},
			expected: false,
		},
		{
			name: "only server url",
			config: SyncConfig{
				ServerURL: "http://localhost:8080",
			},
			expected: false,
		},
		{
			name: "only sync id",
			config: SyncConfig{
				SyncID: "test-id",
			},
			expected: false,
		},
		{
			name: "only user key",
			config: SyncConfig{
				UserKey: "test-key",
			},
			expected: false,
		},
		{
			name: "missing server url",
			config: SyncConfig{
				SyncID:  "test-id",
				UserKey: "test-key",
			},
			expected: false,
		},
		{
			name: "missing sync id",
			config: SyncConfig{
				ServerURL: "http://localhost:8080",
				UserKey:   "test-key",
			},
			expected: false,
		},
		{
			name: "missing user key",
			config: SyncConfig{
				ServerURL: "http://localhost:8080",
				SyncID:    "test-id",
			},
			expected: false,
		},
		{
			name: "fully configured",
			config: SyncConfig{
				ServerURL: "http://localhost:8080",
				SyncID:    "test-id",
				UserKey:   "test-key",
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.IsConfigured()
			if result != tt.expected {
				t.Errorf("IsConfigured() = %v, expected %v", result, tt.expected)
			}
		})
	}
}
