package services

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"

	"github.com/coder/websocket"
	"github.com/ilaziness/vexo/internal/system"
	"github.com/wailsapp/wails/v3/pkg/application"
	"go.uber.org/zap"
)

// WSClient WebSocket 客户端连接
type WSClient struct {
	conn      *websocket.Conn
	stdin     io.WriteCloser
	sessionID string
	done      chan struct{}
	closeOnce sync.Once
}

// WebSocketService WebSocket 服务
type WebSocketService struct {
	app        *application.App
	httpServer *http.Server
	sshService *SSHService
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
	Logger.Info("Starting WebSocket server", zap.String("addr", wsAddr))

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
	Logger.Info("Stopping WebSocket server")
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

	connWS, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		CompressionMode:    websocket.CompressionContextTakeover,
		InsecureSkipVerify: true,
	})
	if err != nil {
		Logger.Error("Failed to accept WebSocket", zap.Error(err))
		return
	}

	// 启动 SSH session（WebSocket 连接建立后再启动）
	err = s.sshService.Start(sessionID, cols, rows)
	if err != nil {
		Logger.Error(errFailedToStartSSHSession, zap.Error(err), zap.String("id", sessionID))
		connWS.Close(websocket.StatusNormalClosure, errFailedToStartSSHSession)
		s.sshService.CloseByID(sessionID)
		http.Error(w, errFailedToStartSSHSession, http.StatusInternalServerError)
		return
	}

	// 获取 SSH 连接
	sshConnAny, ok := s.sshService.SSHConnects.Load(sessionID)
	if !ok {
		Logger.Error("SSH connection not found after start", zap.String("id", sessionID))
		connWS.Close(websocket.StatusNormalClosure, "SSH connection not found")
		http.Error(w, "SSH connection not found", http.StatusNotFound)
		return
	}

	sshConn := sshConnAny.(*SSHConnect)
	if sshConn.stdin == nil {
		Logger.Error("SSH not initialized", zap.String("id", sessionID))
		connWS.Close(websocket.StatusNormalClosure, "SSH stdin not initialized")
		http.Error(w, "SSH not initialized", http.StatusInternalServerError)
		return
	}

	client := &WSClient{
		conn:      connWS,
		stdin:     sshConn.stdin,
		sessionID: sessionID,
		done:      make(chan struct{}),
	}

	// 启动读写协程
	go s.readLoop(client)
	go s.writeLoop(client, sshConn.outputChan)

	// 等待连接关闭
	<-client.done
}

// closeDone 安全地关闭 done 通道（确保只关闭一次）
func (c *WSClient) closeDone() {
	c.closeOnce.Do(func() {
		close(c.done)
	})
}

// readLoop 读取客户端发送的数据（xterm输入）
func (s *WebSocketService) readLoop(client *WSClient) {
	defer func() {
		// 通知连接关闭
		client.closeDone()
		// WebSocket 关闭时，关闭对应的 SSH 连接
		s.closeSSHConnection(client.sessionID)
		client.conn.Close(websocket.StatusNormalClosure, "")
	}()

	ctx := context.Background()
	for {
		_, data, err := client.conn.Read(ctx)
		if err != nil {
			Logger.Debug("WebSocket read error", zap.Error(err), zap.String("sessionID", client.sessionID))
			return
		}

		// 直接写入 SSH stdin
		if client.stdin != nil {
			_, err := client.stdin.Write(data)
			if err != nil {
				Logger.Error("Failed to write to stdin", zap.Error(err), zap.String("sessionID", client.sessionID))
				return
			}
		}
	}
}

// writeLoop 将 SSH output 写入客户端
func (s *WebSocketService) writeLoop(client *WSClient, outputChan chan []byte) {
	defer func() {
		// 通知连接关闭
		client.closeDone()
		// WebSocket 关闭时，关闭对应的 SSH 连接
		s.closeSSHConnection(client.sessionID)
		client.conn.Close(websocket.StatusNormalClosure, "")
	}()

	ctx := context.Background()
	for {
		select {
		case data := <-outputChan:
			err := client.conn.Write(ctx, websocket.MessageBinary, data)
			if err != nil {
				Logger.Debug("WebSocket write error", zap.Error(err), zap.String("sessionID", client.sessionID))
				return
			}

		case <-client.done:
			return
		}
	}
}

// closeSSHConnection 关闭 SSH 连接
func (s *WebSocketService) closeSSHConnection(sessionID string) {
	Logger.Debug("Closing SSH connection due to WebSocket close", zap.String("sessionID", sessionID))
	if err := s.sshService.CloseByID(sessionID); err != nil {
		Logger.Error("Failed to close SSH connection", zap.Error(err), zap.String("sessionID", sessionID))
	}
}
