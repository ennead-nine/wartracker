package common

import (
	"archive/zip"
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"math"

	"golang.org/x/image/draw"
)

type CropRatios struct {
	Rpx float64 `json:"rpx"`
	Rpy float64 `json:"rpy"`
	Rsx float64 `json:"rsx"`
	Rsy float64 `json:"rsy"`
}

func UnzipSS(z []byte) ([]*zip.File, error) {
	a, err := zip.NewReader(bytes.NewReader(z), int64(len(z)))
	if err != nil {
		return nil, fmt.Errorf("failed to create zip reader: %w", err)
	}

	return a.File, nil
}

func ConvertJpgToPng(b []byte) ([]byte, error) {
	img, err := jpeg.Decode(bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("unable to decode jpeg: %w", err)
	}

	buf := new(bytes.Buffer)
	if err := png.Encode(buf, img); err != nil {
		return nil, fmt.Errorf("unable to encode png: %w", err)
	}

	return buf.Bytes(), nil
}

func ResizeImage(b []byte, x int, y int) ([]byte, error) {
	src, err := png.Decode(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	dst := image.NewRGBA(image.Rect(0, 0, x, y))
	draw.NearestNeighbor.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)

	var output bytes.Buffer
	err = png.Encode(&output, dst)
	if err != nil {
		return nil, err
	}

	return output.Bytes(), nil
}

func CropImage(b []byte, cr CropRatios) ([]byte, error) {
	img, err := png.Decode(bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("could not decode image: %w", err)
	}

	r := img.Bounds()

	px := (float64(r.Max.X) - float64(r.Min.X)) * float64(cr.Rpx)
	py := (float64(r.Max.Y) - float64(r.Min.Y)) * float64(cr.Rpy)
	sx := (float64(r.Max.X) - float64(r.Min.X)) * float64(cr.Rsx)
	sy := (float64(r.Max.Y) - float64(r.Min.Y)) * float64(cr.Rsy)

	r = image.Rect(r.Min.X, r.Min.Y, int(math.Round(sx)), int(math.Round(sy)))
	r = r.Add(image.Point{X: int(math.Round(px)), Y: int(math.Round(py))})
	img = img.(SubImager).SubImage(r)

	buf := new(bytes.Buffer)
	err = png.Encode(buf, img)
	if err != nil {
		return nil, fmt.Errorf("could not encode image: %w", err)
	}

	return buf.Bytes(), nil
}

type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}
