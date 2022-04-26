package view

import (
	"github.com/julienschmidt/httprouter"
	"github.com/wutipong/mangaweb/handler"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"
)

func Download(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	item := handler.ParseParam(params, "item")
	item = filepath.FromSlash(item)

	db, err := handler.CreateMetaProvider()
	if err != nil {
		handler.WriteError(w, err)
		return
	}
	defer db.Close()

	m, err := db.Read(item)
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

	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		handler.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Length", strconv.Itoa(len(bytes)))
}
