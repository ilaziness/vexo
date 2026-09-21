package ssh

import (
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/ilaziness/vexo/internal/system"
	"github.com/ilaziness/vexo/internal/utils"
	"go.uber.org/zap"
	cryptossh "golang.org/x/crypto/ssh"
)

const ErrConnectionNotFound = "SSH connection with ID %s not found"

type Endpoint struct {
	Host        string
	Port        int
	User        string
	Password    string
	Key         string
	KeyPassword string
}

func (e Endpoint) Addr() string {
	return fmt.Sprintf("%s:%d", e.Host, e.Port)
}

func (e Endpoint) ClientKey() string {
	return fmt.Sprintf("%s@%s:%d", e.User, e.Host, e.Port)
}

type HostKeyPrompt struct {
	Host        string
	Address     string
	Fingerprint string
	KeyType     string
	KeyBase64   string
}

type HostKeyPrompter interface {
	Prompt(p HostKeyPrompt) error
}

type Session struct {
	ID             string
	ClientKey      string
	client         *cryptossh.Client
	session        *cryptossh.Session
	isClosed       bool
	closeMu        sync.Mutex
	OutputChan     chan []byte
	outputBuffSize int
	outputWg       sync.WaitGroup
	outputOnce     sync.Once
	Stdin          io.WriteCloser
	stopOutput     chan struct{}
	stopOutputOnce sync.Once
	manager        *Manager
}

type hopClient struct {
	target *cryptossh.Client
	jumps  []*cryptossh.Client
}

func (h *hopClient) Close() {
	if h == nil {
		return
	}
	if h.target != nil {
		_ = h.target.Close()
	}
	for i := len(h.jumps) - 1; i >= 0; i-- {
		_ = h.jumps[i].Close()
	}
}

type Manager struct {
	logger               *zap.Logger
	knownHostsPath       string
	knownHostsMu         sync.Mutex
	prompter             HostKeyPrompter
	onClose              func(sessionID string)
	clients              *sync.Map
	sessions             *sync.Map
	clientLocks          sync.Map
	remoteInfoCache      sync.Map
	remoteInfoFetchLocks sync.Map
	hostKey              *hostKeyStore
}

func NewManager(logger *zap.Logger, knownHostsPath string, prompter HostKeyPrompter) *Manager {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &Manager{
		logger:         logger,
		knownHostsPath: knownHostsPath,
		prompter:       prompter,
		clients:        new(sync.Map),
		sessions:       new(sync.Map),
		hostKey:        newHostKeyStore(),
	}
}

func (m *Manager) SetOnClose(fn func(sessionID string)) {
	m.onClose = fn
}

func (m *Manager) GetSession(id string) (*Session, error) {
	v, ok := m.sessions.Load(id)
	if !ok {
		return nil, fmt.Errorf(ErrConnectionNotFound, id)
	}
	return v.(*Session), nil
}

func (m *Manager) GetClient(sessionID string) (*cryptossh.Client, error) {
	sess, err := m.GetSession(sessionID)
	if err != nil {
		return nil, err
	}
	if sess.client == nil {
		return nil, fmt.Errorf(ErrConnectionNotFound, sessionID)
	}
	return sess.client, nil
}

func (m *Manager) HasSession(id string) bool {
	_, ok := m.sessions.Load(id)
	return ok
}

func (m *Manager) clientLock(key string) *sync.Mutex {
	mu, _ := m.clientLocks.LoadOrStore(key, &sync.Mutex{})
	return mu.(*sync.Mutex)
}

func hopsClientKey(hops []Endpoint) string {
	if len(hops) == 0 {
		return ""
	}
	key := hops[len(hops)-1].ClientKey()
	for i := 0; i < len(hops)-1; i++ {
		key += fmt.Sprintf("via:%s@%s", hops[i].User, hops[i].Addr())
	}
	return key
}

