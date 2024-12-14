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
	"fmt"
	"wartracker/pkg/commander"

	"github.com/spf13/cobra"
)

func GetCommander() error {
	var c commander.Commander
	var o []byte
	var err error

	if id != "" {
		err = c.GetById(id)
	} else if notename != "" {
		err = c.GetByNoteName(notename)
	}
	if err != nil {
		return err
	}

	o, err = c.CommanderToJSON()
	if err != nil {
		return err
	}
	fmt.Println(string(o))

	return nil
}

// scanCmd represents the scan command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Scans an alliance screenshot into a json file.",
	Long: `Scan takes an alliance screenshot and Marshals an alliance object 
	into json for cleanup.  Running wartracjer-cli alliance create with the 
	cleaned json will create an alliance object in the database.
	
	Example: wartracker-cli alliance scan -i alliance.png`,
	Run: func(cmd *cobra.Command, args []string) {
		err := GetCommander()
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	commanderCmd.AddCommand(getCmd)

	getCmd.Flags().StringVarP(&id, "id", "i", "", "Commander's wartracker ID")
	getCmd.Flags().StringVarP(&notename, "notename", "n", "", "Commander's in game notename")
	getCmd.MarkFlagsOneRequired("id", "notename")
	getCmd.MarkFlagsMutuallyExclusive("id", "notename")
}
