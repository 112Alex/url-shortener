package sqlite

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite" // init sqlite driver
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

	// Для sqlite используем драйвер modernc.org/sqlite
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

func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {
	const op = "storage.sqlite.SaveURL"

	// Проверяем, что urlToSave не пустой
	if urlToSave == "" {
		return 0, fmt.Errorf("%s: url не указан", op)
	}

	// Проверяем, что alias не пустой
	if alias == "" {
		return 0, fmt.Errorf("%s: alias не указан", op)
	}

	// Сохраняем URL в базе данных
	result, err := s.db.Exec("INSERT INTO url (alias, url) VALUES (?, ?)", alias, urlToSave)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
