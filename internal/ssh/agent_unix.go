//go:build unix

package ssh

import (
	"fmt"
	"net"
	"os"
)

func dialAgentEndpoints() ([]agentEndpoint, error) {
	sock := os.Getenv("SSH_AUTH_SOCK")
	if sock == "" {
		return nil, fmt.Errorf("SSH_AUTH_SOCK is not set")
	}
	conn, err := net.Dial("unix", sock)
	if err != nil {
		return nil, fmt.Errorf("dial SSH_AUTH_SOCK: %w", err)
	}
	return []agentEndpoint{{rw: conn, closer: conn}}, nil
}
