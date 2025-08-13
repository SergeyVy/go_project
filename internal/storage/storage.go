package storage

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type URLRecord struct {
	ID    int64  `json:"id"`
	Alias string `json:"alias"`
	URL   string `json:"url"`
}

type Storage struct {
	DB *sql.DB // теперь доступно в main.go
}

func New(dsn string) (*Storage, error) {
	const op = "storage.New"

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// создаём таблицу, если нет
	if _, err = db.Exec(`
CREATE TABLE IF NOT EXISTS urls (
  id     BIGSERIAL PRIMARY KEY,
  alias  TEXT UNIQUE NOT NULL,
  url    TEXT NOT NULL
);`); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{DB: db}, nil
}

// Получить запись по alias
func (s *Storage) GetByAlias(alias string) (URLRecord, error) {
	const op = "storage.GetByAlias"

	var rec URLRecord
	err := s.DB.QueryRow(
		`SELECT id, alias, url FROM urls WHERE alias = $1`,
		alias,
	).Scan(&rec.ID, &rec.Alias, &rec.URL)

	if errors.Is(err, sql.ErrNoRows) {
		return URLRecord{}, ErrNotFound
	}
	if err != nil {
		return URLRecord{}, fmt.Errorf("%s: %w", op, err)
	}
	return rec, nil
}

// Сохранить новый URL с alias
func (s *Storage) SaveURL(urlToSave, alias string) (int64, error) {
	const op = "storage.SaveURL"

	var id int64
	err := s.DB.QueryRow(
		`INSERT INTO urls(url, alias) VALUES ($1, $2) RETURNING id`,
		urlToSave, alias,
	).Scan(&id)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return 0, ErrURLExists
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

// Получить только URL по alias
func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.GetURL"

	stmt, err := s.DB.Prepare(`SELECT url FROM urls WHERE alias = $1`)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	var resURL string
	err = stmt.QueryRow(alias).Scan(&resURL)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrNotFound
	}
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return resURL, nil
}

// Удалить запись по alias
func (s *Storage) DeleteURL(alias string) error {
	const op = "storage.DeleteURL"

	stmt, err := s.DB.Prepare(`DELETE FROM urls WHERE alias = $1`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(alias)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}
