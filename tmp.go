package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
)

type Entry struct {
	Path string
	Hash []byte
}

var (
	wg      sync.WaitGroup
	counter atomic.Int64
	results map[string][]byte = make(map[string][]byte)
)

const (
	numWorkers = 5000
	pathsSize  = 1000
)

func fileWalker(path string, cancel <-chan struct{}) <-chan string {
	paths := make(chan string, pathsSize)
	go func() {
		err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
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
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Paths closed")
		close(paths)
	}()
	return paths
}

func hashFiles(paths <-chan string, results chan<- Entry, cancel <-chan struct{}) error {
	defer wg.Done()
	for {
		select {
		case <-cancel:
			return nil
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
			log.Println("Pushing to results")
			results <- Entry{path, hash}
		default:
			if paths == nil {
				close(results)
				return nil
			}
		}
	}
}

func Main() {
	if len(os.Args) < 2 {
		fmt.Println("file-sync requires 1 argument")
		return
	}
	path := os.Args[1]

	// db := NewDB("./test.db")
	// defer db.Close()
	log.Println("db created")

	cancel := make(chan struct{})
	paths := fileWalker(path, cancel)
	hashes := make(chan Entry, pathsSize)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go hashFiles(paths, hashes, cancel)
	}

	for hash := range hashes {
		fmt.Println(hash.Path, hash.Hash)
	}
	wg.Wait()
}
