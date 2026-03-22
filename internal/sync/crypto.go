package sync

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

const (
	saltSize   = 16
	nonceSize  = 12
	keySize    = 32
	iterations = 100000
)

// DeriveKey 从 user_key 派生加密密钥
func DeriveKey(userKey string, salt []byte) []byte {
	return pbkdf2.Key([]byte(userKey), salt, iterations, keySize, sha256.New)
}

// clearKey 清除密钥内存
func clearKey(key []byte) {
	for i := range key {
		key[i] = 0
	}
}

// EncryptStream 流式加密
// 格式: salt(16) + nonce(12) + ciphertext + tag(16)
type EncryptStream struct {
	writer   io.Writer
	cipher   cipher.AEAD
	nonce    []byte
	buffer   []byte
}

// NewEncryptStream 创建加密流
func NewEncryptStream(writer io.Writer, userKey string) (*EncryptStream, error) {
	// 生成随机 salt
	salt := make([]byte, saltSize)
	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}

	// 派生密钥
	key := DeriveKey(userKey, salt)
	defer clearKey(key)

	// 创建 AES-GCM
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// 生成随机 nonce
	nonce := make([]byte, nonceSize)
	if _, err := rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// 写入 salt 和 nonce
	if _, err := writer.Write(salt); err != nil {
		return nil, fmt.Errorf("failed to write salt: %w", err)
	}
	if _, err := writer.Write(nonce); err != nil {
		return nil, fmt.Errorf("failed to write nonce: %w", err)
	}

	return &EncryptStream{
		writer: writer,
		cipher: aead,
		nonce:  nonce,
		buffer: make([]byte, 0, 32*1024),
	}, nil
}

// Write 实现 io.Writer 接口
func (e *EncryptStream) Write(p []byte) (n int, err error) {
	e.buffer = append(e.buffer, p...)

	// 当缓冲区足够大时，加密并写入
	if len(e.buffer) >= 32*1024 {
		if err := e.flush(); err != nil {
			return 0, err
		}
	}

	return len(p), nil
}

// Close 关闭加密流，写入剩余数据
func (e *EncryptStream) Close() error {
	if len(e.buffer) > 0 {
		return e.flush()
	}
	return nil
}

// flush 加密并写入缓冲区数据
func (e *EncryptStream) flush() error {
	if len(e.buffer) == 0 {
		return nil
	}

	// 加密数据
	ciphertext := e.cipher.Seal(nil, e.nonce, e.buffer, nil)

	// 写入加密数据
	if _, err := e.writer.Write(ciphertext); err != nil {
		return err
	}

	// 增加 nonce
	incrementNonce(e.nonce)

	// 清空缓冲区
	e.buffer = e.buffer[:0]

	return nil
}

// DecryptStream 流式解密
type DecryptStream struct {
	reader   io.Reader
	cipher   cipher.AEAD
	nonce    []byte
	buffer   []byte
	saltRead bool
	userKey  string
}

// NewDecryptStream 创建解密流
func NewDecryptStream(reader io.Reader, userKey string) (*DecryptStream, error) {
	return &DecryptStream{
		reader:  reader,
		nonce:   make([]byte, nonceSize),
		buffer:  make([]byte, 0, 32*1024+16), // 预留 tag 空间
		userKey: userKey,
	}, nil
}

// Read 实现 io.Reader 接口
func (d *DecryptStream) Read(p []byte) (n int, err error) {
	// 首次读取，读取 salt 和 nonce
	if !d.saltRead {
		salt := make([]byte, saltSize)
		if _, err := io.ReadFull(d.reader, salt); err != nil {
			return 0, fmt.Errorf("failed to read salt: %w", err)
		}

		if _, err := io.ReadFull(d.reader, d.nonce); err != nil {
			return 0, fmt.Errorf("failed to read nonce: %w", err)
		}

		// 派生密钥
		key := DeriveKey(d.userKey, salt)
		defer clearKey(key)

		// 创建 AES-GCM
		block, err := aes.NewCipher(key)
		if err != nil {
			return 0, fmt.Errorf("failed to create cipher: %w", err)
		}

		aead, err := cipher.NewGCM(block)
		if err != nil {
			return 0, fmt.Errorf("failed to create GCM: %w", err)
		}

		d.cipher = aead
		d.saltRead = true
	}

	// 读取加密数据
	buf := make([]byte, 32*1024+16) // 包含 tag
	n, err = d.reader.Read(buf)
	if n > 0 {
		d.buffer = append(d.buffer, buf[:n]...)
	}

	// 尝试解密
	if len(d.buffer) > 16 {
		// 解密数据
		plaintext, err := d.cipher.Open(nil, d.nonce, d.buffer, nil)
		if err == nil {
			// 解密成功
			n = copy(p, plaintext)
			d.buffer = d.buffer[:0]
			incrementNonce(d.nonce)
			return n, nil
		}
	}

	if err == io.EOF && len(d.buffer) > 0 {
		// 最后一块数据
		plaintext, err := d.cipher.Open(nil, d.nonce, d.buffer, nil)
		if err != nil {
			return 0, fmt.Errorf("failed to decrypt final block: %w", err)
		}
		n = copy(p, plaintext)
		return n, io.EOF
	}

	return 0, err
}

// incrementNonce 增加 nonce
func incrementNonce(nonce []byte) {
	for i := len(nonce) - 1; i >= 0; i-- {
		nonce[i]++
		if nonce[i] != 0 {
			break
		}
	}
}


