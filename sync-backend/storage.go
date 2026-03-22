package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Storage 文件存储
type Storage struct {
	baseDir string
}

// NewStorage 创建文件存储
func NewStorage(baseDir string) (*Storage, error) {
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}
	return &Storage{baseDir: baseDir}, nil
}

// getUserDir 获取用户目录路径
func (s *Storage) getUserDir(userID string) string {
	return filepath.Join(s.baseDir, userID)
}

// getFilePath 获取文件路径
func (s *Storage) getFilePath(userID string, version int) string {
	return filepath.Join(s.getUserDir(userID), fmt.Sprintf("v%d.dat", version))
}

// Save 保存文件
func (s *Storage) Save(userID string, version int, reader io.Reader) (string, int64, error) {
	userDir := s.getUserDir(userID)
	if err := os.MkdirAll(userDir, 0755); err != nil {
		return "", 0, fmt.Errorf("failed to create user directory: %w", err)
	}

	filePath := s.getFilePath(userID, version)
	file, err := os.Create(filePath)
	if err != nil {
		return "", 0, fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	size, err := io.Copy(file, reader)
	if err != nil {
		return "", 0, fmt.Errorf("failed to write file: %w", err)
	}

	return filePath, size, nil
}

// Open 打开文件
func (s *Storage) Open(userID string, version int) (io.ReadCloser, error) {
	filePath := s.getFilePath(userID, version)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	return file, nil
}

// Delete 删除文件
func (s *Storage) Delete(userID string, version int) error {
	filePath := s.getFilePath(userID, version)
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

// DeleteUserDir 删除用户目录
func (s *Storage) DeleteUserDir(userID string) error {
	userDir := s.getUserDir(userID)
	if err := os.RemoveAll(userDir); err != nil {
		return fmt.Errorf("failed to delete user directory: %w", err)
	}
	return nil
}
