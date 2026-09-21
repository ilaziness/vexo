package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ilaziness/vexo/internal/ai"
	"github.com/ilaziness/vexo/internal/sync"
	"github.com/ilaziness/vexo/internal/system"
	"github.com/pelletier/go-toml/v2"
	"go.uber.org/zap"
)

const FileName = "config.toml"

var defaultFontFamily = []string{
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

type Config struct {
	General  GeneralConfig  `toml:"general"`
	Terminal TerminalConfig `toml:"terminal"`
	Sync     sync.SyncConfig `toml:"sync"`
	AI       ai.Config      `toml:"ai"`
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

type AppConfig struct {
	General GeneralConfig `toml:"general"`
}

type Store struct {
	Config         *Config
	appConfig      *AppConfig
	userConfigFile string
	appConfigFile  string
	logger         *zap.Logger
}

func Default() *Config {
	execDir := system.GetExecutableDir()
	return &Config{
		General: GeneralConfig{
			UserDataDir: filepath.Join(execDir, "data"),
			Theme:       "dark",
		},
		Terminal: TerminalConfig{
			Font:       strings.Join(defaultFontFamily, ","),
			FontSize:   14,
			LineHeight: 1,
		},
	}
}

func Load(logger *zap.Logger) (*Store, error) {
	execDir := system.GetExecutableDir()
	appConfigPath := filepath.Join(execDir, FileName)
	defaultConfig := Default()

	appConfig := &AppConfig{General: defaultConfig.General}
	if data, err := os.ReadFile(appConfigPath); err == nil {
		_ = toml.Unmarshal(data, appConfig)
	} else if os.IsNotExist(err) {
		if data, err := toml.Marshal(appConfig); err == nil {
			_ = os.WriteFile(appConfigPath, data, 0600)
		}
	}

	userDataDir := appConfig.General.UserDataDir
	if userDataDir == "" {
		userDataDir = defaultConfig.General.UserDataDir
	}
	if err := os.MkdirAll(userDataDir, 0755); err != nil {
		return nil, fmt.Errorf("create user data dir: %w", err)
	}

	userConfigPath := filepath.Join(userDataDir, FileName)
	finalConfig := defaultConfig
	finalConfig.General.UserDataDir = userDataDir

	if data, err := os.ReadFile(userConfigPath); err == nil {
		if err = toml.Unmarshal(data, finalConfig); err != nil {
			logger.Error("unmarshal user config failed", zap.Error(err))
		}
	} else if os.IsNotExist(err) {
		if data, err := toml.Marshal(finalConfig); err == nil {
			_ = os.WriteFile(userConfigPath, data, 0600)
		}
	}

	logger.Debug("config loaded", zap.String("appConfigPath", appConfigPath), zap.String("userConfigPath", userConfigPath))

	return &Store{
		Config:         finalConfig,
		appConfig:      appConfig,
		userConfigFile: userConfigPath,
		appConfigFile:  appConfigPath,
		logger:         logger,
	}, nil
}

func (s *Store) UserDataDir() string {
	return s.Config.General.UserDataDir
}

func (s *Store) Save() error {
	s.logger.Debug("save config to file", zap.String("path", s.userConfigFile))
	data, err := toml.Marshal(s.Config)
	if err != nil {
		s.logger.Error("save config to file failed", zap.Error(err))
		return fmt.Errorf("marshal config failed: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(s.userConfigFile), 0755); err != nil {
		s.logger.Error("create user config dir failed", zap.Error(err))
		return fmt.Errorf("create user config dir failed: %w", err)
	}
	if err := os.WriteFile(s.userConfigFile, data, 0600); err != nil {
		s.logger.Error("save config to file failed", zap.Error(err))
		return fmt.Errorf("write user config file failed: %w", err)
	}
	return nil
}

func (s *Store) SaveAppConfig() error {
	data, err := toml.Marshal(s.appConfig)
	if err != nil {
		return err
	}
	return os.WriteFile(s.appConfigFile, data, 0600)
}

func (s *Store) ApplyUserDataDir(dir string) error {
	if dir == s.Config.General.UserDataDir {
		return nil
	}
	s.appConfig = &AppConfig{General: GeneralConfig{UserDataDir: dir}}
	if err := s.SaveAppConfig(); err != nil {
		return fmt.Errorf("write app config file failed: %w", err)
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create user data dir failed: %w", err)
	}
	s.userConfigFile = filepath.Join(dir, FileName)
	s.Config.General.UserDataDir = dir
	return nil
}
