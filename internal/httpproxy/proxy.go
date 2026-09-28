package httpproxy

import (
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// Client is an HTTP client for update checks and downloads.
// It uses environment proxies first, then the OS system proxy.
// There is no overall timeout so a large download is not cut off mid-body.
func Client() *http.Client {
	base, ok := http.DefaultTransport.(*http.Transport)
	if !ok {
		return &http.Client{Transport: &http.Transport{
			Proxy:                 Proxy,
			ResponseHeaderTimeout: 30 * time.Second,
		}}
	}
	tr := base.Clone()
	tr.Proxy = Proxy
	tr.ResponseHeaderTimeout = 30 * time.Second
	return &http.Client{Transport: tr}
}

// Proxy resolves the proxy for req. An explicit HTTP_PROXY/HTTPS_PROXY wins.
// If those variables are set but this host is excluded (NO_PROXY or localhost),
// the system proxy is not used either.
func Proxy(req *http.Request) (*url.URL, error) {
	if req == nil {
		return nil, nil
	}
	u, err := http.ProxyFromEnvironment(req)
	if err != nil || u != nil {
		return u, err
	}
	if envProxyConfigured(req) {
		return nil, nil
	}
	return cachedSystemProxy(req)
}

// envProxyConfigured reports whether the process has an environment proxy that
// would apply to some host of req's scheme. ProxyFromEnvironment returns nil
// both when no proxy is configured and when the host is excluded.
func envProxyConfigured(req *http.Request) bool {
	scheme := "https"
	if req.URL != nil && req.URL.Scheme != "" {
		scheme = req.URL.Scheme
	}
	probe, err := http.NewRequest(http.MethodGet, scheme+"://proxy-probe.invalid/", nil)
	if err != nil {
		return false
	}
	u, err := http.ProxyFromEnvironment(probe)
	return err == nil && u != nil
}

// systemProxy hits the registry or an external process. Cache it so a download
// and its redirects do not repeat that lookup.
var systemProxyCache struct {
	sync.Mutex
	at  time.Time
	key string
	url *url.URL
	err error
}

func cachedSystemProxy(req *http.Request) (*url.URL, error) {
	scheme, host := "", ""
	if req.URL != nil {
		scheme = req.URL.Scheme
		host = req.URL.Hostname()
	}
	key := scheme + " " + host
	systemProxyCache.Lock()
	if systemProxyCache.key == key && time.Since(systemProxyCache.at) < 30*time.Second {
		u, err := systemProxyCache.url, systemProxyCache.err
		systemProxyCache.Unlock()
		return u, err
	}
	systemProxyCache.Unlock()

	u, err := systemProxy(req)

	systemProxyCache.Lock()
	systemProxyCache.at = time.Now()
	systemProxyCache.key = key
	systemProxyCache.url = u
	systemProxyCache.err = err
	systemProxyCache.Unlock()
	return u, err
}

// parseProxyServer parses a Windows proxy server value.
// A bare host:port is an HTTP proxy. A list like "http=host:port;https=host:port;socks=host:port"
// prefers the matching scheme, then the HTTP proxy (CONNECT), and only then SOCKS.
func parseProxyServer(raw, scheme string) (*url.URL, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}
	if !strings.Contains(raw, "=") {
		return proxyURL("http", raw)
	}
	var httpProxy, httpsProxy, socks string
	for _, part := range strings.Split(raw, ";") {
		key, val, ok := strings.Cut(part, "=")
		if !ok {
			continue
		}
		key = strings.ToLower(strings.TrimSpace(key))
		val = strings.TrimSpace(val)
		if val == "" {
			continue
		}
		switch key {
		case "http":
			httpProxy = val
		case "https":
			httpsProxy = val
		case "socks", "socks5":
			socks = val
		}
	}
	switch scheme {
	case "https":
		if httpsProxy != "" {
			return proxyURL("http", httpsProxy)
		}
		if httpProxy != "" {
			return proxyURL("http", httpProxy)
		}
	default:
		if httpProxy != "" {
			return proxyURL("http", httpProxy)
		}
		if httpsProxy != "" {
			return proxyURL("http", httpsProxy)
		}
	}
	if socks != "" {
		return proxyURL("socks5", socks)
	}
	return nil, nil
}

func proxyURL(kind, hostport string) (*url.URL, error) {
	hostport = strings.TrimSpace(hostport)
	if hostport == "" {
		return nil, nil
	}
	if strings.Contains(hostport, "://") {
		return url.Parse(hostport)
	}
	switch strings.ToLower(kind) {
	case "socks", "socks5":
		return url.Parse("socks5://" + hostport)
	default:
		return url.Parse("http://" + hostport)
	}
}

// bypassHost reports whether host should skip the proxy.
// rules is a semicolon-separated list; "<local>" matches names without a dot.
func bypassHost(host, rules string) bool {
	host = strings.TrimSuffix(strings.ToLower(host), ".")
	if host == "" || rules == "" {
		return false
	}
	if h, _, err := net.SplitHostPort(host); err == nil {
		host = h
	}
	for _, rule := range strings.Split(rules, ";") {
		rule = strings.TrimSpace(strings.ToLower(rule))
		if rule == "" {
			continue
		}
		if rule == "<local>" {
			if !strings.Contains(host, ".") {
				return true
			}
			continue
		}
		rule = strings.TrimPrefix(rule, "*.")
		if host == rule || strings.HasSuffix(host, "."+rule) {
			return true
		}
	}
	return false
}
