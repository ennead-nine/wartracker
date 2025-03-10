package site

import (
	"encoding/json"
	"fmt"
	"os"
)

type Site struct {
	Title string `json:"title" yaml:"title"`
	Icon  string `json:"icon" yaml:"icon"`
}

func NewSite(siteFile string) *Site {
	var s Site

	in, err := os.ReadFile(siteFile)
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal(in, &s)
	if err != nil {
		fmt.Println(err)
	}

	return &s
}
