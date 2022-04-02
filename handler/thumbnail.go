package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
	"github.com/wutipong/mangaweb/meta"
)

func ThumbnailHandler(c echo.Context) error {
	name, err := url.PathUnescape(c.Param("*"))
	if err != nil {
		return err
	}

	provider, err := CreateMetaProvider()
	if err != nil {
		return err
	}
	defer provider.Close()

	m, err := provider.Read(name)
	if errors.Is(err, sql.ErrNoRows) {
		m, _ = meta.NewItem(name)
		err = provider.Write(m)
		if err != nil {
			return err
		}

	} else if err != nil {
		return err
	}

	return c.Blob(http.StatusOK, "image/jpeg", m.Thumbnail)
}
