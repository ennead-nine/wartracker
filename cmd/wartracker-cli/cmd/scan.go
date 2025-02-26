package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

func Scan(f string) {
	data, err := os.Open(f)
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()

	client := &http.Client{}
	url := "http://localhost:3001/vsduel/scan/1/Monday?indent=true"
	req, err := http.NewRequest("POST", url, data)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", string(content))
}

// commanderCmd represents the commander command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan is the parent command for scanning Last War screen shots",
	Long: `The scan command is used to scan screen shots from the Last War 
	game.  There are several data sets that can be scanned.  Use the --help 
	flag to learn more more.`,
	Run: func(cmd *cobra.Command, args []string) {
		Scan(args[0])
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
}
