package ssh

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/pem"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ilaziness/vexo/internal/sshkey"
	cryptossh "golang.org/x/crypto/ssh"
)

func newTestSigner(t *testing.T) cryptossh.Signer {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	signer, err := cryptossh.NewSignerFromKey(key)
	if err != nil {
		t.Fatal(err)
	}
	return signer
}

func methodTypes(methods []cryptossh.AuthMethod) []string {
	out := make([]string, len(methods))
	for i, method := range methods {
		out[i] = fmt.Sprintf("%T", method)
	}
	return out
}

func TestAuthMethodsAgentOnly(t *testing.T) {
	methods, err := authMethods(Endpoint{UseAgent: true}, []cryptossh.Signer{newTestSigner(t)})
	if err != nil {
		t.Fatal(err)
	}
	got := methodTypes(methods)
	if len(got) != 1 || !strings.Contains(got[0], "publicKey") {
		t.Fatalf("methods = %v", got)
	}
}

func TestAuthMethodsAgentWithoutKeys(t *testing.T) {
	methods, err := authMethods(Endpoint{UseAgent: true}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(methods) != 0 {
		t.Fatalf("methods = %v", methodTypes(methods))
	}
}

func TestAuthMethodsAgentAndPassword(t *testing.T) {
	methods, err := authMethods(Endpoint{UseAgent: true, Password: "secret"}, []cryptossh.Signer{newTestSigner(t)})
	if err != nil {
		t.Fatal(err)
	}
	got := methodTypes(methods)
	if len(got) != 2 || !strings.Contains(got[0], "publicKey") || !strings.Contains(got[1], "password") {
		t.Fatalf("methods = %v", got)
	}
}

func TestSerialReadWriterNotCloser(t *testing.T) {
	var rw io.ReadWriter = serialReadWriter{rw: bytes.NewBuffer(nil)}
	if _, ok := rw.(io.Closer); ok {
		t.Fatal("serial wrapper must not implement io.Closer")
	}
}

func TestAuthMethodsEmpty(t *testing.T) {
	methods, err := authMethods(Endpoint{}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(methods) != 0 {
		t.Fatalf("methods = %v", methodTypes(methods))
	}
}

func TestFinalizeAuthMethodsEmpty(t *testing.T) {
	_, err := finalizeAuthMethods(nil, nil)
	if err == nil || !strings.Contains(err.Error(), "没有可用认证") {
		t.Fatalf("err = %v", err)
	}
}

func TestAuthMethodOrder(t *testing.T) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	block, err := cryptossh.MarshalPrivateKey(key, "")
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "id_ecdsa")
	if err := os.WriteFile(path, pem.EncodeToMemory(block), 0o600); err != nil {
		t.Fatal(err)
	}
	methods, err := authMethods(Endpoint{UseAgent: true, Key: path, Password: "secret"}, []cryptossh.Signer{newTestSigner(t)})
	if err != nil {
		t.Fatal(err)
	}
	keyboard := cryptossh.KeyboardInteractive(func(string, string, []string, []bool) ([]string, error) {
		return nil, nil
	})
	methods, err = finalizeAuthMethods(methods, keyboard)
	if err != nil {
		t.Fatal(err)
	}
	got := methodTypes(methods)
	if len(got) != 4 ||
		!strings.Contains(got[0], "publicKey") ||
		!strings.Contains(got[1], "publicKey") ||
		!strings.Contains(got[2], "password") ||
		!strings.Contains(got[3], "KeyboardInteractive") {
		t.Fatalf("methods = %v", got)
	}
}

func TestAuthMethodsKeyPEM(t *testing.T) {
	_, privatePEM, err := sshkey.Generate(sshkey.AlgoEd25519, "test")
	if err != nil {
		t.Fatal(err)
	}
	methods, err := authMethods(Endpoint{KeyPEM: privatePEM, Key: "ignored"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	got := methodTypes(methods)
	if len(got) != 1 || !strings.Contains(got[0], "publicKey") {
		t.Fatalf("methods = %v", got)
	}
}

func TestAgentUnavailableIsFatal(t *testing.T) {
	if !agentUnavailableIsFatal(Endpoint{UseAgent: true}, false) {
		t.Fatal("expected fatal without other auth")
	}
	if agentUnavailableIsFatal(Endpoint{Password: "x"}, false) {
		t.Fatal("password should keep dialing")
	}
	if agentUnavailableIsFatal(Endpoint{Key: "k"}, false) {
		t.Fatal("key file should keep dialing")
	}
	if agentUnavailableIsFatal(Endpoint{KeyPEM: "pem"}, false) {
		t.Fatal("stored key should keep dialing")
	}
	if agentUnavailableIsFatal(Endpoint{}, true) {
		t.Fatal("keyboard-interactive should keep dialing")
	}
}
