package services

import (
	"os"
	"path/filepath"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LogLevel 可由外部修改以调整运行时日志级别
var LogLevel = zap.NewAtomicLevelAt(zap.DebugLevel)

// 包级别 logger 供全局使用
var Logger *zap.Logger

type LogService struct {
	logger *zap.Logger
}

// InitLogger 初始化全局 logger。
// logPath: 日志文件路径（例如 "logs/vexo.log"），会在该目录下按月生成轮转文件。
// level: zap.AtomicLevel，用于运行时调整日志级别。
func InitLogger(logPath string, level zap.AtomicLevel) error {
	dir := filepath.Dir(logPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	// 按月轮转，文件名示例: vexo.log.2025-12
	writer, err := rotatelogs.New(
		logPath+".%Y-%m",
		rotatelogs.WithLinkName(logPath),
		rotatelogs.WithRotationTime(24*time.Hour*30),
		rotatelogs.WithMaxAge(365*24*time.Hour),
	)
	if err != nil {
		return err
	}

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "ts"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	// JSON 文件输出到轮转 writer，同时控制台也输出便于开发时查看
	fileCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.AddSync(writer),
		level,
	)

	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderCfg),
		zapcore.AddSync(os.Stdout),
		level,
	)

	core := zapcore.NewTee(fileCore, consoleCore)
	zlogger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	// 替换全局 logger 与 LogLevel
	Logger = zlogger
	LogLevel = level
	return nil
}

func NewLogService() *LogService {
	if Logger == nil {
		// 若未显式初始化，则使用默认路径和默认级别进行初始化（忽略错误）
		_ = InitLogger("logs/vexo.log", LogLevel)
	}
	return &LogService{
		logger: Logger,
	}
}

func (ls *LogService) Debug(msg string) {
	ls.logger.Debug(msg)
}

func (ls *LogService) Info(msg string) {
	ls.logger.Info(msg)
}

func (ls *LogService) Warn(msg string) {
	ls.logger.Warn(msg)
}

func (ls *LogService) Error(msg string) {
	ls.logger.Error(msg)
}
