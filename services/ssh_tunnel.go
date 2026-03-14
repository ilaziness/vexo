package services

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/ilaziness/vexo/internal/system"
	"github.com/ilaziness/vexo/internal/utils"
	"github.com/things-go/go-socks5"
	"go.uber.org/zap"
)

// ssh tunnel

const (
	tunnelTypeLocal   = "local"
	tunnelTypeRemote  = "remote"
	tunnelTypeDynamic = "dynamic"
)

var (
	sshSessionNotFoundErr = errors.New("ssh connect not found")
)

var sshTunnels sync.Map // map[string]*SSHTunnel，key 是隧道ID
var sshTunnelService *SSHTunnelService

type TunnelList struct {
	TunnelType string       `json:"tunnelType"`
	Tunnels    []TunnelInfo `json:"tunnels"`
}
type TunnelInfo struct {
	ID         string `json:"id"`
	TunnelType string `json:"tunnelType"`
	SessionID  string `json:"sessionID"`
	LocalAddr  string `json:"localAddr"`
	RemoteAddr string `json:"remoteAddr"`
}

type SSHTunnelService struct {
	sshService *SSHService
}

type SSHTunnel struct {
	ID         string
	tunnelType string
	sessionID  string
	LocalAddr  string // 本地监听地址，格式: ip:port
	RemoteAddr string // 远程地址，格式: ip:port

	listener net.Listener
	exitCh   chan struct{}
	wg       sync.WaitGroup
	// For dynamic tunnels
	connLock    sync.Mutex
	activeConns map[net.Conn]struct{}
	socksServer *socks5.Server
}

func NewSSHTunnelService(sshService *SSHService) *SSHTunnelService {
	sshTunnelService = &SSHTunnelService{
		sshService: sshService,
	}
	return sshTunnelService
}

// TunnelList 隧道列表，按类型分组
func (t *SSHTunnelService) TunnelList() []TunnelList {
	tunnelsByType := make(map[string][]TunnelInfo)

	sshTunnels.Range(func(key, value any) bool {
		tunnel, ok := value.(*SSHTunnel)
		if ok {
			tunnelInfo := TunnelInfo{
				ID:         tunnel.ID,
				TunnelType: tunnel.tunnelType,
				SessionID:  tunnel.sessionID,
				LocalAddr:  tunnel.LocalAddr,
				RemoteAddr: tunnel.RemoteAddr,
			}
			tunnelsByType[tunnel.tunnelType] = append(tunnelsByType[tunnel.tunnelType], tunnelInfo)
		}
		return true
	})

	// 将分组后的数据转换为 TunnelList 数组
	var result []TunnelList
	for tunnelType, tunnels := range tunnelsByType {
		result = append(result, TunnelList{
			TunnelType: tunnelType,
			Tunnels:    tunnels,
		})
	}

	return result
}

// StartLocal 本地端口转发（local forwarding）。
// localAddr 格式: ip:port，例如 127.0.0.1:8080
// remoteAddr 格式: host:port，例如 127.0.0.1:3306 或 example.com:443
func (t *SSHTunnelService) StartLocal(sessionID string, localAddr string, remoteAddr string) (string, error) {
	connAny, ok := t.sshService.SSHConnects.Load(sessionID)
	if !ok {
		return "", sshSessionNotFoundErr
	}
	sc := connAny.(*SSHConnect)

	ln, err := net.Listen("tcp", localAddr)
	if err != nil {
		return "", err
	}

	tunnel := &SSHTunnel{
		ID:         utils.GenerateRandomID(),
		tunnelType: tunnelTypeLocal,
		sessionID:  sessionID,
		LocalAddr:  localAddr,
		RemoteAddr: remoteAddr,
		listener:   ln,
		exitCh:     make(chan struct{}),
	}
	sshTunnels.Store(tunnel.ID, tunnel)

	// accept loop
	tunnel.wg.Add(1)
	system.SafeGo(func() {
		defer tunnel.wg.Done()
		for {
			localConn, err := ln.Accept()
			if err != nil {
				// listener closed or error
				Logger.Debug("accept loop quit", zap.String("tunnelID", tunnel.ID), zap.Error(err))
				return
			}
			// handle each connection
			tunnel.wg.Add(1)
			system.SafeGo(func() {
				defer tunnel.wg.Done()
				defer localConn.Close()

				remoteConn, err := sc.client.Dial("tcp", tunnel.RemoteAddr)
				if err != nil {
					Logger.Debug("failed to dial remote", zap.String("tunnelID", tunnel.ID), zap.Error(err))
					return
				}
				defer remoteConn.Close()

				Logger.Debug("SSH tunnel local forwarding connection established")

				// 创建退出信号通道和同步机制
				connExitCh := make(chan struct{})
				var closeOnce sync.Once
				closeConnExitCh := func() {
					closeOnce.Do(func() {
						close(connExitCh)
					})
				}

				// 启动双向数据复制
				tunnel.wg.Add(2)
				system.SafeGo(func() {
					defer tunnel.wg.Done()
					io.Copy(remoteConn, localConn)
					closeConnExitCh()
					Logger.Debug("SSH tunnel copy1 exit")
				})
				system.SafeGo(func() {
					defer tunnel.wg.Done()
					io.Copy(localConn, remoteConn)
					closeConnExitCh()
					Logger.Debug("SSH tunnel copy2 exit")
				})

				select {
				case <-connExitCh:
					// 连接正常结束
					Logger.Debug("SSH tunnel connection closed normally")
				case <-tunnel.exitCh:
					// 隧道被停止，需要中断连接
					Logger.Debug("SSH tunnel stopping, closing connections")
					localConn.Close()
					remoteConn.Close()
					<-connExitCh
				}

				Logger.Debug("SSH tunnel local forwarding connection closed")
			})
		}
	})

	return tunnel.ID, nil
}

