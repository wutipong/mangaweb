package view

import (
	"io"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/wutipong/mangaweb/handler"
	"github.com/wutipong/mangaweb/log"
	"github.com/wutipong/mangaweb/meta"
	"go.uber.org/zap"
)

func Download(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	item := handler.ParseParam(params, "item")
	item = filepath.FromSlash(item)

	log.Get().Info("Download", zap.String("item_name", item))

	m, err := meta.Read(r.Context(), item)
	if err != nil {
		handler.WriteError(w, err)
		return
	}

	reader, err := m.Open()
	if err != nil {
		handler.WriteError(w, err)
		return
	}
	defer reader.Close()

	bytes, err := io.ReadAll(reader)
	if err != nil {
		handler.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Length", strconv.Itoa(len(bytes)))
}
