package database

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

const createTableSQL = `
CREATE TABLE IF NOT EXISTS notifications (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	link TEXT NOT NULL,
	channel_id TEXT NOT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	UNIQUE(link, channel_id)
);
`

type Repository struct {
	db *sql.DB
}

func New(dbPath string) (*Repository, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	if _, err := db.Exec(createTableSQL); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return &Repository{db: db}, nil
}

func (r *Repository) Close() error {
	return r.db.Close()
}

func (r *Repository) Exists(link, channelID string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM notifications WHERE link = ? AND channel_id = ?)",
		link, channelID,
	).Scan(&exists)
	return exists, err
}

func (r *Repository) Save(link, channelID string) error {
	_, err := r.db.Exec(
		"INSERT INTO notifications (link, channel_id) VALUES (?, ?)",
		link, channelID,
	)
	return err
}
