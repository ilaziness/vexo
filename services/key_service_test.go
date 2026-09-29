package services

import (
	"errors"
	"testing"
)

func TestDialogCancelled(t *testing.T) {
	if !dialogCancelled(errors.New("cancelled by user")) {
		t.Fatal("expected windows cancel")
	}
	if !dialogCancelled(errors.New("Canceled")) {
		t.Fatal("expected cancel variant")
	}
	if dialogCancelled(errors.New("disk full")) {
		t.Fatal("unexpected cancel")
	}
	if dialogCancelled(nil) {
		t.Fatal("nil is not cancel")
	}
}
