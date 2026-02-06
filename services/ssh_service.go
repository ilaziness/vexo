package services

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/ilaziness/vexo/internal/system"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/xiwh/zmodem/zmodem"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
)

const (
	SSHEventInput            = "sshInput"
	SSHEventOutput           = "sshOutput"
	SSHEventClose            = "sshClose"
	SFTPEventClose           = "sftpClose"
	ErrSSHConnectionNotFound = "SSH connection with ID %s not found"
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
	SN   int    `json:"sn,omitempty"`
	Data []byte `json:"data"`
}

type SSHConnect struct {
	ID              string
	clientKey       string
	client          *ssh.Client
	sshService      *SSHService
	session         *ssh.Session
	sftpService     *SftpService
	isClosed        bool
	outputChan      chan []byte
	outputBuffSize  int
	inputUnregister func()         // 用于取消输入事件监听器的函数
	outputWg        sync.WaitGroup // 等待输出 goroutine 结束
	outputOnce      sync.Once      // 确保 outputChan 只被关闭一次
	stdin           io.WriteCloser // SSH stdin pipe
}

type SSHService struct {
	clients     *sync.Map // 客户端连接列表，key是host:port字符串, value是*ssh.Client
	SSHConnects *sync.Map // 终端会话列表，key是SSHConnect.ID，value是*sshConnect
}

func NewSSHService() *SSHService {
	return &SSHService{
		clients:     new(sync.Map),
		SSHConnects: new(sync.Map),
	}
}

// dialSSH establishes an SSH connection with the provided credentials and timeout.
func (s *SSHService) dialSSH(host string, port int, user, password, key, keyPassword string, timeout time.Duration) (*ssh.Client, error) {
	if password == "" && key == "" {
		return nil, fmt.Errorf("empty password and key")
	}
	cfg := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         timeout,
	}
	if key != "" {
		keyContent, err := os.ReadFile(key)
		if err != nil {
			return nil, fmt.Errorf("unable to read private key: %v", err)
		}
		var signer ssh.Signer
		if keyPassword == "" {
			signer, err = ssh.ParsePrivateKey(keyContent)
		} else {
			signer, err = ssh.ParsePrivateKeyWithPassphrase(keyContent, []byte(keyPassword))
		}
		if err != nil {
			return nil, err
		}
		cfg.Auth = []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		}
	}
	if password != "" {
		cfg.Auth = append(cfg.Auth, ssh.Password(password))
	}

	Logger.Debug("ssh key", zap.String("file", key))
	addr := net.JoinHostPort(host, fmt.Sprintf("%d", port))
	conn, err := net.DialTimeout("tcp", addr, cfg.Timeout)
	if err != nil {
		Logger.Debug("tcp connect error", zap.Error(err))
		return nil, err
	}
	if err = conn.SetDeadline(time.Now().Add(cfg.Timeout)); err != nil {
		conn.Close()
		return nil, err
	}
	c, chans, reqs, err := ssh.NewClientConn(conn, addr, cfg)
	if err != nil {
		conn.Close()
		Logger.Debug("ssh handshake error", zap.Error(err))
		return nil, err
	}
	if err := conn.SetDeadline(time.Time{}); err != nil {
		c.Close()
		return nil, err
	}
	return ssh.NewClient(c, chans, reqs), nil
}

// Connect establishes an SSH connection to the specified host using the provided credentials.
// return session ID if success
func (s *SSHService) Connect(host string, port int, user, password, key, keyPassword string) (ID string, err error) {
	Logger.Debug("Connecting to SSH server", zap.String("host", host), zap.Int("port", port))

	var client *ssh.Client
	clientKey := fmt.Sprintf("%s@%s:%d", user, host, port)
	clientVal, ok := s.clients.Load(clientKey)
	if ok {
		Logger.Debug("Using existing SSH client", zap.String("clientKey", clientKey))
		client = clientVal.(*ssh.Client)
	} else {
		client, err = s.dialSSH(host, port, user, password, key, keyPassword, time.Second*30)
		if err != nil {
			return "", err
		}
		// Store the new client for reuse
		s.clients.Store(clientKey, client)
		Logger.Debug("ssh connect ok and stored in cache", zap.String("clientKey", clientKey))
	}
	connect := NewSSHConnect(s, clientKey, client)
	s.SSHConnects.Store(connect.ID, connect)
	Logger.Debug("set connect done")
	return connect.ID, nil
}

