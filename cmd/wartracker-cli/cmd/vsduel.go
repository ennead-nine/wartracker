package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"wartracker/pkg/vsduel"

	"github.com/spf13/cobra"
)

var (
	v vsduel.VsDuel

	vw           vsduel.Week
	allianceTags []string

	vd vsduel.Data
)

func VsDuelDataScan(f string) {
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

var vsduelCmd = &cobra.Command{
	Use:   "vsduel",
	Short: "vsduel is the top level command for vsduel operations",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("vsduel called")
	},
}

var vsduelCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "create vsduel objects",
	Long:  `create vsduel objects`,
	Run: func(cmd *cobra.Command, args []string) {
		VsDuelCreate()
	},
}

func VsDuelCreate() {
	var data io.Reader
	var err error

	if Input == "" {
		j, err := json.Marshal(v)
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
	url := "http://localhost:3001/vsduel?indent=true"
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

var vsduelWeekCmd = &cobra.Command{
	Use:   "week",
	Short: "vsduel week object CRUD",
	Long:  `vsduel week object CRUD`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("vsduel week called")
	},
}

var vsduelWeekCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "create vsduel week objects",
	Long:  `create vsduel week objects`,
	Run: func(cmd *cobra.Command, args []string) {
		VsDuelWeekCreate()
	},
}

func VsDuelWeekCreate() {
	var data io.Reader
	var err error

	if Input == "" {
		if vw.AllianceIds == nil {
			log.Fatal("get by alliance tag is not yet implemented")
		}
		j, err := json.Marshal(vw)
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
	url := "http://localhost:3001/vsduel/week?indent=true"
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

var vsduelDataCmd = &cobra.Command{
	Use:   "data",
	Short: "vsduel data object CRUD",
	Long:  `vsduel data object CRUD`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("vsduel data called")
	},
}

var vsduelDataCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "create vsduel data objects",
	Long:  `create vsduel data objects`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("vsduel week create called")
	},
}

func init() {
	rootCmd.AddCommand(vsduelCmd)
	vsduelCmd.Flags().StringVarP(&v.Id, "id", "d", "", "vsduel id")

	vsduelCmd.AddCommand(vsduelCreateCmd)
	vsduelCreateCmd.Flags().StringVarP(&Input, "input", "i", "", "input file")
	vsduelCreateCmd.Flags().StringVarP(&v.Date, "date", "t", "", "vsduel date")
	vsduelCreateCmd.MarkFlagsOneRequired("input", "date")
	vsduelCreateCmd.MarkFlagsMutuallyExclusive("input", "date")
	vsduelCreateCmd.Flags().StringVarP(&v.LeagueLevel, "league-level", "l", "Diamond", "vsduel league level")
	vsduelCreateCmd.MarkFlagsOneRequired("input", "league-level")
	vsduelCreateCmd.MarkFlagsMutuallyExclusive("input", "league-level")
	vsduelCreateCmd.Flags().StringVarP(&v.LeagueId, "league-id", "g", "", "vsduel league id")
	vsduelCreateCmd.MarkFlagsOneRequired("input", "league-id")
	vsduelCreateCmd.MarkFlagsMutuallyExclusive("input", "league-id")
	vsduelCreateCmd.Flags().StringVarP(&v.TournamentId, "tournament-id", "u", "", "vsduel tournament id")
	vsduelCreateCmd.MarkFlagsOneRequired("input", "tournament-id")
	vsduelCreateCmd.MarkFlagsMutuallyExclusive("input", "tournament-id")

	vsduelCmd.AddCommand(vsduelWeekCmd)
	vsduelWeekCmd.PersistentFlags().StringVarP(&vw.Id, "vsduel-week-id", "w", "", "vsduel week id")

	vsduelWeekCmd.AddCommand(vsduelWeekCreateCmd)
	vsduelWeekCreateCmd.Flags().StringVarP(&Input, "input", "i", "", "input file")
	vsduelWeekCreateCmd.Flags().StringVarP(&vw.VsDuelId, "vsduel-id", "d", "", "vsduel id")
	vsduelWeekCreateCmd.Flags().StringArrayVarP(&vw.AllianceIds, "alliance-ids", "a", nil, "alliance ids")
	vsduelWeekCreateCmd.Flags().StringArrayVarP(&allianceTags, "alliance-tags", "m", nil, "alliance tags")
	vsduelWeekCreateCmd.MarkFlagsOneRequired("input", "alliance-ids", "alliance-tags")
	vsduelWeekCreateCmd.MarkFlagsMutuallyExclusive("alliance-ids", "alliance-tags")
	vsduelWeekCreateCmd.MarkFlagsMutuallyExclusive("input", "alliance-ids")
	vsduelWeekCreateCmd.MarkFlagsMutuallyExclusive("input", "alliance-tags")
	vsduelWeekCreateCmd.Flags().IntVarP(&vw.WeekNumber, "week-number", "n", 0, "vsduel week id")
	vsduelWeekCreateCmd.MarkFlagsOneRequired("input", "week-number")
	vsduelWeekCreateCmd.MarkFlagsMutuallyExclusive("input", "week-number")

	vsduelCmd.AddCommand(vsduelDataCmd)
	vsduelDataCmd.Flags().StringVarP(&vd.Id, "data-id", "x", "", "vsduel data id")
	vsduelDataCmd.Flags().StringVarP(&vd.DayOfWeek, "day-of-week", "o", "", "vsduel data day of week")
	vsduelDataCmd.MarkFlagsOneRequired("data-id", "day-of-week-id")
	vsduelDataCmd.MarkFlagsMutuallyExclusive("data-id", "day-of-week-id")
	vsduelDataCmd.Flags().StringVarP(&vd.WeekId, "week-id", "w", "", "vsduel week id")
	vsduelDataCmd.MarkFlagRequired("week-id")

	vsduelDataCmd.AddCommand(vsduelDataCreateCmd)
}
