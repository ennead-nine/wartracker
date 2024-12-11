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
	"os"
	"strconv"
	"strings"
	"time"
	"wartracker/pkg/alliance"
	"wartracker/pkg/scanner"

	"github.com/spf13/cobra"
)

var mainImageFile string
var giftImageFile string
var mainOutputFile string
var mainServer int64

func MainImageMapper(im *scanner.ImageMap, field string) {
	switch field {
	case "tag":
		im.PX = 33
		im.PY = 546
		im.RX = 131
		im.RY = 65
		im.Filter = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
		im.Gray = false
		im.Invert = false
		im.BG = false
	case "name":
		im.PX = 160
		im.PY = 546
		im.RX = 798
		im.RY = 65
		im.Filter = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
		im.Gray = false
		im.Invert = false
		im.BG = false
	case "power":
		im.PX = 155
		im.PY = 680
		im.RX = 333
		im.RY = 60
		im.Filter = "0123456789"
		im.Gray = false
		im.Invert = false
		im.BG = false
	case "giftlevel":
		im.PX = 50
		im.PY = 420
		im.RX = 170
		im.RY = 51
		im.Filter = "0123456789"
		im.Gray = true
		im.Invert = true
		im.BG = true
	case "membercount":
		im.PX = 646
		im.PY = 624
		im.RX = 145
		im.RY = 70
		im.Filter = "0123456789/"
		im.Gray = false
		im.Invert = false
		im.BG = false
	}
}

// ScanAlliance pre-processes the given image file and scans it with tessaract
// into an alliance.Alliance struct
func ScanMainAlliance() (*alliance.Alliance, error) {
	var im scanner.ImageMap
	var a alliance.Alliance
	var d alliance.Data

	// Get the image ready for scanning
	img, err := scanner.SetImageDensity(mainImageFile, 300)
	if err != nil {
		return nil, err
	}
	im.Image = img

	// Setup alliance
	MainImageMapper(&im, "tag")
	tag, err := im.ProcessImageText()
	if err != nil {
		return nil, err
	}

	MainImageMapper(&im, "name")
	name, err := im.ProcessImageText()
	if err != nil {
		return nil, err
	}

	MainImageMapper(&im, "power")
	power, err := im.ProcessImageInt()
	if err != nil {
		return nil, err
	}

	MainImageMapper(&im, "membercount")
	tmembercount, err := im.ProcessImageText()
	if err != nil {
		return nil, err
	}
	tmembercount = strings.Split(tmembercount, "/")[0]
	membercount, err := strconv.Atoi(tmembercount)
	if err != nil {
		return nil, err
	}

	// Gift Level for Main Alliance
	img, err = scanner.SetImageDensity(giftImageFile, 300)
	if err != nil {
		return nil, err
	}

	im.Image = img
	MainImageMapper(&im, "giftlevel")
	giftlevel, err := im.ProcessImageInt()
	if err != nil {
		return nil, err
	}

	d.Date = time.Now().Format(time.DateOnly)
	d.Tag = tag
	d.Name = name
	d.Power = int64(power)
	d.GiftLevel = int64(giftlevel)
	d.MemberCount = int64(membercount)
	a.Data = append(a.Data, d)
	a.Server = mainServer

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
	err = os.WriteFile(mainOutputFile, j, 0644)
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
