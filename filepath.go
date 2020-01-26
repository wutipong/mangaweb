package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
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

func ReadMeta(i item) (meta itemMeta, err error) {
	metaFile := filepath.Join(BaseDirectory, i.Name+".meta")
	f, err := os.OpenFile(metaFile, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return
	}

	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}
	if len(b) == 0 {
		meta = itemMeta{
			CreateTime: time.Now(),
			Favorite:   false,
			Name:       i.Name,
		}

		b, err = json.Marshal(meta)
		if err != nil {
			return
		}
		defer f.Write(b)
	} else {
		err = json.Unmarshal(b, &meta)
	}

	return
}
