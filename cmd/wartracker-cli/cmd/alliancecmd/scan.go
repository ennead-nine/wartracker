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
	"time"
	"wartracker/pkg/alliance"
	"wartracker/pkg/scanner"

	"github.com/spf13/cobra"
)

var ()

// ScanAlliance pre-processes the given image file and scans it with tessaract
// into an alliance.Alliance struct
func ScanAlliance() (*alliance.Alliance, error) {
	var a alliance.Alliance
	var d alliance.Data

	initImageMaps("alliance")

	img, err := scanner.SetImageDensity(infile, 300)
	if err != nil {
		return nil, err
	}

	// Setup alliance
	for k, im := range Imm {
		switch k {
		case "tag":
			d.Tag, err = im.ProcessImageText(img)
		case "name":
			d.Name, err = im.ProcessImageText(img)
		case "power":
			d.Power, err = im.ProcessImageInt(img)
		case "giftlevel":
			d.GiftLevel, err = im.ProcessImageInt(img)
		case "membercount":
			d.MemberCount, err = im.ProcessImageInt(img)
		default:
			return nil, fmt.Errorf("invalid key \"%s\" in map configuration", k)
		}
		if err != nil {
			return nil, err
		}
	}

	d.Date = time.Now().Format(time.DateOnly)
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

	if len(a.Data) > 0 {
		a.Data = a.Data[:1]
	}
	a.Data = append(a.Data, d)

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
	scanCmd.Flags().StringVarP(&infile, "image", "i", "", "image file (PNG) to scan for alliance data")
	scanCmd.Flags().StringVarP(&outfile, "output", "o", "", "JSON file to output alliance data to")
	scanCmd.Flags().Int64VarP(&server, "server", "s", 1, "Alliance's server number")

}
