package view

import (
	"github.com/julienschmidt/httprouter"
	"github.com/wutipong/mangaweb/handler"
	"github.com/wutipong/mangaweb/log"
	"go.uber.org/zap"
	"net/http"
	"path/filepath"
	"strconv"
)

type setFavoriteResponse struct {
	Favorite bool `json:"favorite"`
}

func SetFavoriteHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	item := handler.ParseParam(params, "item")
	item = filepath.FromSlash(item)

	log.Get().Info("Set Favorite Item", zap.String("item_name", item))

	query := r.URL.Query()

	db, err := handler.CreateMetaProvider()
	if err != nil {
		handler.WriteJson(w, err)
		return
	}
	defer db.Close()

	m, err := db.Read(item)
	if err != nil {
		handler.WriteJson(w, err)
		return
	}

	if fav, e := strconv.ParseBool(query.Get("favorite")); e == nil {
		if fav != m.Favorite {
			m.Favorite = fav
			db.Write(m)
		}
	}

	response := setFavoriteResponse{
		Favorite: m.Favorite,
	}

	handler.WriteJson(w, response)
}
