package system

import (
	"strings"
)

// RemoteSystemInfo 远程主机基本信息（不含敏感数据）
type RemoteSystemInfo struct {
	Hostname  string `json:"hostname,omitempty"`
	OSPretty  string `json:"os_pretty,omitempty"`
	OSID      string `json:"os_id,omitempty"`
	OSVersion string `json:"os_version,omitempty"`
	Kernel    string `json:"kernel,omitempty"`
	Arch      string `json:"arch,omitempty"`
	Ready     bool   `json:"ready"`
}

// HasContent 是否至少有一个有效字段
func (r *RemoteSystemInfo) HasContent() bool {
	if r == nil {
		return false
	}
	return r.Hostname != "" || r.OSPretty != "" || r.OSID != "" ||
		r.OSVersion != "" || r.Kernel != "" || r.Arch != ""
}

// NormalizeHost 规范化 host 作为缓存 key
func NormalizeHost(host string) string {
	return strings.ToLower(strings.TrimSpace(host))
}

// ParseRemoteSystemInfo 解析远程采集脚本输出
func ParseRemoteSystemInfo(output string) RemoteSystemInfo {
	info := RemoteSystemInfo{Ready: true}
	lines := strings.Split(output, "\n")

	var scalar []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.Contains(line, "=") {
			key, val, ok := strings.Cut(line, "=")
			if !ok {
				continue
			}
			val = strings.Trim(val, `"`)
			switch key {
			case "PRETTY_NAME":
				info.OSPretty = val
			case "ID":
				info.OSID = val
			case "VERSION_ID":
				info.OSVersion = val
			}
			continue
		}
		if len(scalar) < 3 {
			scalar = append(scalar, line)
		}
	}

	if len(scalar) > 0 {
		info.Hostname = scalar[0]
	}
	if len(scalar) > 1 {
		info.Kernel = scalar[1]
	}
	if len(scalar) > 2 {
		info.Arch = scalar[2]
	}

	return info
}
