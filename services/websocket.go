package services

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/ilaziness/vexo/internal/system"
	"github.com/wailsapp/wails/v3/pkg/application"
	"go.uber.org/zap"
)

// WSClient WebSocket 客户端连接
type WSClient struct {
	conn       *websocket.Conn
	stdin      io.WriteCloser
	sessionID  string
	done       chan struct{}
	closeOnce  sync.Once
	pongMissed int        // 未收到 pong 响应的次数
	pongMu     sync.Mutex // 保护 pongMissed 的互斥锁
}

// WebSocketService WebSocket 服务
type WebSocketService struct {
	app        *application.App
	httpServer *http.Server
	sshService *SSHService
	clients    sync.Map // map[string]*WSClient，存储所有活跃的客户端连接
}

var wsService *WebSocketService
var wsAddr string

// 错误消息常量
const errFailedToStartSSHSession = "Failed to start SSH session"

// NewWebSocketService 创建 WebSocket 服务
func NewWebSocketService(app *application.App, sshSvc *SSHService) *WebSocketService {
	wsService = &WebSocketService{
		app:        app,
		sshService: sshSvc,
	}
	return wsService
}

// registerClient 注册客户端连接
func (s *WebSocketService) registerClient(client *WSClient) {
	s.clients.Store(client.sessionID, client)
	Logger.Debug("WebSocket client registered", zap.String("sessionID", client.sessionID))
}

// unregisterClient 注销客户端连接
func (s *WebSocketService) unregisterClient(sessionID string) {
	s.clients.Delete(sessionID)
	Logger.Debug("WebSocket client unregistered", zap.String("sessionID", sessionID))
}

// GetWebSocketService 获取 WebSocket 服务实例
func GetWebSocketService() *WebSocketService {
	return wsService
}

// getWsAddr 获取 WebSocket 服务器地址
func getWsAddr() string {
	if wsAddr != "" {
		return wsAddr
	}

	base := 10697
	for i := range 100 {
		port := base + i
		addr := fmt.Sprintf("127.0.0.1:%d", port)
		ln, err := net.Listen("tcp", addr)
		if err != nil {
			continue
		}
		ln.Close()
		return addr
	}

	wsAddr = fmt.Sprintf("127.0.0.1:%d", base)
	return wsAddr
}

// Start 启动 WebSocket 服务
func (s *WebSocketService) Start() {
	wsAddr = getWsAddr()
	Logger.Debug("Starting WebSocket server", zap.String("addr", wsAddr))

	mux := http.NewServeMux()
	mux.HandleFunc("/ws/terminal", s.handleWebSocket)

	s.httpServer = &http.Server{
		Addr:    wsAddr,
		Handler: mux,
	}

	go func() {
		defer system.RecoverFromPanic()
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			Logger.Error("WebSocket server error", zap.Error(err))
		}
	}()
}

// Stop 停止 WebSocket 服务
func (s *WebSocketService) Stop() {
	Logger.Debug("Stopping WebSocket server")
	if s.httpServer != nil {
		s.httpServer.Close()
	}
}

