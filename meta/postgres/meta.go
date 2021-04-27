package postgres

import (
	"database/sql"
	"errors"
	"log"
	"mangaweb/meta"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

var databaseURL string

const databaseDriver = "postgres"
const maxAttempt = 10
const waitTime = time.Second * 10

func connectDB() (dbx *sqlx.DB, err error) {
	dbx, err = sqlx.Connect(databaseDriver, databaseURL)
	return
}

func Init(dbAddress string) error {
	var dbx *sqlx.DB
	var err error

	databaseURL = dbAddress
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
			thumbnail bytea,
			read boolean NOT NULL DEFAULT false);`)

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

func isItemExist(db *sqlx.DB, name string) bool {
	_, err := db.Query("SELECT * FROM manga_meta where name = $1", name)
	return err != nil
}

func write(m meta.Item, db *sqlx.DB) error {
	db.NamedExec(`UPDATE manga_meta
		SET favorite = :favorite,
			create_time = :create_time,
			read = :read
		WHERE name = :name
	`, m)
	return nil
}

func newItem(db *sqlx.DB, name string) (m meta.Item, err error) {
	m = meta.Item{
		Name:       name,
		CreateTime: time.Now(),
		Favorite:   false,
		Mutex:      new(sync.Mutex),
	}

	m.GenerateImageIndices()
	m.GenerateThumbnail()

	_, err = db.Exec(`INSERT INTO 
		manga_meta(name, create_time, favorite, file_indices, thumbnail) 
			VALUES($1, $2, $3, $4, $5)`, m.Name, m.CreateTime, m.Favorite, pq.Array(m.FileIndices), m.Thumbnail)

	return
}

func deleteItem(db *sqlx.DB, m meta.Item) error {
	_, err := db.NamedExec(`DELETE from manga_meta
		WHERE name = :name;`, m)
	return err
}

func readItem(db *sqlx.DB, name string) (m meta.Item, err error) {
	row := db.QueryRow("SELECT name, create_time, favorite, file_indices, thumbnail, read from manga_meta where name = $1", name)

	var x []sql.NullInt32
	err = row.Scan(&m.Name, &m.CreateTime, &m.Favorite, pq.Array(&x), &m.Thumbnail, &m.IsRead)

	m.FileIndices = make([]int, len(x))
	for i := range x {
		m.FileIndices[i] = (int)(x[i].Int32)
	}

	return
}

func openItem(db *sqlx.DB, name string) (m meta.Item, err error) {
	m, err = readItem(db, name)
	if errors.Is(err, sql.ErrNoRows) {
		m, err = newItem(db, name)
	}

	return
}

func readAllItems(db *sqlx.DB) (items []meta.Item, err error) {
	rows, err := db.Query(
		`SELECT name, create_time, favorite, file_indices, thumbnail, read 
			FROM manga_meta
			ORDER BY name;`)

	if err != nil {
		return
	}

	for rows.Next() {
		var m meta.Item
		var x []sql.NullInt32

		err = rows.Scan(&m.Name, &m.CreateTime, &m.Favorite, pq.Array(&x), &m.Thumbnail, &m.IsRead)
		if err != nil {
			continue
		}
		m.FileIndices = make([]int, len(x))
		for i := range x {
			m.FileIndices[i] = (int)(x[i].Int32)
		}

		items = append(items, m)
	}

	return
}
