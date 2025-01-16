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
	"archive/zip"
	"database/sql"
	"fmt"
	"image"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
	"wartracker/cmd/wartracker-cli/cmd"
	"wartracker/pkg/commander"
	"wartracker/pkg/scanner"
	"wartracker/pkg/vsduel"

	"github.com/spf13/cobra"
)

var (
	zipfile string
)

func ExtractZip(dir string) ([]string, error) {
	var files []string

	fmt.Printf("Unzipping to %s\n", dir)

	archive, err := zip.OpenReader(zipfile)
	if err != nil {
		return nil, err
	}
	defer archive.Close()

	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return nil, err
	}

	for _, f := range archive.File {
		dstPath := filepath.Join(dir, f.Name)
		if f.FileInfo().IsDir() {
			// Create the directory
			fmt.Println("creating directory...")
			if err := os.MkdirAll(dstPath, os.ModePerm); err != nil {
				return nil, err
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(dstPath), os.ModePerm); err != nil {
			return nil, err
		}
		files = append(files, dstPath)
		dstFile, err := os.OpenFile(dstPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return nil, err
		}
		srcFile, err := f.Open()
		if err != nil {
			return nil, err
		}
		_, err = io.Copy(dstFile, srcFile)
		if err != nil {
			return nil, err
		}
		dstFile.Close()
		srcFile.Close()
	}

	slices.Sort(files)

	return files, nil
}
func ExecScan(img image.Image) error {
	return nil
}

func ScanDuelCommander() (*vsduel.VsDuel, error) {
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

	initImageMaps("vsduelcommander")

	outdir := filepath.Join(cmd.ScratchDir, "vsduel")
	files, err := ExtractZip(outdir)
	if err != nil {
		return nil, err
	}

	n := 0
	last := false
	first := true
	for j, f := range files {
		img, err := scanner.SetImageDensity(f, 300)
		if err != nil {
			return nil, err
		}
		if j == len(files)-1 {
			last = true
		}
		for i := 0; i < 7; i++ {
			var c commander.Commander
			var cd vsduel.VsCommanderData
			for k, im := range Imm {
				if cmd.Debug {
					time.Sleep(5 * time.Second)
				}
				switch k {
				case "rank":
					var imt scanner.ImageMap
					imt = im
					if last {
						imt.PY = imt.PY + 42 + (185 * i)
					} else {
						imt.PY = imt.PY + (174 * i)
					}
					if cmd.Debug {
						fmt.Printf("scanning %s...\n", k)
					}
					if first {
						imt.BG = true
						imt.Invert = true
					}
					cd.Rank, err = imt.ProcessImageInt(img)
					if err != nil {
						cd.Rank = 0
					}
				case "points":
					var imt scanner.ImageMap
					imt = im
					if last {
						imt.PY = imt.PY + 42 + (185 * i)
					} else {
						imt.PY = imt.PY + (174 * i)
					}
					if cmd.Debug {
						fmt.Printf("scanning %s...\n", k)
					}
					cd.Points, err = imt.ProcessImageInt(img)
					if err != nil {
						cd.Points = 0
					}
				case "name":
					var imt scanner.ImageMap
					imt = im
					if last {
						imt.PY = imt.PY + 42 + (185 * i)
					} else {
						imt.PY = imt.PY + (174 * i)
					}
					if cmd.Debug {
						fmt.Printf("scanning %s...\n", k)
					}
					name, err := imt.ProcessImageText(img)
					if err != nil {
						return nil, err
					}
					err = c.GetByNoteName(name)
					if err != nil {
						if err == sql.ErrNoRows {
							cd.CommanderId = ""
						} else {
							return nil, err
						}
					} else {
						cd.CommanderId = c.Id
					}
					cd.Name = strings.ToLower(name)
				default:
					return nil, fmt.Errorf("invalid key \"%s\" in map configuration", k)
				}
				cd.VsDuelDataId = did
				v.VsDuelDataMap[did].VsCommanderDataMap[cd.CommanderId] = cd
			}
			n++
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
var scancommanderCmd = &cobra.Command{
	Use:   "scancommander",
	Short: "Scans a commander screenshot into a json file.",
	Long: `Scan takes an cammander screenshot pareses it into numbers and 
	text, places it in a commander object, then marshals it into json for 
	cleanup.  Running wartracker-cli commander create with the cleaned json 
	will create an commander object in the database.
	
	Example: wartracker-cli commander scan -i commander.png -o commander.json`,
	Run: func(cmd *cobra.Command, args []string) {
		_, err := ScanDuelCommander()
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	vsduelCmd.AddCommand(scancommanderCmd)

	scancommanderCmd.Flags().StringVarP(&dow, "dow", "d", "", "VS Duel Day")
	scancommanderCmd.MarkFlagRequired("dow")
}
