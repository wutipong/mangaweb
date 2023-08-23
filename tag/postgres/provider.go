package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/wutipong/mangaweb/errors"
	"github.com/wutipong/mangaweb/tag"
)

type Provider struct {
	conn *pgx.Conn
}

func Init(ctx context.Context, conn *pgx.Conn) (p *Provider, err error) {
	p = new(Provider)
	p.conn = conn
	return
}

func (p *Provider) Close() error {
	return nil
}
func (p *Provider) Delete(t tag.Tag) error {
	return errors.ErrNotImplemented

}
func (p *Provider) IsTagExist(name string) bool {
	r := p.conn.QueryRow(
		context.Background(),
		`select exists (select 1 from tags where name = $1)`,
		name,
	)

	exists := false
	r.Scan(&exists)

	return exists
}
func (p *Provider) Read(name string) (t tag.Tag, err error) {
	r := p.conn.QueryRow(
		context.Background(),
		`SELECT name, favorite, hidden, thumbnail, version
			FROM manga.tags
			where name = $1;`,
		name,
	)

	err = r.Scan(
		&t.Name,
		&t.Favorite,
		&t.Hidden,
		&t.Version,
	)

	return
}
func (p *Provider) ReadAll() (tags []tag.Tag, err error) {
	rows, err := p.conn.Query(
		context.Background(),
		`SELECT name, favorite, hidden, thumbnail, version
			FROM manga.tags;`,
	)

	if err != nil {
		return
	}

	for {
		var t tag.Tag
		rows.Scan(
			&t.Name,
			&t.Favorite,
			&t.Hidden,
			&t.Version,
		)

		tags = append(tags, t)

		if !rows.Next() {
			break
		}
	}

	return

}
func (p *Provider) Write(t tag.Tag) error {
	_, err := p.conn.Exec(
		context.Background(),
		`INSERT INTO manga.tags(name, favorite, hidden, thumbnail, version)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT(name) DO UPDATE
			SET favorite = $2, 
				hidden = $3,
				thumbnail = $4
				version = $5;`,
		t.Name,
		t.Favorite,
		t.Hidden,
		t.Thumbnail,
		t.Version,
	)
	return err
}
