package storage

import (
	"database/sql"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

const (
	createUrlsTable = `
CREATE TABLE IF NOT EXISTS urls (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	url VARCHAR(256) NOT NULL UNIQUE ,
	status INTEGER DEFAULT 0 -- 0 new, 1 downloaded, 2 failed
)
`
	createVideosTable = `
CREATE TABLE IF NOT EXISTS videos (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	url VARCHAR(256) NOT NULL,
	title VARCHAR(512) NOT NULL,
	duration INTEGER DEFAULT 0,
	thumbnail VARCHAR(256) DEFAULT "",
	author VARCHAR(256) DEFAULT "",
	downloadUrl VARCHAR (256) DEFAULT ""
)
`

	addSql          = `INSERT INTO urls (url) VALUES (?)`
	setUrlStatusSql = `UPDATE urls SET status=? WHERE url=?`
	nextUrlSql      = `SELECT url FROM urls where status=0 ORDER BY id ASC limit 1`

	saveVideoSql = `INSERT INTO videos (url, title, duration, thumbnail) VALUES (?,?,?,?)`
)

type SqliteStorage struct {
	db *sql.DB
	mu sync.Mutex
}

func NewSqliteStorage(dbFile string) (Storage, error) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(createUrlsTable)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(createVideosTable)
	if err != nil {
		return nil, err
	}
	return &SqliteStorage{db: db}, nil
}

func (s *SqliteStorage) Add(url string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec(addSql, url)
	return err
}

func (s *SqliteStorage) setStatus(url string, status int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec(setUrlStatusSql, status, url)
	return err
}

func (s *SqliteStorage) Downloaded(url string) error {
	return s.setStatus(url, URL_STATUS_DOWNLOADED)
}

func (s *SqliteStorage) DownloadFailed(url string) error {
	return s.setStatus(url, URL_STATUS_DOWNLOAD_FAILED)
}

func (s *SqliteStorage) NextUrl() (string, error) {
	row := s.db.QueryRow(nextUrlSql)
	var url string
	err := row.Scan(&url)

	if err != nil && err != sql.ErrNoRows {
		return "", err
	}
	return url, nil
}

func (s *SqliteStorage) SaveVideo(v Video) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, err := s.db.Exec(saveVideoSql, v.Url, v.Title, v.Length, v.Thumbnail)
	return err
}
