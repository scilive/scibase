package images

import (
	"bytes"
	"image"
	"io"
	"math"

	"github.com/disintegration/imaging"
)

// Crop crop image
func Crop(file io.Reader, fileName string, x, y, w, h int) ([]byte, error) {
	format, err := imaging.FormatFromFilename(fileName)
	if err != nil {
		return nil, err
	}
	im, err := imaging.Decode(file)
	if err != nil {
		return nil, err
	}
	im = imaging.Crop(im, image.Rect(x, y, x+w, y+h))
	buff := new(bytes.Buffer)
	err = imaging.Encode(buff, im, format)
	if nil != err {
		return nil, err
	}
	return buff.Bytes(), nil
}

// Resize resize image
func Resize(file io.Reader, fileName string, w, h int) ([]byte, error) {
	format, err := imaging.FormatFromFilename(fileName)
	if err != nil {
		return nil, err
	}
	im, err := imaging.Decode(file)
	if err != nil {
		return nil, err
	}
	if w <= 0 || h <= 0 {
		ow, oh := im.Bounds().Dx(), im.Bounds().Dy()
		if w <= 0 { // resize by height
			w = ow * h / oh
		}
		if h <= 0 { // resize by width
			h = oh * w / ow
		}
	}

	im = imaging.Resize(im, w, h, imaging.Lanczos)
	buff := new(bytes.Buffer)
	err = imaging.Encode(buff, im, format)
	if nil != err {
		return nil, err
	}
	return buff.Bytes(), nil
}

// Thumb resize image and keep original ratio
func Thumbnail(file io.Reader, fileName string, w, h int) ([]byte, error) {
	format, err := imaging.FormatFromFilename(fileName)
	if err != nil {
		return nil, err
	}
	im, err := imaging.Decode(file)
	if err != nil {
		return nil, err
	}
	ow, oh := im.Bounds().Dx(), im.Bounds().Dy()
	nw, nh, x, y := GetGoodResize(ow, oh, w, h)

	im = imaging.Resize(im, nw, nh, imaging.Lanczos)
	im = imaging.Crop(im, image.Rect(x, y, x+w, y+h))
	buff := new(bytes.Buffer)
	err = imaging.Encode(buff, im, format)
	if nil != err {
		return nil, err
	}
	return buff.Bytes(), nil
}

// GetGoodResize get good resize
func GetGoodResize(ow, oh, w, h int) (int, int, int, int) {
	rw := float64(w) / float64(ow)
	rh := float64(h) / float64(oh)
	x, y := 0, 0
	nw, nh := w, h
	if math.Abs(rw-1) < math.Abs(rh-1) { //select the resize smaller ratio
		nh = int(float64(oh) * rw)
		y = (h - nh) / 2
		if y < 0 {
			y = -y
		}
	} else {
		nw = int(float64(ow) * rh)
		x = (w - nw) / 2
		if x < 0 {
			x = -x
		}
	}
	return nw, nh, x, y
}
