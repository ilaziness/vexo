package system

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
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

// RecoverFromPanic 捕获 panic 并记录错误信息
func RecoverFromPanic() {
	if r := recover(); r != nil {
		// 记录详细的错误信息和堆栈跟踪
		fmt.Printf("Recovered from panic: %v\n", r)
		// 打印调用堆栈以便调试
		buf := make([]byte, 4096)
		n := runtime.Stack(buf, false)
		fmt.Printf("Stack trace:\n%s\n", buf[:n])
	}
}