// Start starts the SSH connection with the given ID.
func (s *SSHService) Start(ID string, cols, rows int) error {
	Logger.Debug("Starting SSH connection", zap.String("id", ID))

	conn, ok := s.SSHConnects.Load(ID)
	if !ok {
		return errors.New("ssh connect not found")
	}
	err := conn.(*SSHConnect).Start(cols, rows)
	if err != nil {
		s.CloseByID(ID)
	}
	return err
}

// TestConnectInfo tests SSH connection information without establishing a persistent connection.
func (s *SSHService) TestConnectInfo(host string, port int, user, password, key, keyPassword string) error {
	Logger.Debug("Testing SSH connection", zap.String("host", host), zap.Int("port", port))
	client, err := s.dialSSH(host, port, user, password, key, keyPassword, time.Second*20)
	if err != nil {
		return err
	}
	// Immediately close the test connection
	defer client.Close()
	Logger.Debug("Test SSH connect successful", zap.String("host", host), zap.Int("port", port))
	return nil
}

// StartSftp create sftp service for the SSH connection
func (s *SSHService) StartSftp(ID string) error {
	Logger.Debug("Starting SFTP service", zap.String("id", ID))
	connAny, ok := s.SSHConnects.Load(ID)
	if !ok {
		return fmt.Errorf(ErrSSHConnectionNotFound, ID)
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
		return fmt.Errorf(ErrSSHConnectionNotFound, ID)
	}
	return conn.(*SSHConnect).Resize(cols, rows)
}

// Close closes all SSH connections managed by the service.
func (s *SSHService) Close() {
	Logger.Debug("Close ssh all")
	s.SSHConnects.Range(func(_, connAny any) bool {
		conn := connAny.(*SSHConnect)
		_ = conn.Close()
		return true
	})
	s.clients = new(sync.Map)
	s.SSHConnects = new(sync.Map)
}

// CloseByID closes the SSH connection with the specified ID.
func (s *SSHService) CloseByID(ID string) error {
	Logger.Debug("CloseByID", zap.String("ID", ID))
	connAny, ok := s.SSHConnects.Load(ID)
	if !ok {
		return fmt.Errorf(ErrSSHConnectionNotFound, ID)
	}
	conn := connAny.(*SSHConnect)
	clientKey := conn.clientKey
	_ = conn.Close()
	s.SSHConnects.Delete(ID)
	s.closeClientIfNoConnections(clientKey)
	return nil
}

// closeClientIfNoConnections close ssh client if no other connections use it
func (s *SSHService) closeClientIfNoConnections(clientKey string) {
	hasOtherConnections := false
	s.SSHConnects.Range(func(_, value any) bool {
		conn := value.(*SSHConnect)
		if conn.clientKey == clientKey {
			hasOtherConnections = true
			return false // 停止遍历
		}
		return true
	})
	Logger.Debug("closeClientIfNoConnections", zap.String("clientKey", clientKey), zap.Bool("hasOtherConnections", hasOtherConnections))
	if !hasOtherConnections {
		client, ok := s.clients.LoadAndDelete(clientKey)
		if ok {
			err := client.(*ssh.Client).Close()
			if err != nil {
				Logger.Warn("Failed to close SSH client", zap.Error(err))
			}
		}
	}
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
		outputChan:     make(chan []byte, 10),
		outputBuffSize: 32768, // 32KB
	}
}

// Start 开启交互式终端会话
func (sc *SSHConnect) Start(cols, rows int) error {
	Logger.Debug("Starting SSH session", zap.String("id", sc.ID), zap.String("size", fmt.Sprintf("%dx%d", cols, rows)))
	var err error
	sc.session, err = sc.client.NewSession()
	if err != nil {
		Logger.Error("Failed to create SSH session", zap.Error(err), zap.String("id", sc.ID))
		return errors.New("Start session fail")
	}
	err = sc.session.RequestPty("xterm-256color", rows, cols, ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400, // 输入速度
		ssh.TTY_OP_OSPEED: 14400, // 输出速度
	})
	if err != nil {
		return err
	}
	err = sc.startInput()
	if err != nil {
		return err
	}
	err = sc.startOutput()
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
		defer system.RecoverFromPanic()
		startSN := 1
		for data := range sc.outputChan {
			//Logger.Debug("SSH session output:", zap.ByteString("msg", data), zap.String("id", sc.ID))
			app.Event.Emit(SSHEventOutput, SSHEventData{
				ID:   sc.ID,
				SN:   startSN,
				Data: data,
			})
			startSN++
		}
	}()
}

