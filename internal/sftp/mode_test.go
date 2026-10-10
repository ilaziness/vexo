package sftp

import (
	"os"
	"testing"
)

func TestFileModeUnixRoundTrip(t *testing.T) {
	cases := []struct {
		unix uint32
		name string
	}{
		{0o755, "755"},
		{0o644, "644"},
		{0o4755, "setuid"},
		{0o2755, "setgid"},
		{0o1755, "sticky"},
		{0o6777, "setuid+setgid"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fm := unixToFileMode(tc.unix)
			got := fileModeToUnix(fm)
			if got != tc.unix {
				t.Fatalf("round-trip %o -> FileMode -> %o", tc.unix, got)
			}
		})
	}
}

func TestUnixToFileModeSpecialBits(t *testing.T) {
	m := unixToFileMode(0o4755)
	if m&os.ModeSetuid == 0 {
		t.Fatal("expected ModeSetuid")
	}
	if m.Perm() != 0o755 {
		t.Fatalf("perm=%o", m.Perm())
	}
}
