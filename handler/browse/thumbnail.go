package browse

import (
	"database/sql"
	"errors"
	"github.com/julienschmidt/httprouter"
	"github.com/wutipong/mangaweb/handler"
	"github.com/wutipong/mangaweb/meta"
	"net/http"
	"path/filepath"
)

func ThumbnailHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	item := handler.ParseParam(params, "item")
	item = filepath.FromSlash(item)

	provider, err := handler.CreateMetaProvider()
	if err != nil {
		handler.WriteError(w, err)
		return
	}
	defer provider.Close()

	m, err := provider.Read(item)
	if errors.Is(err, sql.ErrNoRows) {
		m, _ = meta.NewItem(item)
		err = provider.Write(m)
		if err != nil {
			handler.WriteError(w, err)
			return
		}

	} else if err != nil {
		handler.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(m.Thumbnail)
	w.Header().Set("Content-Type", "image/jpeg")
}
