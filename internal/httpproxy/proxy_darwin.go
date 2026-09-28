//go:build darwin

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
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "scutil", "--proxy").Output()
	if err != nil {
		return nil, nil
	}
	vals := map[string]string{}
	for _, line := range strings.Split(string(out), "\n") {
		key, val, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		vals[strings.TrimSpace(key)] = strings.TrimSpace(val)
	}
	scheme := req.URL.Scheme
	if scheme == "https" && vals["HTTPSEnable"] == "1" && vals["HTTPSProxy"] != "" {
		return proxyURL("http", netJoin(vals["HTTPSProxy"], vals["HTTPSPort"]))
	}
	if vals["HTTPEnable"] == "1" && vals["HTTPProxy"] != "" {
		return proxyURL("http", netJoin(vals["HTTPProxy"], vals["HTTPPort"]))
	}
	if vals["SOCKSEnable"] == "1" && vals["SOCKSProxy"] != "" {
		return proxyURL("socks5", netJoin(vals["SOCKSProxy"], vals["SOCKSPort"]))
	}
	return nil, nil
}

func netJoin(host, port string) string {
	if port == "" {
		return host
	}
	if _, err := strconv.Atoi(port); err != nil {
		return host
	}
	return host + ":" + port
}
