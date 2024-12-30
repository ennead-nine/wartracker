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
package vsduelcmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"wartracker/cmd/wartracker-cli/cmd"
	"wartracker/pkg/alliance"
	"wartracker/pkg/scanner"
	"wartracker/pkg/vsduel"

	"github.com/spf13/cobra"
)

// ScanCommander processes the given image file and scans it with tessaract
// into an commander.Commander struct
func ScanVSDuel() (*vsduel.VsDuel, error) {
	var v vsduel.VsDuel
	err := v.GetById(id)
	if err != nil {
		return nil, err
	}

	var did string
	days, err := vsduel.GetDays()
	if err != nil {
		return nil, err
	}
	for id, d := range v.VsDuelDataMap {
		if d.VsDuelDayId == days[dow].Id {
			did = id
		}
	}

	initImageMaps("vsduel")

	img, err := scanner.SetImageDensity(inputFile, 300)
	if err != nil {
		return nil, err
	}

	for k, im := range Imm {
		switch k {
		case "left":
			var a alliance.Alliance
			var ad vsduel.VsAllianceData

			if cmd.Debug {
				fmt.Printf("scanning %s...\n", k)
			}
			np, err := im.ProcessImageText(img)
			if err != nil {
				return nil, err
			}
			names := strings.Split(np, "]")
			tag := names[0][1:]
			var utag string
			fmt.Printf("Correct alliance tag [%s]: ", tag)
			fmt.Scanln(&utag)
			if utag != "" {
				tag = utag
			}
			err = a.GetByTag(tag)
			if err != nil {
				return nil, err
			}
			ad.AllianceId = a.Id
			ad.Tag = tag

			p, err := strconv.Atoi(names[1])
			if err != nil {
				return nil, err
			}
			ad.Points = p
			ad.VsDuelDataId = did

			v.VsDuelDataMap[did].VsAllianceDataMap[a.Id] = ad
		case "right":
			var a alliance.Alliance
			var ad vsduel.VsAllianceData

			if cmd.Debug {
				fmt.Printf("scanning %s...\n", k)
			}
			np, err := im.ProcessImageText(img)
			if err != nil {
				return nil, err
			}
			names := strings.Split(np, "[")
			tag := names[1][:len(names[1])-1]
			var utag string
			fmt.Printf("Correct alliance tag [%s]: ", tag)
			fmt.Scanln(&utag)
			if utag != "" {
				tag = utag
			}
			err = a.GetByTag(tag)
			if err != nil {
				return nil, err
			}
			ad.AllianceId = a.Id
			ad.Tag = tag

			p, err := strconv.Atoi(names[0])
			if err != nil {
				return nil, err
			}
			ad.Points = p

			v.VsDuelDataMap[did].VsAllianceDataMap[a.Id] = ad
		default:
			return nil, fmt.Errorf("invalid key \"%s\" in map configuration", k)
		}
	}

	j, err := v.DuelToJSON()
	if err != nil {
		return nil, err
	}
	err = os.WriteFile(outputFile, j, 0644)
	if err != nil {
		return nil, err
	}

	return &v, nil
}

// scanCmd represents the scan command
var scanAllianceCmd = &cobra.Command{
	Use:   "scanAlliance",
	Short: "Scans an alliance versus duel day screenshot into a json file.",
	Long: `Scan takes an alliance versus deul daily screenshot and pareses it 
	into numbers and text, places it in an aliiance versus data object, then 
	marshals it into json for cleanup.  
	
	Running "wartracker-cli vsduel addAliiance" with the cleaned json 
	will create an alliance versus duel data object in the database.
	
	Example: wartracker-cli commander scan -i commander.png -o commander.json`,
	Run: func(cmd *cobra.Command, args []string) {
		_, err := ScanVSDuel()
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	vsduelCmd.AddCommand(scanAllianceCmd)

	scanAllianceCmd.Flags().StringVarP(&dow, "dow", "d", "", "VS Duel Day")
	scanAllianceCmd.MarkFlagRequired("dow")
}
