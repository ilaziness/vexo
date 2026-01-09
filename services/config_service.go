package services

import (
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

var (
	Version       = "0.0.1"
	versionNumber = 1
	Mode          = "debug"
	ModeDebug     = "debug"
	ModeRelease   = "release"
)

type Config struct {
	General  GeneralConfig  `toml:"general"`
	Terminal TerminalConfig `toml:"terminal"`
}

type GeneralConfig struct {
	UserDataDir string `toml:"user_data_dir"`
}

type TerminalConfig struct {
	Font       string  `toml:"font"`
	FontSize   int     `toml:"font_size"`
	LineHeight float64 `toml:"line_height"`
}

type ConfigService struct {
	Config            *Config
	configPath        string // UserDataDir 下的配置文件路径
	defaultConfigPath string // 可执行目录下的默认配置文件路径
	window            *application.WebviewWindow
}

// GetDefaultConfig 获取默认配置
func GetDefaultConfig() *Config {
	// 获取可执行文件路径
	execPath, err := os.Executable()
	if err != nil {
		// 如果获取失败，使用当前工作目录
		execPath, _ = os.Getwd()
	}

	// 获取可执行文件所在目录
	execDir := filepath.Dir(execPath)
	defaultUserDataDir := filepath.Join(execDir, "data")

	return &Config{
		General: GeneralConfig{
			UserDataDir: defaultUserDataDir,
		},
		Terminal: TerminalConfig{
			Font:       "monospace",
			FontSize:   14,
			LineHeight: 1.2,
		},
	}
}

// mergeConfig 合并配置，userConfig 覆盖 defaultConfig
func mergeConfig(defaultConfig, userConfig *Config) *Config {
	result := &Config{
		General:  defaultConfig.General,
		Terminal: defaultConfig.Terminal,
	}

	// 合并 General 配置
	if userConfig.General.UserDataDir != "" {
		result.General.UserDataDir = userConfig.General.UserDataDir
	}

	// 合并 Terminal 配置
	if userConfig.Terminal.Font != "" {
		result.Terminal.Font = userConfig.Terminal.Font
	}
	if userConfig.Terminal.FontSize > 0 {
		result.Terminal.FontSize = userConfig.Terminal.FontSize
	}
	if userConfig.Terminal.LineHeight > 0 {
		result.Terminal.LineHeight = userConfig.Terminal.LineHeight
	}

	return result
}

func NewConfigService() *ConfigService {
	// 获取可执行文件路径
	execPath, err := os.Executable()
	if err != nil {
		// 如果获取失败，使用当前工作目录
		execPath, _ = os.Getwd()
	}

	// 获取可执行文件所在目录
	execDir := filepath.Dir(execPath)

	// 默认配置文件路径：可执行目录/config/config.toml
	defaultConfigPath := filepath.Join(execDir, "config", "config.toml")

	// 1. 读取或创建默认配置
	defaultConfig := GetDefaultConfig()

	// 确保默认配置目录存在
	defaultConfigDir := filepath.Dir(defaultConfigPath)
	if _, err = os.Stat(defaultConfigDir); os.IsNotExist(err) {
		os.MkdirAll(defaultConfigDir, 0755)
	}

	// 如果默认配置文件不存在，创建并写入
	if _, err = os.Stat(defaultConfigPath); os.IsNotExist(err) {
		if data, err := toml.Marshal(defaultConfig); err == nil {
			os.WriteFile(defaultConfigPath, data, 0644)
		}
	} else {
		// 如果存在，读取默认配置文件
		if data, err := os.ReadFile(defaultConfigPath); err == nil {
			toml.Unmarshal(data, defaultConfig)
		}
	}

	// 确保默认 UserDataDir 存在
	if _, err := os.Stat(defaultConfig.General.UserDataDir); os.IsNotExist(err) {
		os.MkdirAll(defaultConfig.General.UserDataDir, 0755)
	}

	// 2. 检查 UserDataDir 下是否有配置文件
	userConfigPath := filepath.Join(defaultConfig.General.UserDataDir, "config.toml")
	finalConfig := defaultConfig

	if _, err := os.Stat(userConfigPath); err == nil {
		// 读取用户配置并合并
		userConfig := &Config{}
		if data, err := os.ReadFile(userConfigPath); err == nil {
			if err := toml.Unmarshal(data, userConfig); err == nil {
				finalConfig = mergeConfig(defaultConfig, userConfig)

				// 如果合并后的 UserDataDir 发生变化，更新路径
				if finalConfig.General.UserDataDir != defaultConfig.General.UserDataDir {
					// 确保新的 UserDataDir 存在
					if _, err := os.Stat(finalConfig.General.UserDataDir); os.IsNotExist(err) {
						os.MkdirAll(finalConfig.General.UserDataDir, 0755)
					}
					userConfigPath = filepath.Join(finalConfig.General.UserDataDir, "config.toml")
				}
			}
		}
	}

	return &ConfigService{
		Config:            finalConfig,
		configPath:        userConfigPath,
		defaultConfigPath: defaultConfigPath,
	}
}

func (cs *ConfigService) ShowWindow() {
	if cs.window != nil {
		return
	}
	cs.window = app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:            "设置",
		URL:              "/setting",
		Width:            1200,
		Height:           700,
		MinWidth:         1200,
		MinHeight:        700,
		BackgroundColour: application.NewRGB(27, 38, 54),
	})
	cs.window.OnWindowEvent(events.Common.WindowClosing, func(_ *application.WindowEvent) {
		cs.window = nil
	})
	cs.window.Show()
}

func (cs *ConfigService) CloseWindow() {
	if cs.window != nil {
		cs.window.Close()
		cs.window = nil
	}
}

// ReadConfig 读取配置方法，不需要参数
func (cs *ConfigService) ReadConfig() (*Config, error) {
	data, err := os.ReadFile(cs.configPath)
	if err != nil {
		// 如果配置文件不存在，返回默认配置
		if os.IsNotExist(err) {
			// 确保配置文件所在目录存在
			configDir := filepath.Dir(cs.configPath)
			if _, err := os.Stat(configDir); os.IsNotExist(err) {
				os.MkdirAll(configDir, 0755)
			}
			defaultConfig := &Config{
				General: GeneralConfig{
					UserDataDir: filepath.Dir(cs.configPath), // 使用当前目录作为用户数据目录
				},
				Terminal: TerminalConfig{
					Font:       "monospace",
					FontSize:   14,
					LineHeight: 1.2,
				},
			}
			return defaultConfig, nil
		}
		return nil, err
	}

	config := &Config{}
	err = toml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	// 确保用户数据目录存在
	if config.General.UserDataDir != "" {
		if _, err := os.Stat(config.General.UserDataDir); os.IsNotExist(err) {
			os.MkdirAll(config.General.UserDataDir, 0755)
		}
	}

	return config, nil
}

// SaveConfig 保存配置方法，需要传入配置对象，保存到 UserDataDir 指定的目录
func (cs *ConfigService) SaveConfig(config Config) error {
	data, err := toml.Marshal(config)
	if err != nil {
		return err
	}

	// 更新当前配置
	cs.Config = &config

	// 如果 UserDataDir 发生变化，更新配置文件路径
	if config.General.UserDataDir != "" && config.General.UserDataDir != filepath.Dir(cs.configPath) {
		cs.configPath = filepath.Join(config.General.UserDataDir, "config.toml")
	}

	// 确保配置文件所在目录存在
	configDir := filepath.Dir(cs.configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	// 保存到 UserDataDir 下的 config.toml
	return os.WriteFile(cs.configPath, data, 0644)
}
