package scanner

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
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
	Languages   []string
)

type ImageMap struct {
	Rect
	CharFilter
	PreProcess
}

type ImageMaps map[string]ImageMap

type Rect struct {
	PX int `yaml:"px"`
	PY int `yaml:"py"`
	RX int `yaml:"rx"`
	RY int `yaml:"ry"`
}

type CharFilter struct {
	Filter string `yaml:"filter"`
}

type PreProcess struct {
	Gray   bool `yaml:"gray"`
	Invert bool `yaml:"invert"`
	BG     bool `yaml:"bg"`
}

type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}

func ParseAbbvInt(s string) (int, error) {
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
		i := int(math.Round(base * 1000))
		return i, nil
	case "M":
		i := int(math.Round(base * 1000000))
		return i, nil
	case "G":
		i := int(math.Round(base * 1000000000))
		return i, nil
	default:
		i := int(math.Round(base))
		return i, nil
	}
}

func grayImage(img image.Image) (image.Image, error) {
	Process += 1

	if Debug {
		out, err := os.Create(fmt.Sprintf("%s/debug-%d-pregrey.png", ScratchDir, Process))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return img, err
		}
		err = png.Encode(out, img)
		if err != nil {
			return nil, err
		}
		out.Close()
	}

	img = imaging.Grayscale(img)

	if Debug {
		out, err := os.Create(fmt.Sprintf("%s/debug-%d-postgrey.png", ScratchDir, Process))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return img, err
		}
		err = png.Encode(out, img)
		if err != nil {
			return nil, err
		}
		out.Close()
	}

	return img, nil
}

func invertImage(img image.Image) (image.Image, error) {
	Process += 1

	if Debug {
		out, err := os.Create(fmt.Sprintf("%s/debug-%d-preinvert.png", ScratchDir, Process))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return img, err
		}
		err = png.Encode(out, img)
		if err != nil {
			return nil, err
		}
		out.Close()
	}

	img = imaging.Invert(img)

	if Debug {
		out, err := os.Create(fmt.Sprintf("%s/debug-%d-postinvert.png", ScratchDir, Process))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return img, err
		}
		err = png.Encode(out, img)
		if err != nil {
			return nil, err
		}
		out.Close()
	}

	return img, nil
}

func bgImage(img image.Image) (image.Image, error) {
	Process++

	if Debug {
		out, err := os.Create(fmt.Sprintf("%s/debug-%d-prebg.png", ScratchDir, Process))
		if err != nil {
			return img, err
		}
		err = png.Encode(out, img)
		if err != nil {
			return nil, err
		}
		out.Close()
	}

	var minR uint32 = 0xffff
	var minCount int = 0
	for x := img.Bounds().Min.X; x <= img.Bounds().Max.X-1; x++ {
		for y := img.Bounds().Min.Y; y <= img.Bounds().Max.Y-1; y++ {
			c := img.At(x, y)
			r, _, _, _ := c.RGBA()
			if r < minR {
				minR = r
				minCount++
			}

		}
	}
	if Debug {
		fmt.Printf("BG:\n\tmincount for %d: %d\n\tmin for %d: %d\n", Process, minCount, Process, minR)
	}
	cimg := image.NewRGBA(img.Bounds())
	for x := img.Bounds().Min.X; x <= img.Bounds().Max.X-1; x++ {
		for y := img.Bounds().Min.Y; y <= img.Bounds().Max.Y-1; y++ {
			c := img.At(x, y)
			r, _, _, _ := c.RGBA()
			if r == minR {
				cimg.Set(x, y, color.Black)
			} else {
				cimg.Set(x, y, color.White)
			}
		}
	}
	img = cimg

	if Debug {
		out, err := os.Create(fmt.Sprintf("%s/debug-%d-postbg.png", ScratchDir, Process))
		if err != nil {
			return nil, err
		}
		err = png.Encode(out, img)
		if err != nil {
			return nil, err
		}
		out.Close()
	}

	return img, nil
}

// PreProcessImage takes an image.Image object and applies filters for optimal OCR
func PreProcessImage(img image.Image, gray, invert, bg bool) (image.Image, error) {
	var err error
	if gray {
		img, err = grayImage(img)
		if err != nil {
			return nil, err
		}
	}
	if invert {
		img, err = invertImage(img)
		if err != nil {
			return nil, err
		}
	}
	if bg {
		img, err = bgImage(img)
		if err != nil {
			return nil, err
		}
	}

	return img, nil
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
	if len(Languages) > 0 {
		client.SetLanguage(Languages...)
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

func (im *ImageMap) ProcessImage(img image.Image) ([]byte, error) {
	img = GetImageRect(im.PX, im.PY, im.RX, im.RY, img)
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

func (im *ImageMap) ProcessImageInt(img image.Image) (int, error) {
	img = GetImageRect(im.PX, im.PY, im.RX, im.RY, img)
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

	return i, nil
}

func (im *ImageMap) ProcessImageAbbrInt(img image.Image) (int, error) {
	img = GetImageRect(im.PX, im.PY, im.RX, im.RY, img)
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

	return i, nil
}

func (im *ImageMap) ProcessImageText(img image.Image) (string, error) {
	img = GetImageRect(im.PX, im.PY, im.RX, im.RY, img)
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
