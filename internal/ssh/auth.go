package ssh

import (
	"fmt"
	"os"

	cryptossh "golang.org/x/crypto/ssh"
)

// authMethods 按 Agent、密钥文件、密码的顺序组装静态认证方法。
// 没有任何静态凭证时返回空切片，由拨号方决定是否再追加 keyboard-interactive。
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
	return methods, nil
}

// finalizeAuthMethods 在静态方法之后追加 keyboard-interactive。都没有时返回错误。
func finalizeAuthMethods(methods []cryptossh.AuthMethod, keyboard cryptossh.AuthMethod) ([]cryptossh.AuthMethod, error) {
	if keyboard != nil {
		methods = append(methods, keyboard)
	}
	if len(methods) == 0 {
		return nil, fmt.Errorf("没有可用认证")
	}
	return methods, nil
}

func agentUnavailableIsFatal(ep Endpoint, keyboardEnabled bool) bool {
	return ep.Password == "" && ep.Key == "" && !keyboardEnabled
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
