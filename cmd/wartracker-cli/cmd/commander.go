package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

func AddData() error {
	return fmt.Errorf("not implemented")
}

func AddAlias() error {
	data, err := os.Open(Input)
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()

	client := &http.Client{}
	url := "http://localhost:3001/commander/a?indent=true"
	req, err := http.NewRequest("PUT", url, data)
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

	j, err := json.MarshalIndent(content, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(j))

	return nil
}

func GetData() error {
	return fmt.Errorf("not implemented")
}
func GetAlias() error {
	return fmt.Errorf("not implemented")
}
func Add() error {
	return fmt.Errorf("not implemented")
}

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
		var err error

		switch args[0] {
		case "add-data":
			err = AddData()
		case "add-alias":
			err = AddAlias()
		case "get-data":
			err = GetData()
		case "get-alias":
			err = GetAlias()
		case "add":
			err = Add()
		}

		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(commanderCmd)
}
