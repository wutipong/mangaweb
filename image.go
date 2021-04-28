package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	_ "image/png"
	"io/ioutil"
	"mangaweb/image"
	"mangaweb/meta"
	"mangaweb/meta/mongo"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

// GetImage returns an image with specific width/height while retains aspect ratio.
func GetImage(c echo.Context) error {
	p, err := url.PathUnescape(c.Param("*"))
	if err != nil {
		return err
	}

	provider, err := mongo.New()
	if err != nil {
		return err
	}
	defer provider.Close()

	var width, height int64 = 0, 0
	if w, e := strconv.ParseInt(c.QueryParam("width"), 10, 64); e == nil {
		width = w
		height = width
	}

	if h, e := strconv.ParseInt(c.QueryParam("height"), 10, 64); e == nil {
		height = h
	}

	var index = 0
	if i, err := strconv.Atoi(c.QueryParam("i")); err == nil {
		index = i
	}

	m, err := provider.Read(p)
	data, f, err := OpenZipEntry(m, index)

	if width == 0 || height == 0 {
		if err != nil {
			return nil
		}

		var mimetype string
		switch filepath.Ext(strings.ToLower(f)) {
		case ".jpg", ".jpeg":
			mimetype = "image/jpeg"
		case ".png":
			mimetype = "image/png"
		}

		return c.Blob(http.StatusOK, mimetype, data)
	}
	reader := bytes.NewBuffer(data)
	output, err := image.CreateResized(reader, uint(width), uint(height))
	if err != nil {
		return err
	}

	return c.Blob(http.StatusOK, "image/jpeg", output)
}

func OpenZipEntry(m meta.Item, index int) (content []byte, filename string, err error) {
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

	defer reader.Close()
	if content, err = ioutil.ReadAll(reader); err != nil {
		content = nil
		return
	}
	return
}
