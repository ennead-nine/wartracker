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
	"encoding/json"
	"fmt"
	"io"
	"os"

	"wartracker/pkg/alliance"
	"wartracker/pkg/scanner"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var (
	inputFile  string
	outputFile string
	id         string
	tag        string
	server     int

	Imm scanner.ImageMaps
)

// allianceCmd represents the alliance command

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

func ReadAllianceJSON(a *alliance.Alliance) error {
	in, err := os.Open(inputFile)
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