package ssh

import (
	"strings"
	"testing"
)

func TestTruncateUTF8(t *testing.T) {
	if got := truncateUTF8("abc", 10); got != "abc" {
		t.Fatalf("short: %q", got)
	}
	got := truncateUTF8("你好世界", 4)
	if !strings.HasSuffix(got, "…[truncated]") {
		t.Fatalf("expected truncation marker, got %q", got)
	}
	body := strings.TrimSuffix(got, "\n…[truncated]")
	if body != "你" {
		t.Fatalf("expected first complete rune, got %q", body)
	}
}
