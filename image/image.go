package image

import (
	"bytes"
	"image"
	"image/jpeg"
	"io"

	_ "golang.org/x/image/webp"

	"github.com/nfnt/resize"
)

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
