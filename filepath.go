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
func ListDir() (files []string, err error) {

	dir, err := os.Open(BaseDirectory)
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

		if f.IsDir() {
			continue
		}

		ext := strings.ToLower(filepath.Ext(f.Name()))

		if ext == ".zip" || ext == ".cbz" {
			files = append(files, f.Name())
		}
	}
	return
}
