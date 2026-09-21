package transfer

import "testing"

func TestTrackerStopIdempotent(t *testing.T) {
	var n int
	r := NewRegistry(func(ProgressData) { n++ })
	tr := r.New("sess", TypeUpload, "local", "remote", 100)
	tr.Stop(nil)
	tr.Stop(nil)
	if n < 2 {
		t.Fatalf("expected at least create+done events, got %d", n)
	}
}
