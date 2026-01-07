package services

import (
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
	"github.com/wailsapp/wails/v3/pkg/application"
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
	Config     *Config
	configPath string
	window     *application.WebviewWindow
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
	defaultUserDataDir := filepath.Join(execDir, "data")

	// 检查并创建默认用户数据目录
	if _, err := os.Stat(defaultUserDataDir); os.IsNotExist(err) {
		os.MkdirAll(defaultUserDataDir, 0755)
	}

	// 从默认目录或配置中获取用户数据目录
	userDataDir := defaultUserDataDir
	configPath := filepath.Join(userDataDir, "config.toml")

	// 如果配置文件存在，尝试读取它以获取实际的用户数据目录
	config := &Config{}
	if _, err := os.Stat(configPath); err == nil {
		if data, err := os.ReadFile(configPath); err == nil {
			if err := toml.Unmarshal(data, config); err == nil {
				if config.General.UserDataDir != "" {
					userDataDir = config.General.UserDataDir
					// 确保用户数据目录存在
					if _, err := os.Stat(userDataDir); os.IsNotExist(err) {
						os.MkdirAll(userDataDir, 0755)
					}
					configPath = filepath.Join(userDataDir, "config.toml")
				}
			}
		}
	} else {
		// 如果配置文件不存在，使用默认配置
		config = &Config{
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

	return &ConfigService{
		Config:     config,
		configPath: configPath,
	}
}

func (cs *ConfigService) ShowWindow() {
	cs.window = app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:            "设置",
		URL:              "/setting",
		MinWidth:         1200,
		MinHeight:        700,
		BackgroundColour: application.NewRGB(27, 38, 54),
	})
	cs.window.Show()
}

func (cs *ConfigService) CloseWindow() {
	if cs.window != nil {
		cs.window.Close()
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

// SaveConfig 保存配置方法，需要传入配置对象，保存到文件持久化
func (cs *ConfigService) SaveConfig(config Config) error {
	data, err := toml.Marshal(config)
	if err != nil {
		return err
	}

	// 确保配置文件所在目录存在
	configDir := filepath.Dir(cs.configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	return os.WriteFile(cs.configPath, data, 0644)
}
