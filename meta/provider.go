package meta

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wutipong/mangaweb/errors"
)

var pool *pgxpool.Pool

func Init(p *pgxpool.Pool) {
	pool = p
}

func IsItemExist(ctx context.Context, name string) bool {
	r := pool.QueryRow(
		ctx,
		`select exists (select 1 from items where name = $1)`,
		name,
	)

	exists := false
	r.Scan(&exists)

	return exists
}
func Write(ctx context.Context, i Meta) error {
	_, err := pool.Exec(
		ctx,
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
func Delete(ctx context.Context, i Meta) error {
	return errors.ErrNotImplemented
}

func Read(ctx context.Context, name string) (i Meta, err error) {
	r := pool.QueryRow(
		ctx,
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

func ReadAll(ctx context.Context) (items []Meta, err error) {
	rows, err := pool.Query(ctx,
		`SELECT name, create_time, favorite, file_indices, thumbnail, is_read, tags, version
		FROM manga.items;`)

	if err != nil {
		return
	}

	for rows.Next() {
		var i Meta
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
func Search(ctx context.Context, criteria []SearchCriteria, sort SortField, order SortOrder, pageSize int, page int) (items []Meta, err error) {
	// TODO: sanitize the query.
	criteriaStr := make([]string, 0)
	for _, c := range criteria {
		switch c.Field {
		case SearchFieldName:
			criteriaStr = append(criteriaStr, fmt.Sprintf(`items.name LIKE '%%%s%%'`, c.Value.(string)))

		case SearchFieldFavorite:
			criteriaStr = append(criteriaStr, fmt.Sprintf(`items.favorite = %v`, c.Value.(bool)))

		case SearchFieldTag:
			criteriaStr = append(criteriaStr, fmt.Sprintf(`items.tags @> ARRAY['%s']`, c.Value.(string)))
		}
	}

	where := ""
	if len(criteria) > 0 {
		where = fmt.Sprintf(`WHERE %s `, strings.Join(criteriaStr, " AND "))
	}

	sortBy := ""
	switch sort {
	case SortFieldName:
		sortBy = `ORDER BY name`
	case SortFieldCreateTime:
		sortBy = `ORDER BY create_time`
	}

	sortOrder := ""
	switch order {
	case SortOrderAscending:
		sortOrder = "ASC"
	case SortOrderDescending:
		sortOrder = "DESC"

	}

	query := fmt.Sprintf(
		`SELECT name, create_time, favorite, file_indices, thumbnail, is_read, tags, version
		FROM manga.items
		%s 
		%s %s
		LIMIT %d OFFSET %d;`,
		where,
		sortBy,
		sortOrder,
		pageSize,
		pageSize*page,
	)

	rows, err := pool.Query(ctx, query)

	if err != nil {
		return
	}

	for rows.Next() {
		var i Meta
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
func Count(ctx context.Context, criteria []SearchCriteria) (count int64, err error) {
	// TODO: sanitize the query.
	criteriaStr := make([]string, 0)
	for _, c := range criteria {
		switch c.Field {
		case SearchFieldName:
			criteriaStr = append(criteriaStr, fmt.Sprintf(`items.name LIKE '%%%s%%'`, c.Value.(string)))

		case SearchFieldFavorite:
			criteriaStr = append(criteriaStr, fmt.Sprintf(`items.favorite = %v`, c.Value.(bool)))

		case SearchFieldTag:
			criteriaStr = append(criteriaStr, fmt.Sprintf(`items.tags @> ARRAY['%s']`, c.Value.(string)))
		}
	}

	where := ""
	if len(criteria) > 0 {
		where = fmt.Sprintf(`WHERE %s `, strings.Join(criteriaStr, " AND "))
	}

	r := pool.QueryRow(
		ctx,
		fmt.Sprintf(`select count (*) from manga.items %s;`, where),
	)

	r.Scan(&count)
	return
}
