package handler

import (
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
)

func Download(c echo.Context) error {
	filename := c.Param("*")

	db, err := CreateMetaProvider()
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
