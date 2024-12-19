package vsduel

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"wartracker/pkg/db"
	"wartracker/pkg/wtid"

	"gopkg.in/yaml.v3"
)

type Duel struct {
	Id            string          `json:"id" yaml:"id" db:"id"`
	Date          string          `json:"date" yaml:"date" db:"date"`
	League        string          `json:"league" yaml:"league" db:"league"`
	Week          int64           `json:"week" yaml:"week" db:"week"`
	DuelData      []DuelData      `json:"duel-data" yaml:"duelData"`
	CommanderData []CommanderData `json:"commander-data" yaml:"commanderData"`
}
type Day struct {
	Id        string `json:"id" yaml:"id" db:"id"`
	Name      string `json:"name" yaml:"name" db:"name"`
	ShortName string `json:"short-name" yaml:"shortName" db:"short_name"`
	DayOfWeek string `json:"day-of-week" yaml:"dayOfWeek" db:"day_of_week"`
}

type DuelData struct {
	Points     int64  `json:"points" yaml:"points" db:"points"`
	AllianceID string `json:"alliance-id" yaml:"allianceId" db:"alliance_id"`
	DuelId     string `json:"duel-id" yaml:"duelId" db:"duel_id"`
	DayID      string `json:"day-id" yaml:"dayId" db:"day_id"`
}

type CommanderData struct {
	Points      int64  `json:"points" yaml:"points" db:"points"`
	Rank        int64  `json:"rank" yaml:"rank" db:"rank"`
	DuelId      string `json:"duel-id" yaml:"duelId" db:"duel_id"`
	DayID       string `json:"day-id" yaml:"dayId" db:"day_id"`
	CommanderID string `json:"commander-id" yaml:"commanderId" db:"commander_id"`
}

type Days []Day

var DayFile string

func InitDays() error {
	var ds []Day

	df, err := os.ReadFile(DayFile)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(df, &ds)
	if err != nil {
		return err
	}
	if len(ds) != 6 {
		return fmt.Errorf("did not read all 6 days from days config")
	}

	for _, d := range ds {
		var w wtid.WTID
		w.New("wartracker", "vsday", 0)
		d.Id = w.Id

		tx, err := db.Connection.Begin()
		if err != nil {
			return err
		}
		res, err := tx.Exec("INSERT INTO vsduel_day (id, name, short_name, day_of_week) VALUES (?, ?, ?, ?)",
			d.Id,
			d.Name,
			d.ShortName,
			d.DayOfWeek)
		if err != nil {
			return err
		}
		x, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if x != 1 {
			return fmt.Errorf("failed to insert vsduel day")
		}
		err = tx.Commit()
		if err != nil {
			return err
		}
	}

	return nil
}

func GetDays() (Days, error) {
	var ds Days

	rows, err := db.Connection.Queryx("SELECT * FROM vsduel_day")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var d Day
		err = rows.StructScan(&d)
		if err != nil {
			return nil, err
		}
		ds = append(ds, d)
	}

	return ds, nil
}

func (v *Duel) Create() error {
	var w wtid.WTID
	w.New("wartracker", "vsduel", 0)
	v.Id = w.Id

	tx, err := db.Connection.Begin()
	if err != nil {
		return err
	}

	res, err := tx.Exec("INSERT INTO vsduel (id, date, league, week) VALUES (?, ?, ?, ?)",
		v.Id,
		v.Date,
		v.League,
		v.Week)
	if err != nil {
		return err
	}
	x, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if x != 1 {
		return fmt.Errorf("failed to insert vsduel")
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (v *Duel) Update() error {
	if len(v.DuelData) != 2 {
		return fmt.Errorf("length of duel data for [%s] is not 2 ", v.Id)
	}
	tx, err := db.Connection.Begin()
	if err != nil {
		return err
	}
	for i := 0; i < 2; i++ {
		res, err := tx.Exec("INSERT INTO vsduel_data (points, vsduel_day_id, vsduel_id, alliance_id) VALUES (?, ?, ?, ?)",
			v.DuelData[i].Points,
			v.DuelData[i].DayID,
			v.DuelData[i].DuelId,
			v.DuelData[i].AllianceID)
		if err != nil {
			return err
		}
		x, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if x != 1 {
			return fmt.Errorf("failed to insert commander data")
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (v *Duel) AddDuelData() error {
	if len(v.DuelData) != 2 {
		return fmt.Errorf("data for duel [%s] should have 2 entries", v.Date)
	}
	tx, err := db.Connection.Begin()
	if err != nil {
		return err
	}

	for _, d := range v.DuelData {
		res, err := tx.Exec("INSERT INTO vsduel_data (points, vsduel_day_id, vsduel_id, alliance_id) VALUES (?, ?, ?, ?)",
			d.Points,
			d.DayID,
			d.DuelId,
			d.AllianceID)
		if err != nil {
			return err
		}
		x, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if x != 1 {
			return fmt.Errorf("failed to insert duel data")
		}
		err = tx.Commit()
		if err != nil {
			return err
		}
	}

	return nil
}

func (v *Duel) AddCommanderData(d Day) error {
	if len(v.CommanderData) < 1 {
		return fmt.Errorf("data for duel [%s] is empty", v.Date)
	}
	tx, err := db.Connection.Begin()
	if err != nil {
		return err
	}

	for _, d := range v.CommanderData {
		res, err := tx.Exec("INSERT INTO vsduel_commanders (points, rank, vsduel_id, commander_id, vsduel_day_id) VALUES (?, ?, ?, ?, ?)",
			d.Points,
			d.Rank,
			d.DuelId,
			d.CommanderID,
			d.DayID)
		if err != nil {
			return err
		}
		x, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if x != 1 {
			return fmt.Errorf("failed to insert duel data")
		}
		err = tx.Commit()
		if err != nil {
			return err
		}
	}

	return nil
}

func (v *Duel) GetCommanderTotal() error {
	return fmt.Errorf("not yet implemented")
}

func (v *Duel) GetByDate(d string) error {
	return fmt.Errorf("not yet implemented")
}

func (v *Duel) GetById(id string) error {
	err := db.Connection.QueryRowx("SELECT * FROM vsduel WHERE id=?", id).StructScan(v)
	if err != nil {
		return err
	}

	rows, err := db.Connection.Queryx("SELECT * FROM vsduel_data WHERE vsduel_id=?", id)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	for rows.Next() {
		var d DuelData
		err = rows.StructScan(&d)
		if err != nil {
			return err
		}
		v.DuelData = append(v.DuelData, d)
	}

	rows, err = db.Connection.Queryx("SELECT * FROM vsduel_commanders WHERE vsduel_id=?", id)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	for rows.Next() {
		var d CommanderData
		err = rows.StructScan(&d)
		if err != nil {
			return err
		}
		v.CommanderData = append(v.CommanderData, d)
	}

	return err
}

func (d *Duel) DuelToJSON() ([]byte, error) {
	return json.MarshalIndent(d, "", "\t")
}

func (d *Duel) DuelToYAML() ([]byte, error) {
	return yaml.Marshal(d)
}