// StartRemote 远端端口转发（remote forwarding）。
// 在远端（SSH 服务器）监听 127.0.0.1:remotePort，接受连接后转发到本地的 localAddr。
// remotePort: 远端监听端口（远端地址固定 127.0.0.1）
// localAddr 格式: ip:port，例如 127.0.0.1:3306
func (t *SSHTunnelService) StartRemote(sessionID string, remotePort int, localAddr string) (string, error) {
	connAny, ok := t.sshService.SSHConnects.Load(sessionID)
	if !ok {
		return "", sshSessionNotFoundErr
	}
	sc := connAny.(*SSHConnect)

	remoteAddr := fmt.Sprintf("127.0.0.1:%d", remotePort)
	// 在远端通过 ssh client 监听
	ln, err := sc.client.Listen("tcp", remoteAddr)
	if err != nil {
		return "", err
	}

	tunnel := &SSHTunnel{
		ID:         utils.GenerateRandomID(),
		tunnelType: tunnelTypeRemote,
		sessionID:  sessionID,
		LocalAddr:  localAddr,
		RemoteAddr: remoteAddr,
		listener:   ln,
		exitCh:     make(chan struct{}),
	}
	sshTunnels.Store(tunnel.ID, tunnel)

	tunnel.wg.Add(1)
	system.SafeGo(func() {
		defer tunnel.wg.Done()
		for {
			remoteConn, err := ln.Accept()
			if err != nil {
				Logger.Debug("remote accept loop quit", zap.String("tunnelID", tunnel.ID), zap.Error(err))
				return
			}
			tunnel.wg.Add(1)
			system.SafeGo(func() {
				defer tunnel.wg.Done()
				defer remoteConn.Close()

				localConn, err := net.Dial("tcp", tunnel.LocalAddr)
				if err != nil {
					Logger.Debug("failed to dial local", zap.String("tunnelID", tunnel.ID), zap.Error(err))
					return
				}
				defer localConn.Close()

				Logger.Debug("SSH tunnel remote forwarding connection established")

				connExitCh := make(chan struct{})
				var closeOnce sync.Once
				closeConnExitCh := func() {
					closeOnce.Do(func() { close(connExitCh) })
				}

				tunnel.wg.Add(2)
				system.SafeGo(func() {
					defer tunnel.wg.Done()
					io.Copy(localConn, remoteConn)
					closeConnExitCh()
					Logger.Debug("SSH tunnel copy1 exit")
				})
				system.SafeGo(func() {
					defer tunnel.wg.Done()
					io.Copy(remoteConn, localConn)
					closeConnExitCh()
					Logger.Debug("SSH tunnel copy2 exit")
				})

				select {
				case <-connExitCh:
					Logger.Debug("SSH tunnel connection closed normally")
				case <-tunnel.exitCh:
					Logger.Debug("SSH tunnel stopping, closing connections")
					localConn.Close()
					remoteConn.Close()
					<-connExitCh
				}

				Logger.Debug("SSH tunnel remote forwarding connection closed")
			})
		}
	})

	return tunnel.ID, nil
}

