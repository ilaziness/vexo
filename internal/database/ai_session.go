package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// AISession represents an AI chat session
type AISession struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AIMessage represents a message in an AI chat session
type AIMessage struct {
	ID        string    `json:"id"`
	SessionID string    `json:"session_id"`
	Role      string    `json:"role"` // 'user' or 'assistant'
	Content   string    `json:"content"`
	Parts     string    `json:"parts"` // JSON array of ChatMessagePart
	Timestamp time.Time `json:"timestamp"`
}

// AISessionRepository provides data access operations for AI sessions and messages
type AISessionRepository interface {
	CreateSession(ctx context.Context, session *AISession) error
	GetSession(ctx context.Context, id string) (*AISession, error)
	ListSessions(ctx context.Context, limit int) ([]*AISession, error)
	UpdateSession(ctx context.Context, session *AISession) error
	DeleteSession(ctx context.Context, id string) error
	CreateMessage(ctx context.Context, msg *AIMessage) error
	ListMessages(ctx context.Context, sessionID string) ([]*AIMessage, error)
}

// SQLiteAISessionRepository implements AISessionRepository using SQLite
type SQLiteAISessionRepository struct {
	db *sql.DB
}

// NewSQLiteAISessionRepository creates a new SQLite repository
func NewSQLiteAISessionRepository(db *sql.DB) AISessionRepository {
	return &SQLiteAISessionRepository{db: db}
}

// CreateSession creates a new session
func (r *SQLiteAISessionRepository) CreateSession(ctx context.Context, session *AISession) error {
	session.ID = uuid.New().String()
	now := time.Now()
	session.CreatedAt = now
	session.UpdatedAt = now

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO ai_sessions (id, title, created_at, updated_at) VALUES (?, ?, ?, ?)`,
		session.ID, session.Title, session.CreatedAt.Unix(), session.UpdatedAt.Unix())
	if err != nil {
		return fmt.Errorf("create session failed: %w", err)
	}
	return nil
}

// GetSession retrieves a session by ID
func (r *SQLiteAISessionRepository) GetSession(ctx context.Context, id string) (*AISession, error) {
	var session AISession
	var createdAt, updatedAt int64

	err := r.db.QueryRowContext(ctx,
		`SELECT id, title, created_at, updated_at FROM ai_sessions WHERE id = ?`,
		id).Scan(&session.ID, &session.Title, &createdAt, &updatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("session not found")
	}
	if err != nil {
		return nil, fmt.Errorf("get session failed: %w", err)
	}

	session.CreatedAt = time.Unix(createdAt, 0)
	session.UpdatedAt = time.Unix(updatedAt, 0)
	return &session, nil
}

// ListSessions lists sessions ordered by updated_at DESC
func (r *SQLiteAISessionRepository) ListSessions(ctx context.Context, limit int) ([]*AISession, error) {
	query := `SELECT id, title, created_at, updated_at FROM ai_sessions ORDER BY updated_at DESC`
	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list sessions failed: %w", err)
	}
	defer rows.Close()

	var sessions []*AISession
	for rows.Next() {
		var session AISession
		var createdAt, updatedAt int64
		if err := rows.Scan(&session.ID, &session.Title, &createdAt, &updatedAt); err != nil {
			return nil, fmt.Errorf("scan session failed: %w", err)
		}
		session.CreatedAt = time.Unix(createdAt, 0)
		session.UpdatedAt = time.Unix(updatedAt, 0)
		sessions = append(sessions, &session)
	}

	return sessions, nil
}

// UpdateSession updates a session
func (r *SQLiteAISessionRepository) UpdateSession(ctx context.Context, session *AISession) error {
	session.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(ctx,
		`UPDATE ai_sessions SET title = ?, updated_at = ? WHERE id = ?`,
		session.Title, session.UpdatedAt.Unix(), session.ID)
	if err != nil {
		return fmt.Errorf("update session failed: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("session not found")
	}

	return nil
}

// DeleteSession deletes a session (CASCADE delete handles messages)
func (r *SQLiteAISessionRepository) DeleteSession(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM ai_sessions WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete session failed: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("session not found")
	}

	return nil
}

// CreateMessage creates a new message
func (r *SQLiteAISessionRepository) CreateMessage(ctx context.Context, msg *AIMessage) error {
	msg.ID = uuid.New().String()
	if msg.Timestamp.IsZero() {
		msg.Timestamp = time.Now()
	}

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO ai_messages (id, session_id, role, content, parts, timestamp) VALUES (?, ?, ?, ?, ?, ?)`,
		msg.ID, msg.SessionID, msg.Role, msg.Content, msg.Parts, msg.Timestamp.Unix())
	if err != nil {
		return fmt.Errorf("create message failed: %w", err)
	}
	return nil
}

// ListMessages lists messages for a session ordered by timestamp ASC
func (r *SQLiteAISessionRepository) ListMessages(ctx context.Context, sessionID string) ([]*AIMessage, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, session_id, role, content, parts, timestamp FROM ai_messages WHERE session_id = ? ORDER BY timestamp ASC`,
		sessionID)
	if err != nil {
		return nil, fmt.Errorf("list messages failed: %w", err)
	}
	defer rows.Close()

	var messages []*AIMessage
	for rows.Next() {
		var msg AIMessage
		var timestamp int64
		if err := rows.Scan(&msg.ID, &msg.SessionID, &msg.Role, &msg.Content, &msg.Parts, &timestamp); err != nil {
			return nil, fmt.Errorf("scan message failed: %w", err)
		}
		msg.Timestamp = time.Unix(timestamp, 0)
		messages = append(messages, &msg)
	}

	return messages, nil
}
