package services

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/ilaziness/vexo/internal/system"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
)

const remoteSystemInfoScript = `export LC_ALL=C; hostname 2>/dev/null; uname -sr 2>/dev/null; uname -m 2>/dev/null; if [ -r /etc/os-release ]; then . /etc/os-release; echo PRETTY_NAME=$PRETTY_NAME; echo ID=$ID; echo VERSION_ID=$VERSION_ID; fi; exit 0`

// GetOrFetchRemoteSystemInfo 按 host 查缓存；未命中则经 linkID 远程采集并缓存
func (s *SSHService) GetOrFetchRemoteSystemInfo(linkID, host string) *system.RemoteSystemInfo {
	key := system.NormalizeHost(host)
	if key == "" {
		return nil
	}

	if cached, ok := s.remoteInfoCache.Load(key); ok {
		return cached.(*system.RemoteSystemInfo)
	}

	mu := s.remoteInfoFetchMu(key)
	mu.Lock()
	defer mu.Unlock()

	if cached, ok := s.remoteInfoCache.Load(key); ok {
		return cached.(*system.RemoteSystemInfo)
	}

	info := s.fetchRemoteSystemInfo(linkID)
	s.remoteInfoCache.Store(key, info)
	return info
}

func (s *SSHService) remoteInfoFetchMu(key string) *sync.Mutex {
	mu, _ := s.remoteInfoFetchLocks.LoadOrStore(key, &sync.Mutex{})
	return mu.(*sync.Mutex)
}

func (s *SSHService) fetchRemoteSystemInfo(linkID string) *system.RemoteSystemInfo {
	connAny, ok := s.SSHConnects.Load(linkID)
	if !ok {
		return &system.RemoteSystemInfo{Ready: true}
	}
	conn := connAny.(*SSHConnect)
	if conn.client == nil {
		return &system.RemoteSystemInfo{Ready: true}
	}

	output, err := runRemoteCommand(conn.client, remoteSystemInfoScript)
	if err != nil && strings.TrimSpace(output) == "" {
		Logger.Debug("fetch remote system info failed",
			zap.String("linkID", linkID),
			zap.Error(err),
		)
		return &system.RemoteSystemInfo{Ready: true}
	}

	info := system.ParseRemoteSystemInfo(output)
	return &info
}

func runRemoteCommand(client *ssh.Client, command string) (string, error) {
	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	type cmdResult struct {
		out []byte
		err error
	}
	done := make(chan cmdResult, 1)
	go func() {
		out, runErr := session.CombinedOutput(command)
		done <- cmdResult{out: out, err: runErr}
	}()

	select {
	case r := <-done:
		if r.err != nil && len(r.out) == 0 {
			return "", r.err
		}
		return string(r.out), r.err
	case <-ctx.Done():
		_ = session.Close()
		return "", ctx.Err()
	}
}
