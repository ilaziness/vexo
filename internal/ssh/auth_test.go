package ssh

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"io"
	"strings"
	"testing"

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
	_, err := authMethods(Endpoint{UseAgent: true}, nil)
	if err == nil || !strings.Contains(err.Error(), "没有可用密钥") {
		t.Fatalf("err = %v", err)
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
	_, err := authMethods(Endpoint{}, nil)
	if err == nil || !strings.Contains(err.Error(), "empty password and key") {
		t.Fatalf("err = %v", err)
	}
}
