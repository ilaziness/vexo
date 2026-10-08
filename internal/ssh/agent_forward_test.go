package ssh

import (
	"bytes"
	"errors"
	"testing"

	cryptossh "golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

type stubAgent struct {
	keys    []*agent.Key
	signers []cryptossh.Signer
	listErr error
	signErr error
}

func (s *stubAgent) List() ([]*agent.Key, error) {
	if s.listErr != nil {
		return nil, s.listErr
	}
	return s.keys, nil
}

func (s *stubAgent) Sign(key cryptossh.PublicKey, data []byte) (*cryptossh.Signature, error) {
	if s.signErr != nil {
		return nil, s.signErr
	}
	for _, signer := range s.signers {
		if bytes.Equal(signer.PublicKey().Marshal(), key.Marshal()) {
			return signer.Sign(nil, data)
		}
	}
	return nil, errors.New("agent: key not found")
}

func (s *stubAgent) Add(agent.AddedKey) error             { return errors.New("unsupported") }
func (s *stubAgent) Remove(cryptossh.PublicKey) error     { return errors.New("unsupported") }
func (s *stubAgent) RemoveAll() error                     { return errors.New("unsupported") }
func (s *stubAgent) Lock([]byte) error                    { return errors.New("unsupported") }
func (s *stubAgent) Unlock([]byte) error                  { return errors.New("unsupported") }
func (s *stubAgent) Signers() ([]cryptossh.Signer, error) { return s.signers, s.listErr }

func TestMultiAgentListDedup(t *testing.T) {
	a := newTestSigner(t)
	blob := a.PublicKey().Marshal()
	key := &agent.Key{Format: a.PublicKey().Type(), Blob: blob, Comment: "a"}
	dup := &agent.Key{Format: a.PublicKey().Type(), Blob: blob, Comment: "b"}
	m := &multiAgent{agents: []agent.Agent{
		&stubAgent{keys: []*agent.Key{key}},
		&stubAgent{keys: []*agent.Key{dup}},
	}}
	got, err := m.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 {
		t.Fatalf("list len = %d", len(got))
	}
}

func TestMultiAgentSignTriesBackends(t *testing.T) {
	match := newTestSigner(t)
	other := newTestSigner(t)
	m := &multiAgent{agents: []agent.Agent{
		&stubAgent{signErr: errors.New("missing"), signers: []cryptossh.Signer{other}},
		&stubAgent{signers: []cryptossh.Signer{match}},
	}}
	sig, err := m.Sign(match.PublicKey(), []byte("payload"))
	if err != nil {
		t.Fatal(err)
	}
	if sig == nil {
		t.Fatal("expected signature")
	}
}

func TestMultiAgentListAllFail(t *testing.T) {
	m := &multiAgent{agents: []agent.Agent{
		&stubAgent{listErr: errors.New("down")},
	}}
	if _, err := m.List(); err == nil {
		t.Fatal("expected error")
	}
}
