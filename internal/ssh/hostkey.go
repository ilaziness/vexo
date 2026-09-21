package ssh

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	cryptossh "golang.org/x/crypto/ssh"
)

type hostKeyStore struct {
	mu      sync.Mutex
	pending map[string]chan bool
}

func newHostKeyStore() *hostKeyStore {
	return &hostKeyStore{pending: make(map[string]chan bool)}
}

func (s *hostKeyStore) begin(host string) (chan bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.pending[host]; ok {
		return nil, fmt.Errorf("host key prompt already in progress for %s", host)
	}
	ch := make(chan bool, 1)
	s.pending[host] = ch
	return ch, nil
}

func (s *hostKeyStore) finish(host string) {
	s.mu.Lock()
	delete(s.pending, host)
	s.mu.Unlock()
}

func (s *hostKeyStore) decide(host string, accept bool) error {
	s.mu.Lock()
	ch := s.pending[host]
	s.mu.Unlock()
	if ch == nil {
		return fmt.Errorf("no pending host key prompt")
	}
	select {
	case ch <- accept:
	default:
	}
	return nil
}

func (m *Manager) knownHostsFile() string {
	if m.knownHostsPath != "" {
		return m.knownHostsPath
	}
	return filepath.Join("data", "known_hosts")
}

func ensureKnownHostsExists(path string) error {
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0700); err != nil {
			return err
		}
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.WriteFile(path, []byte(""), 0600)
	}
	return nil
}

func appendKnownHost(path, host, keyType, keyBase64 string) error {
	f, err := os.Open(path)
	if err == nil {
		defer f.Close()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			tokens := strings.Fields(line)
			if len(tokens) >= 3 && tokens[0] == host && tokens[1] == keyType && tokens[2] == keyBase64 {
				return nil
			}
		}
		if err := scanner.Err(); err != nil {
			return err
		}
	}
	f2, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f2.Close()
	_, err = f2.WriteString(fmt.Sprintf("%s %s %s\n", host, keyType, keyBase64))
	return err
}

func checkKnownHostInFile(path, host, keyType, keyBase64 string) (bool, bool, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, false, nil
		}
		return false, false, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		tokens := strings.Fields(line)
		if len(tokens) >= 3 && tokens[0] == host {
			if tokens[1] == keyType && tokens[2] == keyBase64 {
				return true, false, nil
			}
			return false, true, nil
		}
	}
	if err := scanner.Err(); err != nil {
		return false, false, err
	}
	return false, false, nil
}

func (m *Manager) hostKeyCallback(host string, remote net.Addr, key cryptossh.PublicKey) error {
	keyType := key.Type()
	keyBase64 := base64.StdEncoding.EncodeToString(key.Marshal())
	knownPath := m.knownHostsFile()
	m.knownHostsMu.Lock()
	if err := ensureKnownHostsExists(knownPath); err != nil {
		m.knownHostsMu.Unlock()
		return err
	}
	found, mismatch, err := checkKnownHostInFile(knownPath, host, keyType, keyBase64)
	m.knownHostsMu.Unlock()
	if err != nil {
		return err
	}
	if found {
		return nil
	}
	if mismatch {
		return fmt.Errorf("host key mismatch for %s", host)
	}
	if m.prompter == nil {
		return fmt.Errorf("host key not trusted")
	}

	ch, err := m.hostKey.begin(host)
	if err != nil {
		return err
	}
	defer m.hostKey.finish(host)

	if err := m.prompter.Prompt(HostKeyPrompt{
		Host:        host,
		Address:     remote.String(),
		Fingerprint: cryptossh.FingerprintSHA256(key),
		KeyType:     keyType,
		KeyBase64:   keyBase64,
	}); err != nil {
		return err
	}

	select {
	case accept := <-ch:
		if accept {
			m.knownHostsMu.Lock()
			err := appendKnownHost(knownPath, host, keyType, keyBase64)
			m.knownHostsMu.Unlock()
			return err
		}
		return fmt.Errorf("host key not trusted")
	case <-time.After(30 * time.Second):
		return fmt.Errorf("host key prompt timeout")
	}
}
