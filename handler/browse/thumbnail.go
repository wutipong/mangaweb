package browse

import (
	"database/sql"
	"errors"
	"net/http"
	"path/filepath"

	"github.com/julienschmidt/httprouter"
	"github.com/wutipong/mangaweb/handler"
	"github.com/wutipong/mangaweb/log"
	"github.com/wutipong/mangaweb/meta"
	"go.uber.org/zap"
)

func ThumbnailHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	item := handler.ParseParam(params, "item")
	item = filepath.FromSlash(item)

	log.Get().Info("Item Thumbnail", zap.String("item_name", item))

	m, err := meta.Read(item)
	if errors.Is(err, sql.ErrNoRows) {
		m, _ = meta.NewItem(item)
		err = meta.Write(m)
		if err != nil {
			handler.WriteError(w, err)
			return
		}

	} else if err != nil {
		handler.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "image/jpeg")
	w.Write(m.Thumbnail)
}
