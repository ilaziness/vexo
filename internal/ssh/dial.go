package ssh

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"go.uber.org/zap"
	cryptossh "golang.org/x/crypto/ssh"
)

func (m *Manager) dialHops(hops []Endpoint, timeout time.Duration) (*hopClient, error) {
	const maxDepth = 5
	if len(hops) == 0 {
		return nil, fmt.Errorf("empty hop list")
	}
	if len(hops) > maxDepth+1 {
		return nil, fmt.Errorf("跳板机层数超过最大限制 (%d)", maxDepth)
	}
	seen := make(map[string]bool)
	var jumps []*cryptossh.Client
	var prev *cryptossh.Client
	closeOpened := func() {
		if prev != nil {
			_ = prev.Close()
		}
		for i := len(jumps) - 1; i >= 0; i-- {
			_ = jumps[i].Close()
		}
	}
	for i, hop := range hops {
		addr := hop.Addr()
		if seen[addr] {
			closeOpened()
			return nil, fmt.Errorf("检测到跳板机循环引用: %s", addr)
		}
		seen[addr] = true
		if i > 0 && hops[i-1].Host == hop.Host && hops[i-1].Port == hop.Port {
			closeOpened()
			return nil, fmt.Errorf("跳板机不能与目标主机相同")
		}
		client, err := m.dialOne(hop, timeout, prev)
		if err != nil {
			closeOpened()
			return nil, err
		}
		if prev != nil {
			jumps = append(jumps, prev)
		}
		prev = client
	}
	return &hopClient{target: prev, jumps: jumps}, nil
}

func (m *Manager) dialOne(ep Endpoint, timeout time.Duration, via *cryptossh.Client) (*cryptossh.Client, error) {
	if ep.Password == "" && ep.Key == "" {
		return nil, fmt.Errorf("empty password and key")
	}
	cfg := &cryptossh.ClientConfig{
		User:            ep.User,
		Auth:            []cryptossh.AuthMethod{},
		HostKeyCallback: m.hostKeyCallback,
		Timeout:         timeout,
	}
	if ep.Key != "" {
		keyContent, err := os.ReadFile(ep.Key)
		if err != nil {
			return nil, fmt.Errorf("unable to read private key: %v", err)
		}
		var signer cryptossh.Signer
		if ep.KeyPassword == "" {
			signer, err = cryptossh.ParsePrivateKey(keyContent)
		} else {
			signer, err = cryptossh.ParsePrivateKeyWithPassphrase(keyContent, []byte(ep.KeyPassword))
		}
		if err != nil {
			return nil, err
		}
		cfg.Auth = []cryptossh.AuthMethod{cryptossh.PublicKeys(signer)}
	}
	if ep.Password != "" {
		cfg.Auth = append(cfg.Auth, cryptossh.Password(ep.Password))
	}

	m.logger.Debug("ssh key", zap.String("file", ep.Key))
	addr := net.JoinHostPort(ep.Host, fmt.Sprintf("%d", ep.Port))

	var conn net.Conn
	var err error
	isProxyConn := via != nil
	if via != nil {
		m.logger.Debug("Connecting via ProxyJump", zap.String("proxyHost", ep.Host), zap.Int("proxyPort", ep.Port))
		dialCtx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		conn, err = via.DialContext(dialCtx, "tcp", addr)
		if err != nil {
			m.logger.Debug("proxy dial tcp error", zap.Error(err))
			return nil, err
		}
	} else {
		conn, err = net.DialTimeout("tcp", addr, timeout)
		if err != nil {
			m.logger.Debug("tcp connect error", zap.Error(err))
			return nil, err
		}
	}

	if !isProxyConn {
		if err = conn.SetDeadline(time.Now().Add(timeout)); err != nil {
			conn.Close()
			m.logger.Debug("set deadline error", zap.Error(err))
			return nil, err
		}
	}
	c, chans, reqs, err := cryptossh.NewClientConn(conn, addr, cfg)
	if err != nil {
		conn.Close()
		m.logger.Debug("ssh handshake error", zap.Error(err))
		return nil, err
	}
	if !isProxyConn {
		if err := conn.SetDeadline(time.Time{}); err != nil {
			c.Close()
			m.logger.Debug("set deadline error", zap.Error(err))
			return nil, err
		}
	}
	return cryptossh.NewClient(c, chans, reqs), nil
}