func (m *Manager) Connect(hops []Endpoint) (string, error) {
	if len(hops) == 0 {
		return "", fmt.Errorf("empty hop list")
	}
	target := hops[len(hops)-1]
	m.logger.Debug("Connecting to SSH server", zap.String("host", target.Host), zap.Int("port", target.Port))

	clientKey := hopsClientKey(hops)
	lock := m.clientLock(clientKey)
	lock.Lock()
	defer lock.Unlock()

	var client *cryptossh.Client
	if v, ok := m.clients.Load(clientKey); ok {
		m.logger.Debug("Using existing SSH client", zap.String("clientKey", clientKey))
		client = v.(*hopClient).target
	} else {
		hc, err := m.dialHops(hops, 30*time.Second)
		if err != nil {
			return "", err
		}
		m.clients.Store(clientKey, hc)
		client = hc.target
		m.logger.Debug("ssh connect ok and stored in cache", zap.String("clientKey", clientKey))
	}

	sess := newSession(m, clientKey, client)
	m.sessions.Store(sess.ID, sess)
	return sess.ID, nil
}

func (m *Manager) TestConnect(hops []Endpoint) error {
	if len(hops) == 0 {
		return fmt.Errorf("empty hop list")
	}
	target := hops[len(hops)-1]
	m.logger.Debug("Testing SSH connection", zap.String("host", target.Host), zap.Int("port", target.Port))
	hc, err := m.dialHops(hops, 20*time.Second)
	if err != nil {
		return err
	}
	defer hc.Close()
	m.logger.Debug("Test SSH connect successful", zap.String("host", target.Host), zap.Int("port", target.Port))
	return nil
}

func (m *Manager) Start(id string, cols, rows int) error {
	m.logger.Debug("Starting SSH connection", zap.String("id", id))
	sess, err := m.GetSession(id)
	if err != nil {
		return err
	}
	if err := sess.Start(cols, rows); err != nil {
		return err
	}
	return nil
}

func (m *Manager) Resize(id string, cols, rows int) error {
	sess, err := m.GetSession(id)
	if err != nil {
		return err
	}
	return sess.Resize(cols, rows)
}

func (m *Manager) SendToSession(sessionID, command string) error {
	sess, err := m.GetSession(sessionID)
	if err != nil {
		return fmt.Errorf("SSH session %s not found", sessionID)
	}
	if sess.Stdin == nil {
		return fmt.Errorf("stdin not available")
	}
	_, err = sess.Stdin.Write([]byte(command + "\n"))
	if err != nil {
		m.logger.Error("Failed to write to stdin", zap.Error(err))
	}
	return err
}

func (m *Manager) ActiveSessions() []map[string]any {
	sessions := make([]map[string]any, 0)
	m.sessions.Range(func(_, value any) bool {
		sess := value.(*Session)
		sessions = append(sessions, map[string]any{
			"id":        sess.ID,
			"clientKey": sess.ClientKey,
		})
		return true
	})
	return sessions
}

func (m *Manager) notifyClosed(id string) {
	if m.onClose != nil {
		m.onClose(id)
		return
	}
	_ = m.CloseSession(id)
}

func (m *Manager) CloseAll() {
	m.logger.Debug("Close ssh all")
	m.sessions.Range(func(_, value any) bool {
		_ = value.(*Session).closePTY()
		return true
	})
	m.clients.Range(func(key, value any) bool {
		value.(*hopClient).Close()
		m.clients.Delete(key)
		return true
	})
	m.sessions = new(sync.Map)
}

func (m *Manager) CloseSession(id string) error {
	m.logger.Debug("CloseByID", zap.String("ID", id))
	sess, err := m.GetSession(id)
	if err != nil {
		return err
	}
	clientKey := sess.ClientKey
	_ = sess.closePTY()
	m.sessions.Delete(id)
	m.closeClientIfNoConnections(clientKey)
	return nil
}

func (m *Manager) closeClientIfNoConnections(clientKey string) {
	lock := m.clientLock(clientKey)
	lock.Lock()
	defer lock.Unlock()
	hasOther := false
	m.sessions.Range(func(_, value any) bool {
		if value.(*Session).ClientKey == clientKey {
			hasOther = true
			return false
		}
		return true
	})
	m.logger.Debug("closeClientIfNoConnections", zap.String("clientKey", clientKey), zap.Bool("hasOtherConnections", hasOther))
	if hasOther {
		return
	}
	if client, ok := m.clients.LoadAndDelete(clientKey); ok {
		client.(*hopClient).Close()
	}
}

func (m *Manager) SetHostKeyDecision(host string, accept bool) error {
	return m.hostKey.decide(host, accept)
}

