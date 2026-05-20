package database

import (
	"database/sql"
	"fmt"

	"go.uber.org/zap"
)

// Migration 表示一个数据库迁移
type Migration struct {
	Version int
	Name    string
	Up      func(*sql.DB) error
}

// migrations 注册所有数据库迁移，按版本号升序排列
var migrations = []Migration{
	{Version: 1, Name: "init schema", Up: migrateInitSchema},
	{Version: 2, Name: "add proxy_jump_id", Up: migrateAddProxyJumpID},
	{Version: 3, Name: "add ai sessions", Up: migrateAddAISessions},
}

// migrateInitSchema 初始化数据库表结构（幂等）
func migrateInitSchema(db *sql.DB) error {
	createGroupTable := `
	CREATE TABLE IF NOT EXISTS bookmark_groups (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	createBookmarkTable := `
	CREATE TABLE IF NOT EXISTS bookmarks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		bookmark_id TEXT UNIQUE NOT NULL,
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

	createUserCommandTable := `
	CREATE TABLE IF NOT EXISTS user_commands (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		category TEXT NOT NULL,
		name TEXT NOT NULL,
		command TEXT NOT NULL,
		description TEXT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	createCommandHistoryTable := `
	CREATE TABLE IF NOT EXISTS command_history (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		command TEXT NOT NULL,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	createIndexes := `
	CREATE UNIQUE INDEX IF NOT EXISTS idx_bookmarks_bookmark_id ON bookmarks(bookmark_id);
	CREATE INDEX IF NOT EXISTS idx_bookmarks_group_id ON bookmarks(group_id);
	CREATE INDEX IF NOT EXISTS idx_user_commands_category ON user_commands(category);
	CREATE INDEX IF NOT EXISTS idx_command_history_timestamp ON command_history(timestamp DESC);
	`

	tables := []string{
		createGroupTable,
		createBookmarkTable,
		createUserCommandTable,
		createCommandHistoryTable,
		createIndexes,
	}

	for _, sql := range tables {
		if _, err := db.Exec(sql); err != nil {
			return fmt.Errorf("exec sql failed: %w", err)
		}
	}

	Logger.Debug("migration: init schema applied")
	return nil
}

// migrateAddProxyJumpID 添加 proxy_jump_id 列（幂等）
func migrateAddProxyJumpID(db *sql.DB) error {
	var columnName string
	err := db.QueryRow(`SELECT name FROM pragma_table_info('bookmarks') WHERE name = 'proxy_jump_id'`).Scan(&columnName)
	if err == sql.ErrNoRows {
		_, err := db.Exec(`ALTER TABLE bookmarks ADD COLUMN proxy_jump_id TEXT DEFAULT ''`)
		if err != nil {
			return fmt.Errorf("add column proxy_jump_id failed: %w", err)
		}
		Logger.Debug("migration: added proxy_jump_id column")
		return nil
	}
	if err != nil {
		return fmt.Errorf("check column proxy_jump_id failed: %w", err)
	}
	return nil
}

// migrateAddAISessions 添加 AI 会话和消息表（幂等）
func migrateAddAISessions(db *sql.DB) error {
	createSessionsTable := `
	CREATE TABLE IF NOT EXISTS ai_sessions (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		created_at INTEGER NOT NULL,
		updated_at INTEGER NOT NULL
	);`

	createMessagesTable := `
	CREATE TABLE IF NOT EXISTS ai_messages (
		id TEXT PRIMARY KEY,
		session_id TEXT NOT NULL,
		role TEXT NOT NULL,
		content TEXT NOT NULL,
		parts TEXT DEFAULT '',
		timestamp INTEGER NOT NULL,
		FOREIGN KEY (session_id) REFERENCES ai_sessions(id) ON DELETE CASCADE
	);`

	createIndexes := `
	CREATE INDEX IF NOT EXISTS idx_ai_messages_session_id ON ai_messages(session_id);
	CREATE INDEX IF NOT EXISTS idx_ai_sessions_updated_at ON ai_sessions(updated_at DESC);
	`

	sqls := []string{
		createSessionsTable,
		createMessagesTable,
		createIndexes,
	}

	for _, sql := range sqls {
		if _, err := db.Exec(sql); err != nil {
			return fmt.Errorf("exec sql failed: %w", err)
		}
	}

	Logger.Debug("migration: added ai_sessions and ai_messages tables")
	return nil
}

// createSchemaMigrationsTable 创建迁移记录表
func createSchemaMigrationsTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("create schema_migrations table failed: %w", err)
	}
	return nil
}

// getAppliedMigrations 获取已执行的迁移版本
func getAppliedMigrations(db *sql.DB) (map[int]bool, error) {
	rows, err := db.Query(`SELECT version FROM schema_migrations`)
	if err != nil {
		return nil, fmt.Errorf("query applied migrations failed: %w", err)
	}
	defer rows.Close()

	applied := make(map[int]bool)
	for rows.Next() {
		var version int
		if err := rows.Scan(&version); err != nil {
			return nil, fmt.Errorf("scan migration version failed: %w", err)
		}
		applied[version] = true
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate migration versions failed: %w", err)
	}

	return applied, nil
}

// runMigrations 执行所有未执行的迁移
func runMigrations(db *sql.DB) error {
	if err := createSchemaMigrationsTable(db); err != nil {
		return err
	}

	applied, err := getAppliedMigrations(db)
	if err != nil {
		return err
	}

	for _, m := range migrations {
		if applied[m.Version] {
			continue
		}

		Logger.Debug("running migration", zap.Int("version", m.Version), zap.String("name", m.Name))

		if err := m.Up(db); err != nil {
			return fmt.Errorf("migration %d (%s) failed: %w", m.Version, m.Name, err)
		}

		_, err := db.Exec(`INSERT INTO schema_migrations (version) VALUES (?)`, m.Version)
		if err != nil {
			return fmt.Errorf("record migration %d failed: %w", m.Version, err)
		}

		Logger.Debug("migration completed", zap.Int("version", m.Version), zap.String("name", m.Name))
	}

	return nil
}
