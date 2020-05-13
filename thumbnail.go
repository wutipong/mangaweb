package main

import (
	"database/sql"
	"errors"
	"net/http"
	"net/url"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

func thumbnail(c echo.Context) error {
	name, err := url.PathUnescape(c.Param("*"))
	if err != nil {
		return err
	}

	db := c.Get("db").(*sqlx.DB)

	var m itemMeta
	err = m.Read(db, name)
	if errors.Is(err, sql.ErrNoRows) {
		m, _ = NewMeta(db, name)
	} else if err != nil {
		return err
	}

	return c.Blob(http.StatusOK, "image/jpeg", m.Thumbnail)
}
