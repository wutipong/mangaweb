package main

import (
	"database/sql"
	"errors"
	"mangaweb/meta/mongo"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
)

func thumbnail(c echo.Context) error {
	name, err := url.PathUnescape(c.Param("*"))
	if err != nil {
		return err
	}

	provider, err := mongo.New()
	if err != nil {
		return err
	}
	defer provider.Close()

	m, err := provider.Read(name)
	if errors.Is(err, sql.ErrNoRows) {
		m, _ = provider.New(name)
	} else if err != nil {
		return err
	}

	return c.Blob(http.StatusOK, "image/jpeg", m.Thumbnail)
}
