package services

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	internalsync "github.com/ilaziness/vexo/internal/sync"
	"github.com/ilaziness/vexo/internal/system"
	"github.com/pelletier/go-toml/v2"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
	"go.uber.org/zap"
)

const (
	EventInputPassword      = "eventInputPassword"
	EventInputPasswordClose = "eventInputPasswordClose"
)

// 用户密码，用于加密/解密私钥密码和API Key
var (
	globalUserPassword string
	passwordMutex      sync.RWMutex
)

func init() {
	application.RegisterEvent[string](EventInputPassword)
	application.RegisterEvent[string](EventInputPasswordClose)
}

var (
	Mode        = "debug"
	ModeDebug   = "debug"
	ModeRelease = "release"
	ConfigFile  = "config.toml"

	Version   = "v1.0.0"
	GitInfo   = "2026-02-06 14:53:25 4cc46033541fe7927208fedf8edea6bf3efafce7"
	BuildTime = "2026-02-06 07:02:20 AM"

	// CurrentDBVersion 当前数据库版本号
	CurrentDBVersion = 2

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
	General  GeneralConfig           `toml:"general"`
	Terminal TerminalConfig          `toml:"terminal"`
	Sync     internalsync.SyncConfig `toml:"sync"`
	AI       AIConfig                `toml:"ai"`
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

// AppConfig 应用配置
type AppConfig struct {
	General   GeneralConfig `toml:"general"`
	DBVersion int           `toml:"db_version,omitempty"`
}

type ConfigService struct {
	Config         *Config
	appConfig      *AppConfig // 应用配置对象
	userConfigFile string     // UserDataDir 下的用户配置文件路径
	appConfigFile  string     // 可执行目录下的应用配置文件路径
	passwordChan   chan struct{}
	chanMutex      sync.Mutex // 保护 passwordChan 的并发访问
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

func (cs *ConfigService) ShowWindow() {
	if AppInstance.SettingWindow == nil {
		AppInstance.SettingWindow = newSettingWindow()
	}
	AppInstance.SettingWindow.OnWindowEvent(events.Common.WindowClosing, func(event *application.WindowEvent) {
		AppInstance.SettingWindow = nil
	})
	AppInstance.SettingWindow.Show()
	AppInstance.SettingWindow.Focus()
}

func (cs *ConfigService) CloseWindow() {
	if AppInstance.SettingWindow != nil {
		AppInstance.SettingWindow.Close()
		AppInstance.SettingWindow = nil
	}
}

func NewConfigService() *ConfigService {
	// 获取可执行文件所在目录
	execDir := system.GetExecutableDir()

	// 应用配置文件路径：可执行目录/config.toml
	appConfigPath := filepath.Join(execDir, ConfigFile)

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
	userConfigPath := filepath.Join(userDataDir, ConfigFile)
	finalConfig := defaultConfig
	finalConfig.General.UserDataDir = userDataDir

	if data, err := os.ReadFile(userConfigPath); err == nil {
		Logger.Debug("user config content", zap.String("content", string(data)))
		if err = toml.Unmarshal(data, finalConfig); err != nil {
			Logger.Error("unmarshal user config failed", zap.Error(err))
		}
	} else if os.IsNotExist(err) {
		// 首次启动，创建用户配置文件
		if data, err := toml.Marshal(finalConfig); err == nil {
			os.WriteFile(userConfigPath, data, 0644)
		}
	}
	Logger.Debug("final config", zap.Any("config", finalConfig), zap.String("appConfigPath", appConfigPath), zap.String("userConfigPath", userConfigPath))

	return &ConfigService{
		Config:         finalConfig,
		appConfig:      appConfig,
		userConfigFile: userConfigPath,
		appConfigFile:  appConfigPath,
	}
}

// SetTheme 设置主题
func (cs *ConfigService) SetTheme(theme string) {
	cs.Config.General.Theme = theme
	cs.saveToFile()
}

// ReadConfig 读取配置方法，返回当前配置
func (cs *ConfigService) ReadConfig() (*Config, error) {
	Logger.Debug("read config", zap.Any("config", cs.Config))
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
			os.WriteFile(cs.appConfigFile, data, 0644)
		}

		// 确保新的 UserDataDir 存在
		if _, err := os.Stat(config.General.UserDataDir); os.IsNotExist(err) {
			os.MkdirAll(config.General.UserDataDir, 0755)
		}

		// 更新用户配置文件路径
		cs.userConfigFile = filepath.Join(config.General.UserDataDir, ConfigFile)
	}

	// 2. 保存完整配置到用户配置文件
	data, err := toml.Marshal(config)
	if err != nil {
		return err
	}

	// 确保用户配置目录存在
	userConfigDir := filepath.Dir(cs.userConfigFile)
	if err := os.MkdirAll(userConfigDir, 0755); err != nil {
		return err
	}

	// 保存到用户配置文件
	if err := os.WriteFile(cs.userConfigFile, data, 0644); err != nil {
		return err
	}

	// 3. 更新内存中的配置
	cs.Config = &config

	return nil
}

