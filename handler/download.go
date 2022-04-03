package handler

import (
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
)

func Download(c echo.Context) error {
	p, err := url.PathUnescape(c.Param("*"))
	if err != nil {
		return err
	}

	db, err := CreateMetaProvider()
	if err != nil {
		return err
	}
	defer db.Close()

	m, err := db.Read(p)
	if err != nil {
		return err
	}

	reader, err := m.Open()
	if err != nil {
		return err
	}
	defer reader.Close()

	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	return c.Blob(http.StatusOK, "application/zip", bytes)
}
