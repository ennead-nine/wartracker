package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"wartracker/pkg/alliance"

	"github.com/spf13/cobra"
)

var (
	a           alliance.Alliance
	warzone     int
	date        string
	name        string
	power       int64
	giftLevel   int
	memberCount int
)

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

var allianceCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		AllianceCreate()
	},
}

func AllianceCreate() {
	var data io.Reader
	var err error
	var ad alliance.Data

	if Input == "" {
		if a.WarzoneId == "" {
			log.Fatal("get by warzone by number not yet implemented")
		}
		ad.Date = date
		ad.Name = name
		ad.GiftLevel = giftLevel
		ad.Power = power
		ad.MemberCount = memberCount

		adm := make(alliance.DataMap)
		adm[date] = ad
		a.DataMap = adm

		j, err := json.Marshal(a)
		if err != nil {
			log.Fatal(err)
		}
		data = bytes.NewReader(j)
	} else {
		data, err = os.Open(Input)
		if err != nil {
			log.Fatal(err)
		}
	}

	client := &http.Client{}
	url := "http://localhost:3001/alliance?indent=true"
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
	rootCmd.AddCommand(allianceCmd)

	allianceCmd.AddCommand(allianceCreateCmd)
	allianceCreateCmd.Flags().StringVarP(&Input, "input", "i", "", "input file")
	allianceCreateCmd.Flags().StringVarP(&a.WarzoneId, "warzone-id", "w", "", "warzone id")
	allianceCreateCmd.Flags().IntVarP(&warzone, "warzone", "z", 0, "warzone number")
	allianceCreateCmd.MarkFlagsOneRequired("input", "warzone-id", "warzone")
	allianceCreateCmd.MarkFlagsMutuallyExclusive("input", "warzone-id")
	allianceCreateCmd.MarkFlagsMutuallyExclusive("input", "warzone")
	allianceCreateCmd.MarkFlagsMutuallyExclusive("warzone-id", "warzone")
	allianceCreateCmd.Flags().StringVarP(&a.Tag, "tag", "t", "", "alliance tag")
	allianceCreateCmd.MarkFlagsOneRequired("input", "tag")
	allianceCreateCmd.MarkFlagsMutuallyExclusive("input", "tag")
	allianceCreateCmd.Flags().StringVarP(&date, "date", "d", "", "alliance data date")
	allianceCreateCmd.MarkFlagsOneRequired("input", "date")
	allianceCreateCmd.MarkFlagsMutuallyExclusive("input", "date")
	allianceCreateCmd.Flags().StringVarP(&name, "name", "n", "", "alliance data name")
	allianceCreateCmd.MarkFlagsOneRequired("input", "date")
	allianceCreateCmd.MarkFlagsMutuallyExclusive("input", "date")
	allianceCreateCmd.Flags().Int64VarP(&power, "power", "p", 0, "alliance data power")
	allianceCreateCmd.MarkFlagsOneRequired("input", "power")
	allianceCreateCmd.MarkFlagsMutuallyExclusive("input", "power")
	allianceCreateCmd.Flags().IntVarP(&giftLevel, "gift-level", "g", 0, "alliance data gift level")
	allianceCreateCmd.MarkFlagsOneRequired("input", "gift-level")
	allianceCreateCmd.MarkFlagsMutuallyExclusive("input", "gift-level")
	allianceCreateCmd.Flags().IntVarP(&memberCount, "member-count", "c", 0, "alliance data member count")
	allianceCreateCmd.MarkFlagsOneRequired("input", "member-count")
	allianceCreateCmd.MarkFlagsMutuallyExclusive("input", "member-count")
}
