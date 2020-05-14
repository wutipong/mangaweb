package main

import (
	"archive/zip"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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
	Name           string          `json:"name" db:"name"`
	CreateTime     time.Time       `json:"create_time" db:"create_time"`
	Favorite       bool            `json:"favorite" db:"favorite"`
	FileIndices    []int           `json:"file_indices"`
	FileIndicesSQL []sql.NullInt32 `db:"file_indices"`
	Thumbnail      []byte          `json:"thumbnail" db:"thumbnail"`
	mutex          *sync.Mutex
}

func connectDB() (dbx *sqlx.DB, err error) {
	dbx, err = sqlx.Connect(databaseDriver, databaseURL)
	return
}

func initDatabase() error {
	var dbx *sqlx.DB
	var err error

	for i := 0; i < maxAttempt; i++ {
		log.Printf("Connecting to Database, attempt #%v.", i)
		dbx, err = connectDB()
		if err == nil {
			break
		}
		time.Sleep(waitTime)
	}

	if err != nil {
		return err
	}

	defer dbx.Close()

	_, err = dbx.Exec(`
		CREATE TABLE IF NOT EXISTS 
		manga_meta (
			name text PRIMARY KEY, 
			create_time timestamp, 
			favorite boolean,
			file_indices integer[],
			thumbnail bytea);`)

	if err != nil {
		return err
	}

	_, err = dbx.Exec(`
		CREATE INDEX IF NOT EXISTS manga_meta_name
		ON manga_meta (name);`)

	if err != nil {
		return err
	}

	_, err = dbx.Exec(`
		CREATE INDEX IF NOT EXISTS manga_meta_name_favorite
		ON manga_meta (name, favorite);`)

	if err != nil {
		return err
	}

	return nil
}

func generateMetaFileName(name string) string {
	return filepath.Join(BaseDirectory, name+".meta")
}

func migrateMeta(db *sqlx.DB) error {
	files, err := ListDir()
	if err != nil {
		return err
	}

	for _, file := range files {
		log.Printf("Processing %s", file)
		metaFile := generateMetaFileName(file)
		f, err := os.Open(metaFile)
		if err != nil {
			continue
		}
		b, err := ioutil.ReadAll(f)

		if err != nil {
			continue
		}

		var m itemMeta
		json.Unmarshal(b, &m)

		m.Write(db)
	}

	return nil
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
			create_time = :create_time
		WHERE name = :name
	`, m)
	return nil
}

func NewMeta(db *sqlx.DB, name string) (meta itemMeta, err error) {
	meta = itemMeta{
		Name:       name,
		CreateTime: time.Now(),
		Favorite:   false,
		mutex:      new(sync.Mutex),
	}

	meta.generateImageIndices()
	meta.generateThumbnail()

	_, err = db.Exec(`INSERT INTO 
		manga_meta(name, create_time, favorite, file_indices, thumbnail) 
			VALUES($1, $2, $3, $4, $5)`, meta.Name, meta.CreateTime, meta.Favorite, pq.Array(meta.FileIndices), meta.Thumbnail)

	return
}

func (m *itemMeta) Read(db *sqlx.DB, name string) error {

	row := db.QueryRow("SELECT name, create_time, favorite, file_indices, thumbnail from manga_meta where name = $1", name)

	var x []sql.NullInt32
	err := row.Scan(&m.Name, &m.CreateTime, &m.Favorite, pq.Array(&x), &m.Thumbnail)

	m.FileIndices = make([]int, len(x))
	for i := range x {
		m.FileIndices[i] = (int)(x[i].Int32)
	}

	return err
}

func (m *itemMeta) updateIndex() {
	m.FileIndices = make([]int, len(m.FileIndicesSQL))
	for i := range m.FileIndicesSQL {
		m.FileIndices[i] = (int)(m.FileIndicesSQL[i].Int32)
	}
}

func OpenMeta(db *sqlx.DB, name string) (meta itemMeta, err error) {
	err = meta.Read(db, name)
	if errors.Is(err, sql.ErrNoRows) {
		meta, err = NewMeta(db, name)
	}

	return
}

func ReadAllMeta(db *sqlx.DB) (meta []itemMeta, err error) {
	rows, err := db.Query(
		`SELECT name, create_time, favorite, file_indices, thumbnail 
			FROM manga_meta
			ORDER BY name;`)

	if err != nil {
		return
	}

	for rows.Next() {
		var m itemMeta
		var x []sql.NullInt32

		err = rows.Scan(&m.Name, &m.CreateTime, &m.Favorite, pq.Array(&x), &m.Thumbnail)
		if err != nil {
			continue
		}
		m.FileIndices = make([]int, len(x))
		for i := range x {
			m.FileIndices[i] = (int)(x[i].Int32)
		}

		meta = append(meta, m)
	}

	return
}

func (m *itemMeta) generateThumbnail() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	fullpath := filepath.Join(BaseDirectory, m.Name)

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

func (m *itemMeta) generateImageIndices() error {
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
