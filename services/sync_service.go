package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ilaziness/vexo/internal/sync"
	"go.uber.org/zap"
)

// ErrSyncNotConfigured 同步未配置错误
var ErrSyncNotConfigured = errors.New("sync not configured")

// SyncService 同步服务
type SyncService struct {
	configService *ConfigService
}

// NewSyncService 创建同步服务
func NewSyncService(configService *ConfigService) *SyncService {
	return &SyncService{
		configService: configService,
	}
}

// GetSyncConfig 获取同步配置
func (s *SyncService) GetSyncConfig() (*sync.SyncConfig, error) {
	config := s.configService.GetSyncConfig()
	return config, nil
}

// SaveSyncConfig 保存同步配置
func (s *SyncService) SaveSyncConfig(config sync.SyncConfig) error {
	return s.configService.SaveSyncConfig(config)
}

// GetSyncProgress 获取同步进度
func (s *SyncService) GetSyncProgress() SyncProgress {
	return GlobalProgressReporter.GetProgress()
}

// UploadSync 上传同步数据
func (s *SyncService) UploadSync() error {
	config := s.configService.GetSyncConfig()
	if !config.IsConfigured() {
		return ErrSyncNotConfigured
	}

	// 重置进度
	GlobalProgressReporter.UpdateProgress("preparing", 0, 0)

	manager := sync.NewSyncManager(config)
	if !manager.IsConfigured() {
		return ErrSyncNotConfigured
	}

	userDataDir := s.configService.Config.General.UserDataDir
	
	// 使用带重试的上传
	err := s.uploadWithRetry(manager, userDataDir, 3)
	if err != nil {
		GlobalProgressReporter.SetError(err.Error())
		return err
	}

	GlobalProgressReporter.SetCompleted()
	return nil
}

// uploadWithRetry 带重试的上传
func (s *SyncService) uploadWithRetry(manager *sync.SyncManager, userDataDir string, maxRetries int) error {
	var lastErr error
	
	for i := 0; i < maxRetries; i++ {
		if i > 0 {
			GlobalProgressReporter.SetStage(fmt.Sprintf("retrying (%d/%d)", i, maxRetries))
			time.Sleep(time.Second * time.Duration(i))
		}

		GlobalProgressReporter.SetStage("packing")
		err := manager.Upload(userDataDir)
		if err == nil {
			return nil
		}
		
		lastErr = err
		Logger.Error("Upload failed, will retry", 
			zap.Int("attempt", i+1), 
			zap.Int("maxRetries", maxRetries), 
			zap.Error(err))
	}
	
	return fmt.Errorf("upload failed after %d attempts: %w", maxRetries, lastErr)
}

// DownloadSync 下载同步数据
// 注意：恢复成功或失败后都需要重启应用，因为数据库连接已关闭
func (s *SyncService) DownloadSync(version int) error {
	config := s.configService.GetSyncConfig()
	if !config.IsConfigured() {
		return ErrSyncNotConfigured
	}

	// 重置进度
	GlobalProgressReporter.UpdateProgress("preparing", 0, 0)

	manager := sync.NewSyncManager(config)
	if !manager.IsConfigured() {
		return ErrSyncNotConfigured
	}

	userDataDir := s.configService.Config.General.UserDataDir
	
	// 使用带重试的下载
	err := s.downloadWithRetry(manager, userDataDir, version, 3)
	if err != nil {
		GlobalProgressReporter.SetError(err.Error())
		// 返回特殊错误标记，表示需要重启（数据库已关闭）
		return fmt.Errorf("%w: 数据恢复失败，请手动重启应用", err)
	}

	GlobalProgressReporter.SetCompleted()
	return nil
}

// downloadWithRetry 带重试的下载
func (s *SyncService) downloadWithRetry(manager *sync.SyncManager, userDataDir string, version int, maxRetries int) error {
	var lastErr error
	
	for i := 0; i < maxRetries; i++ {
		if i > 0 {
			GlobalProgressReporter.SetStage(fmt.Sprintf("retrying (%d/%d)", i, maxRetries))
			time.Sleep(time.Second * time.Duration(i))
		}

		GlobalProgressReporter.SetStage("downloading")
		err := manager.Download(userDataDir, version, DB)
		if err == nil {
			return nil
		}
		
		lastErr = err
		Logger.Error("Download failed, will retry", 
			zap.Int("attempt", i+1), 
			zap.Int("maxRetries", maxRetries), 
			zap.Error(err))
	}
	
	return fmt.Errorf("download failed after %d attempts: %w", maxRetries, lastErr)
}

// ListSyncVersions 列出同步版本（支持分页，limit=0表示不分页）
func (s *SyncService) ListSyncVersions(limit, offset int) ([]sync.VersionInfo, error) {
	config := s.configService.GetSyncConfig()
	if !config.IsConfigured() {
		return nil, ErrSyncNotConfigured
	}

	manager := sync.NewSyncManager(config)
	if !manager.IsConfigured() {
		return nil, ErrSyncNotConfigured
	}

	return manager.ListVersions(limit, offset)
}

// DeleteSyncVersion 删除指定同步版本
func (s *SyncService) DeleteSyncVersion(version int) error {
	config := s.configService.GetSyncConfig()
	if !config.IsConfigured() {
		return ErrSyncNotConfigured
	}

	manager := sync.NewSyncManager(config)
	if !manager.IsConfigured() {
		return ErrSyncNotConfigured
	}

	return manager.DeleteVersion(version)
}

// HealthCheck 健康检查
func (s *SyncService) HealthCheck() error {
	config := s.configService.GetSyncConfig()
	if !config.IsConfigured() {
		return ErrSyncNotConfigured
	}

	manager := sync.NewSyncManager(config)
	if !manager.IsConfigured() {
		return ErrSyncNotConfigured
	}

	return manager.HealthCheck()
}

// CancelSync 取消同步
func (s *SyncService) CancelSync() {
	// 设置取消标志
	GlobalProgressReporter.SetError("cancelled by user")
}

// context for cancellation
type syncContext struct {
	ctx    context.Context
	cancel context.CancelFunc
}

var currentSync *syncContext

// StartSyncContext 开始同步上下文
func (s *SyncService) StartSyncContext() {
	ctx, cancel := context.WithCancel(context.Background())
	currentSync = &syncContext{ctx, cancel}
}

// CancelSyncContext 取消同步上下文
func (s *SyncService) CancelSyncContext() {
	if currentSync != nil && currentSync.cancel != nil {
		currentSync.cancel()
	}
}
