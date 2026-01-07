package services

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigService(t *testing.T) {
	// 创建临时目录用于测试
	tempDir := t.TempDir()
	
	// 临时修改配置文件路径
	originalConfigPath := filepath.Join(tempDir, "config.toml")
	
	// 创建一个 ConfigService 实例用于测试
	configService := &ConfigService{
		configPath: originalConfigPath,
	}
	
	// 测试默认配置
	config, err := configService.ReadConfig()
	if err != nil {
		t.Fatalf("读取配置失败: %v", err)
	}
	
	// 验证默认配置值
	if config.Terminal.Font != "monospace" {
		t.Errorf("期望字体为 monospace，实际为 %s", config.Terminal.Font)
	}
	
	if config.Terminal.FontSize != 14 {
		t.Errorf("期望字体大小为 14，实际为 %d", config.Terminal.FontSize)
	}
	
	if config.Terminal.LineHeight != 1.2 {
		t.Errorf("期望行高为 1.2，实际为 %f", config.Terminal.LineHeight)
	}
	
	// 修改配置
	config.Terminal.Font = "Consolas"
	config.Terminal.FontSize = 16
	config.Terminal.LineHeight = 1.5
	config.General.UserDataDir = tempDir
	
	// 保存配置
	err = configService.SaveConfig(*config)
	if err != nil {
		t.Fatalf("保存配置失败: %v", err)
	}
	
	// 重新读取配置
	newConfig, err := configService.ReadConfig()
	if err != nil {
		t.Fatalf("重新读取配置失败: %v", err)
	}
	
	// 验证保存的配置值
	if newConfig.Terminal.Font != "Consolas" {
		t.Errorf("期望字体为 Consolas，实际为 %s", newConfig.Terminal.Font)
	}
	
	if newConfig.Terminal.FontSize != 16 {
		t.Errorf("期望字体大小为 16，实际为 %d", newConfig.Terminal.FontSize)
	}
	
	if newConfig.Terminal.LineHeight != 1.5 {
		t.Errorf("期望行高为 1.5，实际为 %f", newConfig.Terminal.LineHeight)
	}
}

func TestNewConfigService(t *testing.T) {
	// 测试 NewConfigService 函数
	configService := NewConfigService()
	
	// 验证 configService 不为 nil
	if configService == nil {
		t.Fatal("ConfigService 应该不为 nil")
	}
	
	// 验证 configPath 不为空
	if configService.configPath == "" {
		t.Error("configPath 应该不为空")
	}
	
	// 验证配置文件所在目录存在
	configDir := filepath.Dir(configService.configPath)
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		t.Errorf("配置目录应该存在: %s", configDir)
	}
}