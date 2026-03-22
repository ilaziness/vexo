package main

import (
	"fmt"
	"time"
)

// RateLimiter 频率限制器
type RateLimiter struct {
	service   *SyncService
	minInterval time.Duration
}

// NewRateLimiter 创建频率限制器
func NewRateLimiter(service *SyncService, minInterval time.Duration) *RateLimiter {
	return &RateLimiter{
		service:     service,
		minInterval: minInterval,
	}
}

// Check 检查是否可以同步
func (r *RateLimiter) Check(userID string) error {
	lastSync, err := r.service.GetLastSyncAt(userID)
	if err != nil {
		// 如果获取失败（可能是首次同步），允许同步
		return nil
	}

	elapsed := time.Since(lastSync)
	if elapsed < r.minInterval {
		remaining := r.minInterval - elapsed
		return fmt.Errorf("sync too frequent, please wait %v", remaining.Round(time.Second))
	}

	return nil
}

// DefaultRateLimiter 创建默认频率限制器（1分钟）
func DefaultRateLimiter(service *SyncService) *RateLimiter {
	return NewRateLimiter(service, time.Minute)
}
