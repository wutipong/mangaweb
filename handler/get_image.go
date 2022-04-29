package handler

import (
	"archive/zip"
	"bytes"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/wutipong/mangaweb/log"
	"go.uber.org/zap"
	_ "image/png"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/wutipong/mangaweb/image"
	"github.com/wutipong/mangaweb/meta"
)

// GetImage returns an image with specific width/height while retains aspect ratio.
func GetImage(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	item := ParseParam(params, "item")
	item = filepath.FromSlash(item)

	query := r.URL.Query()

	provider, err := CreateMetaProvider()
	if err != nil {
		WriteError(w, err)
		return
	}
	defer provider.Close()

	var width, height int64 = 0, 0
	if w, e := strconv.ParseInt(query.Get("width"), 10, 64); e == nil {
		width = w
		height = width
	}

	if h, e := strconv.ParseInt(query.Get("height"), 10, 64); e == nil {
		height = h
	}

	var index = 0
	if i, err := strconv.Atoi(query.Get("i")); err == nil {
		index = i
	}

	log.Get().Info("Get image", zap.String("item_name", item), zap.Int("index", index))

	m, err := provider.Read(item)
	if err != nil {
		WriteError(w, err)
		return
	}
	data, f, err := OpenZipEntry(m, index)
	if err != nil {
		WriteError(w, err)
		return
	}

	if width == 0 || height == 0 {
		var contentType string
		switch filepath.Ext(strings.ToLower(f)) {
		case ".jpg", ".jpeg":
			contentType = "image/jpeg"
		case ".png":
			contentType = "image/png"
		case ".webp":
			contentType = "image/webp"
		}

		w.WriteHeader(http.StatusOK)
		w.Write(data)
		w.Header().Set("Content-Type", contentType)

		return
	}

	reader := bytes.NewBuffer(data)

	img, err := image.Create(reader)
	if err != nil {
		WriteError(w, err)
		return
	}

	resized := image.Resize(img, uint(width), uint(height))
	output, err := image.ToJPEG(resized)
	if err != nil {
		WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(output)
	w.Header().Set("Content-Type", "image/jpeg")
}

func OpenZipEntry(m meta.Meta, index int) (content []byte, filename string, err error) {
	if len(m.FileIndices) == 0 {
		err = fmt.Errorf("image file not found")
	}

	fullpath := filepath.Join(meta.BaseDirectory, m.Name)
	r, err := zip.OpenReader(fullpath)
	if err != nil {
		return
	}

	defer r.Close()

	zf := r.File[m.FileIndices[index]]

	if zf == nil {
		err = fmt.Errorf("file not found : %v", index)
		return
	}

	filename = zf.Name
	reader, err := zf.Open()
	if err != nil {
		return
	}
	defer reader.Close()
	if content, err = ioutil.ReadAll(reader); err != nil {
		content = nil
		return
	}
	return
}
