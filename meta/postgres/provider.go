package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/wutipong/mangaweb/errors"
	"github.com/wutipong/mangaweb/meta"
)

type Provider struct {
	conn *pgx.Conn
}

func Init(ctx context.Context, conn *pgx.Conn) (p *Provider, err error) {
	p = new(Provider)
	p.conn = conn
	return
}

func (p *Provider) IsItemExist(name string) bool {
	r := p.conn.QueryRow(
		context.Background(),
		`select exists (select 1 from items where name = $1)`,
		name,
	)

	exists := false
	r.Scan(&exists)

	return exists
}
func (p *Provider) Write(i meta.Meta) error {
	_, err := p.conn.Exec(
		context.Background(),
		`INSERT INTO manga.items(name, create_time, favorite, file_indices, thumbnail, is_read, tags, version)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			ON CONFLICT(name) DO UPDATE
			SET create_time = $2, 
				favorite = $3, 
				file_indices = $4, 
				thumbnail = $5, 
				is_read = $6, 
				tags = $7, 
				version = $8;`,
		i.Name,
		i.CreateTime,
		i.Favorite,
		i.FileIndices,
		i.Thumbnail,
		i.IsRead,
		i.Tags,
		i.Version,
	)
	return err
}
func (p *Provider) Delete(i meta.Meta) error {
	return errors.ErrNotImplemented
}
func (p *Provider) Read(name string) (i meta.Meta, err error) {
	r := p.conn.QueryRow(
		context.Background(),
		`SELECT name, create_time, favorite, file_indices, thumbnail, is_read, tags, version
		FROM manga.items
		WHERE name = $1`,
		name,
	)

	err = r.Scan(
		&i.Name,
		&i.CreateTime,
		&i.Favorite,
		&i.FileIndices,
		&i.Thumbnail,
		&i.IsRead,
		&i.Tags,
		&i.Version,
	)

	return
}

func (p *Provider) ReadAll() (items []meta.Meta, err error) {
	rows, err := p.conn.Query(context.Background(),
		`SELECT name, create_time, favorite, file_indices, thumbnail, is_read, tags, version
		FROM manga.items;`)

	if err != nil {
		return
	}

	for rows.Next() {
		var i meta.Meta
		rows.Scan(
			&i.Name,
			&i.CreateTime,
			&i.Favorite,
			&i.FileIndices,
			&i.Thumbnail,
			&i.IsRead,
			&i.Tags,
			&i.Version)

		items = append(items, i)
	}

	return
}
func (p *Provider) Search(criteria []meta.SearchCriteria, sort meta.SortField, order meta.SortOrder, pageSize int, page int) (items []meta.Meta, err error) {

	return p.ReadAll()
}
func (p *Provider) Count(criteria []meta.SearchCriteria) (count int64, err error) {
	count = 0
	return
}
func (p *Provider) Close() error {
	// return p.conn.Close(context.Background())
	return nil
}
