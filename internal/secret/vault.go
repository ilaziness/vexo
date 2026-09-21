package secret

import (
	"errors"
	"sync"
)

var ErrPasswordNotSet = errors.New("password not entered")

// Vault 进程内保存用户口令，不含任何 UI。
type Vault struct {
	mu       sync.RWMutex
	password string
}

func NewVault() *Vault {
	return &Vault{}
}

func (v *Vault) Set(password string) {
	v.mu.Lock()
	v.password = password
	v.mu.Unlock()
}

func (v *Vault) Get() (string, error) {
	v.mu.RLock()
	defer v.mu.RUnlock()
	if v.password == "" {
		return "", ErrPasswordNotSet
	}
	return v.password, nil
}

func (v *Vault) Clear() {
	v.Set("")
}
