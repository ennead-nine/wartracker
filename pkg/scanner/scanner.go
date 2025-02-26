package scanner

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/png"
	"os"
	"strings"

	vision "cloud.google.com/go/vision/apiv1"
)

type Scanner struct {
	ImageBytes []byte
	CropPoint  image.Point
	CropSize   image.Point
	Debug      bool
	Crop       bool
	ScratchDir string
}

var client *vision.ImageAnnotatorClient

func NewScanner(f []byte, d, c bool, pts ...image.Point) (*Scanner, error) {
	var s Scanner
	var err error

	s.ImageBytes = f
	// TODO: Needs to come from config
	s.ScratchDir = "_scratch"
	s.Debug = d
	s.Crop = c
	if c && len(pts) == 2 {
		s.CropPoint = pts[0]
		s.CropSize = pts[1]
	} else if len(pts) != 2 {
		return nil, fmt.Errorf("crop set to true but crop points passed is not equal to 2")
	}

	ctx := context.Background()

	client, err = vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return &s, err
	}

	return &s, nil
}

func (s *Scanner) SetImageBytes(f []byte) {
	s.ImageBytes = nil
	s.ImageBytes = append(s.ImageBytes, f...)
}

func (s *Scanner) ScanImage() ([]string, error) {
	ctx := context.Background()

	if s.Crop {
		err := s.CropImage()
		if err != nil {
			return nil, fmt.Errorf("could not crop image: %w", err)
		}
	}

	image, err := vision.NewImageFromReader(bytes.NewReader(s.ImageBytes))
	if err != nil {
		return nil, err
	}

	annotations, err := client.DetectTexts(ctx, image, nil, 40)
	if err != nil {
		return nil, fmt.Errorf("vision detecttexts error: %w", err)
	}
	if len(annotations) < 1 {
		return nil, fmt.Errorf("vision detecttexts fround no text in image")
	}

	texts := strings.Split(annotations[0].Description, "\n")

	return texts, nil
}

func (s *Scanner) CropImage() error {
	img, err := png.Decode(bytes.NewReader(s.ImageBytes))
	if err != nil {
		return fmt.Errorf("could not decode image: %w", err)
	}

	r := img.Bounds()
	r = image.Rect(r.Min.X, r.Min.Y, s.CropSize.X, s.CropSize.Y)
	r = r.Add(s.CropPoint)
	img = img.(SubImager).SubImage(r)

	if s.Debug {
		err = s.SaveImage(img)
		if err != nil {
			return fmt.Errorf("could not save image: %w", err)
		}
	}

	buf := new(bytes.Buffer)
	err = png.Encode(buf, img)
	if err != nil {
		return fmt.Errorf("could not encode image: %w", err)
	}

	s.ImageBytes = buf.Bytes()

	return nil
}

func (s *Scanner) SaveImage(img image.Image) error {
	out, err := os.CreateTemp("_scratch", "debug*.png")
	if err != nil {
		return fmt.Errorf("could not open temp file: %w", err)
	}
	defer out.Close()

	err = png.Encode(out, img)
	if err != nil {
		return fmt.Errorf("could not write image to temp file: %w", err)
	}

	return nil
}

type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}
