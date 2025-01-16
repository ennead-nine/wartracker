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
	"os"
	"wartracker/pkg/vsduel"

	"github.com/spf13/cobra"
)

var (
	date   string
	week   int
	league string
)

func CreateDuel() error {
	var d vsduel.VsDuel

	d.Date = date
	d.League = league
	d.Week = week

	err := d.Create()
	if err != nil {
		return err
	}

	j, err := d.DuelToJSON()
	if err != nil {
		return err
	}
	err = os.WriteFile(outputFile, j, 0644)
	if err != nil {
		return err
	}

	return err
}

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := CreateDuel()
		cobra.CheckErr(err)
	},
}

func init() {
	vsduelCmd.AddCommand(createCmd)

	createCmd.Flags().StringVarP(&date, "date", "d", "", "Date the duel start in MM-DD-YYYY format")
	createCmd.Flags().StringVarP(&league, "league", "l", "", "The level of the duel league")
	createCmd.Flags().IntVarP(&week, "week", "w", 0, "Duel league week")
}
