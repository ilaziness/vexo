package services

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/ilaziness/vexo/internal/system"
	"github.com/pelletier/go-toml/v2"
	"github.com/wailsapp/wails/v3/pkg/application"
	"go.uber.org/zap"
)

var (
	Version       = "0.0.1"
	versionNumber = 1
	Mode          = "debug"
	ModeDebug     = "debug"
	ModeRelease   = "release"

	defaultFontFamily = []string{
		"'Noto Sans Mono'",
		"'ui-monospace'",
		"'Cascadia Code'",
		"Consolas",
		"Menlo",
		"Monaco",
		"'DejaVu Sans Mono'",
		"'Ubuntu Mono'",
		"'Liberation Mono'",
		"'Courier New'",
		"'Microsoft YaHei'",
		"'PingFang SC'",
		"'Heiti SC'",
		"'WenQuanYi Micro Hei'",
		"monospace",
	}
)

type Config struct {
	General  GeneralConfig  `toml:"general"`
	Terminal TerminalConfig `toml:"terminal"`
}

type GeneralConfig struct {
	UserDataDir string `toml:"user_data_dir"`
	Theme       string `toml:"theme"`
}

type TerminalConfig struct {
	Font       string  `toml:"font" json:"fontFamily"`
	FontSize   int     `toml:"font_size" json:"fontSize"`
	LineHeight float64 `toml:"line_height" json:"lineHeight"`
}

// AppConfig 应用配置，只保存 user_data_dir
type AppConfig struct {
	General GeneralConfig `toml:"general"`
}

type ConfigService struct {
	Config     *Config
	userConfig string // UserDataDir 下的用户配置文件路径
	appConfig  string // 可执行目录下的应用配置文件路径
	window     *application.WebviewWindow
}

// GetDefaultConfig 获取默认配置
func GetDefaultConfig() *Config {
	// 获取可执行文件所在目录
	execDir := system.GetExecutableDir()
	defaultUserDataDir := filepath.Join(execDir, "data")

	return &Config{
		General: GeneralConfig{
			UserDataDir: defaultUserDataDir,
			Theme:       "dark",
		},
		Terminal: TerminalConfig{
			Font:       strings.Join(defaultFontFamily, ","),
			FontSize:   14,
			LineHeight: 1,
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
	// 获取可执行文件所在目录
	execDir := system.GetExecutableDir()

	// 应用配置文件路径：可执行目录/config.toml
	appConfigPath := filepath.Join(execDir, "config.toml")

	// 1. 获取默认配置
	defaultConfig := GetDefaultConfig()

	// 2. 读取应用配置（如果存在）
	appConfig := &AppConfig{
		General: defaultConfig.General,
	}
	if data, err := os.ReadFile(appConfigPath); err == nil {
		toml.Unmarshal(data, appConfig)
	} else if os.IsNotExist(err) {
		// 首次启动，创建应用配置文件，只保存 user_data_dir
		if data, err := toml.Marshal(appConfig); err == nil {
			os.WriteFile(appConfigPath, data, 0644)
		}
	}

	// 确保 UserDataDir 存在
	userDataDir := appConfig.General.UserDataDir
	if userDataDir == "" {
		userDataDir = defaultConfig.General.UserDataDir
	}
	if _, err := os.Stat(userDataDir); os.IsNotExist(err) {
		os.MkdirAll(userDataDir, 0755)
	}

	// 3. 读取用户配置（如果存在）
	userConfigPath := filepath.Join(userDataDir, "config.toml")
	finalConfig := defaultConfig
	finalConfig.General.UserDataDir = userDataDir

	if data, err := os.ReadFile(userConfigPath); err == nil {
		userConfig := &Config{}
		if err = toml.Unmarshal(data, userConfig); err == nil {
			// 用户配置覆盖默认配置（不包括 UserDataDir，UserDataDir 由应用配置管理）
			finalConfig = mergeConfig(defaultConfig, userConfig)
			finalConfig.General.UserDataDir = userDataDir
		}
	} else if os.IsNotExist(err) {
		// 首次启动，创建用户配置文件
		if data, err := toml.Marshal(finalConfig); err == nil {
			os.WriteFile(userConfigPath, data, 0644)
		}
	}

	return &ConfigService{
		Config:     finalConfig,
		userConfig: userConfigPath,
		appConfig:  appConfigPath,
	}
}

func (cs *ConfigService) ShowWindow() {
	AppInstance.SettingWindow.Show()
	AppInstance.SettingWindow.Focus()
}

func (cs *ConfigService) CloseWindow() {
	AppInstance.SettingWindow.Hide()
}

// SetTheme 设置主题
func (cs *ConfigService) SetTheme(theme string) {
	cs.Config.General.Theme = theme
	cs.saveToFile()
}

// ReadConfig 读取配置方法，返回当前配置
func (cs *ConfigService) ReadConfig() (*Config, error) {
	return cs.Config, nil
}

// SaveConfig 保存配置方法
// 1. 如果 UserDataDir 改变，更新应用配置文件
// 2. 将所有配置保存到用户配置文件
func (cs *ConfigService) SaveConfig(config Config) error {
	// 1. 检查 UserDataDir 是否发生变化
	if config.General.UserDataDir != cs.Config.General.UserDataDir {
		// 更新应用配置文件中的 user_data_dir
		appConfig := &AppConfig{
			General: GeneralConfig{
				UserDataDir: config.General.UserDataDir,
			},
		}
		if data, err := toml.Marshal(appConfig); err == nil {
			os.WriteFile(cs.appConfig, data, 0644)
		}

		// 确保新的 UserDataDir 存在
		if _, err := os.Stat(config.General.UserDataDir); os.IsNotExist(err) {
			os.MkdirAll(config.General.UserDataDir, 0755)
		}

		// 更新用户配置文件路径
		cs.userConfig = filepath.Join(config.General.UserDataDir, "config.toml")
	}

	// 2. 保存完整配置到用户配置文件
	data, err := toml.Marshal(config)
	if err != nil {
		return err
	}

	// 确保用户配置目录存在
	userConfigDir := filepath.Dir(cs.userConfig)
	if err := os.MkdirAll(userConfigDir, 0755); err != nil {
		return err
	}

	// 保存到用户配置文件
	if err := os.WriteFile(cs.userConfig, data, 0644); err != nil {
		return err
	}

	// 3. 更新内存中的配置
	cs.Config = &config

	return nil
}

// saveToFile 保存配置到文件
func (cs *ConfigService) saveToFile() {
	data, err := toml.Marshal(cs.Config)
	if err != nil {
		Logger.Error("save config to file failed", zap.Error(err))
		return
	}
	if err := os.WriteFile(cs.userConfig, data, 0644); err != nil {
		Logger.Error("save config to file failed", zap.Error(err))
		return
	}
}
