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
	"wartracker/pkg/commander"

	"github.com/spf13/cobra"
)

func CreateCommander() error {
	var c commander.Commander

	err := ReadCommanderJSON(&c)
	if err != nil {
		return err
	}

	err = c.GetByNoteName(c.NoteName)
	if err == nil {
		return fmt.Errorf("commander \"%s\" already exists", c.NoteName)
	} else {
		if err == sql.ErrNoRows {
			return c.Create(server)
		} else {
			return err
		}
	}
}

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a commander from a JSON file",
	Long: `Builds a commander object from a JSON file created from "scan".  If 
	the alliance already exists, data is added to the database for today's date.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := CreateCommander()
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	commanderCmd.AddCommand(createCmd)
}
