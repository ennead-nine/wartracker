package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

var (
	w int
)

var warzoneCmd = &cobra.Command{
	Use:   "warzone",
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

var warzoneCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		WarzoneCreate()
	},
}

func WarzoneCreate() {
	var data io.Reader
	var err error

	client := &http.Client{}
	url := fmt.Sprintf("http://localhost:3001/warzone/%d?indent=true", w)
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

func init() {
	rootCmd.AddCommand(warzoneCmd)

	warzoneCmd.AddCommand(warzoneCreateCmd)
	warzoneCreateCmd.Flags().IntVarP(&w, "warzone", "w", 0, "warzone server")
	warzoneCreateCmd.MarkFlagRequired("warzone")
}
