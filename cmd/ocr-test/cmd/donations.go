package cmd

import (
	"bytes"
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"wartracker/pkg/common"

	"github.com/spf13/cobra"
)

var (
	x int
	y int
)

func ScanDonation() error {
	b, err := os.ReadFile(inputFile)
	if err != nil {
		return err
	}

	if filepath.Ext(inputFile) == ".jpg" {
		b, err = common.ConvertJpgToPng(b)
		if err != nil {
			return err
		}
	}

	//	b, err = common.ResizeImage(b, x, y)
	//	if err != nil {
	//		return err
	//	}

	var cr common.CropRatios
	cr.Rpx = .38
	cr.Rpy = .22
	cr.Rsx = .4
	cr.Rsy = .58

	nb, err := common.CropImage(b, cr)
	if err != nil {
		return err
	}

	outputFile, err := os.Create("_scratch/names.png")
	if err != nil {
		return err
	}
	defer outputFile.Close()
	img, err := png.Decode(bytes.NewReader(nb))
	if err != nil {
		return err
	}
	err = png.Encode(outputFile, img)
	if err != nil {
		return err
	}

	cr.Rpx = .75
	cr.Rpy = .22
	cr.Rsx = .22
	cr.Rsy = .58

	pb, err := common.CropImage(b, cr)
	if err != nil {
		return err
	}

	outputFile, err = os.Create("_scratch/points.png")
	if err != nil {
		return err
	}
	defer outputFile.Close()
	img, err = png.Decode(bytes.NewReader(pb))
	if err != nil {
		return err
	}
	err = png.Encode(outputFile, img)

	return err
}

var donationsCmd = &cobra.Command{
	Use:   "donations",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		err := ScanDonation()
		fmt.Println(err)
	},
}

func init() {
	rootCmd.AddCommand(donationsCmd)

	donationsCmd.Flags().IntVarP(&x, "height", "x", 0, "height to resize image to")
	donationsCmd.Flags().IntVarP(&y, "width", "y", 0, "width to resize image to")
}
