package tag

import (
	"github.com/julienschmidt/httprouter"
	"github.com/wutipong/mangaweb/handler"
	"net/http"
	"path/filepath"
	"strconv"
)

type setTagFavoriteResponse struct {
	Favorite bool `json:"favorite"`
}

func SetFavoriteHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	tag := handler.ParseParam(params, "tag")
	tag = filepath.FromSlash(tag)

	query := r.URL.Query()

	db, err := handler.CreateTagProvider()
	if err != nil {
		handler.WriteError(w, err)
		return
	}
	defer db.Close()

	m, err := db.Read(tag)
	if err != nil {
		handler.WriteError(w, err)
		return
	}

	if fav, e := strconv.ParseBool(query.Get("favorite")); e == nil {
		if fav != m.Favorite {
			m.Favorite = fav
			db.Write(m)
		}
	}

	response := setTagFavoriteResponse{
		Favorite: m.Favorite,
	}

	handler.WriteJson(w, response)
}
