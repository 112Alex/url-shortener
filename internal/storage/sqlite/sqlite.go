package sqlite

import (
	"database/sql"
	"fmt"

	//_ "github.com/mattn/go-sqlite3"  register sqlite driver
	_ "modernc.org/sqlite"
)

type Storage struct {
	db *sql.DB
}

func (s *Storage) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	// Проверяем, что путь не пустой
	if storagePath == "" {
		return nil, fmt.Errorf("%s: путь к базе данных не указан", op)
	}

	db, err := sql.Open("sqlite", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS url(
		id INTEGER PRIMARY KEY,
		alias TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL);
	CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}
