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
package cmd

import (
	"fmt"
	"os"
	"wartracker/pkg/db"
	"wartracker/pkg/scanner"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	//Configuration variables
	ConfigFile  string
	DBFile      string
	ScratchDir  string
	Debug       bool
	TessdataDir string
	Languages   []string
	MapDir      string
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "wartracker-cli",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&ConfigFile, "config", "", "config file (default is $HOME/.wartracker-cli.yaml)")
	RootCmd.PersistentFlags().StringVar(&DBFile, "dbfile", "", "database file")
	RootCmd.PersistentFlags().StringVar(&ScratchDir, "scratch", "", "Directory to store scratch files")
	RootCmd.PersistentFlags().BoolVar(&Debug, "debug", false, "Directory to store scratch files")
	RootCmd.PersistentFlags().StringVar(&TessdataDir, "tessdata", "", "Tesseract data directory")
	RootCmd.PersistentFlags().StringSliceVar(&Languages, "languages", nil, "Tesseract languages to use")
	RootCmd.PersistentFlags().StringVar(&MapDir, "mapdir", "", "Directory for command maps")
	cobra.CheckErr(viper.BindPFlag("dbfile", RootCmd.PersistentFlags().Lookup("dbfile")))
	cobra.CheckErr(viper.BindPFlag("scratch", RootCmd.PersistentFlags().Lookup("scratch")))
	cobra.CheckErr(viper.BindPFlag("debug", RootCmd.PersistentFlags().Lookup("debug")))
	cobra.CheckErr(viper.BindPFlag("tessdata", RootCmd.PersistentFlags().Lookup("tessdata")))
	cobra.CheckErr(viper.BindPFlag("languages", RootCmd.PersistentFlags().Lookup("languages")))
	cobra.CheckErr(viper.BindPFlag("mapdir", RootCmd.PersistentFlags().Lookup("mapdir")))
	viper.SetDefault("dbfile", "db/wartracker.db")
	viper.SetDefault("scratch", "_scratch")
	viper.SetDefault("debug", false)
	viper.SetDefault("tessdata", "/Users/erumer/src/github.com/tesseract-ocr/tessdata")
	viper.SetDefault("languages", nil)
	viper.SetDefault("mapdir", "config/scanner/maps")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if ConfigFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(ConfigFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".wartracker-cli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".wartracker-cli")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	initDB()
	initScratch()
	initDebug()
	initLanguages()
	initMapDir()
	initTessdataDir()
}

func initDB() {
	var err error
	db.Connection, err = db.Connect(viper.GetString("dbfile"))
	if err != nil {
		panic(err)
	}
}

func initScratch() {
	Scratch := viper.GetString("scratch")
	err := os.RemoveAll(Scratch)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to initialize scratch directory: ", Scratch)
	}
	err = os.MkdirAll(Scratch, 0755)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to initialize scratch directory: ", Scratch)
	}
}

func initDebug() {
	if viper.GetBool("debug") {
		scanner.Debug = true
		scanner.Process = os.Getpid()
		scanner.ScratchDir = viper.GetString("scratch")
	} else {
		scanner.Debug = false
	}
}

func initTessdataDir() {
	TessdataDir = viper.GetString("tessdata")
	scanner.TessdataDir = TessdataDir
}

func initLanguages() {
	Languages = viper.GetStringSlice("languages")
	scanner.Languages = Languages
}

func initMapDir() {
	MapDir = viper.GetString("mapdir")
}
