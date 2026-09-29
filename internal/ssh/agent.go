package ssh

import (
	"io"

	cryptossh "golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

type agentEndpoint struct {
	rw     io.ReadWriter
	closer io.Closer
}

// serialReadWriter 不实现 io.Closer，让 agent.NewClient 按一次一问一答使用连接。
// Windows 命名管道是同步句柄，客户端的流水线会在两个 goroutine 里同时读写，Close 也停不掉阻塞中的 ReadFile。
type serialReadWriter struct {
	rw io.ReadWriter
}

func (s serialReadWriter) Read(p []byte) (int, error) {
	return s.rw.Read(p)
}

func (s serialReadWriter) Write(p []byte) (int, error) {
	return s.rw.Write(p)
}

// openAgentSigners 连接本机 SSH agent 并返回其签名者。
// close 在握手结束后调用；没有可关闭的连接时为 nil。
func openAgentSigners() ([]cryptossh.Signer, func(), error) {
	eps, err := dialAgentEndpoints()
	cleanup := func() {
		for _, ep := range eps {
			if ep.closer != nil {
				_ = ep.closer.Close()
			}
		}
	}
	if len(eps) == 0 {
		return nil, nil, err
	}

	seen := make(map[string]struct{})
	var signers []cryptossh.Signer
	var signErr error
	for i, ep := range eps {
		list, serr := agent.NewClient(serialReadWriter{rw: ep.rw}).Signers()
		if serr != nil {
			signErr = serr
			if ep.closer != nil {
				_ = ep.closer.Close()
				eps[i].closer = nil
			}
			continue
		}
		for _, signer := range list {
			blob := string(signer.PublicKey().Marshal())
			if _, ok := seen[blob]; ok {
				continue
			}
			seen[blob] = struct{}{}
			signers = append(signers, signer)
		}
	}
	if len(signers) == 0 && signErr != nil {
		cleanup()
		return nil, nil, signErr
	}
	return signers, cleanup, nil
}
