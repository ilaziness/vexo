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
	SFTPEventClose = "sftpClose"
)

func init() {
	application.RegisterEvent[SSHEventData](SSHEventInput)
	application.RegisterEvent[SSHEventData](SSHEventOutput)
	application.RegisterEvent[SSHEventData](SSHEventClose)
	application.RegisterEvent[SSHEventData](SFTPEventClose)
}

// SSHEventData exchange data struct
type SSHEventData struct {
	ID   string `json:"id"`
	Data []byte `json:"data"`
}

type SSHConnect struct {
	ID             string
	clientKey      string
	client         *ssh.Client
	sshService     *SSHService
	session        *ssh.Session
	sftpService    *SftpService
	outputChan     chan []byte
	outputBuffSize int
}

type SSHService struct {
	clients     *sync.Map // 客户端连接列表，key是host:port字符串, value是*ssh.Client
	clientNum   *sync.Map // 客户端session列表，key是host:port字符串，value打开的会话数
	SSHConnects *sync.Map // 终端会话列表，key是SSHConnect.ID，value是*sshConnect
}

func NewSSHService() *SSHService {
	return &SSHService{
		clients:     new(sync.Map),
		clientNum:   new(sync.Map),
		SSHConnects: new(sync.Map),
	}
}

// Connect establishes an SSH connection to the specified host using the provided credentials.
// return session ID if success
func (s *SSHService) Connect(host string, port int, user string, password string, key string) (ID string, err error) {
	Logger.Debug("Connecting to SSH server", zap.String("host", host), zap.Int("port", port))
	var client *ssh.Client
	clientKey := fmt.Sprintf("%s:%d", host, port)
	clientVal, ok := s.clients.Load(clientKey)
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
		s.clients.Store(clientKey, client)
		Logger.Debug("ssh connect ok and stored in cache", zap.String("clientKey", clientKey))
	}
	connect := NewSSHConnect(s, clientKey, client)
	s.SSHConnects.Store(connect.ID, connect)
	s.clientNumAdd(clientKey)
	Logger.Debug("set connect done")
	return connect.ID, nil
}

func (s *SSHService) clientNumAdd(clientKey string) {
	oldNum, ok := s.clientNum.Load(clientKey)
	if ok {
		s.clientNum.Store(clientKey, oldNum.(int)+1)
	} else {
		s.clientNum.Store(clientKey, 1)
	}
}

func (s *SSHService) clientNumSub(clientKey string) {
	oldNum, ok := s.clientNum.Load(clientKey)
	if ok {
		s.clientNum.Store(clientKey, oldNum.(int)-1)
		if oldNum.(int) <= 0 {
			s.closeClient(clientKey)
		}
	}
}

// Start starts the SSH connection with the given ID.
func (s *SSHService) Start(ID string, cols, rows int) error {
	Logger.Debug("Starting SSH connection", zap.String("id", ID))

	conn, ok := s.SSHConnects.Load(ID)
	if !ok {
		return fmt.Errorf("SSH connection with ID %s not found", ID)
	}
	return conn.(*SSHConnect).Start(cols, rows)
}

// Restart restarts the SSH connection with the given ID.
func (s *SSHService) Restart(ID string, cols, rows int) error {
	Logger.Debug("Restarting SSH connection", zap.String("id", ID))
	conn, ok := s.SSHConnects.Load(ID)
	if !ok {
		return fmt.Errorf("SSH connection with ID %s not found", ID)
	}
	_ = s.CloseByID(ID)
	return conn.(*SSHConnect).Start(cols, rows)
}

// StartSftp create sftp service for the SSH connection
func (s *SSHService) StartSftp(ID string) error {
	Logger.Debug("Starting SFTP service", zap.String("id", ID))
	connAny, ok := s.SSHConnects.Load(ID)
	if !ok {
		return fmt.Errorf("SSH connection with ID %s not found", ID)
	}
	conn := connAny.(*SSHConnect)
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
	sftpClient.Store(ID, sftpService)
	Logger.Debug("SFTP service started successfully", zap.String("id", ID))
	return nil
}

// Resize resizes the terminal for the SSH connection with the given ID.
func (s *SSHService) Resize(ID string, cols int, rows int) error {
	conn, ok := s.SSHConnects.Load(ID)
	if !ok {
		return fmt.Errorf("SSH connection with ID %s not found", ID)
	}
	return conn.(*SSHConnect).Resize(cols, rows)
}

// Close closes all SSH connections managed by the service.
func (s *SSHService) Close() {
	s.SSHConnects.Range(func(_, connAny any) bool {
		conn := connAny.(*SSHConnect)
		_ = conn.Close()
		return true
	})
	s.clients = new(sync.Map)
	s.clientNum = new(sync.Map)
	s.SSHConnects = new(sync.Map)
}

