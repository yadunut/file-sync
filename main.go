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
	"time"
)

type Entry struct {
	Path    string
	Hash    []byte
	Deleted bool
}

var (
	wg      sync.WaitGroup
	counter atomic.Int64
	results map[string][]byte = make(map[string][]byte)
)

const (
	numWorkers = 100
)

func worker(db *DB, num int, paths chan string, cancel chan struct{}) error {
	defer func() {
		wg.Done()
		log.Printf("Worker %d finished\n", num)
	}()
	ticker := time.NewTicker(time.Second)
	numProcessed := 0
	log.Printf("Worker %d started\n", num)
	for {
		select {
		case <-cancel:
			return nil
		case <-ticker.C:
			log.Printf("Worker %d processed %d files\n", num, numProcessed)
		case path := <-paths:
			f, err := os.Open(path)
			if err != nil {
				log.Fatal(err)
			}
			h := sha256.New()
			_, err = io.Copy(h, f)
			if err != nil {
				log.Fatal(err)
			}
			f.Close()
			if err != nil {
				log.Fatal(err)
			}
			hash := h.Sum(nil)

			numProcessed++

			hashResult, ok := results[path]
			if !ok {
				// write to db
				db.mu.Lock()
				tx := db.MustBegin()
				insertStmt, err := tx.Preparex(INSERT_STMT)
				if err != nil {
					log.Fatal(err)
				}
				insertStmt.MustExec(path, hash)
				err = tx.Commit()
				if err != nil {
					log.Panic(err)
				}
				db.mu.Unlock()
				continue
			}

			// hashes match
			if bytes.Equal(hashResult, hash) {
				continue
			}
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("file-sync requires 1 argument")
		return
	}
	db := NewDB("./test.db")
	defer db.Close()
	log.Println("db created")

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

	log.Println("Cache Created")

	paths := make(chan string, 1000)
	cancel := make(chan struct{})
	path := os.Args[1]

	go func() {
		log.Println("Filewalker started")
		filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			select {
			case <-cancel:
				return filepath.SkipAll
			default:
			}
			if d.Type().IsRegular() {
				path, err = filepath.Abs(path)
				if err != nil {
					return err
				}
				paths <- path
			}
			return nil
		})
		close(cancel)
		log.Println("Filewalker finished")
	}()

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(db, i, paths, cancel)
	}
	wg.Wait()
}
