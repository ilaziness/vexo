package database

import (
	"database/sql"
	"fmt"
	"time"

	"go.uber.org/zap"
)

// 数据库操作错误消息常量
const (
	errBeginTx     = "begin tx failed: %w"
	errCommitTx    = "commit tx failed: %w"
	errDeleteQuery = "delete %s failed: %w"
	errInsertQuery = "insert %s failed: %w"
	errQuery       = "query %s failed: %w"

	// 数据表名称常量
	tableNameGroups         = "groups"
	tableNameBookmarks      = "bookmarks"
	tableNameBookmarkGroups = "bookmark_groups"
	tableNameUserCommands   = "user commands"
	tableNameCommandHistory = "command history"

	// 错误消息中的表名描述常量
	errTableNameUserCommand = "user command"
)

// BookmarkGroupDB 书签分组数据库模型
type BookmarkGroupDB struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// BookmarkDB 书签数据库模型
type BookmarkDB struct {
	AutoID             int       `json:"auto_id"` // 数据库自增主键
	ID                 string    `json:"id"`      // 字符串业务 ID
	GroupID            int       `json:"group_id"`
	Title              string    `json:"title"`
	Host               string    `json:"host"`
	Port               int       `json:"port"`
	User               string    `json:"user"`
	Password           string    `json:"password"`
	PrivateKey         string    `json:"private_key"`
	PrivateKeyPassword string    `json:"private_key_password"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// UserCommandDB 用户命令数据库模型
type UserCommandDB struct {
	ID          int       `json:"id"`
	Category    string    `json:"category"`
	Name        string    `json:"name"`
	Command     string    `json:"command"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// CommandHistoryDB 命令历史数据库模型
type CommandHistoryDB struct {
	ID        int       `json:"id"`
	Command   string    `json:"command"`
	Timestamp time.Time `json:"timestamp"`
}

// BookmarkRepository 书签数据访问接口
type BookmarkRepository struct {
	db *sql.DB
}

// NewBookmarkRepository 创建书签数据访问实例
func NewBookmarkRepository(db *sql.DB) *BookmarkRepository {
	return &BookmarkRepository{db: db}
}

// GetAllGroups 获取所有分组
func (r *BookmarkRepository) GetAllGroups() ([]*BookmarkGroupDB, error) {
	query := `SELECT id, name, created_at FROM bookmark_groups ORDER BY id`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf(errQuery, tableNameGroups, err)
	}
	defer rows.Close()

	groups := make([]*BookmarkGroupDB, 0)
	for rows.Next() {
		var g BookmarkGroupDB
		if err := rows.Scan(&g.ID, &g.Name, &g.CreatedAt); err != nil {
			Logger.Error("scan group failed", zap.Error(err))
			continue
		}
		groups = append(groups, &g)
	}

	return groups, nil
}

// GetAllBookmarks 获取所有书签
func (r *BookmarkRepository) GetAllBookmarks() ([]*BookmarkDB, error) {
	query := `SELECT id, bookmark_id, group_id, title, host, port, user, password, private_key, private_key_password, created_at, updated_at 
			  FROM bookmarks ORDER BY id`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf(errQuery, tableNameBookmarks, err)
	}
	defer rows.Close()

	bookmarks := make([]*BookmarkDB, 0)
	for rows.Next() {
		var b BookmarkDB
		if err := rows.Scan(
			&b.AutoID, &b.ID, &b.GroupID, &b.Title, &b.Host, &b.Port,
			&b.User, &b.Password, &b.PrivateKey, &b.PrivateKeyPassword,
			&b.CreatedAt, &b.UpdatedAt,
		); err != nil {
			Logger.Error("scan bookmark failed", zap.Error(err))
			continue
		}
		bookmarks = append(bookmarks, &b)
	}

	return bookmarks, nil
}