// CloseByID closes the SSH connection with the specified ID.
func (s *SSHService) CloseByID(ID string) error {
	connAny, ok := s.SSHConnects.Load(ID)
	if !ok {
		return fmt.Errorf("SSH connection with ID %s not found", ID)
	}
	conn := connAny.(*SSHConnect)
	_ = conn.Close()
	s.SSHConnects.Delete(ID)
	s.clientNumSub(conn.clientKey)
	return nil
}

// closeClient close ssh client
func (s *SSHService) closeClient(clientKey string) {
	s.SSHConnects.Range(func(_, value any) bool {
		conn := value.(*SSHConnect)
		if conn.clientKey == clientKey {
			client, ok := s.clients.LoadAndDelete(clientKey)
			if ok {
				err := client.(*ssh.Client).Close()
				if err != nil {
					Logger.Warn("Failed to close SSH client", zap.Error(err))
				}
			}
			return false
		}
		return true
	})
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

// generateID generates a unique ID based on the current timestamp.
func generateConnectID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func NewSSHConnect(sshService *SSHService, clientKey string, client *ssh.Client) *SSHConnect {
	return &SSHConnect{
		sshService:     sshService,
		clientKey:      clientKey,
		client:         client,
		ID:             generateConnectID(),
		outputChan:     make(chan []byte, 1024),
		outputBuffSize: 32768, // 32KB - 业界标准缓冲区大小
	}
}

// Start 开启交互式终端会话
func (sc *SSHConnect) Start(cols, rows int) error {
	Logger.Debug("Starting SSH session", zap.String("id", sc.ID), zap.String("size", fmt.Sprintf("%dx%d", cols, rows)))
	var err error
	sc.session, err = sc.client.NewSession()
	if err != nil {
		return err
	}
	err = sc.session.RequestPty("xterm-256color", rows, cols, ssh.TerminalModes{
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
			Logger.Debug("SSH session ended with error:", zap.String("msg", err.Error()), zap.String("id", sc.ID))
		} else {
			Logger.Info("SSH session ended", zap.String("ID", sc.ID))
		}
		_ = sc.Close()
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
		buf := make([]byte, sc.outputBuffSize)
		for {
			n, err := stdout.Read(buf)
			if err != nil {
				if errors.Is(err, io.EOF) {
					return
				}
				Logger.Error("Error reading from stdout:", zap.Error(err))
				sc.outputChan <- []byte(err.Error())
				return
			}
			data := make([]byte, n)
			copy(data, buf[:n])
			sc.outputChan <- data
		}
	}()
	go func() {
		buf := make([]byte, sc.outputBuffSize)
		for {

			n, err := stderr.Read(buf)
			if err != nil {
				if errors.Is(err, io.EOF) {
					return
				}
				Logger.Error("Error reading from stderr:", zap.Error(err))
				sc.outputChan <- []byte(err.Error())
				return
			}
			data := make([]byte, n)
			copy(data, buf[:n])
			sc.outputChan <- data

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
		inData := event.Data.(SSHEventData)
		if inData.ID != sc.ID {
			return
		}
		Logger.Debug("Received sshInput event", zap.Any("data", event.Data))
		_, err = stdin.Write(inData.Data)
		if err != nil {
			Logger.Warn("Error writing to stdin:", zap.Error(err))
			return
		}
	})
	return nil
}

// Close terminates the SSH connection and session.
func (sc *SSHConnect) Close() error {
	Logger.Debug("Closing SSH connection", zap.String("ID", sc.ID))
	// Close SFTP service if exists
	if sc.sftpService != nil {
		sc.sftpService.Close()
		sc.sftpService = nil
		Logger.Debug("SFTP service closed", zap.String("ID", sc.ID))
	}

	if sc.session != nil {
		_ = sc.session.Signal(ssh.SIGTERM)
		_ = sc.session.Close()
	}

	sc.sendCloseEvent()
	sc.sshService.clientNumSub(sc.clientKey)
	sc.sshService.SSHConnects.Delete(sc.ID)
	Logger.Info("SSH connection closed", zap.String("ID", sc.ID))
	return nil
}

// sendCloseEvent 发送会话关闭消息
func (sc *SSHConnect) sendCloseEvent() {
	app.Event.Emit(SSHEventClose, SSHEventData{ID: sc.ID})
}

// Resize sets the remote pty size. cols/rows come from frontend terminal.
func (sc *SSHConnect) Resize(cols int, rows int) error {
	Logger.Debug("Resizing SSH session", zap.String("id", sc.ID), zap.Int("cols", cols), zap.Int("rows", rows))
	if sc.session == nil {
		return fmt.Errorf("no active session")
	}
	// WindowChange takes height, width
	return sc.session.WindowChange(rows, cols)
}
