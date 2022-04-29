package handler

import (
	"github.com/julienschmidt/httprouter"
	"github.com/wutipong/mangaweb/log"
	"github.com/wutipong/mangaweb/scheduler"
	"net/http"
)

type RescanLibraryResponse struct {
	Result bool `json:"result"`
}

func RescanLibraryHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	log.Get().Info("Rescan library")

	scheduler.ScheduleScanLibrary()

	response := RescanLibraryResponse{
		Result: true,
	}

	WriteJson(w, response)
}