// startOutput 启动读取远程输出的协程，并集成 ZModem 支持
func (sc *SSHConnect) startOutput() error {
	if sc.stdin == nil {
		return fmt.Errorf("stdin not initialized; call startInput() before startOutput()")
	}
	// stdout, err := sc.session.StdoutPipe()
	// if err != nil {
	// 	return err
	// }
	stderr, err := sc.session.StderrPipe()
	if err != nil {
		return err
	}
	sc.sendOutput()

	zc := &ZModemConsumer{sshConnect: sc}
	zm := zmodem.New(zmodem.ZModemConsumer{
		OnUploadSkip:    zc.OnUploadSkip,
		OnUpload:        zc.OnUpload,
		OnCheckDownload: zc.OnCheckDownload,
		OnDownload:      zc.OnDownload,
		Writer:          sc.stdin,
		EchoWriter:      &ZModemEchoWriter{sshConnect: sc},
	})
	sc.session.Stdout = zm

	// 启动 stdout 读取 goroutine
	//sc.outputWg.Go(sc.readFromPipe(stdout, "stdout"))

	// 启动 stderr 读取 goroutine
	sc.outputWg.Go(sc.readFromPipe(stderr, "stderr"))
	return nil
}

// readFromPipe 从管道读取数据的通用方法
func (sc *SSHConnect) readFromPipe(pipe io.Reader, pipeName string) func() {
	return func() {
		defer system.RecoverFromPanic()
		buf := make([]byte, sc.outputBuffSize)
		for {
			n, err := pipe.Read(buf)
			if err != nil {
				if errors.Is(err, io.EOF) {
					Logger.Debug(fmt.Sprintf("%s EOF reached", pipeName), zap.String("id", sc.ID))
					return
				}
				Logger.Error(fmt.Sprintf("Error reading from %s:", pipeName), zap.Error(err), zap.String("id", sc.ID))
				return
			}
			if n > 0 {
				data := make([]byte, n)
				copy(data, buf[:n])
				sc.outputChan <- data
			}
		}
	}
}

// closeOutputChan 安全地关闭输出 channel（只关闭一次）
func (sc *SSHConnect) closeOutputChan() {
	sc.outputOnce.Do(func() {
		close(sc.outputChan)
	})
}

