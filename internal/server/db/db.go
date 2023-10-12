package db

import (
	"sync"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const (
	INSERT_STMT = "INSERT INTO files(path, hash) VALUES (?, ?)"
	UPDATE_STMT = "UPDATE files SET hash = ?, deleted = ? WHERE path = ?"
	QUERY_STMT  = "SELECT path, hash, deleted FROM files WHERE path = ?"
	SCHEMA      = `CREATE TABLE IF NOT EXISTS files
  (path TEXT NOT NULL PRIMARY KEY,
  hash BLOB NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted BOOLEAN NOT NULL DEFAULT FALSE
  )`
)

type DB struct {
	*sqlx.DB
	mu sync.Mutex
}

func NewDB(path string) *DB {
	db := sqlx.MustConnect("sqlite3", path)
	db.MustExec(SCHEMA)
	return &DB{db, sync.Mutex{}}
}
