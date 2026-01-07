package services

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
)

const (
	SSHEventInput  = "sshInput"
	SSHEventOutput = "sshOutput"
	SSHEventClose  = "sshClose"
)

func init() {
	application.RegisterEvent[SSHEventData](SSHEventInput)
	application.RegisterEvent[SSHEventData](SSHEventOutput)
	application.RegisterEvent[SSHEventData](SSHEventClose)
}

// sshClient 客户端连接列表，key是host:port字符串
var sshClient = new(sync.Map)

// sshSession 终端会话列表，key是SSHConnect.ID
var sshSession = new(sync.Map)

type SSHEventData struct {
	ID   string `json:"id"`
	Data []byte `json:"data"`
}

type SSHConnect struct {
	ID          string
	sshService  *SSHService
	clientKey   string
	client      *ssh.Client
	session     *ssh.Session
	sftpService *SftpService
	outputChan  chan []byte
	closeSig    chan struct{}
	isClosed    bool
}

type SSHService struct {
	SSHConnects map[string]*SSHConnect // key是SSHConnect.ID
	mutex       sync.Mutex
}

func NewSSHService() *SSHService {
	return &SSHService{
		SSHConnects: make(map[string]*SSHConnect),
	}
}

// generateID generates a unique ID based on the current timestamp.
func generateConnectID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// Connect establishes an SSH connection to the specified host using the provided credentials.
// return session ID if success
func (s *SSHService) Connect(host string, port int, user string, password string, key string) (ID string, err error) {
	Logger.Debug("Connecting to SSH server", zap.String("host", host), zap.Int("port", port))
	var client *ssh.Client
	clientKey := fmt.Sprintf("%s:%d", host, port)
	clientVal, ok := sshClient.Load(clientKey)
	if ok {
		Logger.Debug("Using existing SSH client", zap.String("clientKey", clientKey))
		client = clientVal.(*ssh.Client)
	} else {
		cfg := &ssh.ClientConfig{
			User:            user,
			Auth:            []ssh.AuthMethod{},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         time.Second * 60,
		}
		if key != "" {
			keyContent, err := os.ReadFile(key)
			if err != nil {
				log.Fatalf("unable to read private key: %v", err)
				return "", err
			}
			signer, err := ssh.ParsePrivateKey(keyContent)
			if err != nil {
				return "", err
			}
			cfg.Auth = []ssh.AuthMethod{
				ssh.PublicKeys(signer),
			}
		}
		if password != "" {
			cfg.Auth = append(cfg.Auth, ssh.Password(password))
		}
		client, err = ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, port), cfg)
		if err != nil {
			return "", err
		}
		// Store the new client for reuse
		sshClient.Store(clientKey, client)
		Logger.Debug("ssh connect ok and stored in cache", zap.String("clientKey", clientKey))
	}
	connect := NewSSHConnect(s, clientKey, client)
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.SSHConnects[connect.ID] = connect
	// Add to sshSession sync.Map
	sshSession.Store(connect.ID, connect)
	Logger.Debug("set connect done")
	return connect.ID, nil
}

// Start starts the SSH connection with the given ID.
func (s *SSHService) Start(ID string) error {
	Logger.Debug("Starting SSH connection", zap.String("id", ID))
	s.mutex.Lock()
	defer s.mutex.Unlock()
	conn, ok := s.SSHConnects[ID]
	if !ok {
		return fmt.Errorf("SSH connection with ID %s not found", ID)
	}
	return conn.Start()
}

// StartSftp create sftp service for the SSH connection
func (s *SSHService) StartSftp(ID string) error {
	Logger.Debug("Starting SFTP service", zap.String("id", ID))
	s.mutex.Lock()
	defer s.mutex.Unlock()
	conn, ok := s.SSHConnects[ID]
	if !ok {
		return fmt.Errorf("SSH connection with ID %s not found", ID)
	}
	// If SFTP service already exists, return without error
	if conn.sftpService != nil {
		Logger.Debug("SFTP service already exists for this connection", zap.String("id", ID))
		return nil
	}
	sftpService := NewSftpService()
	err := sftpService.Connect(conn.client)
	if err != nil {
		return err
	}
	conn.sftpService = sftpService
	Logger.Debug("SFTP service started successfully", zap.String("id", ID))
	return nil
}

// Resize resizes the terminal for the SSH connection with the given ID.
func (s *SSHService) Resize(ID string, cols int, rows int) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	conn, ok := s.SSHConnects[ID]
	if !ok {
		return fmt.Errorf("SSH connection with ID %s not found", ID)
	}
	return conn.Resize(cols, rows)
}

// Close closes all SSH connections managed by the service.
func (s *SSHService) Close() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, conn := range s.SSHConnects {
		conn.Close(false)
	}
	s.SSHConnects = make(map[string]*SSHConnect)
}

// CloseByID closes the SSH connection with the specified ID.
func (s *SSHService) CloseByID(ID string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	conn, ok := s.SSHConnects[ID]
	if !ok {
		return fmt.Errorf("SSH connection with ID %s not found", ID)
	}
	conn.Close(false)
	delete(s.SSHConnects, ID)
	return nil
}

// GetSftpService returns the SFTP service for the specified SSH connection ID
func (s *SSHService) GetSftpService(ID string) (*SftpService, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	conn, ok := s.SSHConnects[ID]
	if !ok {
		return nil, fmt.Errorf("SSH connection with ID %s not found", ID)
	}
	if conn.sftpService == nil {
		return nil, fmt.Errorf("SFTP service not started for connection %s", ID)
	}
	return conn.sftpService, nil
}

