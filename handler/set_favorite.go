package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type setFavoriteResponse struct {
	Favorite bool `json:"favorite"`
}

func SetFavoriteHandler(c echo.Context) error {
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

	if fav, e := strconv.ParseBool(c.QueryParam("favorite")); e == nil {
		if fav != m.Favorite {
			m.Favorite = fav
			db.Write(m)
		}
	}

	response := setFavoriteResponse{
		Favorite: m.Favorite,
	}
	return c.JSON(http.StatusOK, response)
}
