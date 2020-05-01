package main

import (
	"archive/zip"
	"os"
)

type Page struct {
	Index int
	Name  string
}

func ListPages(file string) (pages []Page, err error) {
	var meta itemMeta
	err = meta.Read(file)
	if err != nil {
		return
	}

	if len(meta.FileIndices) == 0 {
		meta.GenerateImageIndices()
		meta.Write()
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
