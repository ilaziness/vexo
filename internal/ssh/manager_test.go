package ssh

import (
	"testing"
	"time"
)

func TestHopsClientKey(t *testing.T) {
	hops := []Endpoint{
		{Host: "jump", Port: 22, User: "j"},
		{Host: "target", Port: 22, User: "t"},
	}
	got := hopsClientKey(hops)
	want := "t@target:22via:j@jump:22"
	if got != want {
		t.Fatalf("hopsClientKey = %q, want %q", got, want)
	}
}

func TestHopClientCloseNil(t *testing.T) {
	var h *hopClient
	h.Close()
	(&hopClient{}).Close()
}

func TestHostKeyDecideIndependentHosts(t *testing.T) {
	s := newHostKeyStore()
	if err := s.decide("a", true); err == nil {
		t.Fatal("expected error when no pending prompt")
	}
	chA, err := s.begin("a")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.begin("a"); err == nil {
		t.Fatal("expected duplicate host prompt to fail")
	}
	chB, err := s.begin("b")
	if err != nil {
		t.Fatal(err)
	}
	if err := s.decide("a", true); err != nil {
		t.Fatal(err)
	}
	if err := s.decide("b", false); err != nil {
		t.Fatal(err)
	}
	select {
	case v := <-chA:
		if !v {
			t.Fatal("expected accept for host a")
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for host a")
	}
	select {
	case v := <-chB:
		if v {
			t.Fatal("expected reject for host b")
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for host b")
	}
	s.finish("a")
	s.finish("b")
	if _, err := s.begin("a"); err != nil {
		t.Fatal(err)
	}
}