// handleWebSocket 处理 WebSocket 连接
func (s *WebSocketService) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// 获取 URL 中的 session ID 参数
	sessionID := r.URL.Query().Get("id")
	if sessionID == "" {
		Logger.Error("Session ID is required")
		http.Error(w, "Session ID is required", http.StatusBadRequest)
		return
	}

	// 获取 cols 和 rows 参数
	cols := 80
	rows := 24
	if colsStr := r.URL.Query().Get("cols"); colsStr != "" {
		if _, err := fmt.Sscanf(colsStr, "%d", &cols); err != nil {
			cols = 80
		}
	}
	if rowsStr := r.URL.Query().Get("rows"); rowsStr != "" {
		if _, err := fmt.Sscanf(rowsStr, "%d", &rows); err != nil {
			rows = 24
		}
	}

	Logger.Debug("New WebSocket connection", zap.String("sessionID", sessionID), zap.Int("cols", cols), zap.Int("rows", rows))

	// 启动 SSH session
	err := s.sshService.Start(sessionID, cols, rows)
	if err != nil {
		Logger.Error(errFailedToStartSSHSession, zap.Error(err), zap.String("id", sessionID))
		s.sshService.CloseByID(sessionID)
		http.Error(w, errFailedToStartSSHSession, http.StatusInternalServerError)
		return
	}

	// 获取 SSH 连接
	sshConnAny, ok := s.sshService.SSHConnects.Load(sessionID)
	if !ok {
		Logger.Error("SSH connection not found after start", zap.String("id", sessionID))
		http.Error(w, "SSH connection not found", http.StatusNotFound)
		return
	}

	sshConn := sshConnAny.(*SSHConnect)
	if sshConn.stdin == nil {
		Logger.Error("SSH not initialized", zap.String("id", sessionID))
		http.Error(w, "SSH not initialized", http.StatusInternalServerError)
		return
	}

	client := &WSClient{
		done:      make(chan struct{}),
		stdin:     sshConn.stdin,
		sessionID: sessionID,
	}
	connWS, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		CompressionMode:    websocket.CompressionContextTakeover,
		InsecureSkipVerify: true,
		OnPongReceived: func(_ context.Context, _ []byte) {
			client.pongMu.Lock()
			client.pongMissed = 0
			client.pongMu.Unlock()
		},
	})
	if err != nil {
		Logger.Error("Failed to accept WebSocket", zap.Error(err))
		return
	}
	client.conn = connWS

	// 注册客户端
	s.registerClient(client)

	// 启动读写协程
	go client.readLoop()
	go client.writeLoop(sshConn.outputChan)
	client.ping()

	client.closeDone()
	s.unregisterClient(client.sessionID)
	s.closeSSHConnection(client.sessionID)
	client.conn.Close(websocket.StatusNormalClosure, "")
	Logger.Debug("WebSocket connection closed", zap.String("sessionID", client.sessionID))
}

// closeDone 安全地关闭
func (c *WSClient) closeDone() {
	c.closeOnce.Do(func() {
		close(c.done)
	})
}

// ping 发送 ping 消息保持连接活跃，5秒发送一次，累计2次未收到响应则关闭连接
func (c *WSClient) ping() {
	ctx := context.Background()
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.pongMu.Lock()
			c.pongMissed++
			if c.pongMissed > 2 {
				c.pongMu.Unlock()
				Logger.Warn("WebSocket pong timeout, closing connection", zap.String("sessionID", c.sessionID))
				return
			}
			c.pongMu.Unlock()

			err := c.conn.Ping(ctx)
			if err != nil {
				Logger.Error("WebSocket ping error", zap.Error(err), zap.String("sessionID", c.sessionID))
				return
			}
		case <-c.done:
			return
		}
	}
}

// readLoop 读取客户端发送的数据
func (c *WSClient) readLoop() {
	defer c.closeDone()
	ctx := context.Background()
	for {
		_, data, err := c.conn.Read(ctx)
		if err != nil {
			Logger.Debug("WebSocket read error", zap.Error(err), zap.String("sessionID", c.sessionID))
			return
		}

		// 直接写入 SSH stdin
		if c.stdin != nil {
			_, err := c.stdin.Write(data)
			if err != nil {
				Logger.Error("Failed to write to stdin", zap.Error(err), zap.String("sessionID", c.sessionID))
				return
			}
		}
	}
}

// writeLoop 将 SSH output 写入客户端
func (c *WSClient) writeLoop(outputChan chan []byte) {
	defer c.closeDone()
	ctx := context.Background()
	for {
		select {
		case data := <-outputChan:
			err := c.conn.Write(ctx, websocket.MessageBinary, data)
			if err != nil {
				Logger.Debug("WebSocket write error", zap.Error(err), zap.String("sessionID", c.sessionID))
				return
			}

		case <-c.done:
			return
		}
	}
}

// closeSSHConnection 关闭 SSH 连接
func (s *WebSocketService) closeSSHConnection(sessionID string) {
	Logger.Debug("Closing SSH connection due to WebSocket close", zap.String("sessionID", sessionID))
	if err := s.sshService.CloseByID(sessionID); err != nil {
		Logger.Debug("Failed to close SSH connection", zap.Error(err), zap.String("sessionID", sessionID))
	}
}

// CloseClient 关闭指定客户端连接
func (s *WebSocketService) CloseClient(sessionID string) {
	Logger.Debug("Closing WebSocket client connection", zap.String("sessionID", sessionID))

	clientAny, ok := s.clients.Load(sessionID)
	if !ok {
		Logger.Debug("WebSocket client not found", zap.String("sessionID", sessionID))
		return
	}

	client := clientAny.(*WSClient)
	client.closeDone()
}
