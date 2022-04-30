package handler

import (
	"encoding/json"
	"github.com/wutipong/mangaweb/errors"
	"net/http"
)

func WriteJson(w http.ResponseWriter, v any) {
	if err, ok := v.(error); ok {
		w.WriteHeader(http.StatusInternalServerError)
		if _, ok := err.(errors.Error); !ok {
			v = errors.ErrUnknown.Wrap(err)
		}

	} else {
		w.WriteHeader(http.StatusOK)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	b, _ := json.Marshal(v)
	w.Write(b)
}

func WriteHtml(w http.ResponseWriter, content string) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(content))
}