// StartDynamic 动态端口转发（SOCKS5）。
// 本地监听 SOCKS5 请求，并通过 SSH 客户端发起出站连接。
// localAddr 格式: ip:port，例如 127.0.0.1:1080
func (t *SSHTunnelService) StartDynamic(sessionID string, localAddr string) (string, error) {
	connAny, ok := t.sshService.SSHConnects.Load(sessionID)
	if !ok {
		return "", sshSessionNotFoundErr
	}
	sc := connAny.(*SSHConnect)

	ln, err := net.Listen("tcp", localAddr)
	if err != nil {
		return "", err
	}

	tunnel := &SSHTunnel{
		ID:          utils.GenerateRandomID(),
		tunnelType:  tunnelTypeDynamic,
		sessionID:   sessionID,
		LocalAddr:   localAddr,
		RemoteAddr:  "",
		listener:    ln,
		exitCh:      make(chan struct{}),
		activeConns: make(map[net.Conn]struct{}),
	}
	sshTunnels.Store(tunnel.ID, tunnel)

	opts := []socks5.Option{
		socks5.WithDial(func(ctx context.Context, network, address string) (net.Conn, error) {
			to, cancel := context.WithTimeout(ctx, time.Second*20)
			defer cancel()
			defer func() {
				if to.Err() != nil {
					Logger.Debug(to.Err().Error(), zap.String("network", network), zap.String("address", address))
				}
			}()
			return sc.client.DialContext(to, network, address)
		}),
	}
	if Mode != ModeRelease {
		opts = append(opts, socks5.WithLogger(socks5.NewLogger(log.New(os.Stdout, "socks5: ", log.LstdFlags))))
	}
	tunnel.socksServer = socks5.NewServer(opts...)
	tunnel.wg.Add(1)
	system.SafeGo(func() {
		defer tunnel.wg.Done()
		t.handleDynamicConn(tunnel, ln)
	})

	return tunnel.ID, nil
}

// handleDynamicConn 处理动态转发的 SOCKS5 连接
func (t *SSHTunnelService) handleDynamicConn(tunnel *SSHTunnel, ln net.Listener) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			tunnel.connLock.Lock()
			for actConn := range tunnel.activeConns {
				_ = actConn.Close()
			}
			tunnel.connLock.Unlock()
			Logger.Debug("socks5 accept loop quit", zap.Error(err))
			return
		}
		tunnel.connLock.Lock()
		tunnel.activeConns[conn] = struct{}{}
		tunnel.connLock.Unlock()

		tunnel.wg.Add(1)
		system.SafeGo(func() {
			Logger.Debug("socks5 conn started")
			defer tunnel.wg.Done()
			defer func() {
				conn.Close()
				tunnel.connLock.Lock()
				delete(tunnel.activeConns, conn)
				tunnel.connLock.Unlock()
			}()

			_ = tunnel.socksServer.ServeConn(conn)
			Logger.Debug("socks5 conn stopped")
		})
	}
}

// StopByID 根据隧道 ID 停止隧道
func (t *SSHTunnelService) StopByID(tunnelID string) error {
	tunnelAny, ok := sshTunnels.Load(tunnelID)
	if !ok {
		return fmt.Errorf("tunnel not found")
	}
	tunnel := tunnelAny.(*SSHTunnel)
	_ = tunnel.listener.Close()
	close(tunnel.exitCh)
	tunnel.wg.Wait()
	sshTunnels.Delete(tunnelID)
	Logger.Debug("SSH tunnel stopped", zap.String("remoteAddr", tunnel.RemoteAddr))
	return nil
}

// StopAll 停止所有隧道（local/remote/dynamic）
func (t *SSHTunnelService) StopAll() {
	sshTunnels.Range(func(key, value any) bool {
		tunnel, ok := value.(*SSHTunnel)
		if ok {
			_ = tunnel.listener.Close()
			close(tunnel.exitCh)
			tunnel.wg.Wait()
			sshTunnels.Delete(key)
		}
		return true
	})
	Logger.Debug("All SSH tunnels stopped")
}

// StopAllBySession 停止指定 sessionID 的所有隧道
func (t *SSHTunnelService) StopAllBySession(sessionID string) {
	sshTunnels.Range(func(key, value any) bool {
		tunnel, ok := value.(*SSHTunnel)
		if ok && tunnel.sessionID == sessionID {
			_ = tunnel.listener.Close()
			close(tunnel.exitCh)
			tunnel.wg.Wait()
			sshTunnels.Delete(key)
		}
		return true
	})
	Logger.Debug("SSH tunnels stopped for session", zap.String("sessionID", sessionID))
}
