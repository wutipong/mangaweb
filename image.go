package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/nfnt/resize"
)

// GetImage returns an image with specific width/height while retains aspect ratio.
func GetImage(c echo.Context) error {
	p, err := url.PathUnescape(c.Param("*"))
	if err != nil {
		return err
	}
	fullpath := BaseDirectory + string(os.PathSeparator) + p

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

	r, err := zip.OpenReader(fullpath)
	if err != nil {
		return err
	}
	defer r.Close()

	if index < 0 || index > len(r.File) {
		return fmt.Errorf("index out of range : %v", index)
	}

	f := r.File[index]

	if !filter(f.Name) {
		return fmt.Errorf("invalid format : %v", f.Name)
	}

	reader, err := f.Open()
	if err != nil {
		return nil
	}
	defer reader.Close()

	if width == 0 || height == 0 {
		data, err := ioutil.ReadAll(reader)

		if err != nil {
			return nil
		}

		var mimetype string
		switch filepath.Ext(strings.ToLower(f.Name)) {
		case ".jpg", ".jpeg":
			mimetype = "image/jpeg"
		case ".png":
			mimetype = "image/png"
		}

		c.Blob(http.StatusOK, mimetype, data)
	}

	img, _, err := image.Decode(reader)

	if err != nil {
		return err
	}

	resized := resize.Thumbnail(uint(width), uint(height), img, resize.MitchellNetravali)
	buffer := bytes.Buffer{}

	jpeg.Encode(&buffer, resized, nil)
	return c.Blob(http.StatusOK, "image/jpeg", buffer.Bytes())
}
