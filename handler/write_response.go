package handler

import (
	"encoding/json"
	"github.com/labstack/gommon/log"
	"net/http"
)

func WriteJson(w http.ResponseWriter, v any) {
	if _, ok := v.(error); ok {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	b, _ := json.Marshal(v)
	w.Write(b)
}

func WriteError(w http.ResponseWriter, err error) {
	log.Error(err)
	WriteJson(w, err)
}

func WriteHtml(w http.ResponseWriter, content string) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(content))
}
