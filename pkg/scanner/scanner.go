package scanner

import (
	"fmt"
	"wartracker/pkg/vsduel"
)

var Scores []vsduel.CommanderData

func VSDuelCommanders(day string, zipfile string) error {
	fmt.Printf("VSDuelCommanders called with Day: %s and Zipfile: %s.\n", day, zipfile)
	return nil
}
