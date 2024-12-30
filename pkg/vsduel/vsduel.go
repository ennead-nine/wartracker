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

type VsDuel struct {
	Id            string `json:"id" yaml:"id" db:"id"`
	Date          string `json:"date" yaml:"date" db:"date"`
	League        string `json:"league" yaml:"league" db:"league"`
	Week          int    `json:"week" yaml:"week" db:"week"`
	VsDuelDataMap `json:"vsduel-data" yaml:"vsDuelData"`
}
type VsDay struct {
	Id        string `json:"id" yaml:"id" db:"id"`
	Name      string `json:"name" yaml:"name" db:"name"`
	ShortName string `json:"short-name" yaml:"shortName" db:"short_name"`
	DayOfWeek string `json:"day-of-week" yaml:"dayOfWeek" db:"day_of_week"`
}

type VsDuelData struct {
	Id                 string `json:"id" yaml:"id" db:"id"`
	VsDuelId           string `json:"vsduel-id" yaml:"vsDuelId" db:"vsduel_id"`
	VsDuelDayId        string `json:"vsduel-day-id" yaml:"vsDuelDayId" db:"vsduel-day-id"`
	VsAllianceDataMap  `json:"vsduel-alliance-data" yaml:"vsDuelAllianceData"`
	VsCommanderDataMap `json:"vsduel-commander-data" yaml:"vsDuelCommanderData"`
}

type VsAllianceData struct {
	Points       int    `json:"points" yaml:"points" db:"points"`
	Tag          string `json:"tag" yaml:"tag" db:"tag"`
	AllianceId   string `json:"alliance-id" yaml:"allianceId" db:"alliance_id"`
	VsDuelDataId string `json:"vsduel-data-id" yaml:"vsDuelDataId" db:"vsduel-data-id"`
}

type VsCommanderData struct {
	Points       int    `json:"points" yaml:"points" db:"points"`
	Rank         int    `json:"rank" yaml:"rank" db:"rank"`
	Name         string `json:"name" yaml:"name" db:"name"`
	CommanderId  string `json:"commander-id" yaml:"commanderId" db:"commander_id"`
	VsDuelDataId string `json:"vsduel-data-id" yaml:"vsDuelDataId" db:"vsduel-data-id"`
}

type VsDays map[string]VsDay
type VsDuelDataMap map[string]VsDuelData
type VsCommanderDataMap map[string]VsCommanderData
type VsAllianceDataMap map[string]VsAllianceData

var DayFile string

var (
	ErrDuelDataInsert = fmt.Errorf("failed to insert duel data")
	ErrNumDays        = fmt.Errorf("number of versus days is not 6")
)

func initDays() error {
	var ds []VsDay

	df, err := os.ReadFile(DayFile)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(df, &ds)
	if err != nil {
		return err
	}
	if len(ds) != 6 {
		return ErrNumDays
	}

	for _, d := range ds {
		var w wtid.WTID
		w.New("wartracker", "vsday", 0)

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
			return ErrDuelDataInsert
		}
		err = tx.Commit()
		if err != nil {
			return err
		}
	}

	return nil
}

func GetDays() (VsDays, error) {
	ds := make(VsDays)

	rows, err := db.Connection.Queryx("SELECT * FROM vsduel_day")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var d VsDay
		err = rows.StructScan(&d)
		if err != nil {
			return nil, err
		}
		ds[d.DayOfWeek] = d
	}

	if len(ds) != 6 {
		return nil, ErrNumDays
	}

	return ds, nil
}

func (v *VsDuel) Create() error {
	var w wtid.WTID
	w.New("wartracker", "vsduel", 0)
	v.Id = string(w.Id)

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
		return ErrDuelDataInsert
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	ds, err := GetDays()
	if err != nil && err == sql.ErrNoRows {
		err = initDays()
		if err != nil {
			return err
		}
		ds, err = GetDays()
	}
	if err != nil {
		return err
	}

	return v.initVsDuelData(ds)
}

