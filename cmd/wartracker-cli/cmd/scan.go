package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// commanderCmd represents the commander command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan is the parent command for scanning Last War screen shots",
	Long: `The scan command is used to scan screen shots from the Last War 
	game.  There are several data sets that can be scanned.  Use the --help 
	flag to learn more more.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("scan called")
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
}
