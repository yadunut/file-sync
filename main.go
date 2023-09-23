package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
)

var wg sync.WaitGroup

func walk(dir string) {
	defer wg.Done()
	filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() && dir != path {
			wg.Add(1)
			go walk(path)
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
      h.Sum(nil)
			// fmt.Printf("%x\n", h.Sum(nil))
		}
		return nil
	})
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("file-sync requires 1 argument")
		return
	}
	path := os.Args[1]
	wg.Add(1)
	go walk(path)
	wg.Wait()
}
