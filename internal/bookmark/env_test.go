package bookmark

import (
	"testing"
)

func TestParseFormatEnvVars(t *testing.T) {
	raw, err := FormatEnvVars(map[string]string{"FOO": "bar", "TERM": "xterm"})
	if err != nil {
		t.Fatal(err)
	}
	m, err := ParseEnvVars(raw)
	if err != nil {
		t.Fatal(err)
	}
	if m["FOO"] != "bar" {
		t.Fatalf("FOO=%q", m["FOO"])
	}
	if _, ok := m["TERM"]; ok {
		t.Fatal("TERM should be stripped")
	}
	empty, err := ParseEnvVars("")
	if err != nil || len(empty) != 0 {
		t.Fatalf("empty: %v %v", empty, err)
	}
	if _, err := ParseEnvVars("{bad"); err == nil {
		t.Fatal("expected error")
	}
}
