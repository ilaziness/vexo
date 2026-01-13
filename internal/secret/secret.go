package secret

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/argon2"
)

const (
	// Argon2 parameters
	time    = 3
	memory  = 64 * 1024 // 64 KB
	threads = 4
	keyLen  = 32
	saltLen = 16
	nonceLen = 12
)

// Encrypt encrypts plaintext using Argon2 + AES-GCM with the given password
func Encrypt(password, plaintext string) (string, error) {
	// Generate random salt
	salt := make([]byte, saltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %v", err)
	}

	// Derive key using Argon2
	key := argon2.IDKey([]byte(password), salt, time, uint32(memory), uint8(threads), keyLen)

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create AES cipher: %v", err)
	}

	// Generate nonce
	nonce := make([]byte, nonceLen)
	if _, err := rand.Read(nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %v", err)
	}

	// Create GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %v", err)
	}

	// Encrypt
	ciphertext := gcm.Seal(nil, nonce, []byte(plaintext), nil)

	// Combine salt + nonce + ciphertext
	combined := make([]byte, 0, saltLen+nonceLen+len(ciphertext))
	combined = append(combined, salt...)
	combined = append(combined, nonce...)
	combined = append(combined, ciphertext...)

	// Base64 encode
	return base64.StdEncoding.EncodeToString(combined), nil
}

// Decrypt decrypts encrypted text using Argon2 + AES-GCM with the given password
func Decrypt(password, encrypted string) (string, error) {
	// Base64 decode
	combined, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %v", err)
	}

	if len(combined) < saltLen+nonceLen {
		return "", fmt.Errorf("encrypted data too short")
	}

	// Extract salt, nonce, ciphertext
	salt := combined[:saltLen]
	nonce := combined[saltLen : saltLen+nonceLen]
	ciphertext := combined[saltLen+nonceLen:]

	// Derive key using Argon2
	key := argon2.IDKey([]byte(password), salt, time, uint32(memory), uint8(threads), keyLen)

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create AES cipher: %v", err)
	}

	// Create GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %v", err)
	}

	// Decrypt
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: %v", err)
	}

	return string(plaintext), nil
}
