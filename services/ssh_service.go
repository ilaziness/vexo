package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/ilaziness/vexo/internal/bookmark"
	"github.com/ilaziness/vexo/internal/sftp"
	"github.com/ilaziness/vexo/internal/ssh"
	"github.com/ilaziness/vexo/internal/termws"
	"github.com/ilaziness/vexo/internal/tunnel"
	"github.com/wailsapp/wails/v3/pkg/application"
)

const EventHostKeyPrompt = "eventHostKeyPrompt"

func init() {
	application.RegisterEvent[string](EventHostKeyPrompt)
}

type hostKeyPrompter struct {
	app *application.App
}

func (p *hostKeyPrompter) Prompt(hp ssh.HostKeyPrompt) error {
	data, err := json.Marshal(map[string]string{
		"host": hp.Host, "address": hp.Address, "fingerprint": hp.Fingerprint,
		"key_type": hp.KeyType, "key_base64": hp.KeyBase64,
	})
	if err != nil {
		return err
	}
	p.app.Event.Emit(EventHostKeyPrompt, string(data))
	return nil
}

type SSHService struct {
	app       *application.App
	mgr       *ssh.Manager
	sftp      *sftp.Manager
	tunnels   *tunnel.Manager
	termws    *termws.Server
	bookmarks *bookmark.Service
	closing   sync.Map
}

func NewSSHService(app *application.App, mgr *ssh.Manager, sftpMgr *sftp.Manager, tunnels *tunnel.Manager) *SSHService {
	return &SSHService{app: app, mgr: mgr, sftp: sftpMgr, tunnels: tunnels}
}

func (s *SSHService) bind(term *termws.Server, bookmarks *bookmark.Service) {
	s.termws = term
	s.bookmarks = bookmarks
}

func (s *SSHService) hops(host string, port int, user, password, key, keyPassword, proxyJumpID string) ([]ssh.Endpoint, error) {
	target := ssh.Endpoint{Host: host, Port: port, User: user, Password: password, Key: key, KeyPassword: keyPassword}
	if proxyJumpID == "" {
		return []ssh.Endpoint{target}, nil
	}
	if s.bookmarks == nil {
		return nil, errors.New("bookmark service not initialized")
	}
	return s.bookmarks.ResolveHops(target, proxyJumpID)
}

func (s *SSHService) connectHops(hops []ssh.Endpoint) (string, error) {
	return s.mgr.Connect(hops)
}

func (s *SSHService) testHops(hops []ssh.Endpoint) error {
	return s.mgr.TestConnect(hops)
}

func (s *SSHService) hasSession(id string) bool {
	return s.mgr.HasSession(id)
}

func (s *SSHService) exec(ctx context.Context, linkID, command string, timeout time.Duration, maxOut int) (string, error) {
	if s == nil || s.mgr == nil {
		return "", fmt.Errorf("ssh manager unavailable")
	}
	return s.mgr.Exec(ctx, linkID, command, timeout, maxOut)
}

func (s *SSHService) annotateSession(linkID, notice string) error {
	if s == nil || s.mgr == nil {
		return nil
	}
	return s.mgr.AnnotateSession(linkID, notice)
}

func (s *SSHService) Connect(host string, port int, user, password, key, keyPassword, proxyJumpID string) (string, error) {
	hops, err := s.hops(host, port, user, password, key, keyPassword, proxyJumpID)
	if err != nil {
		return "", err
	}
	return s.connectHops(hops)
}

func (s *SSHService) Start(ID string, cols, rows int) error {
	if err := s.mgr.Start(ID, cols, rows); err != nil {
		_ = s.CloseByID(ID)
		return err
	}
	return nil
}

func (s *SSHService) TestConnectInfo(host string, port int, user, password, key, keyPassword, proxyJumpID string) error {
	hops, err := s.hops(host, port, user, password, key, keyPassword, proxyJumpID)
	if err != nil {
		return err
	}
	return s.testHops(hops)
}

func (s *SSHService) StartSftp(ID string) error {
	return s.sftp.Connect(ID)
}

func (s *SSHService) Resize(ID string, cols, rows int) error {
	return s.mgr.Resize(ID, cols, rows)
}

func (s *SSHService) Close() {
	if s.termws != nil {
		s.termws.CloseAllClients()
	}
	s.tunnels.StopAllAfter(func() {
		s.sftp.CloseAll()
		s.mgr.CloseAll()
	})
}

func (s *SSHService) CloseByID(ID string) error {
	if _, loaded := s.closing.LoadOrStore(ID, struct{}{}); loaded {
		return nil
	}
	defer s.closing.Delete(ID)
	if s.termws != nil {
		s.termws.CloseClient(ID)
	}
	var err error
	s.tunnels.StopSessionAfter(ID, func() {
		s.sftp.CloseSession(ID)
		err = s.mgr.CloseSession(ID)
	})
	return err
}

func (s *SSHService) SelectKeyFile() (string, error) {
	return s.app.Dialog.OpenFile().SetTitle("选择私钥文件").PromptForSingleSelection()
}

func (s *SSHService) GetActiveSessions() []map[string]any {
	return s.mgr.ActiveSessions()
}

func (s *SSHService) SendToSession(sessionID, command string) error {
	return s.mgr.SendToSession(sessionID, command)
}

func (s *SSHService) SetHostKeyDecision(host string, accept bool) error {
	return s.mgr.SetHostKeyDecision(host, accept)
}

func (s *SSHService) GetOrFetchRemoteSystemInfo(linkID, host string) *ssh.RemoteSystemInfo {
	return s.mgr.GetOrFetchRemoteSystemInfo(linkID, host)
}
