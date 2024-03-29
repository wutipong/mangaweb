package meta

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/disintegration/imaging"
	"github.com/facette/natsort"
	"github.com/wutipong/mangaweb/tag"

	_ "golang.org/x/image/webp"
)

// Meta the metadata for each manga item.
// Do not change the field type nor names. Add new field when necessary.
// Also, when update the structure, if the new field is required, increment the CurrentVersion by one
// and create a migration function.
type Meta struct {
	Name        string    `json:"name" db:"name" bson:"name"`
	CreateTime  time.Time `json:"create_time" db:"create_time" bson:"create_time"`
	Favorite    bool      `json:"favorite" db:"favorite" bson:"favorite"`
	FileIndices []int     `json:"file_indices" bson:"file_indices"`
	Thumbnail   []byte    `json:"thumbnail" db:"thumbnail" bson:"thumbnail"`
	IsRead      bool      `json:"is_read" db:"read" bson:"is_read"`
	Tags        []string  `json:"tags" db:"tags" bson:"tags"`

	Version int `json:"version" db:"version" bson:"version"`
}

type ProviderFactory func() (p Provider, err error)

// CurrentVersion the current version of `Meta` structure.
const CurrentVersion = 1

func NewItem(name string) (i Meta, err error) {
	createTime := time.Now()

	if stat, e := fs.Stat(os.DirFS(BaseDirectory), name); e == nil {
		createTime = stat.ModTime()
	}

	i = Meta{
		Name:       name,
		CreateTime: createTime,
		Favorite:   false,
		Version:    CurrentVersion,
	}

	i.GenerateImageIndices()
	i.GenerateThumbnail(0)
	i.PopulateTags()

	return
}

func (m *Meta) Open() (reader io.ReadCloser, err error) {
	mutex := new(sync.Mutex)
	mutex.Lock()
	defer mutex.Unlock()

	fullpath := filepath.Join(BaseDirectory, m.Name)

	reader, err = os.Open(fullpath)
	return
}

func (m *Meta) GenerateThumbnail(fileIndex int) error {
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

	stream, err := r.File[m.FileIndices[fileIndex]].Open()
	if err != nil {
		return err
	}

	defer stream.Close()

	img, err := imaging.Decode(stream, imaging.AutoOrientation(true))
	if err != nil {
		return err
	}

	const thumbnailSize = 200
	if img.Bounds().Dx() > thumbnailSize || img.Bounds().Dy() > thumbnailSize {
		resized := imaging.Fit(img, thumbnailSize, thumbnailSize, imaging.MitchellNetravali)
		img = resized
	}

	buffer := bytes.Buffer{}
	imaging.Encode(&buffer, img, imaging.JPEG, imaging.JPEGQuality(75))

	m.Thumbnail = buffer.Bytes()

	return nil
}

func (m *Meta) GenerateImageIndices() error {
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
		return natsort.Compare(fileNames[i].FileName, fileNames[j].FileName)
	})

	m.FileIndices = make([]int, len(fileNames))
	for i, p := range fileNames {
		m.FileIndices[i] = p.Index
	}

	return nil
}

func (m *Meta) PopulateTags() {
	m.Tags = tag.ParseTag(m.Name)
}