func newSession(m *Manager, clientKey string, client *cryptossh.Client) *Session {
	return &Session{
		manager:        m,
		ClientKey:      clientKey,
		client:         client,
		ID:             utils.GenerateRandomID(),
		OutputChan:     make(chan []byte, 200),
		stopOutput:     make(chan struct{}),
		outputBuffSize: 1024 * 10,
	}
}

func (sc *Session) Start(cols, rows int) error {
	sc.manager.logger.Debug("Starting SSH session", zap.String("id", sc.ID), zap.String("size", fmt.Sprintf("%dx%d", cols, rows)))
	var err error
	sc.session, err = sc.client.NewSession()
	if err != nil {
		return fmt.Errorf("Start session fail")
	}
	if err = sc.session.RequestPty("xterm-256color", rows, cols, cryptossh.TerminalModes{
		cryptossh.ECHO:          1,
		cryptossh.TTY_OP_ISPEED: 14400,
		cryptossh.TTY_OP_OSPEED: 14400,
	}); err != nil {
		return sc.failStart(err)
	}
	if err = sc.startInput(); err != nil {
		return sc.failStart(err)
	}
	if err = sc.startOutput(); err != nil {
		return sc.failStart(err)
	}
	if err = sc.session.Shell(); err != nil {
		return sc.failStart(err)
	}
	pty := sc.session
	go func() {
		waitErr := pty.Wait()
		if waitErr != nil {
			sc.manager.logger.Debug("SSH session ended with error:", zap.String("msg", waitErr.Error()), zap.String("id", sc.ID))
		} else {
			sc.manager.logger.Debug("SSH session ended", zap.String("ID", sc.ID))
		}
		sc.manager.notifyClosed(sc.ID)
	}()
	return nil
}

func (sc *Session) failStart(err error) error {
	if sc.session != nil {
		_ = sc.session.Close()
		sc.session = nil
	}
	sc.Stdin = nil
	return err
}

func (sc *Session) startOutput() error {
	stdout, err := sc.session.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := sc.session.StderrPipe()
	if err != nil {
		return err
	}
	sc.outputWg.Go(sc.readFromPipe(stdout, "stdout"))
	sc.outputWg.Go(sc.readFromPipe(stderr, "stderr"))
	return nil
}

func (sc *Session) readFromPipe(pipe io.Reader, pipeName string) func() {
	return func() {
		defer system.RecoverFromPanic()
		buf := make([]byte, sc.outputBuffSize)
		for {
			n, err := pipe.Read(buf)
			if err != nil {
				if err == io.EOF {
					sc.manager.logger.Debug(fmt.Sprintf("%s EOF reached", pipeName), zap.String("id", sc.ID))
					return
				}
				sc.manager.logger.Error(fmt.Sprintf("Error reading from %s:", pipeName), zap.Error(err), zap.String("id", sc.ID))
				return
			}
			if n > 0 {
				data := make([]byte, n)
				copy(data, buf[:n])
				select {
				case sc.OutputChan <- data:
				case <-sc.stopOutput:
					return
				}
			}
		}
	}
}

func (sc *Session) closeOutputChan() {
	sc.outputOnce.Do(func() {
		close(sc.OutputChan)
	})
}

func (sc *Session) startInput() error {
	stdin, err := sc.session.StdinPipe()
	if err != nil {
		return err
	}
	sc.Stdin = stdin
	return nil
}

func (sc *Session) closePTY() error {
	sc.closeMu.Lock()
	defer sc.closeMu.Unlock()
	sc.manager.logger.Debug("Closing SSH connection", zap.String("ID", sc.ID))
	if sc.isClosed {
		return nil
	}
	sc.isClosed = true
	sc.stopOutputOnce.Do(func() { close(sc.stopOutput) })
	if sc.session != nil {
		_ = sc.session.Signal(cryptossh.SIGTERM)
		_ = sc.session.Close()
		sc.session = nil
	}
	go func() {
		sc.outputWg.Wait()
		sc.closeOutputChan()
		sc.manager.logger.Debug("Output channel closed in Close()", zap.String("ID", sc.ID))
	}()
	return nil
}

func (sc *Session) Resize(cols, rows int) error {
	sc.manager.logger.Debug("Resizing SSH session", zap.String("id", sc.ID), zap.Int("cols", cols), zap.Int("rows", rows))
	if sc.session == nil {
		return fmt.Errorf("no active session")
	}
	return sc.session.WindowChange(rows, cols)
}
