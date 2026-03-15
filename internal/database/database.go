package database

import (
	"database/sql"
	"fmt"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
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
}

// NewDatabase 创建数据库实例
func NewDatabase(userDataDir string) *Database {
	dbPath := filepath.Join(userDataDir, dbFileName)
	return &Database{
		dbPath: dbPath,
	}
}

// Initialize 初始化数据库，创建表结构
// skipCreateTables: 如果为 true，则跳过表结构的创建和变更
func (d *Database) Initialize(skipCreateTables bool) error {
	var err error
	d.db, err = sql.Open("sqlite3", d.dbPath)
	if err != nil {
		return fmt.Errorf("open db failed: %w", err)
	}

	d.db.SetMaxOpenConns(25)
	d.db.SetMaxIdleConns(5)

	if err := d.db.Ping(); err != nil {
		return fmt.Errorf("ping db failed: %w", err)
	}

	Logger.Debug("db opened successfully", zap.String("dbPath", d.dbPath))

	if !skipCreateTables {
		if err := d.createTables(); err != nil {
			return fmt.Errorf("create tables failed: %w", err)
		}
	}

	Logger.Debug("db initialized successfully")

	d.BookmarkRepo = NewBookmarkRepository(d.db)
	d.UserCommandRepo = NewUserCommandRepository(d.db)
	d.CommandHistoryRepo = NewCommandHistoryRepository(d.db)

	return nil
}

// createTables 创建数据表
func (d *Database) createTables() error {
	// 书签分组表
	createGroupTable := `
	CREATE TABLE IF NOT EXISTS bookmark_groups (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	// 书签表
	createBookmarkTable := `
	CREATE TABLE IF NOT EXISTS bookmarks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,  -- 自增主键
		bookmark_id TEXT UNIQUE NOT NULL,      -- 字符串业务 ID（唯一索引）
		group_id INTEGER NOT NULL,
		title TEXT NOT NULL,
		host TEXT NOT NULL,
		port INTEGER NOT NULL,
		user TEXT,
		password TEXT,
		private_key TEXT,
		private_key_password TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (group_id) REFERENCES bookmark_groups(id) ON DELETE CASCADE
	);`

	// 用户自定义命令表
	createUserCommandTable := `
	CREATE TABLE IF NOT EXISTS user_commands (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		category TEXT NOT NULL,
		name TEXT NOT NULL,
		command TEXT NOT NULL,
		description TEXT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	// 命令历史表
	createCommandHistoryTable := `
	CREATE TABLE IF NOT EXISTS command_history (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		command TEXT NOT NULL,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	// 创建索引
	createIndexes := `
	CREATE UNIQUE INDEX IF NOT EXISTS idx_bookmarks_bookmark_id ON bookmarks(bookmark_id);
	CREATE INDEX IF NOT EXISTS idx_bookmarks_group_id ON bookmarks(group_id);
	CREATE INDEX IF NOT EXISTS idx_user_commands_category ON user_commands(category);
	CREATE INDEX IF NOT EXISTS idx_command_history_timestamp ON command_history(timestamp DESC);
	`

	// 执行建表 SQL
	tables := []string{
		createGroupTable,
		createBookmarkTable,
		createUserCommandTable,
		createCommandHistoryTable,
		createIndexes,
	}

	for _, table := range tables {
		if _, err := d.db.Exec(table); err != nil {
			return fmt.Errorf("exec sql failed: %w", err)
		}
	}

	Logger.Debug("tables created successfully")
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
		Logger.Debug("db closed successfully")
	}
	return nil
}
