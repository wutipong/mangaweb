package main

import (
	"archive/zip"
	"log"
	"os"
	"path/filepath"
	"sort"
)

type Page struct {
	Index int
	Name  string
}

func ListPages(file string) (pages []Page, err error) {
	fullpath := BaseDirectory + string(os.PathSeparator) + file

	r, err := zip.OpenReader(fullpath)
	if err != nil {
		return
	}
	defer r.Close()

	var fileNames []string
	for _, f := range r.File {
		if filter(f.Name) {
			fileNames = append(fileNames, filepath.Base(f.Name))
		}
	}

	sort.Strings(fileNames)

	for _, f := range fileNames {
		log.Println(f)
	}

	pages = make([]Page, len(fileNames))
	for i, f := range fileNames {
		pages[i] = Page{
			Name:  f,
			Index: i,
		}
	}

	return
}
