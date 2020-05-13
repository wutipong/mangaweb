package main

import (
	"archive/zip"
	"os"

	"github.com/jmoiron/sqlx"
)

type Page struct {
	Index int
	Name  string
}

func ListPages(db *sqlx.DB, file string) (pages []Page, err error) {
	var meta itemMeta
	err = meta.Read(db, file)
	if err != nil {
		return
	}

	if len(meta.FileIndices) == 0 {
		return
	}

	fullpath := BaseDirectory + string(os.PathSeparator) + file
	r, err := zip.OpenReader(fullpath)
	if err != nil {
		return
	}

	defer r.Close()

	pages = make([]Page, len(meta.FileIndices))
	for i, f := range meta.FileIndices {
		pages[i] = Page{
			Name:  r.File[f].Name,
			Index: i,
		}
	}

	return
}
