package tag

import (
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/wutipong/mangaweb/handler"
	"github.com/wutipong/mangaweb/log"
	"github.com/wutipong/mangaweb/tag"
	"go.uber.org/zap"
)

type setTagFavoriteResponse struct {
	Favorite bool `json:"favorite"`
}

func SetFavoriteHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	tagStr := handler.ParseParam(params, "tag")
	tagStr = filepath.FromSlash(tagStr)

	log.Get().Info("Set favorite tag", zap.String("tag", tagStr))

	query := r.URL.Query()

	m, err := tag.Read(r.Context(), tagStr)
	if err != nil {
		handler.WriteJson(w, err)
		return
	}

	if fav, e := strconv.ParseBool(query.Get("favorite")); e == nil {
		if fav != m.Favorite {
			m.Favorite = fav
			tag.Write(r.Context(), m)
		}
	}

	response := setTagFavoriteResponse{
		Favorite: m.Favorite,
	}

	handler.WriteJson(w, response)
}
