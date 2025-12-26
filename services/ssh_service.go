package services

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
	"unicode/utf8"

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
	ID         string
	sshService *SSHService
	clientKey  string
	client     *ssh.Client
	session    *ssh.Session
	outputChan chan []byte
	closeSig   chan struct{}

	stdoutBuf []byte
	stderrBuf []byte
}

type SSHService struct {
	app         *application.App
	SSHConnects map[string]*SSHConnect // key是SSHConnect.ID
	mutex       sync.Mutex
}

func NewSSHService(app *application.App) *SSHService {
	return &SSHService{
		app:         app,
		SSHConnects: make(map[string]*SSHConnect),
	}
}

// generateID generates a unique ID based on the current timestamp.
func generateConnectID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// Connect establishes an SSH connection to the specified host using the provided credentials.
// 连接成功返回数据ID
func (s *SSHService) Connect(host string, port int, user string, password string, key string) (ID string, err error) {
	var client *ssh.Client
	clientKey := fmt.Sprintf("%s:%d", host, port)
	clientVal, ok := sshClient.Load(clientKey)
	if ok {
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
		time.Sleep(time.Second * 1)
		client, err = ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, port), cfg)
		if err != nil {
			return "", err
		}
	}
	connect := NewSSHConnect(s, clientKey, client)
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.SSHConnects[connect.ID] = connect
	// Add to sshSession sync.Map
	sshSession.Store(connect.ID, connect)
	return connect.ID, nil
}

// Start starts the SSH connection with the given ID.
func (s *SSHService) Start(ID string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	conn, ok := s.SSHConnects[ID]
	if !ok {
		return fmt.Errorf("SSH connection with ID %s not found", ID)
	}
	return conn.Start()
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
	f, err := s.app.Dialog.OpenFile().SetTitle("选择私钥文件").PromptForSingleSelection()
	if err != nil {
		return "", err
	}
	return f, nil
}

//-----------------------------------------------------------------------------

// 返回可以安全发送的前缀长度（即最后一个完整 UTF-8 字符的位置）
func safeTruncateUTF8(buf []byte) int {
	if len(buf) == 0 {
		return 0
	}

	// 从末尾开始检查，最多回退 3 个字节（UTF-8 最大 4 字节，但首字节不会错）
	for i := 0; i < 4 && i < len(buf); i++ {
		if utf8.Valid(buf[:len(buf)-i]) {
			return len(buf) - i
		}
	}
	// 如果全都不合法（极端情况），至少发送第一个字节避免死锁
	return 1
}

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
			sc.sshService.app.Event.Emit(SSHEventOutput, SSHEventData{
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
	stderror, err := sc.session.StderrPipe()
	if err != nil {
		return err
	}
	sc.sendOutput()
	go func() {
		buf := make([]byte, 1024)
		for {
			select {
			case <-sc.closeSig:
				Logger.Debug("SSH stdout reading goroutine received close signal")
				return
			default:
				n, err := stdout.Read(buf)
				if err != nil {
					Logger.Error("Error reading from stdout:", zap.Error(err))
					// 发送剩余缓冲（即使不完整，EOF 时也发）
					if len(sc.stdoutBuf) > 0 {
						sc.outputChan <- sc.stdoutBuf
						sc.stdoutBuf = nil
					}
					return
				}
				// 追加新数据
				sc.stdoutBuf = append(sc.stdoutBuf, buf[:n]...)
				// 找到可安全发送的长度
				safeLen := safeTruncateUTF8(sc.stdoutBuf)
				if safeLen > 0 {
					sc.outputChan <- sc.stdoutBuf[:safeLen]
					sc.stdoutBuf = sc.stdoutBuf[safeLen:]
				}
			}
		}
	}()
	go func() {
		buf := make([]byte, 1024)
		for {
			select {
			case <-sc.closeSig:
				Logger.Debug("SSH stderr reading goroutine received close signal")
				return
			default:
				n, err := stderror.Read(buf)
				if err != nil {
					Logger.Error("Error reading from stderr:", zap.Error(err))
					// 发送剩余缓冲（即使不完整，EOF 时也发）
					if len(sc.stderrBuf) > 0 {
						sc.outputChan <- sc.stderrBuf
						sc.stderrBuf = nil
					}
					return
				}
				// 追加新数据
				sc.stderrBuf = append(sc.stderrBuf, buf[:n]...)
				// 找到可安全发送的长度
				safeLen := safeTruncateUTF8(sc.stderrBuf)
				if safeLen > 0 {
					sc.outputChan <- sc.stderrBuf[:safeLen]
					sc.stderrBuf = sc.stderrBuf[safeLen:]
				}
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
	sc.sshService.app.Event.On(SSHEventInput, func(event *application.CustomEvent) {
		indata := event.Data.(SSHEventData)
		if indata.ID != sc.ID {
			return
		}
		Logger.Debug("Received sshInput event", zap.Any("data", event.Data))
		stdin.Write([]byte(indata.Data))
	})
	return nil
}

// Close terminates the SSH connection and session.
// hasExit 会话是否已经结束，true已经结束，不再发送退出信号
func (sc *SSHConnect) Close(hasExit bool) error {
	if !hasExit {
		sc.sshService.app.Event.Emit(SSHEventInput, SSHEventData{
			ID:   sc.ID,
			Data: []byte("exit\r"),
		})
		<-sc.closeSig
	}
	if sc.session != nil {
		_ = sc.session.Close()
	}
	sc.outputChan <- []byte("connect closed")
	close(sc.outputChan)
	sc.sendColseEvent()
	sc.sshService.closeClient(sc.clientKey)
	// Remove from sshSession sync.Map
	sshSession.Delete(sc.ID)
	Logger.Info("SSH connection closed", zap.String("ID", sc.ID))
	return nil
}

// sendColseEvent 发送会话关闭消息
func (sc *SSHConnect) sendColseEvent() {
	sc.sshService.app.Event.Emit(SSHEventClose, SSHEventData{ID: sc.ID})
}

// Resize sets the remote pty size. cols/rows come from frontend terminal.
func (sc *SSHConnect) Resize(cols int, rows int) error {
	if sc.session == nil {
		return fmt.Errorf("no active session")
	}
	// WindowChange takes height, width
	return sc.session.WindowChange(rows, cols)
}

// GetSSHSession retrieves an SSH session by its ID from the sshSession sync.Map
func GetSSHSession(id string) (*SSHConnect, bool) {
	session, ok := sshSession.Load(id)
	if !ok {
		return nil, false
	}
	return session.(*SSHConnect), true
}
