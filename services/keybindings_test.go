package services

import (
	"runtime"
	"testing"
)

func TestNormalizeAccel(t *testing.T) {
	got := normalizeAccel("CmdOrCtrl+Shift+C")
	wantMod := "Ctrl"
	if runtime.GOOS == "darwin" {
		wantMod = "Cmd"
	}
	want := wantMod + "+Shift+C"
	if got != want {
		t.Fatalf("normalizeAccel(CmdOrCtrl+Shift+C) = %q, want %q", got, want)
	}

	got = normalizeAccel("CmdOrCtrl+F")
	want = wantMod + "+F"
	if got != want {
		t.Fatalf("normalizeAccel(CmdOrCtrl+F) = %q, want %q", got, want)
	}

	got = normalizeAccel("CmdOrCtrl+,")
	want = wantMod + "+,"
	if got != want {
		t.Fatalf("normalizeAccel(CmdOrCtrl+,) = %q, want %q", got, want)
	}
}

func TestAccelAliasesWindowsComma(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("OEM_COMMA alias is Windows-only")
	}
	aliases := accelAliases("Ctrl+,")
	if len(aliases) != 1 || aliases[0] != "Ctrl+OEM_COMMA" {
		t.Fatalf("accelAliases(Ctrl+,) = %v, want [Ctrl+OEM_COMMA]", aliases)
	}
}

func TestDisplayKeys(t *testing.T) {
	got := displayKeys("CmdOrCtrl+Shift+C")
	if runtime.GOOS == "darwin" {
		if got != "⌘⇧C" {
			t.Fatalf("displayKeys on darwin = %q, want ⌘⇧C", got)
		}
		return
	}
	if got != "Ctrl+Shift+C" {
		t.Fatalf("displayKeys = %q, want Ctrl+Shift+C", got)
	}
}
