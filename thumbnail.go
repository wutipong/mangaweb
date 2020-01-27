package main

import (
	"errors"
	"net/http"
	"net/url"
	"os"

	"github.com/labstack/echo/v4"
)

func thumbnail(c echo.Context) error {
	name, err := url.PathUnescape(c.Param("*"))
	if err != nil {
		return err
	}

	var m itemMeta
	err = m.Read(name)
	if errors.Is(err, os.ErrNotExist) {
	} else if err != nil {
		return err
	}
	if m.Thumbnail == nil {
		defer m.Write()
		m.GenerateThumbnail()
	}

	return c.Blob(http.StatusOK, "image/jpeg", m.Thumbnail)
}
