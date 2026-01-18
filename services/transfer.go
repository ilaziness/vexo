package services

import (
	"context"
	"errors"
	"io"
	"math"
	"sync"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"

	"github.com/ilaziness/vexo/internal/utils"
)

const (
	EventProgress = "eventProgress"

	TransferTypeUpload   = "upload"
	TransferTypeDownload = "download"
)

func init() {
	application.RegisterEvent[ProgressData](EventProgress)
}

var activeTransfers = new(sync.Map) // map[transferID]*transferTracker

type ProgressData struct {
	ID           string // 传输ID
	SessionID    string
	TransferType string // 传输类型
	LocalFile    string // 本地文件/目录，包含完整路径
	RemoteFile   string // 远程文件/目录，包含完整路径
	TotalSize    int64
	Rate         float64
	Done         bool
	Error        string
}

// copyWithContext copies from src to dst using the provided context for cancellation
func copyWithContext(dst io.Writer, src io.Reader, ctx context.Context) (written int64, err error) {
	buf := make([]byte, 4*1024)
	for {
		select {
		case <-ctx.Done():
			if errors.Is(ctx.Err(), context.Canceled) {
				return written, errors.New("user cancelled")
			}
			return written, ctx.Err()
		default:
		}
		n, err := src.Read(buf)
		if n > 0 {
			if nw, err := dst.Write(buf[:n]); err != nil {
				return written + int64(nw), err
			}
			written += int64(n)
		}
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return written, err
		}
	}
}

type transferTracker struct {
	sessionID    string
	id           string
	transferType string
	localFile    string
	remoteFile   string
	total        int64
	transferred  int64
	mutex        sync.Mutex
	ticker       *time.Ticker
	done         chan struct{}
	ctx          context.Context
	cancelFunc   context.CancelFunc
}

func newTransferTracker(sessionID, transferType, localFile, remoteFile string, total int64) *transferTracker {
	ctx, cancelFunc := context.WithCancel(context.Background())
	tracker := &transferTracker{
		sessionID:    sessionID,
		id:           utils.GenerateRandomID(),
		transferType: transferType,
		localFile:    localFile,
		remoteFile:   remoteFile,
		total:        total,
		done:         make(chan struct{}),
		ctx:          ctx,
		cancelFunc:   cancelFunc,
	}
	activeTransfers.Store(tracker.id, tracker)
	tracker.startProgress()
	return tracker
}

func (t *transferTracker) update(n int64) {
	t.mutex.Lock()
	t.transferred += n
	t.mutex.Unlock()
}

func (t *transferTracker) getRate() float64 {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	if t.total == 0 {
		return 100.0
	}
	// 计算进度百分比
	rate := float64(t.transferred) * 100.0 / float64(t.total)
	if rate > 100.0 {
		rate = 100.0
	}
	// 精确到2位小数
	return math.Round(rate*100) / 100
}

func (t *transferTracker) startProgress() {
	// Emit initial progress with rate 0
	app.Event.Emit(EventProgress, ProgressData{
		ID:           t.id,
		SessionID:    t.sessionID,
		TransferType: t.transferType,
		LocalFile:    t.localFile,
		RemoteFile:   t.remoteFile,
		TotalSize:    t.total,
		Rate:         0.0,
	})

	t.ticker = time.NewTicker(500 * time.Millisecond)
	go func() {
		defer t.ticker.Stop()
		for {
			select {
			case <-t.ticker.C:
				rate := t.getRate()
				app.Event.Emit(EventProgress, ProgressData{
					ID:           t.id,
					SessionID:    t.sessionID,
					TransferType: t.transferType,
					LocalFile:    t.localFile,
					RemoteFile:   t.remoteFile,
					TotalSize:    t.total,
					Rate:         rate,
				})
				if rate >= 100.0 {
					close(t.done)
					return
				}
			case <-t.done:
				return
			}
		}
	}()
}

func (t *transferTracker) stopProgress(err error) {
	if t.ticker != nil {
		t.ticker.Stop()
	}
	select {
	case <-t.done:
		// already closed
	default:
		close(t.done)
	}
	// Send final progress with Done: true
	rate := t.getRate()
	progressData := ProgressData{
		ID:           t.id,
		SessionID:    t.sessionID,
		TransferType: t.transferType,
		LocalFile:    t.localFile,
		RemoteFile:   t.remoteFile,
		TotalSize:    t.total,
		Rate:         rate,
		Done:         true,
	}
	if err != nil {
		progressData.Error = err.Error()
	}
	app.Event.Emit(EventProgress, progressData)
	activeTransfers.Delete(t.id)
}

type progressWriter struct {
	writer  io.Writer
	tracker *transferTracker
}

func (pw *progressWriter) Write(p []byte) (int, error) {
	n, err := pw.writer.Write(p)
	if n > 0 {
		pw.tracker.update(int64(n))
	}
	return n, err
}
