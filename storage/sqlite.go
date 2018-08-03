package storage

import (
	_ "github.com/mattn/go-sqlite3"
	"database/sql"
	"sync"
)

const (
	createSql = `
CREATE TABLE IF NOT EXISTS urls (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	url VARCHAR(256) NOT NULL,
	downloaded INTEGER DEFAULT 0
)

`
	addSql        = `INSERT INTO urls (url) VALUES (?)`
	downloadedSql = `UPDATE urls SET downloaded=1 WHERE url=?`
	nextUrlSql    = `SELECT url FROM urls where downloaded=0 ORDER BY id ASC limit 1`
)

type SqliteQueue struct {
	db *sql.DB
	mu sync.Mutex
}

func NewSqliteStorage(dbFile string) (Storage, error) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(createSql)
	if err != nil {
		return nil, err
	}
	return &SqliteQueue{db: db}, nil
}

func (q *SqliteQueue) Add(url string) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	_, err := q.db.Exec(addSql, url)
	return err
}

func (q *SqliteQueue) Downloaded(url string) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	_, err := q.db.Exec(downloadedSql, url)
	return err
}

func (q *SqliteQueue) NextUrl() (string, error) {
	row := q.db.QueryRow(nextUrlSql)
	var url string
	err := row.Scan(&url)
	if err != nil {
		return "", err
	}
	return url, nil
}
