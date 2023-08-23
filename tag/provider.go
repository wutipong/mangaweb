package tag

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wutipong/mangaweb/errors"
)

var pool *pgxpool.Pool = nil

func Init(p *pgxpool.Pool) {
	pool = p
}

func Delete(ctx context.Context, t Tag) error {
	return errors.ErrNotImplemented
}

func IsTagExist(ctx context.Context, name string) bool {
	r := pool.QueryRow(
		ctx,
		`select exists (select 1 from tags where name = $1)`,
		name,
	)

	exists := false
	r.Scan(&exists)

	return exists
}

func Read(ctx context.Context, name string) (t Tag, err error) {
	r := pool.QueryRow(
		ctx,
		`SELECT name, favorite, hidden, thumbnail, version
			FROM manga.tags
			where name = $1;`,
		name,
	)

	err = r.Scan(
		&t.Name,
		&t.Favorite,
		&t.Hidden,
		&t.Thumbnail,
		&t.Version,
	)

	return
}

func ReadAll(ctx context.Context) (tags []Tag, err error) {
	rows, err := pool.Query(
		ctx,
		`SELECT name, favorite, hidden, thumbnail, version
			FROM manga.tags;`,
	)

	if err != nil {
		return
	}

	for rows.Next() {
		var t Tag
		rows.Scan(
			&t.Name,
			&t.Favorite,
			&t.Hidden,
			&t.Thumbnail,
			&t.Version,
		)

		tags = append(tags, t)
	}

	return
}

func Write(ctx context.Context, t Tag) error {
	_, err := pool.Exec(
		ctx,
		`INSERT INTO manga.tags(name, favorite, hidden, thumbnail, version)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT(name) DO UPDATE
			SET favorite = $2, 
				hidden = $3,
				thumbnail = $4,
				version = $5;`,
		t.Name,
		t.Favorite,
		t.Hidden,
		t.Thumbnail,
		t.Version,
	)
	return err
}
