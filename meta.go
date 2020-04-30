package main

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

type itemMeta struct {
	Name       string    `json:"name"`
	CreateTime time.Time `json:"create_time"`
	Favorite   bool      `json:"favorite"`
	Pages      []string  `json:"pages"`
	Thumbnail  []byte    `json:"thumbnail"`
	mutex      sync.Mutex
}

func generateMetaFileName(name string) string {
	return filepath.Join(BaseDirectory, name+".meta")
}

func (m *itemMeta) Write() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	metaFile := generateMetaFileName(m.Name)
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}

	f, err := os.Create(metaFile)
	if err != nil {
		return err
	}

	defer f.Close()
	f.Write(b)

	return nil
}

func NewMeta(name string) itemMeta {
	return itemMeta{
		Name:       name,
		CreateTime: time.Now(),
		Favorite:   false,
	}
}

func (m *itemMeta) Read(name string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	metaFile := generateMetaFileName(name)
	f, err := os.Open(metaFile)
	if err != nil {
		return err
	}

	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, &m)
}

func ReadMeta(name string) (meta itemMeta, err error) {
	err = meta.Read(name)
	if errors.Is(err, os.ErrNotExist) {
		meta = NewMeta(name)
		err = meta.Write()
	}

	return
}

func (m *itemMeta) GenerateThumbnail() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	fullpath := BaseDirectory + string(os.PathSeparator) + m.Name

	r, err := zip.OpenReader(fullpath)
	if err != nil {
		return err
	}
	defer r.Close()

	if len(m.Pages) == 0 {
		m.GeneratePages()
	}

	var reader io.ReadCloser
	for _, zf := range r.File {
		if zf.Name == m.Pages[0] {
			reader, err = zf.Open()
			break
		}
	}

	if err != nil {
		return err
	}

	defer reader.Close()

	thumbnail, err := CreateResized(reader, 200, 200)
	if err != nil {
		return err
	}

	m.Thumbnail = thumbnail

	return nil
}

func (m *itemMeta) GeneratePages() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	fullpath := BaseDirectory + string(os.PathSeparator) + m.Name

	r, err := zip.OpenReader(fullpath)
	if err != nil {
		return err
	}
	defer r.Close()

	var fileNames []string
	for _, f := range r.File {
		if filter(f.Name) {
			fileNames = append(fileNames, f.Name)
		}
	}

	sort.Strings(fileNames)

	m.Pages = fileNames

	return nil
}
