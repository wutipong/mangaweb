package tag

import (
	"github.com/julienschmidt/httprouter"
	"github.com/wutipong/mangaweb/handler"
	"github.com/wutipong/mangaweb/log"
	"go.uber.org/zap"
	"net/http"
	"path/filepath"
)

func ThumbnailHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	tagStr := handler.ParseParam(params, "tag")
	tagStr = filepath.FromSlash(tagStr)

	log.Get().Info("Tag thumbnail image", zap.String("tag", tagStr))

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
