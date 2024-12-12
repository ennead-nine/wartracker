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
	"fmt"
	"io"
	"os"

	"wartracker/cmd/wartracker-cli/cmd"
	"wartracker/pkg/scanner"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// allianceCmd represents the alliance command
var allianceCmd = &cobra.Command{
	Use:   "alliance",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("alliance called")
	},
}

func init() {
	cmd.RootCmd.AddCommand(allianceCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// allianceCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// allianceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

var (
	// Common flag variables
	infile  string
	outfile string
	server  int64
	id      string
	tag     string

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
