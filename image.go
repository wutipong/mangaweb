package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"io"
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

	reader, f, err := OpenZipEntry(p, index)
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
		switch filepath.Ext(strings.ToLower(f)) {
		case ".jpg", ".jpeg":
			mimetype = "image/jpeg"
		case ".png":
			mimetype = "image/png"
		}

		return c.Blob(http.StatusOK, mimetype, data)
	}
	output, err := CreateResized(reader, uint(width), uint(height))
	if err != nil {
		return err
	}

	return c.Blob(http.StatusOK, "image/jpeg", output)
}

func OpenZipEntry(name string, index int) (reader io.ReadCloser, filename string, err error) {
	fullpath := BaseDirectory + string(os.PathSeparator) + name
	r, err := zip.OpenReader(fullpath)
	if err != nil {
		return
	}
	if index < 0 || index > len(r.File) {
		err = fmt.Errorf("index out of range : %v", index)
		return
	}
	f := r.File[index]
	if !filter(f.Name) {
		err = fmt.Errorf("invalid format : %v", f.Name)
		return
	}

	filename = f.Name
	reader, err = f.Open()
	return
}

func CreateResized(reader io.Reader, width uint, height uint) (output []byte, err error) {
	img, _, err := image.Decode(reader)
	if err != nil {
		return
	}

	resized := resize.Thumbnail(width, height, img, resize.MitchellNetravali)
	buffer := bytes.Buffer{}

	err = jpeg.Encode(&buffer, resized, nil)
	if err != nil {
		return
	}
	output = buffer.Bytes()
	return
}