func (v *VsDuel) initVsDuelData(ds VsDays) error {
	tx, err := db.Connection.Begin()
	if err != nil {
		return err
	}
	for _, d := range ds {
		var w wtid.WTID
		var dd VsDuelData
		w.New("wartracker", "vsdueldata", 0)
		dd.Id = string(w.Id)
		dd.VsDuelDayId = d.Id
		dd.VsDuelId = v.Id

		res, err := tx.Exec("INSERT INTO vsduel_data (id, vsduel_day_id, vsduel_id) VALUES (?, ?, ?)",
			dd.Id,
			dd.VsDuelDayId,
			dd.VsDuelId)
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

	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (v *VsDuel) UpsertAllianceData(did string) error {
	d := v.VsDuelDataMap[did]

	tx, err := db.Connection.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec("DELETE FROM vsduel_alliance WHERE vsduel_data_id=?", did)
	if err != nil {
		return nil
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	tx, err = db.Connection.Begin()
	if err != nil {
		return err
	}
	for _, d := range d.VsAllianceDataMap {
		res, err := tx.Exec("INSERT INTO vsduel_alliance (points, tag, alliance_id, vsduel_data_id) VALUES (?, ?, ?, ?)",
			d.Points,
			d.Tag,
			d.AllianceId,
			did)
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
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (v *VsDuel) UpsertCommanderData(did string) error {
	d := v.VsDuelDataMap[did]

	tx, err := db.Connection.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec("DELETE FROM vsduel_commander WHERE vsduel_data_id=?", did)
	if err != nil {
		return nil
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	tx, err = db.Connection.Begin()
	if err != nil {
		return err
	}
	for _, d := range d.VsCommanderDataMap {
		res, err := tx.Exec("INSERT INTO vsduel_commander (points, rank, name, commander_id, vsduel_data_id) VALUES (?, ?, ?, ?, ?)",
			d.Points,
			d.Rank,
			d.Name,
			d.CommanderId,
			d.VsDuelDataId)
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
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (v *VsDuel) GetCommanderTotal() error {
	return fmt.Errorf("not yet implemented")
}

func (v *VsDuel) GetWeekByDate(d string) error {
	return fmt.Errorf("not yet implemented")
}

func (v *VsDuel) GetById(id string) error {
	err := db.Connection.QueryRowx("SELECT * FROM vsduel WHERE id=?", id).StructScan(v)
	if err != nil {
		return err
	}

	rows, err := db.Connection.Queryx("SELECT * FROM vsduel_data WHERE vsduel_id=?", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		} else {
			return err
		}
	}
	for rows.Next() {
		var d VsDuelData
		err = rows.StructScan(&d)
		if err != nil {
			return err
		}
		v.VsDuelDataMap[d.Id] = d
	}

	for _, d := range v.VsDuelDataMap {
		rows, err = db.Connection.Queryx("SELECT * FROM vsduel_alliance WEHRE vsduel_data_id=?", d.Id)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil
			} else {
				return err
			}
		}
		for rows.Next() {
			var a VsAllianceData
			err = rows.StructScan(&a)
			if err != nil {
				return err
			}
			d.VsAllianceDataMap[d.Id] = a
		}
	}

	for _, d := range v.VsDuelDataMap {
		rows, err = db.Connection.Queryx("SELECT * FROM vsduel_commander WEHRE vsduel_data_id=?", d.Id)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil
			} else {
				return err
			}
		}
		for rows.Next() {
			var c VsCommanderData
			err = rows.StructScan(&c)
			if err != nil {
				return err
			}
			d.VsCommanderDataMap[d.Id] = c
		}
	}

	return nil
}

func (d *VsDuel) DuelToJSON() ([]byte, error) {
	return json.MarshalIndent(d, "", "\t")
}

func (d *VsDuel) DuelToYAML() ([]byte, error) {
	return yaml.Marshal(d)
}
