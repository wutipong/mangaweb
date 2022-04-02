package view

import (
	"archive/zip"
	"path/filepath"

	"github.com/wutipong/mangaweb/meta"
)

type Page struct {
	Index int
	Name  string
}

func ListPages(m meta.Item) (pages []Page, err error) {
	if len(m.FileIndices) == 0 {
		return
	}

	fullpath := filepath.Join(meta.BaseDirectory, m.Name)
	r, err := zip.OpenReader(fullpath)
	if err != nil {
		return
	}

	defer r.Close()

	pages = make([]Page, len(m.FileIndices))
	for i, f := range m.FileIndices {
		pages[i] = Page{
			Name:  r.File[f].Name,
			Index: i,
		}
	}

	return
}
