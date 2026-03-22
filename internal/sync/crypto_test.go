package sync

import (
	"bytes"
	"testing"
)

func TestDeriveKey(t *testing.T) {
	userKey := "test-user-key"
	salt := []byte("test-salt-123456")

	key1 := DeriveKey(userKey, salt)
	key2 := DeriveKey(userKey, salt)

	if len(key1) != keySize {
		t.Errorf("expected key size %d, got %d", keySize, len(key1))
	}

	// 相同的输入应该产生相同的密钥
	if !bytes.Equal(key1, key2) {
		t.Error("same input should produce same key")
	}

	// 不同的盐应该产生不同的密钥
	differentSalt := []byte("different-salt-12")
	key3 := DeriveKey(userKey, differentSalt)
	if bytes.Equal(key1, key3) {
		t.Error("different salt should produce different key")
	}
}

func TestEncryptDecryptStream(t *testing.T) {
	userKey := "test-encryption-key"
	originalData := []byte("Hello, World! This is a test message for encryption and decryption.")

	// 加密
	var encryptedBuffer bytes.Buffer
	encryptStream, err := NewEncryptStream(&encryptedBuffer, userKey)
	if err != nil {
		t.Fatalf("failed to create encrypt stream: %v", err)
	}

	_, err = encryptStream.Write(originalData)
	if err != nil {
		t.Fatalf("failed to write data: %v", err)
	}

	err = encryptStream.Close()
	if err != nil {
		t.Fatalf("failed to close encrypt stream: %v", err)
	}

	encryptedData := encryptedBuffer.Bytes()
	if len(encryptedData) == 0 {
		t.Fatal("encrypted data should not be empty")
	}

	// 解密
	decryptStream, err := NewDecryptStream(bytes.NewReader(encryptedData), userKey)
	if err != nil {
		t.Fatalf("failed to create decrypt stream: %v", err)
	}

	var decryptedBuffer bytes.Buffer
	_, err = decryptedBuffer.ReadFrom(decryptStream)
	if err != nil {
		t.Fatalf("failed to read decrypted data: %v", err)
	}

	decryptedData := decryptedBuffer.Bytes()
	if !bytes.Equal(originalData, decryptedData) {
		t.Errorf("decrypted data mismatch: expected %s, got %s", originalData, decryptedData)
	}
}

func TestClearKey(t *testing.T) {
	key := []byte{0x01, 0x02, 0x03, 0x04, 0x05}
	originalKey := make([]byte, len(key))
	copy(originalKey, key)

	clearKey(key)

	// 验证密钥已被清零
	for i, b := range key {
		if b != 0 {
			t.Errorf("key byte at index %d should be 0, got %d", i, b)
		}
	}
}

func TestIncrementNonce(t *testing.T) {
	nonce := []byte{0x00, 0x00, 0x00, 0x00}
	incrementNonce(nonce)

	expected := []byte{0x00, 0x00, 0x00, 0x01}
	if !bytes.Equal(nonce, expected) {
		t.Errorf("nonce increment failed: expected %v, got %v", expected, nonce)
	}

	// 测试溢出
	nonce = []byte{0xFF, 0xFF, 0xFF, 0xFF}
	incrementNonce(nonce)
	expected = []byte{0x00, 0x00, 0x00, 0x00}
	if !bytes.Equal(nonce, expected) {
		t.Errorf("nonce overflow failed: expected %v, got %v", expected, nonce)
	}
}
