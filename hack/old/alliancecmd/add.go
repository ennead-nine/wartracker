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
	"errors"
	"fmt"
	"wartracker/pkg/alliance"

	"github.com/spf13/cobra"
)

var (
	ErrNoData = errors.New("alliancecmd: no data to add to resource")
)

func AddAllianceData() error {
	var a alliance.Alliance

	err := ReadAllianceJSON(&a)
	if err != nil {
		return err
	}

	for date := range a.DataMap {
		return a.AddData(date)
	}

	return ErrNoData
}

// createCmd represents the create command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Update commander data from a JSON file",
	Long: `Builds an object for an existing commander from a JSON file created 
	from "scan" and adds the data to the database.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := AddAllianceData()
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	allianceCmd.AddCommand(addCmd)
}
