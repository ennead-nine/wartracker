package scanner

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/otiai10/gosseract/v2"
)

var (
	Debug      bool
	Process    int
	ScratchDir string
)

type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}

// PreProcessImage takes an image.Image object and applies filters for optimal OCR
func PreProcessImage(img image.Image, gray, invert, bg bool) (image.Image, error) {
	var ppimg = img
	if gray {
		ppimg = imaging.Grayscale(img)
	}
	if invert {
		ppimg = imaging.Invert(ppimg)
	}
	if bg {
		cimg := image.NewRGBA(ppimg.Bounds())
		draw.Draw(cimg, ppimg.Bounds(), ppimg, image.Point{}, draw.Over)
		fullrect := ppimg.Bounds()
		for x := fullrect.Min.X; x <= fullrect.Max.X; x++ {
			for y := fullrect.Min.Y; y <= fullrect.Max.Y; y++ {
				r1, _, _, _ := cimg.At(x, y).RGBA()

				if r1 >= 50 {
					cimg.Set(x, y, color.RGBA{255, 255, 255, 255})
				} else {
					cimg.Set(x, y, color.RGBA{0, 0, 0, 255})
				}
			}
		}
		ppimg = cimg
	}

	if Debug {
		out, err := os.Create(fmt.Sprintf("%s/debug-%d.png", ScratchDir, Process))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return ppimg, err
		}
		err = png.Encode(out, ppimg)
		if err != nil {
			return nil, err
		}
		Process += 1
	}

	return ppimg, nil
}

func GetImageRect(px, py, rx, ry int, img image.Image) image.Image {
	var newimg image.Image
	p := image.Point{px, py}
	r := image.Rect(0, 0, rx, ry)
	r = r.Add(p)
	newimg = img.(SubImager).SubImage(r)

	fmt.Println(img.Bounds().String())

	return newimg
}

func GetImageText(img image.Image, w ...string) (string, error) {
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		return "", err
	}

	client := gosseract.NewClient()
	defer client.Close()

	client.SetPageSegMode(6)
	if w != nil {
		client.SetWhitelist(w[0])
	}
	client.SetTessdataPrefix("/Users/erumer/src/github.com/tesseract-ocr/tessdata")
	client.SetImageFromBytes(buf.Bytes())
	text, err := client.Text()
	if err != nil {
		return "", err
	}

	fmt.Printf("GetImageText: FOUND: %s\n", text)
	text = strings.Replace(text, "<", "", -1)
	text = strings.Replace(text, ">", "", -1)
	return text, nil
}
