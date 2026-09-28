package httpproxy

import "testing"

func TestParseProxyServer(t *testing.T) {
	u, err := parseProxyServer("127.0.0.1:7890", "https")
	if err != nil || u == nil || u.String() != "http://127.0.0.1:7890" {
		t.Fatalf("bare host: %v %v", u, err)
	}
	u, err = parseProxyServer("http=127.0.0.1:7890;https=127.0.0.1:7891", "https")
	if err != nil || u == nil || u.Host != "127.0.0.1:7891" {
		t.Fatalf("https entry: %v %v", u, err)
	}
	u, err = parseProxyServer("http=127.0.0.1:7890;socks=127.0.0.1:1080", "https")
	if err != nil || u == nil || u.String() != "http://127.0.0.1:7890" {
		t.Fatalf("http before socks: %v %v", u, err)
	}
	u, err = parseProxyServer("socks=127.0.0.1:1080", "https")
	if err != nil || u == nil || u.Scheme != "socks5" || u.Host != "127.0.0.1:1080" {
		t.Fatalf("socks: %v %v", u, err)
	}
}

func TestBypassHost(t *testing.T) {
	if !bypassHost("localhost", "<local>;example.com") {
		t.Fatal("local")
	}
	if bypassHost("github.com", "<local>") {
		t.Fatal("github is not local")
	}
	if !bypassHost("api.example.com", "*.example.com") {
		t.Fatal("suffix")
	}
}
