package main

import (
	"bytes"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"

	_ "github.com/mattn/go-sqlite3"
)

type Entry struct {
	path    string
	hash    []byte
	deleted bool
}

const (
	INSERT_STMT = "INSERT INTO files(path, hash) VALUES (?, ?)"
	UPDATE_STMT = "UPDATE files SET hash = ?, deleted = ? WHERE path = ?"
	QUERY_STMT  = "SELECT path, hash, deleted FROM files WHERE path = ?"
)

var (
	wg      sync.WaitGroup
	counter atomic.Int64
)

type DB struct {
	db *sql.DB
	sync.Mutex
}

func walk(dir string, db *DB) {
	defer wg.Done()

	db.Lock()
	defer db.Unlock()
	tx, err := db.db.Begin()
	if err != nil {
		log.Panic(err)
	}

	insertStmt, err := tx.Prepare(INSERT_STMT)
	if err != nil {
		log.Panic(err)
	}
	defer insertStmt.Close()

	updateStmt, err := tx.Prepare(UPDATE_STMT)
	if err != nil {
		log.Panic(err)
	}
	defer updateStmt.Close()

	queryStmt, err := tx.Prepare(QUERY_STMT)
	if err != nil {
		log.Panic(err)
	}
	defer queryStmt.Close()

	filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() && dir != path {
			wg.Add(1)
			go walk(path, db)
			return filepath.SkipDir
		}

		if d.Type().IsRegular() {
			path, err = filepath.Abs(path)
			if err != nil {
				return nil
			}
			f, err := os.Open(path)
			if err != nil {
				return nil
			}

			h := sha256.New()
			io.Copy(h, f)
			hash := h.Sum(nil)

			// check if file exists in the database
			var e *Entry
			if err := queryStmt.QueryRow(path).Scan(e); err != nil {
				if err == sql.ErrNoRows {
					// file doesn't exist in db
					_, err := insertStmt.Exec(path, hash)
					if err != nil {
						log.Panic(err)
					}
					return nil
				}
				if err != sql.ErrNoRows {
					log.Panic(err)
				}
			}

			// check if hash is same as e if same, do nothing
			if bytes.Equal(e.hash, hash) {
				log.Println("found match")
				return nil
			}
			// if not nil, update
			// not doing this
		}
		return nil
	})

	err = tx.Commit()
	if err != nil {
		log.Panic(err)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("file-sync requires 1 argument")
		return
	}

	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	createTable := `CREATE TABLE IF NOT EXISTS files 
  (path TEXT NOT NULL PRIMARY KEY, 
  hash BLOB NOT NULL, 
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, 
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, 
  deleted BOOLEAN NOT NULL DEFAULT FALSE
  )`

	_, err = db.Exec(createTable)
	if err != nil {
		log.Panic(err)
	}

	path := os.Args[1]
	wg.Add(1)
	go walk(path, &DB{db: db})
	wg.Wait()
}
