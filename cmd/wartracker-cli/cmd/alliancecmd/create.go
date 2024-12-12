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
	"io"
	"os"
	"wartracker/pkg/alliance"

	"github.com/spf13/cobra"
)

func ReadAllianceJSON(a *alliance.Alliance) error {
	in, err := os.Open(infile)
	if err != nil {
		return err
	}
	defer in.Close()

	jf, err := io.ReadAll(in)
	if err != nil {
		return err
	}

	return json.Unmarshal(jf, a)
}

func CreateAlliance() error {
	var a alliance.Alliance

	err := ReadAllianceJSON(&a)
	if err != nil {
		return err
	}

	err = a.GetByTag(a.Data[0].Tag)
	if err != sql.ErrNoRows {
		return fmt.Errorf("alliance [%s] already exists", a.Data[0].Tag)
	}

	return a.Add(a.Server)
}

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an alliance from a JSON file",
	Long: `Builds an alliance object from a JSON file created from "scan".  If 
	the alliance already exists, data is added to the database for today's date.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := CreateAlliance()
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	allianceCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	createCmd.Flags().StringVarP(&infile, "inputfile", "i", "", "JSON file to create an allaince")
}
