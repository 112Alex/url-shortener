package sqlite

import (
	"database/sql"
	"fmt"
	"strings"
	"url-shortener/internal/storage"

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

	stmt, err := s.db.Prepare("INSERT INTO url(alias, url) VALUES (?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(alias, urlToSave)
	if err != nil {
		// Проверяем ошибку на нарушение уникальности
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrURLExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to insert id: %w", op, err)
	}

	return id, nil
}
