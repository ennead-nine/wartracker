package wtconfig

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Configuration struct {
	Wartracker `yaml:"warTracker"`
	Scanner    `yaml:"Scanner"`
	CLI        `yaml:"cli"`
	UI         `yaml:"ui"`
	API        `yaml:"api"`
}

type CLI struct {
	APIUrl string `yaml:"apiUrl"`
}

type UI struct {
	ListenAddr string `yaml:"listenAddr"`
	APIUrl     string `yaml:"apiUrl"`
}

type API struct {
	ListenAddr string `yaml:"listenAddr"`
}

type Wartracker struct {
	ScratchDir string `yaml:"scratchDir"`
}

type Scanner struct {
	VsDuelNamesCrop  VsDuelCrop `yaml:"vsDuelNamesCrop"`
	VsDuelPointsCrop VsDuelCrop `yaml:"vsDuelNamesCrop"`
}

type VsDuelCrop struct {
	Px int `yaml:"px"`
	Py int `yaml:"py"`
	H  int `yaml:"h"`
	W  int `yaml:"w"`
}

var Config Configuration

func Load(f string) error {
	cf, err := os.ReadFile(f)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(cf, &Config)
}
