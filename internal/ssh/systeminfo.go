package ssh

import (
	"context"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	cryptossh "golang.org/x/crypto/ssh"
)

const remoteSystemInfoScript = `export LC_ALL=C; hostname 2>/dev/null; uname -sr 2>/dev/null; uname -m 2>/dev/null; if [ -r /etc/os-release ]; then . /etc/os-release; echo PRETTY_NAME=$PRETTY_NAME; echo ID=$ID; echo VERSION_ID=$VERSION_ID; fi; exit 0`

func (m *Manager) GetOrFetchRemoteSystemInfo(linkID, host string) *RemoteSystemInfo {
	key := NormalizeHost(host)
	if key == "" {
		return nil
	}
	if cached, ok := m.remoteInfoCache.Load(key); ok {
		return cached.(*RemoteSystemInfo)
	}
	mu := m.remoteInfoFetchMu(key)
	mu.Lock()
	defer mu.Unlock()
	if cached, ok := m.remoteInfoCache.Load(key); ok {
		return cached.(*RemoteSystemInfo)
	}
	info := m.fetchRemoteSystemInfo(linkID)
	m.remoteInfoCache.Store(key, info)
	return info
}

func (m *Manager) remoteInfoFetchMu(key string) *sync.Mutex {
	mu, _ := m.remoteInfoFetchLocks.LoadOrStore(key, &sync.Mutex{})
	return mu.(*sync.Mutex)
}

func (m *Manager) fetchRemoteSystemInfo(linkID string) *RemoteSystemInfo {
	client, err := m.GetClient(linkID)
	if err != nil {
		return &RemoteSystemInfo{Ready: true}
	}
	output, err := runRemoteCommand(client, remoteSystemInfoScript)
	if err != nil && strings.TrimSpace(output) == "" {
		m.logger.Debug("fetch remote system info failed", zap.String("linkID", linkID), zap.Error(err))
		return &RemoteSystemInfo{Ready: true}
	}
	info := ParseRemoteSystemInfo(output)
	return &info
}

func runRemoteCommand(client *cryptossh.Client, command string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return runRemoteCommandCtx(ctx, client, command, DefaultExecMaxOut)
}