// saveToFile 保存配置到文件
func (cs *ConfigService) saveToFile() {
	Logger.Debug("save config to file", zap.Any("config", cs.Config))
	data, err := toml.Marshal(cs.Config)
	if err != nil {
		Logger.Error("save config to file failed", zap.Error(err))
		return
	}
	if err := os.WriteFile(cs.userConfigFile, data, 0644); err != nil {
		Logger.Error("save config to file failed", zap.Error(err))
		return
	}
}

// UpdateAppConfig 更新应用配置并持久化到文件
func (cs *ConfigService) UpdateAppConfig() error {
	if data, err := toml.Marshal(cs.appConfig); err != nil {
		return err
	} else {
		return os.WriteFile(cs.appConfigFile, data, 0644)
	}
}

// CheckAndUpdateDBVersion 检查是否需要更新数据库版本
// 返回 true 表示需要执行数据库初始化（版本号不匹配或首次），false 表示跳过
func (cs *ConfigService) CheckAndUpdateDBVersion() bool {
	// 判断是否需要初始化
	needInit := cs.appConfig.DBVersion < CurrentDBVersion || cs.appConfig.DBVersion == 0

	// 如果需要初始化，更新配置对象和文件中的版本号
	if needInit {
		cs.appConfig.DBVersion = CurrentDBVersion
		if err := cs.UpdateAppConfig(); err != nil {
			Logger.Error("update app config db version failed", zap.Error(err))
		}
	}

	return needInit
}

// GetSyncConfig 获取同步配置
func (cs *ConfigService) GetSyncConfig() *internalsync.SyncConfig {
	return &cs.Config.Sync
}

// SaveSyncConfig 保存同步配置
func (cs *ConfigService) SaveSyncConfig(syncConfig internalsync.SyncConfig) error {
	cs.Config.Sync = syncConfig
	Logger.Debug("save sync config", zap.Any("syncConfig", syncConfig))
	cs.saveToFile()
	return nil
}

// SetUserPassword 设置用户密码用于加密/解密
func (cs *ConfigService) SetUserPassword(password string) {
	Logger.Debug("set user password")

	// 使用写锁保护全局密码变量
	passwordMutex.Lock()
	globalUserPassword = password
	passwordMutex.Unlock()

	// 使用 chanMutex 保护 channel 操作
	cs.chanMutex.Lock()
	if cs.passwordChan != nil {
		select {
		case cs.passwordChan <- struct{}{}:
		default:
			// channel 已满或已关闭，跳过
		}
		cs.passwordChan = nil
	}
	cs.chanMutex.Unlock()

	Logger.Debug("user password set successfully")
}

// waitForPassword 触发事件并等待用户输入密码，超时60秒
func (cs *ConfigService) waitForPassword(reason string) error {
	Logger.Debug("waiting for user password input", zap.String("reason", reason))

	// 使用读锁检查密码
	passwordMutex.RLock()
	if globalUserPassword != "" {
		passwordMutex.RUnlock()
		return nil
	}
	passwordMutex.RUnlock()

	// 创建新的 channel
	cs.chanMutex.Lock()
	cs.passwordChan = make(chan struct{}, 1)
	cs.chanMutex.Unlock()

	app.Event.Emit(EventInputPassword, reason)

	// 等待密码输入完成，超时60秒
	select {
	case <-cs.passwordChan:
		// 使用读锁检查密码是否已设置
		passwordMutex.RLock()
		pwd := globalUserPassword
		passwordMutex.RUnlock()
		if pwd == "" {
			return errors.New("password not entered")
		}
		app.Event.Emit(EventInputPasswordClose, "")
		return nil
	case <-time.After(60 * time.Second):
		// 超时，安全地关闭 channel：先设为 nil 阻止新发送，再关闭
		cs.chanMutex.Lock()
		ch := cs.passwordChan
		cs.passwordChan = nil
		cs.chanMutex.Unlock()
		if ch != nil {
			close(ch)
		}
		return errors.New("password input timeout")
	}
}

// GetUserPassword 获取当前用户密码
func (cs *ConfigService) GetUserPassword() string {
	passwordMutex.RLock()
	defer passwordMutex.RUnlock()
	return globalUserPassword
}

// getPasswordWithPrompt 获取密码，如未设置则提示用户输入
func (cs *ConfigService) getPasswordWithPrompt(reason string) (string, error) {
	password := cs.GetUserPassword()
	if password == "" {
		if err := cs.waitForPassword(reason); err != nil {
			return "", err
		}
		password = cs.GetUserPassword()
	}
	return password, nil
}
