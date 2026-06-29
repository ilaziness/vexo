package system

import "testing"

func TestParseRemoteSystemInfo(t *testing.T) {
	output := `myserver
Linux 5.15.0-91-generic
x86_64
PRETTY_NAME=Ubuntu 22.04.3 LTS
ID=ubuntu
VERSION_ID=22.04
`
	info := ParseRemoteSystemInfo(output)
	if !info.Ready {
		t.Fatal("expected ready")
	}
	if info.Hostname != "myserver" {
		t.Errorf("hostname: got %q", info.Hostname)
	}
	if info.Kernel != "Linux 5.15.0-91-generic" {
		t.Errorf("kernel: got %q", info.Kernel)
	}
	if info.Arch != "x86_64" {
		t.Errorf("arch: got %q", info.Arch)
	}
	if info.OSPretty != "Ubuntu 22.04.3 LTS" {
		t.Errorf("os pretty: got %q", info.OSPretty)
	}
	if info.OSID != "ubuntu" {
		t.Errorf("os id: got %q", info.OSID)
	}
	if info.OSVersion != "22.04" {
		t.Errorf("os version: got %q", info.OSVersion)
	}
	if !info.HasContent() {
		t.Fatal("expected content")
	}
}

func TestParseRemoteSystemInfoSkipsEmptyLeadingLines(t *testing.T) {
	output := `

host1
Linux 6.1.0
aarch64
`
	info := ParseRemoteSystemInfo(output)
	if info.Hostname != "host1" {
		t.Errorf("hostname: got %q", info.Hostname)
	}
	if info.Kernel != "Linux 6.1.0" {
		t.Errorf("kernel: got %q", info.Kernel)
	}
}

func TestNormalizeHost(t *testing.T) {
	if got := NormalizeHost(" 192.168.1.1 "); got != "192.168.1.1" {
		t.Errorf("got %q", got)
	}
}
