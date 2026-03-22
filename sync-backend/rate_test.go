package main

import (
	"testing"
	"time"
)

func TestNewRateLimiter(t *testing.T) {
	// 创建 mock service
	service := &SyncService{}
	limiter := NewRateLimiter(service, time.Minute)

	if limiter == nil {
		t.Fatal("rate limiter should not be nil")
	}
	if limiter.minInterval != time.Minute {
		t.Errorf("expected interval 1m, got %v", limiter.minInterval)
	}
}

func TestRateLimiterCheck(t *testing.T) {
	// 注意：这个测试需要数据库连接，这里仅作示例
	// 实际测试应该使用 mock
	t.Skip("Skipping test that requires database connection")
}

func TestDefaultRateLimiter(t *testing.T) {
	service := &SyncService{}
	limiter := DefaultRateLimiter(service)

	if limiter == nil {
		t.Fatal("default rate limiter should not be nil")
	}
	if limiter.minInterval != time.Minute {
		t.Errorf("expected default interval 1m, got %v", limiter.minInterval)
	}
}
