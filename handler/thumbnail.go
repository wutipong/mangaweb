package handler

import (
	"database/sql"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/wutipong/mangaweb/meta"
	"net/http"
	"path/filepath"
)

func ThumbnailHandler(c echo.Context) error {
	filename := c.Param("*")
	filename = filepath.FromSlash(filename)

	provider, err := CreateMetaProvider()
	if err != nil {
		return err
	}
	defer provider.Close()

	m, err := provider.Read(filename)
	if errors.Is(err, sql.ErrNoRows) {
		m, _ = meta.NewItem(filename)
		err = provider.Write(m)
		if err != nil {
			return err
		}

	} else if err != nil {
		return err
	}

	return c.Blob(http.StatusOK, "image/jpeg", m.Thumbnail)
}
