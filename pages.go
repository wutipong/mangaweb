package main

import (
	"archive/zip"
	"os"
	"path/filepath"
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

	for i, f := range r.File {
		if filter(f.Name) {
			pages = append(pages, Page{
				Name:  filepath.Base(f.Name),
				Index: i,
			})
		}
	}
	return
}
