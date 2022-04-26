package view

import (
	"github.com/labstack/echo/v4"
	"github.com/wutipong/mangaweb/handler"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

func Download(c echo.Context) error {
	filename := c.Param("*")
	filename = filepath.FromSlash(filename)

	db, err := handler.CreateMetaProvider()
	if err != nil {
		return err
	}
	defer db.Close()

	m, err := db.Read(filename)
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
