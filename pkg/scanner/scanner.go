package scanner

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/otiai10/gosseract/v2"
	"gopkg.in/gographics/imagick.v3/imagick"
)

var (
	Debug       bool
	Process     int
	ScratchDir  string
	TessdataDir string
)

type ImageMap struct {
	Rect
	CharFilter
	PreProcess
	Image image.Image
}

type Rect struct {
	PX int
	PY int
	RX int
	RY int
}

type CharFilter struct {
	Filter string
}

type PreProcess struct {
	Gray   bool
	Invert bool
	BG     bool
}

type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}

func ParseAbbvInt(s string) (int64, error) {
	if len(s) < 2 {
		return 0, fmt.Errorf("length of string \"%s\" is less than 2", s)
	}
	trim := strings.TrimSpace(s)
	unit := trim[len(trim)-1:]
	unit = strings.ToUpper(unit)
	tbase := strings.TrimRight(trim, "KMG")
	base, err := strconv.ParseFloat(tbase, 64)
	if err != nil {
		return 0, err
	}

	switch unit {
	case "K":
		i := int64(math.Round(base * 1000))
		return i, nil
	case "M":
		i := int64(math.Round(base * 1000000))
		return i, nil
	case "G":
		i := int64(math.Round(base * 1000000000))
		return i, nil
	default:
		i := int64(math.Round(base))
		return i, nil
	}
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

func SetImageDensity(inFile string, d int) (image.Image, error) {
	imagick.Initialize()
	// Schedule cleanup
	defer imagick.Terminate()
	var err error

	mw := imagick.NewMagickWand()

	i, err := os.OpenFile(inFile, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer i.Close()

	err = mw.ReadImageFile(i)
	if err != nil {
		return nil, err
	}

	err = mw.SetResolution(float64(d), float64(d))
	if err != nil {
		return nil, err
	}

	outFile := ScratchDir + "/" + filepath.Base(inFile)
	err = mw.WriteImage(outFile)
	if err != nil {
		return nil, err
	}

	imgfile, err := os.Open(outFile)
	if err != nil {
		return nil, err
	}
	defer imgfile.Close()
	img, err := png.Decode(imgfile)
	if err != nil {
		return nil, err
	}

	return img, nil
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
		if w[0] != "" {
			client.SetWhitelist(w[0])
		}
	}
	client.SetTessdataPrefix(TessdataDir)
	client.SetImageFromBytes(buf.Bytes())
	text, err := client.Text()
	if err != nil {
		return "", err
	}

	if Debug {
		fmt.Printf("GetImageText: FOUND: %s\n", text)
	}
	return text, nil
}

func (im *ImageMap) ProcessImage() ([]byte, error) {
	img := GetImageRect(im.PX, im.PY, im.RX, im.RY, im.Image)
	img, err := PreProcessImage(img, im.Gray, im.Invert, im.BG)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	err = png.Encode(&buf, img)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (im *ImageMap) ProcessImageInt() (int64, error) {
	img := GetImageRect(im.PX, im.PY, im.RX, im.RY, im.Image)
	img, err := PreProcessImage(img, im.Gray, im.Invert, im.BG)
	if err != nil {
		return 0, err
	}

	t, err := GetImageText(img, im.Filter)
	if err != nil {
		return 0, err
	}

	i, err := strconv.Atoi(t)
	if err != nil {
		return 0, err
	}

	return int64(i), nil
}

func (im *ImageMap) ProcessImageAbbrInt() (int64, error) {
	img := GetImageRect(im.PX, im.PY, im.RX, im.RY, im.Image)
	img, err := PreProcessImage(img, im.Gray, im.Invert, im.BG)
	if err != nil {
		return 0, err
	}

	t, err := GetImageText(img, im.Filter)
	if err != nil {
		return 0, err
	}

	i, err := ParseAbbvInt(t)
	if err != nil {
		return 0, err
	}

	return int64(i), nil
}

func (im *ImageMap) ProcessImageText() (string, error) {
	img := GetImageRect(im.PX, im.PY, im.RX, im.RY, im.Image)
	img, err := PreProcessImage(img, im.Gray, im.Invert, im.BG)
	if err != nil {
		return "", err
	}

	t, err := GetImageText(img, im.Filter)
	if err != nil {
		return "", err
	}

	return t, nil
}
