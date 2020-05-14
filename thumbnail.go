package main

import (
	"database/sql"
	"errors"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
)

func thumbnail(c echo.Context) error {
	name, err := url.PathUnescape(c.Param("*"))
	if err != nil {
		return err
	}

	db, err := connectDB()
	if err != nil {
		return err
	}
	defer db.Close()

	var m itemMeta
	err = m.Read(db, name)
	if errors.Is(err, sql.ErrNoRows) {
		m, _ = NewMeta(db, name)
	} else if err != nil {
		return err
	}

	return c.Blob(http.StatusOK, "image/jpeg", m.Thumbnail)
}