func (s *SSHService) closeClient(clientKey string) {
	for _, conn := range s.SSHConnects {
		if conn.clientKey == clientKey {
			return
		}
	}
	client, ok := sshClient.LoadAndDelete(clientKey)
	if !ok {
		return
	}
	client.(*ssh.Client).Close()
}

// SelectKeyFile 选择私钥文件
func (s *SSHService) SelectKeyFile() (string, error) {
	f, err := app.Dialog.OpenFile().SetTitle("选择私钥文件").PromptForSingleSelection()
	if err != nil {
		return "", err
	}
	return f, nil
}

// -----------------------------------------------------------------------------

func NewSSHConnect(sshService *SSHService, clientKey string, client *ssh.Client) *SSHConnect {
	return &SSHConnect{
		sshService: sshService,
		clientKey:  clientKey,
		client:     client,
		ID:         generateConnectID(),
		outputChan: make(chan []byte, 1024),
		closeSig:   make(chan struct{}, 1),
	}
}

// Start 开启交互式终端会话
func (sc *SSHConnect) Start() error {
	Logger.Debug("Starting SSH session", zap.String("id", sc.ID))
	var err error
	sc.session, err = sc.client.NewSession()
	if err != nil {
		return err
	}
	err = sc.session.RequestPty("xterm-256color", 80, 40, ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400, // 输入速度
		ssh.TTY_OP_OSPEED: 14400, // 输出速度
	})
	if err != nil {
		return err
	}
	err = sc.startOutput()
	if err != nil {
		return err
	}
	err = sc.startInput()
	if err != nil {
		return err
	}
	err = sc.session.Shell()
	if err != nil {
		return err
	}
	go func() {
		err = sc.session.Wait()
		if err != nil {
			Logger.Warn("SSH session ended with error:", zap.String("msg", err.Error()), zap.String("id", sc.ID))
		} else {
			Logger.Info("SSH session ended", zap.String("ID", sc.ID))
		}
		sc.closeSig <- struct{}{}
		close(sc.closeSig)
		_ = sc.Close(true)
	}()
	return nil
}

func (sc *SSHConnect) sendOutput() {
	go func() {
		for data := range sc.outputChan {
			app.Event.Emit(SSHEventOutput, SSHEventData{
				ID:   sc.ID,
				Data: data,
			})
		}
	}()
}

// startOutput 启动读取远程输出的协程
func (sc *SSHConnect) startOutput() error {
	stdout, err := sc.session.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := sc.session.StderrPipe()
	if err != nil {
		return err
	}
	sc.sendOutput()
	go func() {
		buf := make([]byte, 10240)
		for {
			select {
			case <-sc.closeSig:
				Logger.Debug("SSH stdout reading goroutine received close signal")
				return
			default:
				n, err := stdout.Read(buf)
				if err != nil {
					if errors.Is(err, io.EOF) {
						return
					}
					Logger.Error("Error reading from stdout:", zap.Error(err))
					return
				}
				data := make([]byte, n)
				copy(data, buf[:n])
				sc.outputChan <- data
			}
		}
	}()
	go func() {
		buf := make([]byte, 10240)
		for {
			select {
			case <-sc.closeSig:
				Logger.Debug("SSH stderr reading goroutine received close signal")
				return
			default:
				n, err := stderr.Read(buf)
				if err != nil {
					if errors.Is(err, io.EOF) {
						return
					}
					Logger.Error("Error reading from stderr:", zap.Error(err))
					return
				}
				data := make([]byte, n)
				copy(data, buf[:n])
				sc.outputChan <- data
			}
		}
	}()
	return nil
}

// startInput 启动监听前端输入的协程
func (sc *SSHConnect) startInput() error {
	stdin, err := sc.session.StdinPipe()
	if err != nil {
		return err
	}
	app.Event.On(SSHEventInput, func(event *application.CustomEvent) {
		indata := event.Data.(SSHEventData)
		if indata.ID != sc.ID {
			return
		}
		Logger.Debug("Received sshInput event", zap.Any("data", event.Data))
		stdin.Write(indata.Data)
	})
	return nil
}

// Close terminates the SSH connection and session.
// hasExit 会话是否已经结束，true已经结束，不再发送退出信号
func (sc *SSHConnect) Close(hasExit bool) error {
	Logger.Debug("Closing SSH connection", zap.String("ID", sc.ID))
	// Close SFTP service if exists
	if sc.sftpService != nil {
		sc.sftpService.Close()
		sc.sftpService = nil
		Logger.Debug("SFTP service closed", zap.String("ID", sc.ID))
	}
	if !hasExit {
		app.Event.Emit(SSHEventInput, SSHEventData{
			ID:   sc.ID,
			Data: []byte("exit\r"),
		})
		<-sc.closeSig
	}
	if sc.session != nil {
		_ = sc.session.Close()
	}
	_, ok := <-sc.outputChan
	if ok {
		sc.outputChan <- []byte("connect closed")
	}
	close(sc.outputChan)
	sc.sendCloseEvent()
	sc.sshService.closeClient(sc.clientKey)
	// Remove from sshSession sync.Map
	sshSession.Delete(sc.ID)
	Logger.Info("SSH connection closed", zap.String("ID", sc.ID))
	return nil
}

// sendCloseEvent 发送会话关闭消息
func (sc *SSHConnect) sendCloseEvent() {
	app.Event.Emit(SSHEventClose, SSHEventData{ID: sc.ID})
}

// Resize sets the remote pty size. cols/rows come from frontend terminal.
func (sc *SSHConnect) Resize(cols int, rows int) error {
	if sc.session == nil {
		return fmt.Errorf("no active session")
	}
	// WindowChange takes height, width
	return sc.session.WindowChange(rows, cols)
}
