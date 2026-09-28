package tools

import (
	"strings"
	"testing"
)

func TestConvertChmodRoundTrip(t *testing.T) {
	s := New()
	cases := []struct {
		octal    string
		wantSym  string
		wantBack string
	}{
		{"755", "rwxr-xr-x", "755"},
		{"644", "rw-r--r--", "644"},
		{"4755", "rwsr-xr-x", "4755"},
		{"1777", "rwxrwxrwt", "1777"},
		{"0", "---------", "000"},
	}
	for _, tc := range cases {
		sym := s.ConvertChmod(tc.octal, "toSymbolic")
		if !sym.Success || sym.Symbolic != tc.wantSym {
			t.Fatalf("toSymbolic(%s): success=%v symbolic=%q err=%s", tc.octal, sym.Success, sym.Symbolic, sym.Error)
		}
		back := s.ConvertChmod(sym.Symbolic, "toOctal")
		if !back.Success || back.Octal != tc.wantBack {
			t.Fatalf("toOctal(%s): success=%v octal=%q err=%s", sym.Symbolic, back.Success, back.Octal, back.Error)
		}
	}
}

func TestConvertChmodRejectsWrongCharPosition(t *testing.T) {
	s := New()
	res := s.ConvertChmod("wxr------", "toOctal")
	if res.Success {
		t.Fatalf("expected failure for mis-ordered rwx, got octal=%s", res.Octal)
	}
}

func TestConvertChmodFromLsListing(t *testing.T) {
	s := New()
	res := s.ConvertChmod("drwxr-xr-x", "toOctal")
	if !res.Success || res.Octal != "755" {
		t.Fatalf("got success=%v octal=%s err=%s", res.Success, res.Octal, res.Error)
	}
}

func TestConvertBasePrefixOnlyForMatchingBase(t *testing.T) {
	s := New()
	fail := s.ConvertBase("0x10", 10, 10)
	if fail.Success {
		t.Fatalf("expected failure parsing 0x10 as base 10, got %s", fail.Result)
	}
	ok := s.ConvertBase("0x10", 16, 10)
	if !ok.Success || ok.Result != "16" {
		t.Fatalf("0x10 base16->10: success=%v result=%s", ok.Success, ok.Result)
	}
	bin := s.ConvertBase("0b1010", 2, 10)
	if !bin.Success || bin.Result != "10" {
		t.Fatalf("0b1010 base2->10: success=%v result=%s", bin.Success, bin.Result)
	}
}

func TestCalculateCIDRKeepsResultOnBadCheckIP(t *testing.T) {
	s := New()
	res := s.CalculateCIDR("192.168.1.0/24", "not-an-ip")
	if !res.Success {
		t.Fatal("expected Success=true with network still computed")
	}
	if res.Network != "192.168.1.0/24" {
		t.Fatalf("network=%s", res.Network)
	}
	if res.Error == "" {
		t.Fatal("expected Error explaining invalid check IP")
	}
	if res.Contains != nil {
		t.Fatal("Contains should be nil when check IP invalid")
	}
}

func TestCalculateCIDRFamilyMismatch(t *testing.T) {
	s := New()
	res := s.CalculateCIDR("192.168.1.0/24", "2001:db8::1")
	if !res.Success || res.Error == "" {
		t.Fatalf("expected success with family mismatch error, got success=%v err=%s", res.Success, res.Error)
	}
}

func TestDiffTextLimit(t *testing.T) {
	s := New()
	big := strings.Repeat("x\n", maxDiffLines+1)
	res := s.DiffText(big, "a")
	if res.Success {
		t.Fatal("expected failure for oversized input")
	}
}

func TestDiffTextBasic(t *testing.T) {
	s := New()
	res := s.DiffText("a\nb\nc", "a\nx\nc")
	if !res.Success || len(res.Lines) != 4 {
		t.Fatalf("success=%v lines=%d", res.Success, len(res.Lines))
	}
}

func TestDecodeJWTRejectsEmptyParts(t *testing.T) {
	s := New()
	res := s.DecodeJWT("aaa..bbb")
	if res.Success {
		t.Fatal("expected failure for empty payload")
	}
}

func TestDecodeJWTBearerPrefix(t *testing.T) {
	s := New()
	token := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxIn0.sig"
	res := s.DecodeJWT(token)
	if !res.Success {
		t.Fatalf("expected success, err=%s", res.Error)
	}
}

func TestGenerateRandomUnknownCharset(t *testing.T) {
	s := New()
	res := s.GenerateRandom("string", 8, "nope", "")
	if res.Success {
		t.Fatal("expected failure for unknown charset")
	}
}

func TestConvertChmodLeadingZeroOctal(t *testing.T) {
	s := New()
	res := s.ConvertChmod("0755", "toSymbolic")
	if !res.Success || res.Symbolic != "rwxr-xr-x" {
		t.Fatalf("got success=%v symbolic=%q err=%s", res.Success, res.Symbolic, res.Error)
	}
}

func TestConvertTimestampDateTimeLocal(t *testing.T) {
	s := New()
	res := s.ConvertTimestamp("2026-04-01T15:04", true)
	if !res.Success {
		t.Fatalf("datetime-local parse failed: %s", res.Error)
	}
}

func TestConvertTimestampMillis(t *testing.T) {
	s := New()
	// 2026-04-01 00:00:00 UTC roughly — use a known ms value
	res := s.ConvertTimestamp("1711929600000", false) // 2024-04-01 00:00:00 UTC
	if !res.Success || res.Datetime == "" {
		t.Fatalf("ms timestamp failed: success=%v err=%s", res.Success, res.Error)
	}
	// Should not be year 57000+
	if strings.HasPrefix(res.Datetime, "57") || len(res.Datetime) < 10 {
		t.Fatalf("unexpected datetime for ms input: %s", res.Datetime)
	}
}
