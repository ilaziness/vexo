package services

import (
	"sync"
)

// SyncProgress 同步进度
type SyncProgress struct {
	Stage       string  `json:"stage"`       // 当前阶段: packing, encrypting, uploading, downloading, decrypting, unpacking
	Progress    float64 `json:"progress"`    // 进度百分比 0-100
	TotalBytes  int64   `json:"totalBytes"`  // 总字节数
	DoneBytes   int64   `json:"doneBytes"`   // 已完成字节数
	Speed       float64 `json:"speed"`       // 速度 bytes/s
	ETA         int     `json:"eta"`         // 预计剩余时间(秒)
	IsCompleted bool    `json:"isCompleted"` // 是否完成
	Error       string  `json:"error"`       // 错误信息
}

// ProgressReporter 进度报告器
type ProgressReporter struct {
	mu        sync.RWMutex
	progress  SyncProgress
	onUpdate  func(SyncProgress)
}

// NewProgressReporter 创建进度报告器
func NewProgressReporter(onUpdate func(SyncProgress)) *ProgressReporter {
	return &ProgressReporter{
		onUpdate: onUpdate,
	}
}

// UpdateProgress 更新进度
func (p *ProgressReporter) UpdateProgress(stage string, done, total int64) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.progress.Stage = stage
	p.progress.DoneBytes = done
	p.progress.TotalBytes = total
	if total > 0 {
		p.progress.Progress = float64(done) / float64(total) * 100
	}

	if p.onUpdate != nil {
		p.onUpdate(p.progress)
	}
}

// SetStage 设置阶段
func (p *ProgressReporter) SetStage(stage string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.progress.Stage = stage
	if p.onUpdate != nil {
		p.onUpdate(p.progress)
	}
}

// SetCompleted 设置完成
func (p *ProgressReporter) SetCompleted() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.progress.IsCompleted = true
	p.progress.Progress = 100
	if p.onUpdate != nil {
		p.onUpdate(p.progress)
	}
}

// SetError 设置错误
func (p *ProgressReporter) SetError(err string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.progress.Error = err
	if p.onUpdate != nil {
		p.onUpdate(p.progress)
	}
}

// GetProgress 获取当前进度
func (p *ProgressReporter) GetProgress() SyncProgress {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.progress
}

// GlobalProgressReporter 全局进度报告器
var GlobalProgressReporter = NewProgressReporter(nil)

// SetProgressCallback 设置进度回调
func SetProgressCallback(callback func(SyncProgress)) {
	GlobalProgressReporter.onUpdate = callback
}
