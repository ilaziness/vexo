package services

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
	"golang.org/x/crypto/ssh"
)

const (
	EventHostKeyPrompt = "eventHostKeyPrompt"
)

func init() {
	application.RegisterEvent[string](EventHostKeyPrompt)
}

// known hosts helpers
func (s *SSHService) knownHostsPath() string {
	if ConfigSvc != nil && ConfigSvc.Config != nil {
		return filepath.Join(ConfigSvc.Config.General.UserDataDir, "known_hosts")
	}
	return filepath.Join("data", "known_hosts")
}

func (s *SSHService) ensureKnownHostsExists(path string) error {
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.WriteFile(path, []byte(""), 0644); err != nil {
			return err
		}
	}
	return nil
}

func (s *SSHService) appendKnownHost(path, host, keyType, keyBase64 string) error {
	s.hostKeyMu.Lock()
	defer s.hostKeyMu.Unlock()
	f, err := os.Open(path)
	if err == nil {
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			tokens := strings.Fields(line)
			if len(tokens) >= 3 && tokens[0] == host && tokens[1] == keyType && tokens[2] == keyBase64 {
				f.Close()
				return nil
			}
		}
		f.Close()
	}
	f2, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f2.Close()
	line := fmt.Sprintf("%s %s %s\n", host, keyType, keyBase64)
	if _, err := f2.WriteString(line); err != nil {
		return err
	}
	return nil
}

// checkKnownHostInFile scans known_hosts file and returns (found, mismatch, error)
func (s *SSHService) checkKnownHostInFile(path, host, keyType, keyBase64 string) (bool, bool, error) {
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

// prepareHostKeyPrompt sets up prompt state and emits event to frontend
func (s *SSHService) prepareHostKeyPrompt(host string, remote net.Addr, fingerprint, keyType, keyBase64 string) error {
	s.hostKeyMu.Lock()
	if s.hostKeyChan != nil {
		s.hostKeyMu.Unlock()
		return fmt.Errorf("another host key prompt in progress")
	}
	s.hostKeyChan = make(chan bool, 1)
	s.pendingHost = host
	s.hostKeyMu.Unlock()

	payload := map[string]string{
		"host":        host,
		"address":     remote.String(),
		"fingerprint": fingerprint,
		"key_type":    keyType,
		"key_base64":  keyBase64,
	}
	data, _ := json.Marshal(payload)
	if app != nil {
		app.Event.Emit(EventHostKeyPrompt, string(data))
	}
	return nil
}

// waitForHostDecision waits for frontend decision, handles cleanup and append
func (s *SSHService) waitForHostDecision(host, knownPath, keyType, keyBase64 string) error {
	select {
	case accept := <-s.hostKeyChan:
		s.hostKeyMu.Lock()
		s.hostKeyChan = nil
		s.pendingHost = ""
		s.hostKeyMu.Unlock()
		if accept {
			return s.appendKnownHost(knownPath, host, keyType, keyBase64)
		}
		return fmt.Errorf("host key not trusted")
	case <-time.After(30 * time.Second):
		s.hostKeyMu.Lock()
		s.hostKeyChan = nil
		s.pendingHost = ""
		s.hostKeyMu.Unlock()
		return fmt.Errorf("host key prompt timeout")
	}
}

// hostKeyCallback prompts frontend when unknown host key encountered
func (s *SSHService) hostKeyCallback(host string, remote net.Addr, key ssh.PublicKey) error {
	keyType := key.Type()
	keyBase64 := base64.StdEncoding.EncodeToString(key.Marshal())
	knownPath := s.knownHostsPath()
	if err := s.ensureKnownHostsExists(knownPath); err != nil {
		return err
	}

	found, mismatch, err := s.checkKnownHostInFile(knownPath, host, keyType, keyBase64)
	if err != nil {
		return err
	}
	if found {
		return nil
	}
	if mismatch {
		return fmt.Errorf("host key mismatch for %s", host)
	}

	// unknown host - prompt frontend
	if err := s.prepareHostKeyPrompt(host, remote, ssh.FingerprintSHA256(key), keyType, keyBase64); err != nil {
		return err
	}
	return s.waitForHostDecision(host, knownPath, keyType, keyBase64)
}

// SetHostKeyDecision is called from frontend to respond to HostKeyPrompt
func (s *SSHService) SetHostKeyDecision(host string, accept bool) error {
	s.hostKeyMu.Lock()
	defer s.hostKeyMu.Unlock()
	if s.pendingHost == "" || s.hostKeyChan == nil {
		return fmt.Errorf("no pending host key prompt")
	}
	if s.pendingHost != host {
		return fmt.Errorf("pending host mismatch")
	}
	select {
	case s.hostKeyChan <- accept:
	default:
	}
	return nil
}
