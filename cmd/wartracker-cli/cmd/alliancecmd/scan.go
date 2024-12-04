/*
Copyright © 2024 P4K Ennead  <ennead.tbc@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package alliancecmd

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"strconv"
	"strings"
	"time"
	"wartracker/pkg/alliance"

	"github.com/disintegration/imaging"
	"github.com/otiai10/gosseract/v2"
	"github.com/spf13/cobra"
)

var imageFile string
var outputFile string
var server int64

type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}

// PreProcessImage takes an image.Image object and applies filters for optimal OCR
func PreProcessImage(img image.Image, gray, invert, bg bool) image.Image {
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

	return ppimg
}

func GetImageText(img image.Image, w string) (string, error) {
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		return "", err
	}

	client := gosseract.NewClient()
	defer client.Close()

	client.SetPageSegMode(6)
	client.SetWhitelist(w)
	client.SetTessdataPrefix("../../../../tesseract-ocr/tessdata")
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

// GetAllianceTagImage gets the alliance tag text from an alliance screenshot
func GetAllianceTagText(img image.Image) (string, error) {
	// Alliance tag rect
	pointX := 157
	pointY := 292
	rectX := 48
	rectY := 20

	p := image.Point{pointX, pointY}
	r := image.Rect(0, 0, rectX, rectY)
	r = r.Add(p)
	img = img.(SubImager).SubImage(r)

	img = PreProcessImage(img, false, false, false)

	outf, _ := os.Create(fmt.Sprintf("%s-tag.png", imageFile))
	defer outf.Close()
	_ = png.Encode(outf, img)

	return GetImageText(img, "<>0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
}

// GetAlliancePowerText gets the alliance power text from an alliance screenshot
func GetAlliancePowerText(img image.Image) (int, error) {
	// Alliance tag rect
	pointX := 280
	pointY := 317
	rectX := 96
	rectY := 18

	p := image.Point{pointX, pointY}
	r := image.Rect(0, 0, rectX, rectY)
	r = r.Add(p)
	img = img.(SubImager).SubImage(r)

	img = PreProcessImage(img, true, true, true)

	outf, _ := os.Create(fmt.Sprintf("%s-power.png", imageFile))
	defer outf.Close()
	_ = png.Encode(outf, img)

	tpower, err := GetImageText(img, "0123456789")
	if err != nil {
		return 0, err
	}

	power, err := strconv.Atoi(tpower)
	if err != nil {
		return 0, err
	}

	return power, nil
}

// GetAllianceGiftImage gets the alliance gift level text from an alliance screenshot
func GetAllianceGiftText(img image.Image) (int, error) {
	// Alliance tag rect
	pointX := 356
	pointY := 351
	rectX := 19
	rectY := 15

	p := image.Point{pointX, pointY}
	r := image.Rect(0, 0, rectX, rectY)
	r = r.Add(p)
	img = img.(SubImager).SubImage(r)

	img = PreProcessImage(img, true, true, true)

	outf, _ := os.Create(fmt.Sprintf("%s-gift.png", imageFile))
	defer outf.Close()
	_ = png.Encode(outf, img)

	tgift, err := GetImageText(img, "0123456789")
	if err != nil {
		return 0, err
	}

	gift, err := strconv.Atoi(tgift)
	if err != nil {
		return 0, err
	}

	return gift, nil
}

// GetAllianceGiftImage gets the alliance gift level text from an alliance screenshot
func GetAllianceMemberText(img image.Image) (int, error) {
	// Alliance tag rect
	pointX := 316
	pointY := 366
	rectX := 28
	rectY := 16

	p := image.Point{pointX, pointY}
	r := image.Rect(0, 0, rectX, rectY)
	r = r.Add(p)
	img = img.(SubImager).SubImage(r)

	img = PreProcessImage(img, true, true, true)

	outf, _ := os.Create(fmt.Sprintf("%s-member.png", imageFile))
	defer outf.Close()
	_ = png.Encode(outf, img)

	tmemcount, err := GetImageText(img, "0123456789")
	if err != nil {
		return 0, err
	}

	memcount, err := strconv.Atoi(tmemcount)
	if err != nil {
		return 0, err
	}

	return memcount, nil
}

// ScanAlliance pre-processes the given image file and scans it with tessaract
// into an alliance.Alliance struct
func ScanAlliance() (*alliance.Alliance, error) {
	var a alliance.Alliance
	var d alliance.Data

	// Open image file
	imgfile, err := os.Open(imageFile)
	if err != nil {
		panic(err)
	}
	defer imgfile.Close()
	img, err := png.Decode(imgfile)
	if err != nil {
		panic(err)
	}

	// Setup alliance
	tag, err := GetAllianceTagText(img)
	if err != nil {
		return nil, err
	}
	power, err := GetAlliancePowerText(img)
	if err != nil {
		return nil, err
	}
	gift, err := GetAllianceGiftText(img)
	if err != nil {
		return nil, err
	}
	memcount, err := GetAllianceMemberText(img)
	if err != nil {
		return nil, err
	}

	d.Date = time.Now().Format(time.DateOnly)
	d.Tag = tag
	d.Power = int64(power)
	d.GiftLevel = int64(gift)
	d.MemberCount = int64(memcount)
	a.Data = append(a.Data, d)
	a.Server = server
	err = a.GetByTag(d.Tag)
	if err != nil {
		err = a.Add(server)
		if err != nil {
			return nil, err
		}
	}

	fmt.Printf("Aliiance:\n%#v\n", a)

	return &a, nil
}

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scans an alliance screenshot into a json file.",
	Long: `Scan takes an alliance screenshot and Marshals an alliance object 
	into json for cleanup.  Running wartracjer-cli alliance create with the 
	cleaned json will create an alliance object in the database.
	
	Example: wartracker-cli alliance scan -i alliance.png`,
	Run: func(cmd *cobra.Command, args []string) {
		_, err := ScanAlliance()
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	allianceCmd.AddCommand(scanCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// canmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	scanCmd.Flags().StringVarP(&imageFile, "image", "i", "", "image file (PNG) to scan for alliance data")
	scanCmd.Flags().StringVarP(&outputFile, "output", "o", "", "JSON file to output alliance data to")
	scanCmd.Flags().Int64VarP(&server, "server", "s", 1, "Alliance's server number")
}
