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
	"database/sql"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"os"
	"strconv"
	"time"
	"wartracker/pkg/alliance"
	"wartracker/pkg/scanner"

	"github.com/spf13/cobra"
)

var mainImageFile string
var giftImageFile string
var mainOutputFile string
var mainServer int64

// GetAllianceTagImage gets the alliance tag text from an alliance screenshot
func GetMainAllianceTagText(img image.Image) (string, error) {
	// Alliance tag rect
	px := 27
	py := 546
	rx := 140
	ry := 60

	img = scanner.GetImageRect(px, py, rx, ry, img)
	img, err := scanner.PreProcessImage(img, false, false, false)
	if err != nil {
		return "", err
	}

	return scanner.GetImageText(img, "<>0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
}

// GetAllianceNameImage gets the alliance name text from an alliance screenshot
func GetMainAllianceNameText(img image.Image) (string, error) {
	// Alliance name rect
	px := 160
	py := 546
	rx := 500
	ry := 60

	img = scanner.GetImageRect(px, py, rx, ry, img)
	img, err := scanner.PreProcessImage(img, false, false, false)
	if err != nil {
		return "", err
	}

	return scanner.GetImageText(img)
}

// GetAlliancePowerText gets the alliance power text from an alliance screenshot
func GetMainAlliancePowerText(img image.Image) (int, error) {
	// Alliance tag rect
	px := 155
	py := 680
	rx := 333
	ry := 60

	img = scanner.GetImageRect(px, py, rx, ry, img)
	img, err := scanner.PreProcessImage(img, false, false, false)
	if err != nil {
		return 0, err
	}

	tpower, err := scanner.GetImageText(img, "0123456789")
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
func GetMainAllianceGiftText(img image.Image) (int, error) {
	// Alliance tag rect
	px := 146
	py := 427
	rx := 47
	ry := 34

	img = scanner.GetImageRect(px, py, rx, ry, img)
	img, err := scanner.PreProcessImage(img, true, true, true)
	if err != nil {
		return 0, err
	}

	tgift, err := scanner.GetImageText(img, "0123456789")
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
func GetMainAllianceMemberText(img image.Image) (int, error) {
	// Alliance tag rect
	px := 648
	py := 642
	rx := 42
	ry := 42

	img = scanner.GetImageRect(px, py, rx, ry, img)
	img, err := scanner.PreProcessImage(img, false, false, false)
	if err != nil {
		return 0, err
	}

	tmemcount, err := scanner.GetImageText(img, "0123456789")
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
func ScanMainAlliance() (*alliance.Alliance, error) {
	var a alliance.Alliance
	var d alliance.Data

	// Open image file
	imgfile, err := os.Open(mainImageFile)
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
	name, err := GetAllianceNameText(img)
	if err != nil {
		return nil, err
	}
	power, err := GetAlliancePowerText(img)
	if err != nil {
		return nil, err
	}
	memcount, err := GetAllianceMemberText(img)
	if err != nil {
		return nil, err
	}

	// Gift Level for Main Alliance

	// Open image file
	giftfile, err := os.Open(giftImageFile)
	if err != nil {
		panic(err)
	}
	defer giftfile.Close()
	img, err = png.Decode(giftfile)
	if err != nil {
		panic(err)
	}

	gift, err := GetAllianceGiftText(img)
	if err != nil {
		return nil, err
	}

	d.Date = time.Now().Format(time.DateOnly)
	d.Tag = tag
	d.Name = name
	d.Power = int64(power)
	d.GiftLevel = int64(gift)
	d.MemberCount = int64(memcount)
	a.Data = append(a.Data, d)
	a.Server = server

	err = a.GetByTag(d.Tag)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
		fmt.Printf("A new alliance will need to be created from this data.  Please run 'wartracker-cli alliance new -o %s' after verifying the data\n", outputFile)
	} else {
		fmt.Printf("This alliance already exists. To add the new data run 'wartracker-cli alliance add -o %s' to add the new data.\n", outputFile)
	}

	j, err := json.MarshalIndent(a, "", "\t")
	if err != nil {
		return nil, err
	}
	err = os.WriteFile(outputFile, j, 0644)
	if err != nil {
		return nil, err
	}

	return &a, nil
}

// scanCmd represents the scan command
var scanmainCmd = &cobra.Command{
	Use:   "scanmain",
	Short: "Scans an alliance screenshot into a json file.",
	Long: `Scan takes an alliance screenshot and Marshals an alliance object 
	into json for cleanup.  Running wartracjer-cli alliance create with the 
	cleaned json will create an alliance object in the database.
	
	Example: wartracker-cli alliance scan -i alliance.png`,
	Run: func(cmd *cobra.Command, args []string) {
		_, err := ScanMainAlliance()
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	allianceCmd.AddCommand(scanmainCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// canmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	scanmainCmd.Flags().StringVarP(&mainImageFile, "mainImage", "m", "", "Image file with main alliance info (PNG) to scan for alliance data")
	scanmainCmd.Flags().StringVarP(&giftImageFile, "giftImage", "g", "", "Image file with alliance gift info (PNG) to scan for alliance data")
	scanmainCmd.Flags().StringVarP(&mainOutputFile, "output", "o", "", "JSON file to output alliance data to")
	scanmainCmd.Flags().Int64VarP(&mainServer, "server", "s", 1, "Alliance's server number")
}
