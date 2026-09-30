package ssh

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCheckKnownHostSameTypeMismatch(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "known_hosts")
	if err := os.WriteFile(path, []byte("h:22 ssh-ed25519 OLDKEY\n"), 0600); err != nil {
		t.Fatal(err)
	}
	found, mismatch, oldType, oldKey, err := checkKnownHostInFile(path, "h:22", "ssh-ed25519", "NEWKEY")
	if err != nil {
		t.Fatal(err)
	}
	if found || !mismatch {
		t.Fatalf("found=%v mismatch=%v, want mismatch", found, mismatch)
	}
	if oldType != "ssh-ed25519" || oldKey != "OLDKEY" {
		t.Fatalf("old=%s %s", oldType, oldKey)
	}
}

func TestCheckKnownHostDifferentTypeIsUnknown(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "known_hosts")
	if err := os.WriteFile(path, []byte("h:22 ssh-ed25519 OLDKEY\n"), 0600); err != nil {
		t.Fatal(err)
	}
	found, mismatch, _, _, err := checkKnownHostInFile(path, "h:22", "ssh-rsa", "RSAKEY")
	if err != nil {
		t.Fatal(err)
	}
	if found || mismatch {
		t.Fatalf("found=%v mismatch=%v, want unknown", found, mismatch)
	}
}

func TestReplaceKnownHost(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "known_hosts")
	content := "a:22 ssh-ed25519 AAA\nb:22 ssh-ed25519 BBB\na:22 ssh-rsa CCC\n"
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	if err := replaceKnownHost(path, "a:22", "ssh-ed25519", "NEW"); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	got := string(data)
	wantParts := []string{
		"a:22 ssh-ed25519 NEW",
		"b:22 ssh-ed25519 BBB",
		"a:22 ssh-rsa CCC",
	}
	for _, part := range wantParts {
		if !containsLine(got, part) {
			t.Fatalf("missing %q in:\n%s", part, got)
		}
	}
}

func TestDeleteKnownHost(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "known_hosts")
	content := "a:22 ssh-ed25519 AAA\na:22 ssh-rsa BBB\nb:22 ssh-ed25519 CCC\n"
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	if err := deleteKnownHost(path, "a:22", "ssh-ed25519"); err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(path)
	got := string(data)
	if containsLine(got, "a:22 ssh-ed25519 AAA") {
		t.Fatal("ed25519 entry should be deleted")
	}
	if !containsLine(got, "a:22 ssh-rsa BBB") {
		t.Fatal("rsa entry should remain")
	}
}

func containsLine(data, line string) bool {
	for _, l := range splitLines(data) {
		if l == line {
			return true
		}
	}
	return false
}

func splitLines(s string) []string {
	var out []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			if i > start {
				out = append(out, s[start:i])
			}
			start = i + 1
		}
	}
	if start < len(s) {
		out = append(out, s[start:])
	}
	return out
}
