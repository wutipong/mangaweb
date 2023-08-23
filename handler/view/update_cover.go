package view

import (
	_ "image/png"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/wutipong/mangaweb/handler"
	"github.com/wutipong/mangaweb/log"
	"github.com/wutipong/mangaweb/meta"
	"go.uber.org/zap"
)

// UpdateCover a handler to update the cover to specific image
func UpdateCover(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	item := handler.ParseParam(params, "item")
	item = filepath.FromSlash(item)

	query := r.URL.Query()

	var index = 0
	if i, err := strconv.Atoi(query.Get("i")); err == nil {
		index = i
	}

	log.Get().Info("Update Cover", zap.String("item_name", item), zap.Int("index", index))

	m, err := meta.Read(r.Context(), item)
	if err != nil {
		handler.WriteJson(w, err)
		return
	}

	entryIndex := m.FileIndices[index]

	err = m.GenerateThumbnail(entryIndex)
	if err != nil {
		handler.WriteJson(w, err)
		return
	}

	err = meta.Write(r.Context(), m)
	if err != nil {
		handler.WriteJson(w, err)
		return
	}

	handler.WriteJson(w, "success")
}
