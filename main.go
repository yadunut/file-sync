package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Entry struct {
	Path    string
	Hash    []byte
	Deleted bool
}

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

var (
	wg      sync.WaitGroup
	counter atomic.Int64
	results map[string][]byte = make(map[string][]byte)
)

type DB struct {
	db *sqlx.DB
	sync.Mutex
}

func walk(dir string, db *DB) {
	defer wg.Done()

	db.Lock()
	defer db.Unlock()
	tx := db.db.MustBegin()

	insertStmt, err := tx.Preparex(INSERT_STMT)
	if err != nil {
		log.Panic(err)
	}
	defer insertStmt.Close()

	updateStmt, err := tx.Preparex(UPDATE_STMT)
	if err != nil {
		log.Panic(err)
	}
	defer updateStmt.Close()

	// queryStmt, err := tx.Preparex(QUERY_STMT)
	// if err != nil {
	// 	log.Panic(err)
	// }
	// defer queryStmt.Close()

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

			hashResult, ok := results[path]
			if !ok {
				_, err := insertStmt.Exec(path, hash)
				if err != nil {
					log.Panic(err)
				}
				return nil
			} 

			// check if hash is same as e if same, do nothing
			if bytes.Equal(hashResult, hash) {
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

	db := sqlx.MustOpen("sqlite3", "./test.db")
	defer db.Close()

	db.MustExec(SCHEMA)

	// instead of using sqlite for querying, what if we query all the files and store into a map?

	rows, err := db.Queryx("SELECT path, hash FROM files")
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var path string
		var hash []byte
		err = rows.Scan(&path, &hash)
		if err != nil {
			log.Fatal(err)
		}
		results[path] = hash
	}

	path := os.Args[1]
	wg.Add(1)
	go walk(path, &DB{db: db})
	wg.Wait()
}
