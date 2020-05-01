package main

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"bitbucket.org/zombiezen/cardcpx/natsort"
)

type itemMeta struct {
	Name        string    `json:"name"`
	CreateTime  time.Time `json:"create_time"`
	Favorite    bool      `json:"favorite"`
	FileIndices []int     `json:"file_indices"`
	Thumbnail   []byte    `json:"thumbnail"`
	mutex       sync.Mutex
}

func generateMetaFileName(name string) string {
	return filepath.Join(BaseDirectory, name+".meta")
}

func isMetaFileExist(name string) bool {
	metaFile := generateMetaFileName(name)

	if _, err := os.Stat(metaFile); err == nil {
		return true
	}

	return false
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

	if len(m.FileIndices) == 0 {
		return fmt.Errorf("file list is empty")
	}

	reader, err := r.File[m.FileIndices[0]].Open()
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

func (m *itemMeta) GenerateImageIndices() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	fullpath := BaseDirectory + string(os.PathSeparator) + m.Name

	r, err := zip.OpenReader(fullpath)
	if err != nil {
		return err
	}
	defer r.Close()

	type fileIndexPair struct {
		Index    int
		FileName string
	}

	var fileNames []fileIndexPair
	for i, f := range r.File {
		if filter(f.Name) {
			fileNames = append(fileNames,
				fileIndexPair{
					i, f.Name,
				})
		}
	}

	sort.Slice(fileNames, func(i, j int) bool {
		return natsort.Less(fileNames[i].FileName, fileNames[j].FileName)
	})

	m.FileIndices = make([]int, len(fileNames))
	for i, p := range fileNames {
		m.FileIndices[i] = p.Index
	}

	return nil
}
