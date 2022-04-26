package tag

import (
	"github.com/wutipong/mangaweb/handler"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/labstack/echo/v4"
)

type setTagFavoriteResponse struct {
	Favorite bool `json:"favorite"`
}

func SetFavoriteHandler(c echo.Context) error {
	filename := c.Param("*")
	filename = filepath.FromSlash(filename)

	db, err := handler.CreateTagProvider()
	if err != nil {
		return err
	}
	defer db.Close()

	m, err := db.Read(filename)
	if err != nil {
		return err
	}

	if fav, e := strconv.ParseBool(c.QueryParam("favorite")); e == nil {
		if fav != m.Favorite {
			m.Favorite = fav
			db.Write(m)
		}
	}

	response := setTagFavoriteResponse{
		Favorite: m.Favorite,
	}
	return c.JSON(http.StatusOK, response)
}
