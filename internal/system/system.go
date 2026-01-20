package system

import (
	"os"
	"path/filepath"
)

// GetExecutableDir 获取当前可执行文件所在目录
// 如果获取失败，使用当前工作目录作为fallback
func GetExecutableDir() string {
	execPath, err := os.Executable()
	if err != nil {
		// 如果获取失败，使用当前工作目录
		execPath, _ = os.Getwd()
	}

	return filepath.Dir(execPath)
}
