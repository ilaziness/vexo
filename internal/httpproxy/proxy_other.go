//go:build !windows && !darwin

package httpproxy

import (
	"context"
	"net/http"
	"net/url"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func systemProxy(req *http.Request) (*url.URL, error) {
	if req == nil || req.URL == nil {
		return nil, nil
	}
	if _, err := exec.LookPath("gsettings"); err != nil {
		return nil, nil
	}
	mode := gsettings("org.gnome.system.proxy", "mode")
	if mode != "manual" {
		return nil, nil
	}
	ignore := gsettings("org.gnome.system.proxy", "ignore-hosts")
	if bypassHost(req.URL.Hostname(), strings.NewReplacer("'", "", "[", "", "]", "", ",", ";", " ", "").Replace(ignore)) {
		return nil, nil
	}
	key := "org.gnome.system.proxy.http"
	if req.URL.Scheme == "https" {
		key = "org.gnome.system.proxy.https"
	}
	host, port := proxyEndpoint(key)
	if !validEndpoint(host, port) && req.URL.Scheme == "https" {
		host, port = proxyEndpoint("org.gnome.system.proxy.http")
	}
	if validEndpoint(host, port) {
		return proxyURL("http", host+":"+port)
	}
	host, port = proxyEndpoint("org.gnome.system.proxy.socks")
	if !validEndpoint(host, port) {
		return nil, nil
	}
	return proxyURL("socks5", host+":"+port)
}

func proxyEndpoint(schema string) (string, string) {
	return gsettings(schema, "host"), gsettings(schema, "port")
}

func validEndpoint(host, port string) bool {
	if host == "" || host == "0" || port == "" || port == "0" {
		return false
	}
	_, err := strconv.Atoi(port)
	return err == nil
}

func gsettings(schema, key string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "gsettings", "get", schema, key).Output()
	if err != nil {
		return ""
	}
	s := strings.TrimSpace(string(out))
	s = strings.Trim(s, "'")
	if fields := strings.Fields(s); len(fields) == 2 {
		s = strings.Trim(fields[1], "'")
	}
	if s == "" || s == "''" || s == "@as []" {
		return ""
	}
	return s
}
