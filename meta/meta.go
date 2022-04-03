package meta

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"bitbucket.org/zombiezen/cardcpx/natsort"
	"github.com/wutipong/mangaweb/image"
)

// Item the meta data for each manga item.
// Do not change the field type nor names. Add new field when necessary.
// Also, when update the structure, if the new field is required, increment the CurrentItemVersion by one
// and create a migration function.
type Item struct {
	Name        string    `json:"name" db:"name" bson:"name"`
	CreateTime  time.Time `json:"create_time" db:"create_time" bson:"create_time"`
	Favorite    bool      `json:"favorite" db:"favorite" bson:"favorite"`
	FileIndices []int     `json:"file_indices" bson:"file_indices"`
	Thumbnail   []byte    `json:"thumbnail" db:"thumbnail" bson:"thumbnail"`
	IsRead      bool      `json:"is_read" db:"read" bson:"is_read"`
	Version     int       `json:"version" db:"version" bson:"version"`
}

type MetaProviderFactory func() (p Provider, err error)

//CurrentItemVersion the current version of `Item` structure.
const CurrentItemVersion = 0

func NewItem(name string) (i Item, err error) {

	createTime := time.Now()

	if stat, e := fs.Stat(os.DirFS(BaseDirectory), name); e == nil {
		createTime = stat.ModTime()
	}

	i = Item{
		Name:       name,
		CreateTime: createTime,
		Favorite:   false,
		Version:    CurrentItemVersion,
	}

	i.GenerateImageIndices()
	i.GenerateThumbnail(0)

	return
}

func (m *Item) Open() (reader io.ReadCloser, err error) {
	mutex := new(sync.Mutex)
	mutex.Lock()
	defer mutex.Unlock()

	fullpath := filepath.Join(BaseDirectory, m.Name)

	reader, err = os.Open(fullpath)
	return
}

func (m *Item) GenerateThumbnail(fileIndex int) error {
	mutex := new(sync.Mutex)
	mutex.Lock()
	defer mutex.Unlock()

	fullpath := filepath.Join(BaseDirectory, m.Name)

	r, err := zip.OpenReader(fullpath)
	if err != nil {
		return err
	}
	defer r.Close()

	if len(m.FileIndices) == 0 {
		return fmt.Errorf("file list is empty")
	}

	img, err := image.CreateCover(fileIndex, r)
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
	mutex := new(sync.Mutex)
	mutex.Lock()
	defer mutex.Unlock()

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
