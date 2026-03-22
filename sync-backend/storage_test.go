package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestNewStorage(t *testing.T) {
	tmpDir := t.TempDir()
	storage, err := NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	if storage == nil {
		t.Fatal("storage should not be nil")
	}
	if storage.baseDir != tmpDir {
		t.Errorf("expected baseDir %s, got %s", tmpDir, storage.baseDir)
	}
}

func TestStorageSaveAndOpen(t *testing.T) {
	tmpDir := t.TempDir()
	storage, err := NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}

	userID := "test-user"
	version := 1
	data := []byte("test data content")

	// 测试保存
	filePath, size, err := storage.Save(userID, version, bytes.NewReader(data))
	if err != nil {
		t.Fatalf("failed to save file: %v", err)
	}
	if size != int64(len(data)) {
		t.Errorf("expected size %d, got %d", len(data), size)
	}

	// 验证文件存在
	expectedPath := filepath.Join(tmpDir, userID, "v1.dat")
	if filePath != expectedPath {
		t.Errorf("expected path %s, got %s", expectedPath, filePath)
	}

	// 测试打开
	reader, err := storage.Open(userID, version)
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}
	defer reader.Close()

	content, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	if !bytes.Equal(content, data) {
		t.Errorf("content mismatch: expected %s, got %s", data, content)
	}
}

func TestStorageDelete(t *testing.T) {
	tmpDir := t.TempDir()
	storage, err := NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}

	userID := "test-user"
	version := 1
	data := []byte("test data")

	// 保存文件
	_, _, err = storage.Save(userID, version, bytes.NewReader(data))
	if err != nil {
		t.Fatalf("failed to save file: %v", err)
	}

	// 删除文件
	err = storage.Delete(userID, version)
	if err != nil {
		t.Fatalf("failed to delete file: %v", err)
	}

	// 验证文件不存在
	_, err = storage.Open(userID, version)
	if err == nil {
		t.Error("should not be able to open deleted file")
	}
}

func TestStorageDeleteUserDir(t *testing.T) {
	tmpDir := t.TempDir()
	storage, err := NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}

	userID := "test-user"
	data := []byte("test data")

	// 保存多个版本
	for i := 1; i <= 3; i++ {
		_, _, err = storage.Save(userID, i, bytes.NewReader(data))
		if err != nil {
			t.Fatalf("failed to save file: %v", err)
		}
	}

	// 删除用户目录
	err = storage.DeleteUserDir(userID)
	if err != nil {
		t.Fatalf("failed to delete user dir: %v", err)
	}

	// 验证目录不存在
	userDir := filepath.Join(tmpDir, userID)
	if _, err := os.Stat(userDir); !os.IsNotExist(err) {
		t.Error("user directory should be deleted")
	}
}
