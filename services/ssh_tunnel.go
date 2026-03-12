package services

import (
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/ilaziness/vexo/internal/system"
	"github.com/ilaziness/vexo/internal/utils"
	"go.uber.org/zap"
)

// ssh tunnel

const (
	tunnelTypeLocal   = "local"
	tunnelTypeRemote  = "remote"
	tunnelTypeDynamic = "dynamic"
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
	LocalPort  int    `json:"localPort"`
	RemoteAddr string `json:"remoteAddr"`
}

type SSHTunnelService struct {
	sshService *SSHService
}

type SSHTunnel struct {
	ID         string
	tunnelType string
	sessionID  string // 关联的 SSHConnect 的 sessionID
	LocalPort  int
	RemoteAddr string

	listener net.Listener
	exitCh   chan struct{}
	wg       sync.WaitGroup
}

func NewSSHTunnelService(sshService *SSHService) *SSHTunnelService {
	sshTunnelService = &SSHTunnelService{
		sshService: sshService,
	}
	return sshTunnelService
}

func (t *SSHTunnelService) TunnelList() []TunnelList {
	// 按隧道类型分组
	tunnelsByType := make(map[string][]TunnelInfo)

	sshTunnels.Range(func(key, value any) bool {
		tunnel, ok := value.(*SSHTunnel)
		if ok {
			tunnelInfo := TunnelInfo{
				ID:         tunnel.ID,
				TunnelType: tunnel.tunnelType,
				SessionID:  tunnel.sessionID,
				LocalPort:  tunnel.LocalPort,
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
func (t *SSHTunnelService) StartLocal(sessionID string, LocalPort int, RemoteAddr string) (string, error) {
	connAny, ok := t.sshService.SSHConnects.Load(sessionID)
	if !ok {
		return "", fmt.Errorf("ssh connect not found: %s", sessionID)
	}
	sc := connAny.(*SSHConnect)

	// 默认将本地端口转发到远端的 127.0.0.1:LocalPort
	addr := fmt.Sprintf("127.0.0.1:%d", LocalPort)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return "", err
	}

	tunnel := &SSHTunnel{
		ID:         utils.GenerateRandomID(),
		tunnelType: tunnelTypeLocal,
		sessionID:  sessionID,
		LocalPort:  LocalPort,
		RemoteAddr: RemoteAddr,
		listener:   ln,
		exitCh:     make(chan struct{}),
	}
	// 使用隧道ID 作为 key 存储，这样 StopLocalByID(tunnelID) 能正确找到对应隧道
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

// StopLocal 停止本地端口转发
func (t *SSHTunnelService) StopLocal(sessionID string) error {
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
	return nil
}

// StopLocalByID 根据隧道 ID 停止本地隧道
func (t *SSHTunnelService) StopLocalByID(tunnelID string) error {
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

// StartRemote 远端端口转发（remote forwarding）。
// 在远端（SSH 服务器）监听 127.0.0.1:remotePort，接受连接后转发到本地的 127.0.0.1:localPort。
func (t *SSHTunnelService) StartRemote(sessionID string, remotePort int, localPort int) (string, error) {
	connAny, ok := t.sshService.SSHConnects.Load(sessionID)
	if !ok {
		return "", fmt.Errorf("ssh connect not found: %s", sessionID)
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
		LocalPort:  localPort,
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

				localAddr := fmt.Sprintf("127.0.0.1:%d", tunnel.LocalPort)
				localConn, err := net.Dial("tcp", localAddr)
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

// StopRemote 停止远端端口转发
func (t *SSHTunnelService) StopRemote(sessionID string) error {
	sshTunnels.Range(func(key, value any) bool {
		tunnel, ok := value.(*SSHTunnel)
		if ok && tunnel.tunnelType == tunnelTypeRemote && tunnel.sessionID == sessionID {
			_ = tunnel.listener.Close()
			close(tunnel.exitCh)
			tunnel.wg.Wait()
			sshTunnels.Delete(key)
		}
		return true
	})
	Logger.Debug("SSH remote tunnels stopped for session", zap.String("sessionID", sessionID))
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
