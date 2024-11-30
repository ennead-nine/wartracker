package alliance

import (
	"lastwar/pkg/commander"
)

type Alliance struct {
	Id     string `json:"id" yaml:"id"`
	Server int64  `json:"server" yaml:"server"`
	AllianceData
}

type AllianceData struct {
	Date        string              `json:"date" yaml:"date"`
	Name        string              `json:"name" yaml:"name"`
	Tag         string              `json:"tag" yaml:"tag"`
	Power       int64               `json:"power" yaml:"power"`
	GiftLevel   int64               `json:"gift-level" yaml:"giftLevel"`
	MemberCount int64               `json:"member-count" yaml:"memberCount"`
	R5          commander.Commander `json:"r5" yaml:"r5"`
}

// Adds and alliance to the database
func (a *Alliance) AddAlliance() error {
	return nil
}

// Populates an Alliance struct from the database by ID.  Uses the latest AllianceData
func (a *Alliance) GetAllianceById(id string) error {
	return nil
}

// Populates and Alliance structure using AllianceData from the specified date
func (a *Alliance) GetAllianceDataByDate(d string) error {
	return nil
}

// Populates and Alliance structure using the latest AllianceData where the name matches
func (a *Alliance) GetAllianceDataByName(n string) error {
	return nil
}

// Returns an Alliance's entire data history
func (a *Alliance) GetAllianceDataHistory() ([]AllianceData, error) {
	return nil, nil
}
