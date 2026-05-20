package database

import (
	"database/sql"
	"fmt"
	"path/filepath"

	"go.uber.org/zap"
	_ "modernc.org/sqlite"
)

// Logger 全局日志实例
var Logger *zap.Logger
var dbFileName = "vexo.db"

// Database 数据库管理实例
type Database struct {
	db                 *sql.DB
	dbPath             string
	BookmarkRepo       *BookmarkRepository
	UserCommandRepo    *UserCommandRepository
	CommandHistoryRepo *CommandHistoryRepository
	AISessionRepo      AISessionRepository
}

// NewDatabase 创建数据库实例
func NewDatabase(userDataDir string) *Database {
	dbPath := filepath.Join(userDataDir, dbFileName)
	return &Database{
		dbPath: dbPath,
	}
}

// Initialize 初始化数据库，执行迁移
func (d *Database) Initialize() error {
	var err error
	d.db, err = sql.Open("sqlite", d.dbPath)
	if err != nil {
		return fmt.Errorf("open db failed: %w", err)
	}

	d.db.SetMaxOpenConns(25)
	d.db.SetMaxIdleConns(5)

	if err := d.db.Ping(); err != nil {
		return fmt.Errorf("ping db failed: %w", err)
	}

	Logger.Debug("db opened successfully", zap.String("dbPath", d.dbPath))

	if err := runMigrations(d.db); err != nil {
		return fmt.Errorf("run migrations failed: %w", err)
	}

	Logger.Debug("db initialized successfully")

	d.BookmarkRepo = NewBookmarkRepository(d.db)
	d.UserCommandRepo = NewUserCommandRepository(d.db)
	d.CommandHistoryRepo = NewCommandHistoryRepository(d.db)
	d.AISessionRepo = NewSQLiteAISessionRepository(d.db)

	return nil
}

// GetDB 获取数据库连接实例
func (d *Database) GetDB() *sql.DB {
	return d.db
}

// Close 关闭数据库连接
func (d *Database) Close() error {
	if d.db != nil {
		if err := d.db.Close(); err != nil {
			return fmt.Errorf("close db failed: %w", err)
		}
		d.db = nil
		Logger.Debug("db closed successfully")
	}
	return nil
}
