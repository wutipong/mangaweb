package tag

import (
	"github.com/julienschmidt/httprouter"
	"github.com/wutipong/mangaweb/handler"
	"net/http"
	"path/filepath"
)

func ThumbnailHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	tagStr := handler.ParseParam(params, "tag")
	tagStr = filepath.FromSlash(tagStr)

	provider, err := handler.CreateTagProvider()
	if err != nil {
		handler.WriteError(w, err)
		return
	}
	defer provider.Close()

	m, err := provider.Read(tagStr)
	if err != nil {
		handler.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(m.Thumbnail)
	w.Header().Set("Content-Type", "image/jpeg")
}