// SaveBookmarks 保存书签分组和书签（全量覆盖，保留用于批量导入场景）
func (r *BookmarkRepository) SaveBookmarks(groups []*BookmarkGroupDB, bookmarks []*BookmarkDB) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("begin tx failed: %w", err)
	}
	defer tx.Rollback()

	// 清空现有数据
	if _, err := tx.Exec(`DELETE FROM bookmarks`); err != nil {
		return fmt.Errorf(errDeleteQuery, tableNameBookmarks, err)
	}
	if _, err := tx.Exec(`DELETE FROM bookmark_groups`); err != nil {
		return fmt.Errorf(errDeleteQuery, tableNameBookmarkGroups, err)
	}

	// 插入分组并获取 ID
	groupIDMap := make(map[int]int) // 原索引 -> 新 ID
	for i, group := range groups {
		result, err := tx.Exec(`INSERT INTO bookmark_groups (name) VALUES (?)`, group.Name)
		if err != nil {
			return fmt.Errorf(errInsertQuery, "group", err)
		}
		id, _ := result.LastInsertId()
		groupIDMap[i] = int(id)
	}

	// 插入书签（按组 ID 分组）
	groupBookmarks := make(map[int][]*BookmarkDB)
	for _, b := range bookmarks {
		groupBookmarks[b.GroupID] = append(groupBookmarks[b.GroupID], b)
	}

	for groupIdx := range groups {
		groupID := groupIDMap[groupIdx]
		for _, bookmark := range groupBookmarks[groupID] {
			_, err := tx.Exec(
				`INSERT INTO bookmarks (bookmark_id, group_id, title, host, port, user, password, private_key, private_key_password, created_at, updated_at) 
				 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
				bookmark.ID, groupID, bookmark.Title, bookmark.Host, bookmark.Port,
				bookmark.User, bookmark.Password, bookmark.PrivateKey, bookmark.PrivateKeyPassword,
				bookmark.CreatedAt, bookmark.UpdatedAt,
			)
			if err != nil {
				return fmt.Errorf(errInsertQuery, "bookmark", err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf(errCommitTx, err)
	}

	Logger.Debug("bookmarks saved to db")
	return nil
}

// GetBookmarkByID 按字符串 ID 查询书签
func (r *BookmarkRepository) GetBookmarkByID(id string) (*BookmarkDB, error) {
	query := `SELECT id, bookmark_id, group_id, title, host, port, user, password, private_key, private_key_password, created_at, updated_at 
			  FROM bookmarks WHERE bookmark_id = ?`
	row := r.db.QueryRow(query, id)

	var b BookmarkDB
	err := row.Scan(&b.AutoID, &b.ID, &b.GroupID, &b.Title, &b.Host, &b.Port,
		&b.User, &b.Password, &b.PrivateKey, &b.PrivateKeyPassword,
		&b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("bookmark not found")
		}
		return nil, fmt.Errorf(errQuery, "bookmark by id", err)
	}

	return &b, nil
}

// GetBookmarkByAutoID 按自增 ID 查询书签
func (r *BookmarkRepository) GetBookmarkByAutoID(id int) (*BookmarkDB, error) {
	query := `SELECT id, bookmark_id, group_id, title, host, port, user, password, private_key, private_key_password, created_at, updated_at 
			  FROM bookmarks WHERE id = ?`
	row := r.db.QueryRow(query, id)

	var b BookmarkDB
	err := row.Scan(&b.AutoID, &b.ID, &b.GroupID, &b.Title, &b.Host, &b.Port,
		&b.User, &b.Password, &b.PrivateKey, &b.PrivateKeyPassword,
		&b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("bookmark not found")
		}
		return nil, fmt.Errorf(errQuery, "bookmark by auto id", err)
	}

	return &b, nil
}

// GetBookmarksByGroupID 按分组 ID 查询书签
func (r *BookmarkRepository) GetBookmarksByGroupID(groupID int) ([]*BookmarkDB, error) {
	query := `SELECT id, bookmark_id, group_id, title, host, port, user, password, private_key, private_key_password, created_at, updated_at 
			  FROM bookmarks WHERE group_id = ? ORDER BY id`
	rows, err := r.db.Query(query, groupID)
	if err != nil {
		return nil, fmt.Errorf(errQuery, "bookmarks by group id", err)
	}
	defer rows.Close()

	bookmarks := make([]*BookmarkDB, 0)
	for rows.Next() {
		var b BookmarkDB
		if err := rows.Scan(
			&b.AutoID, &b.ID, &b.GroupID, &b.Title, &b.Host, &b.Port,
			&b.User, &b.Password, &b.PrivateKey, &b.PrivateKeyPassword,
			&b.CreatedAt, &b.UpdatedAt,
		); err != nil {
			Logger.Error("scan bookmark failed", zap.Error(err))
			continue
		}
		bookmarks = append(bookmarks, &b)
	}

	return bookmarks, nil
}

// InsertBookmark 插入书签
func (r *BookmarkRepository) InsertBookmark(bookmark *BookmarkDB) error {
	query := `INSERT INTO bookmarks (bookmark_id, group_id, title, host, port, user, password, private_key, private_key_password, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query,
		bookmark.ID, bookmark.GroupID, bookmark.Title, bookmark.Host, bookmark.Port,
		bookmark.User, bookmark.Password, bookmark.PrivateKey, bookmark.PrivateKeyPassword,
		bookmark.CreatedAt, bookmark.UpdatedAt)
	if err != nil {
		return fmt.Errorf(errInsertQuery, "bookmark", err)
	}

	Logger.Debug("bookmark inserted", zap.String("id", bookmark.ID))
	return nil
}

// UpdateBookmark 更新书签（根据字符串 ID）
func (r *BookmarkRepository) UpdateBookmark(bookmark *BookmarkDB) error {
	query := `UPDATE bookmarks 
			  SET title = ?, host = ?, port = ?, user = ?, password = ?, 
			      private_key = ?, private_key_password = ?, updated_at = ? 
			  WHERE bookmark_id = ?`
	_, err := r.db.Exec(query,
		bookmark.Title, bookmark.Host, bookmark.Port,
		bookmark.User, bookmark.Password, bookmark.PrivateKey, bookmark.PrivateKeyPassword,
		bookmark.UpdatedAt, bookmark.ID)
	if err != nil {
		return fmt.Errorf(errInsertQuery, "update bookmark", err)
	}

	Logger.Debug("bookmark updated", zap.String("id", bookmark.ID))
	return nil
}

// DeleteBookmark 删除书签（根据字符串 ID）
func (r *BookmarkRepository) DeleteBookmark(id string) error {
	_, err := r.db.Exec(`DELETE FROM bookmarks WHERE bookmark_id = ?`, id)
	if err != nil {
		return fmt.Errorf(errDeleteQuery, "bookmark", err)
	}

	Logger.Debug("bookmark deleted", zap.String("id", id))
	return nil
}

// DeleteBookmarkByAutoID 删除书签（根据自增 ID）
func (r *BookmarkRepository) DeleteBookmarkByAutoID(id int) error {
	_, err := r.db.Exec(`DELETE FROM bookmarks WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf(errDeleteQuery, "bookmark by auto id", err)
	}

	Logger.Debug("bookmark deleted by auto id", zap.Int("id", id))
	return nil
}

// InsertGroup 插入分组
func (r *BookmarkRepository) InsertGroup(group *BookmarkGroupDB) error {
	query := `INSERT INTO bookmark_groups (name) VALUES (?)`
	_, err := r.db.Exec(query, group.Name)
	if err != nil {
		return fmt.Errorf(errInsertQuery, "group", err)
	}

	Logger.Debug("group inserted", zap.String("name", group.Name))
	return nil
}

// GetGroupByID 按ID查询分组
func (r *BookmarkRepository) GetGroupByID(id int) (*BookmarkGroupDB, error) {
	query := `SELECT id, name, created_at FROM bookmark_groups WHERE id = ?`
	row := r.db.QueryRow(query, id)

	var g BookmarkGroupDB
	err := row.Scan(&g.ID, &g.Name, &g.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("group not found")
		}
		return nil, fmt.Errorf(errQuery, "group by id", err)
	}

	return &g, nil
}

// GetGroupByName 按名称查询分组
func (r *BookmarkRepository) GetGroupByName(name string) (*BookmarkGroupDB, error) {
	query := `SELECT id, name, created_at FROM bookmark_groups WHERE name = ?`
	row := r.db.QueryRow(query, name)

	var g BookmarkGroupDB
	err := row.Scan(&g.ID, &g.Name, &g.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("group not found")
		}
		return nil, fmt.Errorf(errQuery, "group by name", err)
	}

	return &g, nil
}

// UpdateGroupName 更新分组名称
func (r *BookmarkRepository) UpdateGroupName(oldName, newName string) error {
	query := `UPDATE bookmark_groups SET name = ? WHERE name = ?`
	_, err := r.db.Exec(query, newName, oldName)
	if err != nil {
		return fmt.Errorf(errInsertQuery, "update group name", err)
	}

	Logger.Debug("group name updated", zap.String("old", oldName), zap.String("new", newName))
	return nil
}

// GetGroupBookmarkCount 获取分组下的书签数量
func (r *BookmarkRepository) GetGroupBookmarkCount(groupName string) (int, error) {
	query := `SELECT COUNT(b.id) FROM bookmarks b
			  INNER JOIN bookmark_groups g ON b.group_id = g.id
			  WHERE g.name = ?`
	var count int
	err := r.db.QueryRow(query, groupName).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf(errQuery, "group bookmark count", err)
	}
	return count, nil
}

// DeleteGroup 删除分组
func (r *BookmarkRepository) DeleteGroup(name string) error {
	_, err := r.db.Exec(`DELETE FROM bookmark_groups WHERE name = ?`, name)
	if err != nil {
		return fmt.Errorf(errDeleteQuery, "group", err)
	}

	Logger.Debug("group deleted", zap.String("name", name))
	return nil
}

// UserCommandRepository 用户命令数据访问接口
type UserCommandRepository struct {
	db *sql.DB
}

// NewUserCommandRepository 创建用户命令数据访问实例
func NewUserCommandRepository(db *sql.DB) *UserCommandRepository {
	return &UserCommandRepository{db: db}
}

// GetAllCommands 获取所有用户命令
func (r *UserCommandRepository) GetAllCommands() ([]*UserCommandDB, error) {
	query := `SELECT id, category, name, command, description, created_at FROM user_commands ORDER BY category, name`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf(errQuery, tableNameUserCommands, err)
	}
	defer rows.Close()

	commands := make([]*UserCommandDB, 0)
	for rows.Next() {
		var cmd UserCommandDB
		if err := rows.Scan(&cmd.ID, &cmd.Category, &cmd.Name, &cmd.Command, &cmd.Description, &cmd.CreatedAt); err != nil {
			Logger.Error("scan user command failed", zap.Error(err))
			continue
		}
		commands = append(commands, &cmd)
	}

	return commands, nil
}

// SaveCommands 保存用户命令（全量覆盖，保留用于批量导入场景）
func (r *UserCommandRepository) SaveCommands(commands []*UserCommandDB) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf(errBeginTx, err)
	}
	defer tx.Rollback()

	// 清空现有数据
	if _, err := tx.Exec(`DELETE FROM user_commands`); err != nil {
		return fmt.Errorf(errDeleteQuery, tableNameUserCommands, err)
	}

	// 插入新数据
	for _, cmd := range commands {
		_, err := tx.Exec(
			`INSERT INTO user_commands (category, name, command, description, created_at) VALUES (?, ?, ?, ?, ?)`,
			cmd.Category, cmd.Name, cmd.Command, cmd.Description, cmd.CreatedAt,
		)
		if err != nil {
			return fmt.Errorf(errInsertQuery, errTableNameUserCommand, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf(errCommitTx, err)
	}

	Logger.Debug("user commands saved to db")
	return nil
}

// InsertCommand 插入用户命令
func (r *UserCommandRepository) InsertCommand(command *UserCommandDB) (int64, error) {
	query := `INSERT INTO user_commands (category, name, command, description, created_at) VALUES (?, ?, ?, ?, ?)`
	result, err := r.db.Exec(query, command.Category, command.Name, command.Command, command.Description, command.CreatedAt)
	if err != nil {
		return 0, fmt.Errorf(errInsertQuery, errTableNameUserCommand, err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get last insert id failed: %w", err)
	}

	Logger.Debug("user command inserted", zap.String("name", command.Name))
	return id, nil
}

// UpdateCommand 更新用户命令
func (r *UserCommandRepository) UpdateCommand(command *UserCommandDB) error {
	query := `UPDATE user_commands SET category = ?, name = ?, command = ?, description = ? WHERE id = ?`
	_, err := r.db.Exec(query, command.Category, command.Name, command.Command, command.Description, command.ID)
	if err != nil {
		return fmt.Errorf(errInsertQuery, "update user command", err)
	}

	Logger.Debug("user command updated", zap.String("name", command.Name))
	return nil
}

// DeleteCommand 删除用户命令
func (r *UserCommandRepository) DeleteCommand(id int) error {
	_, err := r.db.Exec(`DELETE FROM user_commands WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf(errDeleteQuery, errTableNameUserCommand, err)
	}

	Logger.Debug("user command deleted", zap.Int("id", id))
	return nil
}

// GetCommandByCategoryAndName 按分类和名称查询命令
func (r *UserCommandRepository) GetCommandByCategoryAndName(category, name string) (*UserCommandDB, error) {
	query := `SELECT id, category, name, command, description, created_at FROM user_commands WHERE category = ? AND name = ?`
	row := r.db.QueryRow(query, category, name)

	var cmd UserCommandDB
	err := row.Scan(&cmd.ID, &cmd.Category, &cmd.Name, &cmd.Command, &cmd.Description, &cmd.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("command not found")
		}
		return nil, fmt.Errorf(errQuery, "command by category and name", err)
	}

	return &cmd, nil
}

// CommandHistoryRepository 命令历史数据访问接口
type CommandHistoryRepository struct {
	db *sql.DB
}

// NewCommandHistoryRepository 创建命令历史数据访问实例
func NewCommandHistoryRepository(db *sql.DB) *CommandHistoryRepository {
	return &CommandHistoryRepository{db: db}
}

// GetHistory 获取命令历史（最新 100 条）
func (r *CommandHistoryRepository) GetHistory() ([]*CommandHistoryDB, error) {
	query := `SELECT id, command, timestamp FROM command_history ORDER BY timestamp DESC LIMIT 100`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf(errQuery, tableNameCommandHistory, err)
	}
	defer rows.Close()

	history := make([]*CommandHistoryDB, 0)
	for rows.Next() {
		var h CommandHistoryDB
		if err := rows.Scan(&h.ID, &h.Command, &h.Timestamp); err != nil {
			Logger.Error("scan command history failed", zap.Error(err))
			continue
		}
		history = append(history, &h)
	}

	return history, nil
}

