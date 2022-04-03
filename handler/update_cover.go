package handler

import (
	_ "image/png"
	"net/http"
	"net/url"
	"strconv"

	"github.com/labstack/echo/v4"
)

// UpdateCover a handler to update the cover to specific image
func UpdateCover(c echo.Context) error {
	p, err := url.PathUnescape(c.Param("*"))
	if err != nil {
		return err
	}

	provider, err := CreateMetaProvider()
	if err != nil {
		return err
	}
	defer provider.Close()

	var index = 0
	if i, err := strconv.Atoi(c.QueryParam("i")); err == nil {
		index = i
	}

	m, err := provider.Read(p)
	if err != nil {
		return err
	}

	entryIndex := m.FileIndices[index]
	err = m.GenerateThumbnail(entryIndex)

	if err != nil {
		return err
	}

	err = provider.Write(m)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, "success")
}
