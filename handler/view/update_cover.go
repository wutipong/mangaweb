package view

import (
	"github.com/wutipong/mangaweb/handler"
	_ "image/png"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/labstack/echo/v4"
)

// UpdateCover a handler to update the cover to specific image
func UpdateCover(c echo.Context) error {
	filename := c.Param("*")
	filename = filepath.FromSlash(filename)

	provider, err := handler.CreateMetaProvider()
	if err != nil {
		return err
	}
	defer provider.Close()

	var index = 0
	if i, err := strconv.Atoi(c.QueryParam("i")); err == nil {
		index = i
	}

	m, err := provider.Read(filename)
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
