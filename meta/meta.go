package meta

import (
	"archive/zip"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/wutipong/mangaweb/image"

	"bitbucket.org/zombiezen/cardcpx/natsort"
)

type Item struct {
	Name        string      `json:"name" db:"name" bson:"name"`
	CreateTime  time.Time   `json:"create_time" db:"create_time" bson:"create_time"`
	Favorite    bool        `json:"favorite" db:"favorite" bson:"favorite"`
	FileIndices []int       `json:"file_indices" bson:"file_indices"`
	Thumbnail   []byte      `json:"thumbnail" db:"thumbnail" bson:"thumbnail"`
	IsRead      bool        `json:"is_read" db:"read" bson:"is_read"`
	Mutex       *sync.Mutex `json:"-" db:"-" bson:"-"`
}

func (m *Item) GenerateThumbnail() error {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	fullpath := filepath.Join(BaseDirectory, m.Name)

	r, err := zip.OpenReader(fullpath)
	if err != nil {
		return err
	}
	defer r.Close()

	if len(m.FileIndices) == 0 {
		return fmt.Errorf("file list is empty")
	}

	img, err := image.CreateCover(m.FileIndices, r)
	if err != nil {
		return err
	}

	resized := image.CreateThumbnail(img)
	jpeg, err := image.ToJPEG(resized)
	if err != nil {
		return err
	}

	m.Thumbnail = jpeg

	return nil
}

func (m *Item) GenerateImageIndices() error {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

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
