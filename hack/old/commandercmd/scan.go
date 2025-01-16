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
	"time"
	"wartracker/cmd/wartracker-cli/cmd"
	"wartracker/pkg/alliance"
	"wartracker/pkg/commander"
	"wartracker/pkg/scanner"

	"github.com/spf13/cobra"
)

// ScanCommander processes the given image file and scans it with tessaract
// into an commander.Commander struct
func ScanCommander() (*commander.Commander, error) {
	var c commander.Commander
	var d commander.Data
	var a alliance.Alliance
	var tag string

	initImageMaps("commander")

	img, err := scanner.SetImageDensity(inputFile, 600)
	if err != nil {
		return nil, err
	}

	for k, im := range Imm {
		switch k {
		case "pfp":
		//			if cmd.Debug {
		//				fmt.Printf("scanning %s...\n", k)
		//			}
		//			d.PFP, err = im.ProcessImage(img)
		case "hqlevel":
			if cmd.Debug {
				fmt.Printf("scanning %s...\n", k)
			}
			d.HQLevel, err = im.ProcessImageInt(img)
			if err != nil {
				d.HQLevel = 0
			}
		case "nametag":
			if cmd.Debug {
				fmt.Printf("scanning %s...\n", k)
			}
			nt, err := im.ProcessImageText(img)
			if err != nil {
				return nil, err
			}
			names := strings.Split(nt, "]")
			if len(names) < 2 {
				tag = ""
				c.NoteName = strings.ToLower(names[0])
				d.AllianceId = ""
			} else {
				tag = names[0][1:]
				err = a.GetByTag(tag)
				if err != nil && err != sql.ErrNoRows {
					return nil, err
				}
				d.AllianceId = a.Id
				c.NoteName = strings.ToLower(names[1])
			}
		case "hqpower":
			if cmd.Debug {
				fmt.Printf("scanning %s...\n", k)
			}
			d.HQPower, err = im.ProcessImageAbbrInt(img)
			if err != nil {
				d.HQPower = 0
			}
		case "kills":
			if cmd.Debug {
				fmt.Printf("scanning %s...\n", k)
			}
			d.Kills, err = im.ProcessImageAbbrInt(img)
			if err != nil {
				d.Kills = 0
			}
		case "proflevel":
			if cmd.Debug {
				fmt.Printf("scanning %s...\n", k)
			}
			d.ProfessionLevel, err = im.ProcessImageInt(img)
			if err != nil {
				d.ProfessionLevel = 0
			}
		case "likes":
			if cmd.Debug {
				fmt.Printf("scanning %s...\n", k)
			}
			d.Likes, err = im.ProcessImageInt(img)
			if err != nil {
				d.Likes = 0
			}
		default:
			return nil, fmt.Errorf("invalid key \"%s\" in map configuration", k)
		}
	}

	if a.Id == "" {
		fmt.Printf("Commander's alliance [%s] does not exist in the database\n", tag)
	} else {
		fmt.Printf("Associating commander with [%s]\n", tag)
	}

	err = c.GetByNoteName(c.NoteName)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("Commander does not exist.  Please use \"wartracker-cli commander create -i %s -s [server]\" to create.\n", outputFile)
		} else {
			return nil, err
		}
	} else {
		fmt.Printf("Commander does exists.  Please use \"wartracker-cli commander update -i %s -s [server]\" to update.\n", outputFile)
	}

	d.Date = time.Now().Format(time.DateOnly)
	c.Data[d.Date] = d

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
}
