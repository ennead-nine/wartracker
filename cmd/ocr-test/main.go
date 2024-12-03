package main

import (
	"fmt"
	"image"
	"image/png"
	"os"

	"github.com/disintegration/imaging"
	"github.com/otiai10/gosseract/v2"
)

var Ranks [7]image.Image

type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}

func GetLines(f string) [7]image.Image {
	var imgs [7]image.Image

	baseX := 75
	baseY := 546
	lineX := 110
	lineY := 90

	oimgFile, err := os.Open(f)
	if err != nil {
		panic(err)
	}
	defer oimgFile.Close()

	oimg, err := png.Decode(oimgFile)
	if err != nil {
		panic(err)
	}

	i := 0
	for i < 7 {
		point := image.Point{baseX, baseY}
		rect := image.Rect(0, 0, lineX, lineY)
		rect = rect.Add(point)
		img := oimg.(SubImager).SubImage(rect)
		gray := imaging.Grayscale(img)
		invert := imaging.Invert(gray)

		imgs[i] = invert

		baseY += 180
		i += 1
	}

	return imgs
}

func main() {
	client := gosseract.NewClient()
	defer client.Close()

	f := "/Users/erumer/Downloads/Photos-001/Screenshot_20241129-062231.png"

	imgs := GetLines(f)

	i := 0
	for i < 7 {
		outf, err := os.Create(fmt.Sprintf("rank-line%d.png", i+15))
		if err != nil {
			panic(err)
		}
		defer outf.Close()
		if err := png.Encode(outf, imgs[i]); err != nil {
			panic(err)
		}

		d1 := []byte(fmt.Sprintf("%d", i+15))
		err = os.WriteFile(fmt.Sprintf("rank-line%d.gt.txt", i+15), d1, 0644)
		if err != nil {
			panic(err)
		}

		client.SetLanguage("wartracker", "eng")
		client.SetImage(fmt.Sprintf("rank%d.png", i))
		text, _ := client.Text()
		fmt.Println(text)

		i += 1
	}
}
