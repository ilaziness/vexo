package ssh

import (
	"bufio"
	"context"
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/proxy"
)

// ProxyConfig is an already-resolved dial proxy (HTTP CONNECT or SOCKS5).
type ProxyConfig struct {
	Type     string
	Host     string
	Port     int
	User     string
	Password string
}

func (p ProxyConfig) Enabled() bool {
	typ := strings.ToLower(strings.TrimSpace(p.Type))
	return (typ == "http" || typ == "socks5") &&
		strings.TrimSpace(p.Host) != "" &&
		p.Port >= 1 && p.Port <= 65535
}

// Validate checks type/host/port. label prefixes error text.
func (p ProxyConfig) Validate(label string) error {
	if label == "" {
		label = "代理"
	}
	if !p.Enabled() {
		typ := strings.ToLower(strings.TrimSpace(p.Type))
		if typ != "http" && typ != "socks5" {
			return fmt.Errorf("%s类型无效", label)
		}
		if strings.TrimSpace(p.Host) == "" {
			return fmt.Errorf("%s主机不能为空", label)
		}
		return fmt.Errorf("%s端口必须在 1–65535 之间", label)
	}
	return nil
}

// cacheKey identifies the proxy path for SSH client reuse (no password).
func (p ProxyConfig) cacheKey() string {
	if !p.Enabled() {
		return "direct"
	}
	return fmt.Sprintf("%s|%s|%d|%s", strings.ToLower(strings.TrimSpace(p.Type)), p.Host, p.Port, p.User)
}

func (p ProxyConfig) addr() string {
	return net.JoinHostPort(p.Host, fmt.Sprintf("%d", p.Port))
}

func dialTCP(addr string, timeout time.Duration, p ProxyConfig) (net.Conn, error) {
	if !p.Enabled() {
		return net.DialTimeout("tcp", addr, timeout)
	}
	if strings.ToLower(strings.TrimSpace(p.Type)) == "socks5" {
		return dialSOCKS5(addr, timeout, p)
	}
	return dialHTTPConnect(addr, timeout, p)
}

func dialSOCKS5(addr string, timeout time.Duration, p ProxyConfig) (net.Conn, error) {
	var auth *proxy.Auth
	if p.User != "" || p.Password != "" {
		auth = &proxy.Auth{User: p.User, Password: p.Password}
	}
	base := &net.Dialer{Timeout: timeout}
	d, err := proxy.SOCKS5("tcp", p.addr(), auth, base)
	if err != nil {
		return nil, fmt.Errorf("socks5 proxy: %w", err)
	}
	cd, ok := d.(proxy.ContextDialer)
	if !ok {
		return nil, fmt.Errorf("socks5 proxy does not support context dial")
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return cd.DialContext(ctx, "tcp", addr)
}

func dialHTTPConnect(addr string, timeout time.Duration, p ProxyConfig) (net.Conn, error) {
	conn, err := net.DialTimeout("tcp", p.addr(), timeout)
	if err != nil {
		return nil, fmt.Errorf("http proxy dial: %w", err)
	}
	if err := conn.SetDeadline(time.Now().Add(timeout)); err != nil {
		_ = conn.Close()
		return nil, err
	}
	req := &http.Request{
		Method: http.MethodConnect,
		URL:    &url.URL{Opaque: addr},
		Host:   addr,
		Header: make(http.Header),
	}
	if p.User != "" || p.Password != "" {
		token := base64.StdEncoding.EncodeToString([]byte(p.User + ":" + p.Password))
		req.Header.Set("Proxy-Authorization", "Basic "+token)
	}
	if err := req.Write(conn); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("http proxy CONNECT write: %w", err)
	}
	br := bufio.NewReader(conn)
	resp, err := http.ReadResponse(br, req)
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("http proxy CONNECT read: %w", err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		_ = conn.Close()
		return nil, fmt.Errorf("http proxy CONNECT failed: %s", resp.Status)
	}
	if err := conn.SetDeadline(time.Time{}); err != nil {
		_ = conn.Close()
		return nil, err
	}
	if br.Buffered() > 0 {
		return &bufferedConn{Conn: conn, r: br}, nil
	}
	return conn, nil
}

type bufferedConn struct {
	net.Conn
	r *bufio.Reader
}

func (c *bufferedConn) Read(b []byte) (int, error) {
	return c.r.Read(b)
}