// startInput 启动监听前端输入的协程
func (sc *SSHConnect) startInput() error {
	stdin, err := sc.session.StdinPipe()
	if err != nil {
		return err
	}
	sc.stdin = stdin

	// 保存取消注册函数，以便在连接关闭时取消监听器
	sc.inputUnregister = app.Event.On(SSHEventInput, func(event *application.CustomEvent) {
		inData := event.Data.(SSHEventData)
		if inData.ID != sc.ID {
			return
		}
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
	if sc.isClosed {
		return nil
	}

	// 取消输入事件监听器，防止重复注册导致的输入重复问题
	if sc.inputUnregister != nil {
		sc.inputUnregister()
		sc.inputUnregister = nil
		Logger.Debug("Input event listener unregistered", zap.String("ID", sc.ID))
	}

	// Close SFTP service if exists
	if sc.sftpService != nil {
		sc.sftpService.Close()
		sc.sftpService = nil
		Logger.Debug("SFTP service closed", zap.String("ID", sc.ID))
	}

	// 关闭 SSH session，这会触发 stdout/stderr 的 EOF，导致读取 goroutine 退出
	if sc.session != nil {
		_ = sc.session.Signal(ssh.SIGTERM)
		_ = sc.session.Close()
		sc.session = nil
	}

	// 等待所有输出读取 goroutine 完成，然后关闭 channel
	// 注意：如果 startOutput 还没有被调用，outputWg.Wait() 会立即返回（计数为0）
	go func() {
		sc.outputWg.Wait()
		sc.closeOutputChan()
		Logger.Debug("Output channel closed in Close()", zap.String("ID", sc.ID))
	}()

	sc.sendCloseEvent()
	sc.isClosed = true
	sc.sshService.CloseByID(sc.ID)
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

// outputZModemMessage 通过 outputChan 发送消息到前端
func (sc *SSHConnect) outputZModemMessage(message string) {
	var (
		Reset  = "\033[0m"
		Yellow = "\033[33m"
	)
	colored := Yellow + "[ZMODEM] " + message + Reset
	data := []byte(colored + "\r\n")
	sc.outputChan <- data
}

// -----------------------------------------------------------------------------

// ZModemConsumer 实现 zmodem.ZModemConsumer 接口
type ZModemConsumer struct {
	sshConnect  *SSHConnect
	downloadDir string
}

// OnUploadSkip 处理跳过上传文件
func (c *ZModemConsumer) OnUploadSkip(file *zmodem.ZModemFile) {
	c.sshConnect.outputZModemMessage(fmt.Sprintf("文件已存在，跳过: %s", file.Filename))
}

// OnUpload 处理文件上传请求
func (c *ZModemConsumer) OnUpload() *zmodem.ZModemFile {
	filePath, err := AppInstance.SelectFile()
	if err != nil || filePath == "" {
		c.sshConnect.outputZModemMessage("文件选择已取消")
		return nil
	}

	fileName := filepath.Base(filePath)
	c.sshConnect.outputZModemMessage(fmt.Sprintf("开始上传文件: %s", fileName))

	// 创建 ZModem 文件对象
	uploadFile, err := zmodem.NewZModemLocalFile(filePath)
	if err != nil {
		c.sshConnect.outputZModemMessage(fmt.Sprintf("创建文件对象失败: %v", err))
		return nil
	}
	return uploadFile
}

// OnCheckDownload 检查下载文件是否已存在
func (c *ZModemConsumer) OnCheckDownload(file *zmodem.ZModemFile) {
	downloadDir, err := AppInstance.SelectDirectory()
	if err != nil {
		c.sshConnect.outputZModemMessage(fmt.Sprintf("选择下载目录失败: %v", err))
		file.Skip()
		return
	}
	// 检查文件是否已存在，如果存在则跳过
	targetPath := filepath.Join(downloadDir, file.Filename)
	if _, err := os.Stat(targetPath); err == nil {
		c.sshConnect.outputZModemMessage(fmt.Sprintf("文件已存在，跳过: %s", file.Filename))
		file.Skip()
		return
	}
	c.downloadDir = downloadDir
}

// OnDownload 处理文件下载
func (c *ZModemConsumer) OnDownload(file *zmodem.ZModemFile, reader io.ReadCloser) error {
	fileName := file.Filename
	c.sshConnect.outputZModemMessage(fmt.Sprintf("开始下载文件: %s", fileName))

	// 创建下载目录
	downloadDir := c.downloadDir
	if err := os.MkdirAll(downloadDir, os.ModePerm); err != nil {
		c.sshConnect.outputZModemMessage(fmt.Sprintf("创建下载目录失败: %v", err))
		return err
	}

	// 打开目标文件
	targetPath := filepath.Join(downloadDir, fileName)
	f, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		c.sshConnect.outputZModemMessage(fmt.Sprintf("创建目标文件失败: %v", err))
		return err
	}
	defer f.Close()

	// 复制数据并显示进度
	copied := int64(0)
	buffer := make([]byte, 100*1024) // 32KB buffer

	for {
		n, err := reader.Read(buffer)
		if n > 0 {
			_, writeErr := f.Write(buffer[:n])
			if writeErr != nil {
				c.sshConnect.outputZModemMessage(fmt.Sprintf("写入文件失败: %v", writeErr))
				return writeErr
			}

			copied += int64(n)
			c.sshConnect.outputZModemMessage(fmt.Sprintf("已接收: %d bytes", copied))
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			c.sshConnect.outputZModemMessage(fmt.Sprintf("读取数据失败: %v", err))
			return err
		}
	}

	c.sshConnect.outputZModemMessage(fmt.Sprintf("文件下载完成: %s (%d bytes)", targetPath, copied))
	return nil
}

// ZModemEchoWriter 实现 zmodem.EchoWriter 接口
type ZModemEchoWriter struct {
	sshConnect *SSHConnect
}

func (w *ZModemEchoWriter) Write(p []byte) (n int, err error) {
	w.sshConnect.outputChan <- p
	return len(p), nil
}
