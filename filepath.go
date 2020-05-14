package main

import (
	"os"
	"path/filepath"
	"strings"
)

//BaseDirectory the data directory
var BaseDirectory string
var filter func(path string) bool

func init() {
	filter = func(path string) bool {
		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".jpeg" {
			return true
		}
		if ext == ".jpg" {
			return true
		}
		if ext == ".png" {
			return true
		}
		return false
	}
}

// ListDir returns a list of content of a directory.
func ListDir(path string) (files []string, err error) {

	actualPath := filepath.Join(BaseDirectory)
	dir, err := os.Open(actualPath)
	if err != nil {
		return
	}
	children, err := dir.Readdir(0)
	if err != nil {
		return
	}

	for _, f := range children {
		if strings.HasPrefix(f.Name(), ".") {
			continue
		}

		name := filepath.Join(path, f.Name())

		if f.IsDir() {
			continue

			subFiles, e := ListDir(name)
			if e != nil {
				err = e
				return
			}
			files = append(files, subFiles...)
		}

		ext := strings.ToLower(filepath.Ext(f.Name()))

		if ext == ".zip" || ext == ".cbz" {

			files = append(files, name)
		}
	}
	return
}
