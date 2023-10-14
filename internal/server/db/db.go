package db

import (
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const (
	DIRECTORIES_SCHEMA = `CREATE TABLE IF NOT EXISTS DIRECTORIES (
    id INTEGER PRIMARY KEY UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    path TEXT UNIQUE NOT NULL,
    deleted BOOLEAN NOT NULL DEFAULT FALSE
);`
)

const (
	DIRECTORY_INSERT_STMT = "INSERT INTO directories(path) VALUES (?)"
	DIRECTORY_REMOVE_STMT = "DELETE FROM directories WHERE path = ?"
	UPDATE_STMT           = "UPDATE files SET hash = ?, deleted = ? WHERE path = ?"
	QUERY_STMT            = "SELECT path, hash, deleted FROM files WHERE path = ?"
)

type DB struct {
	*sqlx.DB
	mu sync.Mutex
}

type Base struct {
	ID        uint64    `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Directory struct {
	Path    string `db:"path"`
	Deleted bool   `db:"deleted"`
	Base
}

func NewDB(path string) *DB {
	db := sqlx.MustConnect("sqlite3", path)
	db.MustExec(DIRECTORIES_SCHEMA)
	return &DB{db, sync.Mutex{}}
}

func (db *DB) AddDirectory(path string) error {
	tx := db.MustBegin()
	_, err := tx.Exec(DIRECTORY_INSERT_STMT, path)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (db *DB) GetDirectories() ([]Directory, error) {
	var directories []Directory
	err := db.Select(&directories, "SELECT * FROM directories")
	if err != nil {
		return nil, err
	}
	return directories, nil
}

func (db *DB) RemoveDirectory(path string) error {
	tx := db.MustBegin()
	_, err := tx.Exec(DIRECTORY_REMOVE_STMT, path)
	if err != nil {
		return err
	}
	return tx.Commit()
}
