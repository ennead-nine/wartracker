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
	"encoding/json"
	"fmt"
	"io"
	"os"
	"wartracker/cmd/wartracker-cli/cmd"
	"wartracker/pkg/commander"
	"wartracker/pkg/scanner"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// commanderCmd represents the commander command
var commanderCmd = &cobra.Command{
	Use:   "commander",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("commander called")
	},
}

func init() {
	cmd.RootCmd.AddCommand(commanderCmd)

	commanderCmd.PersistentFlags().StringVar(&inputFile, "inputfile", "", "File with vs day data")
	commanderCmd.MarkPersistentFlagFilename("inputfile")
	commanderCmd.PersistentFlags().StringVar(&outputFile, "outputfile", "", "YAML file with vs day data")
	commanderCmd.MarkPersistentFlagFilename("outputfile")
	commanderCmd.PersistentFlags().StringVar(&id, "id", "", "Resource ID")
	commanderCmd.PersistentFlags().StringVar(&name, "name", "", "Resource ID")
	commanderCmd.PersistentFlags().IntVar(&server, "server", 0, "Resource ID")
}

var (
	inputFile  string
	outputFile string
	id         string
	name       string
	server     int

	Imm scanner.ImageMaps
)

func initImageMaps(m string) {
	Imm = make(scanner.ImageMaps)

	in, err := os.Open(fmt.Sprintf("%s/%s.yaml", viper.GetString("mapdir"), m))
	if err != nil {
		fmt.Println(err)
	}
	defer in.Close()

	yf, err := io.ReadAll(in)
	if err != nil {
		fmt.Println(err)
	}

	err = yaml.Unmarshal(yf, Imm)
	if err != nil {
		fmt.Println(err)
	}
}

func ReadCommanderJSON(c *commander.Commander) error {
	in, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer in.Close()

	jf, err := io.ReadAll(in)
	if err != nil {
		return err
	}

	return json.Unmarshal(jf, c)
}
