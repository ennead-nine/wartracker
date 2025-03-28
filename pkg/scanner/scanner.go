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
	Debug      bool
	ScratchDir string
}

var client *vision.ImageAnnotatorClient

func NewScanner(f []byte, debug bool) (*Scanner, error) {
	var s Scanner
	var err error

	s.ImageBytes = f
	// TODO: Needs to come from config
	s.ScratchDir = "_scratch"
	s.Debug = debug

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