// SaveHistory 保存命令历史（全量覆盖，保留用于批量导入场景）
func (r *CommandHistoryRepository) SaveHistory(history []*CommandHistoryDB) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf(errBeginTx, err)
	}
	defer tx.Rollback()

	// 清空现有数据
	if _, err := tx.Exec(`DELETE FROM command_history`); err != nil {
		return fmt.Errorf(errDeleteQuery, tableNameCommandHistory, err)
	}

	// 插入新数据
	for _, h := range history {
		_, err := tx.Exec(
			`INSERT INTO command_history (command, timestamp) VALUES (?, ?)`,
			h.Command, h.Timestamp,
		)
		if err != nil {
			return fmt.Errorf(errInsertQuery, tableNameCommandHistory, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf(errCommitTx, err)
	}

	Logger.Debug("command history saved to db")
	return nil
}

// InsertHistory 插入命令历史
func (r *CommandHistoryRepository) InsertHistory(history *CommandHistoryDB) (int64, error) {
	query := `INSERT INTO command_history (command, timestamp) VALUES (?, ?)`
	result, err := r.db.Exec(query, history.Command, history.Timestamp)
	if err != nil {
		return 0, fmt.Errorf(errInsertQuery, "command history", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get last insert id failed: %w", err)
	}

	Logger.Debug("command history inserted", zap.String("command", history.Command))
	return id, nil
}

// ClearHistory 清空命令历史
func (r *CommandHistoryRepository) ClearHistory() error {
	_, err := r.db.Exec(`DELETE FROM command_history`)
	if err != nil {
		return fmt.Errorf(errDeleteQuery, tableNameCommandHistory, err)
	}

	Logger.Debug("command history cleared")
	return nil
}
