//go:build windows

package ssh

import (
	"encoding/binary"
	"io"
	"testing"

	"golang.org/x/sys/windows"
)

func TestPageantConnFraming(t *testing.T) {
	pc := &pageantConn{
		query: func(req []byte) ([]byte, error) {
			if len(req) == 0 {
				t.Fatal("empty request")
			}
			return []byte{0, 0, 0, 3, 'a', 'b', 'c'}, nil
		},
	}
	if _, err := pc.Write([]byte{0, 0, 0, 1, 11}); err != nil {
		t.Fatal(err)
	}
	var size [4]byte
	if _, err := io.ReadFull(pc, size[:]); err != nil {
		t.Fatal(err)
	}
	if binary.BigEndian.Uint32(size[:]) != 3 {
		t.Fatalf("size = %d", binary.BigEndian.Uint32(size[:]))
	}
	body := make([]byte, 3)
	if _, err := io.ReadFull(pc, body); err != nil {
		t.Fatal(err)
	}
	if string(body) != "abc" {
		t.Fatalf("body = %q", body)
	}
}

func TestCreatePageantMapping(t *testing.T) {
	sa, err := currentUserMappingSA()
	if err != nil {
		t.Fatal(err)
	}
	mapping, name, err := createPageantMapping(sa)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = windows.CloseHandle(mapping) }()
	if mapping == 0 || name == "" {
		t.Fatalf("mapping = %d, name = %q", mapping, name)
	}
}

func TestNormalizePipe(t *testing.T) {
	got := normalizePipe(`//./pipe/openssh-ssh-agent`)
	want := `\\.\pipe\openssh-ssh-agent`
	if got != want {
		t.Fatalf("normalizePipe = %q, want %q", got, want)
	}
	if !isPipePath(got) {
		t.Fatal("expected pipe path")
	}
}
