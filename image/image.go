package image

import (
	"archive/zip"
	"bytes"
	"image"
	"io"

	_ "golang.org/x/image/webp"

	"github.com/disintegration/imaging"
)

func Create(reader io.Reader) (output image.Image, err error) {
	output, err = imaging.Decode(reader, imaging.AutoOrientation(true))
	return
}

func CreateCover(fileIndex int, r *zip.ReadCloser) (output image.Image, err error) {
	reader, err := r.File[fileIndex].Open()
	if err != nil {
		return
	}
	output, err = Create(reader)
	return
}

// TODO: Add resize function that preserve aspect ratio.
func Resize(img image.Image, width uint, height uint) image.Image {
	return imaging.Resize(img, int(width), int(height), imaging.MitchellNetravali)
}

func CreateThumbnail(img image.Image) image.Image {
	const thumbnailHeight = 200
	return imaging.Resize(img, 0, thumbnailHeight, imaging.MitchellNetravali)
}

func ToJPEG(img image.Image) (output []byte, err error) {
	buffer := bytes.Buffer{}

	err = imaging.Encode(&buffer, img, imaging.JPEG, imaging.JPEGQuality(60))
	output = buffer.Bytes()

	return
}
