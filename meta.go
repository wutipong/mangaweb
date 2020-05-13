package main

import (
	"archive/zip"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"bitbucket.org/zombiezen/cardcpx/natsort"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

const databaseURL = "user=user password=password host=db dbname=manga port=5432 sslmode=disable"
const databaseDriver = "postgres"
const maxAttempt = 10
const waitTime = time.Second * 10

type itemMeta struct {
	Name        string    `json:"name" db:"name"`
	CreateTime  time.Time `json:"create_time" db:"create_time"`
	Favorite    bool      `json:"favorite" db:"favorite"`
	FileIndices []int     `json:"file_indices" db:"file_indices"`
	Thumbnail   []byte    `json:"thumbnail" db:"thumbnail"`
	mutex       *sync.Mutex
}

func initDatabase() (dbx *sqlx.DB, err error) {
	for i := 0; i < maxAttempt; i++ {
		log.Printf("Connecting to Database, attempt #%v.", i)
		dbx, err = sqlx.Connect(databaseDriver, databaseURL)
		if err == nil {
			break
		}
		time.Sleep(waitTime)
	}

	if err != nil {
		return
	}

	_, err = dbx.Exec(`
		CREATE TABLE IF NOT EXISTS 
		manga_meta (
			name text, 
			create_time timestamp, 
			favorite boolean,
			file_indices integer[],
			thumbnail bytea);`)

	return
}

func generateMetaFileName(name string) string {
	return filepath.Join(BaseDirectory, name+".meta")
}

func isMetaFileExist(db *sqlx.DB, name string) bool {

	_, err := db.Query("SELECT * FROM manga_meta where name = $1", name)
	if err != nil {
		return true
	}

	return false
}

func (m *itemMeta) Write(db *sqlx.DB) error {
	db.NamedExec(`UPDATE manga_meta
		SET favorite = :favorite,
		WHERE name = :name
	`, m)
	return nil
}

func NewMeta(db *sqlx.DB, name string) itemMeta {
	meta := itemMeta{
		Name:       name,
		CreateTime: time.Now(),
		Favorite:   false,
		mutex:      new(sync.Mutex),
	}

	meta.GenerateImageIndices()
	meta.GenerateThumbnail()

	_, e := db.Exec(`INSERT INTO 
		manga_meta(name, create_time, favorite, file_indices, thumbnail) 
			VALUES($1, $2, $3, $4, $5)`, meta.Name, meta.CreateTime, meta.Favorite, pq.Array(meta.FileIndices), meta.Thumbnail)

	log.Printf("%v", e)
	return meta
}

func (m *itemMeta) Read(db *sqlx.DB, name string) error {

	row := db.DB.QueryRow("SELECT name, create_time, favorite, file_indices, thumbnail from manga_meta where name = $1", name)

	var x []sql.NullInt32
	err := row.Scan(&m.Name, &m.CreateTime, &m.Favorite, pq.Array(&x /*m.FileIndices*/), &m.Thumbnail)

	m.FileIndices = make([]int, len(x))
	for i := range x {
		m.FileIndices[i] = (int)(x[i].Int32)
	}
	return err
}

func ReadMeta(db *sqlx.DB, name string) (meta itemMeta, err error) {
	err = meta.Read(db, name)
	if errors.Is(err, sql.ErrNoRows) {
		meta = NewMeta(db, name)
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
