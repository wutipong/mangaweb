package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/wutipong/mangaweb/scheduler"
	"net/http"
)

type RescanLibraryResponse struct {
	Result bool `json:"result"`
}

func RescanLibraryHandler(c echo.Context) error {
	scheduler.ScheduleScanLibrary()

	response := RescanLibraryResponse{
		Result: true,
	}
	return c.JSON(http.StatusOK, response)
}
