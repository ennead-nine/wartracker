package vsduel

import (
	"lastwar/pkg/alliance"
	"lastwar/pkg/commander"
)

type Day struct {
	Id        string `json:"id" yaml:"id"`
	Name      string `json:"name" yaml:"name"`
	ShortName string `json:"short-name" yaml:"shortName"`
	DayOfWeek string `json:"day-of-week" yaml:"dayOfWeek"`
}

type Duel struct {
	Id        string `json:"id" yaml:"id"`
	Date      string `json:"date" yaml:"date"`
	League    string `json:"league" yaml:"league"`
	Week      int64  `json:"week" yaml:"week"`
	Alliance1 alliance.Alliance
	Alliance2 alliance.Alliance
}

type DuelData struct {
	Alliance1Points int64 `json:"alliance1-points" yaml:"alliance1Points"`
	Alliance2Points int64 `json:"alliance2-points" yaml:"alliance2Points"`
	Day             Day
	Duel
}

type CommanderData struct {
	Points int64 `json:"points" yaml:"points"`
	Rank   int64 `json:"rank" yaml:"rank"`
	Duel
	Day       Day
	Commander commander.Commander
}
