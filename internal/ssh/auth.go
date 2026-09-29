package ssh

import (
	"fmt"
	"os"

	cryptossh "golang.org/x/crypto/ssh"
)

// authMethods 按 Agent、密钥文件、密码的顺序组装认证方法。
func authMethods(ep Endpoint, agentSigners []cryptossh.Signer) ([]cryptossh.AuthMethod, error) {
	var methods []cryptossh.AuthMethod
	if ep.UseAgent && len(agentSigners) > 0 {
		methods = append(methods, cryptossh.PublicKeys(agentSigners...))
	}
	if ep.Key != "" {
		signer, err := signerFromKeyFile(ep.Key, ep.KeyPassword)
		if err != nil {
			return nil, err
		}
		methods = append(methods, cryptossh.PublicKeys(signer))
	}
	if ep.Password != "" {
		methods = append(methods, cryptossh.Password(ep.Password))
	}
	if len(methods) == 0 {
		if ep.UseAgent {
			return nil, fmt.Errorf("SSH agent 中没有可用密钥")
		}
		return nil, fmt.Errorf("empty password and key")
	}
	return methods, nil
}

func signerFromKeyFile(path, passphrase string) (cryptossh.Signer, error) {
	keyContent, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("unable to read private key: %v", err)
	}
	if passphrase == "" {
		return cryptossh.ParsePrivateKey(keyContent)
	}
	return cryptossh.ParsePrivateKeyWithPassphrase(keyContent, []byte(passphrase))
}
