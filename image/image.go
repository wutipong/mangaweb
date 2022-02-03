package image

import (
	"archive/zip"
	"bytes"
	"image"
	"image/jpeg"
	"io"

	_ "golang.org/x/image/webp"

	"github.com/nfnt/resize"
)

func Create(reader io.Reader) (output image.Image, err error) {
	output, _, err = image.Decode(reader)
	return
}

func CreateCover(items []int, r *zip.ReadCloser) (output image.Image, err error) {
	reader, err := r.File[0].Open()
	if err != nil {
		return
	}
	output, err = Create(reader)
	return
}

func Resize(img image.Image, width uint, height uint) image.Image {
	return resize.Thumbnail(width, height, img, resize.MitchellNetravali)
}

func CreateThumbnail(img image.Image) image.Image {
	const thumbnailHeight = 200
	return resize.Resize(0, thumbnailHeight, img, resize.MitchellNetravali)
}

func ToJPEG(img image.Image) (output []byte, err error) {
	buffer := bytes.Buffer{}

	err = jpeg.Encode(&buffer, img, nil)
	if err != nil {
		return
	}
	output = buffer.Bytes()
	return
}
