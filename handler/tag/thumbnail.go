package tag

import (
	"github.com/labstack/echo/v4"
	"github.com/wutipong/mangaweb/handler"
	"net/http"
	"path/filepath"
)

func ThumbnailHandler(c echo.Context) error {
	tagStr := c.Param("*")
	tagStr = filepath.FromSlash(tagStr)

	provider, err := handler.CreateTagProvider()
	if err != nil {
		return err
	}
	defer provider.Close()

	m, err := provider.Read(tagStr)
	if err != nil {
		return err
	}

	return c.Blob(http.StatusOK, "image/jpeg", m.Thumbnail)
}
