package ssh

import (
	"os"
	"path/filepath"
	"testing"

	"go.uber.org/zap"
)

func testSession(t *testing.T, id string) (*Manager, *Session) {
	t.Helper()
	m := NewManager(zap.NewNop(), "", nil, nil)
	sess := &Session{
		ID:         id,
		manager:    m,
		OutputChan: make(chan []byte, 8),
		stopOutput: make(chan struct{}),
	}
	m.sessions.Store(id, sess)
	return m, sess
}

func TestStartStopSessionLog(t *testing.T) {
	m, _ := testSession(t, "s1")
	dir := t.TempDir()
	path := filepath.Join(dir, "out.log")

	if err := m.StartSessionLog("s1", path); err != nil {
		t.Fatal(err)
	}
	if !m.IsSessionLogging("s1") {
		t.Fatal("expected logging")
	}
	if err := m.StartSessionLog("s1", path); err == nil {
		t.Fatal("expected error on duplicate start")
	}
	if err := m.StopSessionLog("s1"); err != nil {
		t.Fatal(err)
	}
	if m.IsSessionLogging("s1") {
		t.Fatal("expected not logging")
	}
	if err := m.StopSessionLog("s1"); err != nil {
		t.Fatal(err)
	}
	if err := m.StopSessionLog("missing"); err != nil {
		t.Fatal(err)
	}
}

func TestTeeOutputWritesAndTruncates(t *testing.T) {
	m, sess := testSession(t, "s2")
	dir := t.TempDir()
	path := filepath.Join(dir, "out.log")

	if err := m.StartSessionLog("s2", path); err != nil {
		t.Fatal(err)
	}
	sess.teeOutput([]byte("hello"))
	sess.teeOutput([]byte(" world"))
	if err := m.StopSessionLog("s2"); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "hello world" {
		t.Fatalf("got %q", data)
	}

	if err := m.StartSessionLog("s2", path); err != nil {
		t.Fatal(err)
	}
	sess.teeOutput([]byte("new"))
	_ = m.StopSessionLog("s2")
	data, err = os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "new" {
		t.Fatalf("truncate got %q", data)
	}
}

func TestTeeOutputWriteErrorStops(t *testing.T) {
	m, sess := testSession(t, "s3")
	var gotID, gotReason string
	m.SetOnSessionLogStopped(func(id, reason string) {
		gotID, gotReason = id, reason
	})
	dir := t.TempDir()
	path := filepath.Join(dir, "out.log")
	if err := m.StartSessionLog("s3", path); err != nil {
		t.Fatal(err)
	}
	sess.logMu.Lock()
	_ = sess.logFile.Close()
	sess.logMu.Unlock()

	sess.teeOutput([]byte("x"))
	if m.IsSessionLogging("s3") {
		t.Fatal("expected logging stopped after write error")
	}
	if gotID != "s3" || gotReason != SessionLogStopReasonWriteError {
		t.Fatalf("callback = %q %q", gotID, gotReason)
	}
}

func TestClosePTYStopsLog(t *testing.T) {
	m, sess := testSession(t, "s5")
	var gotReason string
	var sawCloseMu bool
	m.SetOnSessionLogStopped(func(_, reason string) {
		gotReason = reason
		// Callback must not run while closeMu is held (would deadlock on re-entry).
		unlocked := sess.closeMu.TryLock()
		sawCloseMu = unlocked
		if unlocked {
			sess.closeMu.Unlock()
		}
	})
	path := filepath.Join(t.TempDir(), "out.log")
	if err := m.StartSessionLog("s5", path); err != nil {
		t.Fatal(err)
	}
	_ = sess.closePTY()
	if m.IsSessionLogging("s5") {
		t.Fatal("expected not logging after close")
	}
	if gotReason != SessionLogStopReasonClosed {
		t.Fatalf("reason = %q", gotReason)
	}
	if !sawCloseMu {
		t.Fatal("expected closeMu free during log-stopped callback")
	}
}

func TestStartSessionLogRejectsAfterClose(t *testing.T) {
	m, sess := testSession(t, "s6")
	_ = sess.closePTY()
	err := m.StartSessionLog("s6", filepath.Join(t.TempDir(), "x.log"))
	if err == nil {
		t.Fatal("expected error after close")
	}
}
