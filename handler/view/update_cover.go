package view

import (
	"github.com/julienschmidt/httprouter"
	"github.com/wutipong/mangaweb/handler"
	"github.com/wutipong/mangaweb/log"
	"go.uber.org/zap"
	_ "image/png"
	"net/http"
	"path/filepath"
	"strconv"
)

// UpdateCover a handler to update the cover to specific image
func UpdateCover(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	item := handler.ParseParam(params, "item")
	item = filepath.FromSlash(item)

	query := r.URL.Query()

	provider, err := handler.CreateMetaProvider()
	if err != nil {
		handler.WriteError(w, err)
		return
	}
	defer provider.Close()

	var index = 0
	if i, err := strconv.Atoi(query.Get("i")); err == nil {
		index = i
	}

	log.Get().Info("Update Cover", zap.String("item_name", item), zap.Int("index", index))

	m, err := provider.Read(item)
	if err != nil {
		handler.WriteError(w, err)
		return
	}

	entryIndex := m.FileIndices[index]
	err = m.GenerateThumbnail(entryIndex)

	if err != nil {
		handler.WriteError(w, err)
		return
	}

	err = provider.Write(m)
	if err != nil {
		handler.WriteError(w, err)
		return
	}

	handler.WriteJson(w, "success")
}
