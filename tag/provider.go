package tag

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wutipong/mangaweb/errors"
)

var pool *pgxpool.Pool = nil

func Init(ctx context.Context, p *pgxpool.Pool) {
	pool = p
}

func Close() error {
	return nil
}

func Delete(t Tag) error {
	return errors.ErrNotImplemented
}

func IsTagExist(name string) bool {
	r := pool.QueryRow(
		context.Background(),
		`select exists (select 1 from tags where name = $1)`,
		name,
	)

	exists := false
	r.Scan(&exists)

	return exists
}

func Read(name string) (t Tag, err error) {
	r := pool.QueryRow(
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
		&t.Thumbnail,
		&t.Version,
	)

	return
}

func ReadAll() (tags []Tag, err error) {
	rows, err := pool.Query(
		context.Background(),
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

func Write(t Tag) error {
	_, err := pool.Exec(
		context.Background(),
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
