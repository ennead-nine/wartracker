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
package commandercmd

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"wartracker/pkg/alliance"
	"wartracker/pkg/commander"
	"wartracker/pkg/scanner"

	"github.com/spf13/cobra"
)

// Variables for flags
var imageFile string
var outputFile string
var server int64

// ImageMapper builds a scanner.ImageMap object with data for the parts of a
// commander image that need to be parsed.
func ImageMapper(im *scanner.ImageMap, field string) {
	switch field {
	case "pfp":
		im.PX = 88
		im.PY = 526
		im.RX = 450
		im.RY = 450
		im.Filter = ""
		im.Gray = false
		im.Invert = false
		im.BG = false
	case "hqlevel":
		im.PX = 177
		im.PY = 1122
		im.RX = 131
		im.RY = 70
		im.Filter = "0123456789"
		im.Gray = true
		im.Invert = true
		im.BG = true
	case "nametag":
		im.PX = 300
		im.PY = 1114
		im.RX = 538
		im.RY = 86
		im.Filter = "[]0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
		im.Gray = true
		im.Invert = true
		im.BG = true
	case "hqpower":
		im.PX = 177
		im.PY = 1237
		im.RX = 178
		im.RY = 78
		im.Filter = "0123456789.KMG"
		im.Gray = false
		im.Invert = false
		im.BG = false
	case "kills":
		im.PX = 498
		im.PY = 1237
		im.RX = 173
		im.RY = 71
		im.Filter = "0123456789.KMG"
		im.Gray = false
		im.Invert = false
		im.BG = false
	case "proflevel":
		im.PX = 798
		im.PY = 1234
		im.RX = 173
		im.RY = 71
		im.Filter = "0123456789"
		im.Gray = false
		im.Invert = false
		im.BG = false
	}
}

// ScanCommander processes the given image file and scans it with tessaract
// into an commander.Commander struct
func ScanCommander() (*commander.Commander, error) {
	var im scanner.ImageMap
	var c commander.Commander
	var d commander.Data
	var a alliance.Alliance

	// Get the image ready for scanning
	img, err := scanner.SetImageDensity(imageFile, 300)
	if err != nil {
		return nil, err
	}
	im.Image = img

	// Get commander's profile pic
	ImageMapper(&im, "pfp")
	pfp, err := im.ProcessImage()
	if err != nil {
		return nil, err
	}

	// Get commander's HQ level
	ImageMapper(&im, "hqlevel")
	hqlevel, err := im.ProcessImageInt()
	if err != nil {
		return nil, err
	}

	// Get commander's name and tag.  Scanning returns "[TAG]Name" and gets
	// parsed into tag and name
	ImageMapper(&im, "nametag")
	nametag, err := im.ProcessImageText()
	if err != nil {
		return nil, err
	}
	names := strings.Split(nametag, "]")
	tag := names[0][1:]
	name := names[1]

	// Get commander's HQ power
	ImageMapper(&im, "hqpower")
	hqpower, err := im.ProcessImageAbbrInt()
	if err != nil {
		return nil, err
	}

	// Get commander's kills
	ImageMapper(&im, "kills")
	kills, err := im.ProcessImageAbbrInt()
	if err != nil {
		return nil, err
	}

	// Get commander's profession level
	ImageMapper(&im, "proflevel")
	proflevel, err := im.ProcessImageInt()
	if err != nil {
		return nil, err
	}

	// Buildout the commander.Commander struct to save to JSON
	d.PFP = pfp
	d.HQLevel = hqlevel
	c.NoteName = name
	d.HQPower = hqpower
	d.Kills = kills
	d.ProfessionLevel = proflevel
	err = a.GetByTag(tag)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
		fmt.Printf("Commander's alliance [%s] does not exist in the database\n", tag)
	} else {
		fmt.Printf("Associating commander with [%s]%s", tag, a.Data[0].Name)
	}
	d.AllianceID = a.Id

	c.Data = append(c.Data, d)

	j, err := c.CommanderToJSON()
	if err != nil {
		return nil, err
	}
	err = os.WriteFile(outputFile, j, 0644)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scans a commander screenshot into a json file.",
	Long: `Scan takes an cammander screenshot pareses it into numbers and 
	text, places it in a commander object, then marshals it into json for 
	cleanup.  Running wartracker-cli commander create with the cleaned json 
	will create an commander object in the database.
	
	Example: wartracker-cli commander scan -i commander.png -o commander.json`,
	Run: func(cmd *cobra.Command, args []string) {
		_, err := ScanCommander()
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	commanderCmd.AddCommand(scanCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// canmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	scanCmd.Flags().StringVarP(&imageFile, "image", "i", "", "image file (PNG) to scan for commander data")
	scanCmd.MarkFlagRequired("image")
	scanCmd.MarkFlagFilename("image")
	scanCmd.Flags().StringVarP(&outputFile, "output", "o", "", "JSON file to output commander data to")
	scanCmd.MarkFlagRequired("output")
	scanCmd.MarkFlagFilename("output")
	scanCmd.Flags().Int64VarP(&server, "server", "s", 1, "Commander's server number")
	scanCmd.MarkFlagRequired("server")
}
