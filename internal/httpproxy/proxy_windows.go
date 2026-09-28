//go:build windows

package httpproxy

import (
	"net/http"
	"net/url"

	"golang.org/x/sys/windows/registry"
)

func systemProxy(req *http.Request) (*url.URL, error) {
	if req == nil || req.URL == nil {
		return nil, nil
	}
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Internet Settings`, registry.QUERY_VALUE)
	if err != nil {
		return nil, nil
	}
	defer k.Close()

	enable, _, err := k.GetIntegerValue("ProxyEnable")
	if err != nil || enable == 0 {
		return nil, nil
	}
	server, _, err := k.GetStringValue("ProxyServer")
	if err != nil || server == "" {
		return nil, nil
	}
	if override, _, err := k.GetStringValue("ProxyOverride"); err == nil && bypassHost(req.URL.Hostname(), override) {
		return nil, nil
	}
	return parseProxyServer(server, req.URL.Scheme)
}
