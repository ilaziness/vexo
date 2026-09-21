package database

import (
	"database/sql"
	"fmt"
	"path/filepath"

	"go.uber.org/zap"
	_ "modernc.org/sqlite"
)

var dbFileName = "vexo.db"

type Database struct {
	db                 *sql.DB
	dbPath             string
	logger             *zap.Logger
	BookmarkRepo       *BookmarkRepository
	UserCommandRepo    *UserCommandRepository
	CommandHistoryRepo *CommandHistoryRepository
	AISessionRepo      AISessionRepository
}

func NewDatabase(userDataDir string, logger *zap.Logger) *Database {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &Database{
		dbPath: filepath.Join(userDataDir, dbFileName),
		logger: logger,
	}
}

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
	d.logger.Debug("db opened successfully", zap.String("dbPath", d.dbPath))
	if err := runMigrations(d.db, d.logger); err != nil {
		return fmt.Errorf("run migrations failed: %w", err)
	}
	d.logger.Debug("db initialized successfully")
	d.BookmarkRepo = NewBookmarkRepository(d.db, d.logger)
	d.UserCommandRepo = NewUserCommandRepository(d.db, d.logger)
	d.CommandHistoryRepo = NewCommandHistoryRepository(d.db, d.logger)
	d.AISessionRepo = NewSQLiteAISessionRepository(d.db)
	return nil
}

func (d *Database) Close() error {
	if d.db != nil {
		if err := d.db.Close(); err != nil {
			return fmt.Errorf("close db failed: %w", err)
		}
		d.db = nil
		d.logger.Debug("db closed successfully")
	}
	return nil
}
