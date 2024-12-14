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
	"wartracker/cmd/wartracker-cli/cmd"
	"wartracker/pkg/alliance"
	"wartracker/pkg/scanner"

	"github.com/spf13/cobra"
)

var (
	giftfile string
)

// ScanAlliance pre-processes the given image file and scans it with tessaract
// into an alliance.Alliance struct
func ScanMainAlliance() (*alliance.Alliance, error) {
	var a alliance.Alliance
	var d alliance.Data

	initImageMaps("alliancemain")

	img, err := scanner.SetImageDensity(infile, 300)
	if err != nil {
		return nil, err
	}

	gimg, err := scanner.SetImageDensity(giftfile, 300)
	if err != nil {
		return nil, err
	}

	// Setup alliance
	for k, im := range Imm {
		switch k {
		case "tag":
			if cmd.Debug {
				fmt.Printf("scanning %s...\n", k)
			}
			d.Tag, err = im.ProcessImageText(img)
		case "name":
			if cmd.Debug {
				fmt.Printf("scanning %s...\n", k)
			}
			d.Name, err = im.ProcessImageText(img)
		case "power":
			if cmd.Debug {
				fmt.Printf("scanning %s...\n", k)
			}
			d.Power, err = im.ProcessImageInt(img)
		case "giftlevel":
			if cmd.Debug {
				fmt.Printf("scanning %s...\n", k)
			}
			d.GiftLevel, err = im.ProcessImageInt(gimg)
		case "membercount":
			if cmd.Debug {
				fmt.Printf("scanning %s...\n", k)
			}
			tmc, err := im.ProcessImageText(img)
			if err != nil {
				return nil, err
			}
			tmc = strings.Split(tmc, "/")[0]
			mc, err := strconv.Atoi(tmc)
			if err != nil {
				return nil, err
			}
			d.MemberCount = int64(mc)
		default:
			return nil, fmt.Errorf("invalid key \"%s\" in map configuration", k)
		}
		if err != nil {
			return nil, err
		}
	}
	d.Date = time.Now().Format(time.DateOnly)
	a.Data = append(a.Data, d)
	a.Server = server

	err = a.GetByTag(d.Tag)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
		fmt.Printf("A new alliance will need to be created from this data.  Please run 'wartracker-cli alliance new -o %s' after verifying the data\n", outfile)
	} else {
		fmt.Printf("This alliance already exists. To add the new data run 'wartracker-cli alliance add -o %s' to add the new data.\n", outfile)
	}

	a.Data = a.Data[:1]

	j, err := json.MarshalIndent(a, "", "\t")
	if err != nil {
		return nil, err
	}
	err = os.WriteFile(outfile, j, 0644)
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

	scanmainCmd.Flags().StringVarP(&infile, "mainimage", "i", "", "Image file with main alliance info (PNG) to scan for alliance data")
	scanmainCmd.Flags().StringVarP(&giftfile, "giftimage", "g", "", "Image file with alliance gift info (PNG) to scan for alliance data")
	scanmainCmd.Flags().StringVarP(&outfile, "output", "o", "", "JSON file to output alliance data to")
	scanmainCmd.Flags().Int64VarP(&server, "server", "s", 0, "Alliance's server number")

}
