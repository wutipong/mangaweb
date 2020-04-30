package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type itemMeta struct {
	Name       string    `json:"name"`
	CreateTime time.Time `json:"create_time"`
	Favorite   bool      `json:"favorite"`
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

	data, _, err := OpenZipEntry(m.Name, 0)
	if err != nil {
		return err
	}
	reader := bytes.NewBuffer(data)

	thumbnail, err := CreateResized(reader, 200, 200)
	if err != nil {
		return err
	}

	m.Thumbnail = thumbnail

	return nil
}
