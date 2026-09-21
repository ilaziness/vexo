package database

import (
	"database/sql"
	"fmt"
	"time"

	"go.uber.org/zap"
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
	ProxyJumpID        string    `json:"proxy_jump_id"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// BookmarkRepository 书签数据访问接口
type BookmarkRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewBookmarkRepository(db *sql.DB, logger *zap.Logger) *BookmarkRepository {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &BookmarkRepository{db: db, logger: logger}
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
			r.logger.Error("scan group failed", zap.Error(err))
			continue
		}
		groups = append(groups, &g)
	}

	return groups, nil
}

// GetAllBookmarks 获取所有书签
func (r *BookmarkRepository) GetAllBookmarks() ([]*BookmarkDB, error) {
	query := `SELECT id, bookmark_id, group_id, title, host, port, user, password, private_key, private_key_password, proxy_jump_id, created_at, updated_at 
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
			&b.User, &b.Password, &b.PrivateKey, &b.PrivateKeyPassword, &b.ProxyJumpID,
			&b.CreatedAt, &b.UpdatedAt,
		); err != nil {
			r.logger.Error("scan bookmark failed", zap.Error(err))
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
				`INSERT INTO bookmarks (bookmark_id, group_id, title, host, port, user, password, private_key, private_key_password, proxy_jump_id, created_at, updated_at) 
				 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
				bookmark.ID, groupID, bookmark.Title, bookmark.Host, bookmark.Port,
				bookmark.User, bookmark.Password, bookmark.PrivateKey, bookmark.PrivateKeyPassword, bookmark.ProxyJumpID,
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

	r.logger.Debug("bookmarks saved to db")
	return nil
}

// GetBookmarkByID 按字符串 ID 查询书签
func (r *BookmarkRepository) GetBookmarkByID(id string) (*BookmarkDB, error) {
	query := `SELECT id, bookmark_id, group_id, title, host, port, user, password, private_key, private_key_password, proxy_jump_id, created_at, updated_at 
			  FROM bookmarks WHERE bookmark_id = ?`
	row := r.db.QueryRow(query, id)

	var b BookmarkDB
	err := row.Scan(&b.AutoID, &b.ID, &b.GroupID, &b.Title, &b.Host, &b.Port,
		&b.User, &b.Password, &b.PrivateKey, &b.PrivateKeyPassword, &b.ProxyJumpID,
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
	query := `SELECT id, bookmark_id, group_id, title, host, port, user, password, private_key, private_key_password, proxy_jump_id, created_at, updated_at 
			  FROM bookmarks WHERE id = ?`
	row := r.db.QueryRow(query, id)

	var b BookmarkDB
	err := row.Scan(&b.AutoID, &b.ID, &b.GroupID, &b.Title, &b.Host, &b.Port,
		&b.User, &b.Password, &b.PrivateKey, &b.PrivateKeyPassword, &b.ProxyJumpID,
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
	query := `SELECT id, bookmark_id, group_id, title, host, port, user, password, private_key, private_key_password, proxy_jump_id, created_at, updated_at 
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
			&b.User, &b.Password, &b.PrivateKey, &b.PrivateKeyPassword, &b.ProxyJumpID,
			&b.CreatedAt, &b.UpdatedAt,
		); err != nil {
			r.logger.Error("scan bookmark failed", zap.Error(err))
			continue
		}
		bookmarks = append(bookmarks, &b)
	}

	return bookmarks, nil
}

// GetBookmarkByTitleAndGroup 按标题和分组 ID 查询书签（用于检查重复）
func (r *BookmarkRepository) GetBookmarkByTitleAndGroup(title string, groupID int) (*BookmarkDB, error) {
	query := `SELECT id, bookmark_id, group_id, title, host, port, user, password, private_key, private_key_password, proxy_jump_id, created_at, updated_at 
			  FROM bookmarks WHERE title = ? AND group_id = ?`
	var b BookmarkDB
	err := r.db.QueryRow(query, title, groupID).Scan(
		&b.AutoID, &b.ID, &b.GroupID, &b.Title, &b.Host, &b.Port,
		&b.User, &b.Password, &b.PrivateKey, &b.PrivateKeyPassword, &b.ProxyJumpID,
		&b.CreatedAt, &b.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("bookmark not found")
		}
		return nil, fmt.Errorf(errQuery, "bookmark by title and group", err)
	}

	return &b, nil
}

// InsertBookmark 插入书签
func (r *BookmarkRepository) InsertBookmark(bookmark *BookmarkDB) error {
	query := `INSERT INTO bookmarks (bookmark_id, group_id, title, host, port, user, password, private_key, private_key_password, proxy_jump_id, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query,
		bookmark.ID, bookmark.GroupID, bookmark.Title, bookmark.Host, bookmark.Port,
		bookmark.User, bookmark.Password, bookmark.PrivateKey, bookmark.PrivateKeyPassword, bookmark.ProxyJumpID,
		bookmark.CreatedAt, bookmark.UpdatedAt)
	if err != nil {
		return fmt.Errorf(errInsertQuery, "bookmark", err)
	}

	r.logger.Debug("bookmark inserted", zap.String("id", bookmark.ID))
	return nil
}

// UpdateBookmark 更新书签（根据字符串 ID）
func (r *BookmarkRepository) UpdateBookmark(bookmark *BookmarkDB) error {
	query := `UPDATE bookmarks 
			  SET title = ?, host = ?, port = ?, user = ?, password = ?, 
			      private_key = ?, private_key_password = ?, proxy_jump_id = ?, updated_at = ? 
			  WHERE bookmark_id = ?`
	_, err := r.db.Exec(query,
		bookmark.Title, bookmark.Host, bookmark.Port,
		bookmark.User, bookmark.Password, bookmark.PrivateKey, bookmark.PrivateKeyPassword, bookmark.ProxyJumpID,
		bookmark.UpdatedAt, bookmark.ID)
	if err != nil {
		return fmt.Errorf(errInsertQuery, "update bookmark", err)
	}

	r.logger.Debug("bookmark updated", zap.String("id", bookmark.ID))
	return nil
}

// DeleteBookmark 删除书签（根据字符串 ID）
func (r *BookmarkRepository) DeleteBookmark(id string) error {
	_, err := r.db.Exec(`DELETE FROM bookmarks WHERE bookmark_id = ?`, id)
	if err != nil {
		return fmt.Errorf(errDeleteQuery, "bookmark", err)
	}

	r.logger.Debug("bookmark deleted", zap.String("id", id))
	return nil
}

// DeleteBookmarkByAutoID 删除书签（根据自增 ID）
func (r *BookmarkRepository) DeleteBookmarkByAutoID(id int) error {
	_, err := r.db.Exec(`DELETE FROM bookmarks WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf(errDeleteQuery, "bookmark by auto id", err)
	}

	r.logger.Debug("bookmark deleted by auto id", zap.Int("id", id))
	return nil
}

// InsertGroup 插入分组
func (r *BookmarkRepository) InsertGroup(group *BookmarkGroupDB) error {
	query := `INSERT INTO bookmark_groups (name) VALUES (?)`
	_, err := r.db.Exec(query, group.Name)
	if err != nil {
		return fmt.Errorf(errInsertQuery, "group", err)
	}

	r.logger.Debug("group inserted", zap.String("name", group.Name))
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

	r.logger.Debug("group name updated", zap.String("old", oldName), zap.String("new", newName))
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

	r.logger.Debug("group deleted", zap.String("name", name))
	return nil
}

