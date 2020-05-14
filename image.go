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

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/nfnt/resize"
)

// GetImage returns an image with specific width/height while retains aspect ratio.
func GetImage(c echo.Context) error {
	p, err := url.PathUnescape(c.Param("*"))
	if err != nil {
		return err
	}

	db, err := connectDB()
	if err != nil {
		return err
	}
	defer db.Close()

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

	data, f, err := OpenZipEntry(db, p, index)

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
	output, err := CreateResized(reader, uint(width), uint(height))
	if err != nil {
		return err
	}

	return c.Blob(http.StatusOK, "image/jpeg", output)
}

func OpenZipEntry(db *sqlx.DB, name string, index int) (content []byte, filename string, err error) {
	var meta itemMeta
	meta.Read(db, name)

	if len(meta.FileIndices) == 0 {
		err = fmt.Errorf("image file not found")
	}

	fullpath := BaseDirectory + string(os.PathSeparator) + name
	r, err := zip.OpenReader(fullpath)
	if err != nil {
		return
	}

	defer r.Close()

	zf := r.File[meta.FileIndices[index]]

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
