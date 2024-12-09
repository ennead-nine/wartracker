package alliance

import (
	"encoding/json"
	"fmt"
	"wartracker/pkg/db"
	"wartracker/pkg/wtid"

	"gopkg.in/yaml.v3"
)

type Alliance struct {
	Id     string `json:"id" yaml:"id" db:"id"`
	Server int64  `json:"server" yaml:"server" db:"server"`
	Data   []Data
}

type Data struct {
	Date        string `json:"date" yaml:"date" db:"date"`
	Name        string `json:"name" yaml:"name" db:"name"`
	Tag         string `json:"tag" yaml:"tag" db:"tag"`
	Power       int64  `json:"power" yaml:"power" db:"power"`
	GiftLevel   int64  `json:"gift-level" yaml:"giftLevel" db:"gift_level"`
	MemberCount int64  `json:"member-count" yaml:"memberCount" db:"member_count"`
	R5Id        string `json:"r5-id" yaml:"r5Id" db:"r5_id"`
	AllianceID  string `json:"alliance-id" yaml:"allianceId" db:"alliance_id"`
}

// Adds and alliance to the database
func (a *Alliance) Add(server int64) error {
	var w wtid.WTID
	w.New("wartracker", "alliance", server)
	a.Id = w.Id
	a.Data[0].AllianceID = a.Id

	tx, err := db.Connection.Begin()
	if err != nil {
		return err
	}
	res, err := tx.Exec("INSERT INTO alliance (id, server) VALUES (?, ?)",
		a.Id,
		a.Server)
	if err != nil {
		return err
	}
	x, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if x != 1 {
		return fmt.Errorf("failed to insert alliance")
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	tx, err = db.Connection.Begin()
	if err != nil {
		return err
	}
	res, err = tx.Exec("INSERT INTO alliance_data (name, tag, date, power, gift_level, member_count, r5_id, alliance_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		a.Data[0].Name,
		a.Data[0].Tag,
		a.Data[0].Date,
		a.Data[0].Power,
		a.Data[0].GiftLevel,
		a.Data[0].MemberCount,
		a.Data[0].R5Id,
		a.Id)
	if err != nil {
		return err
	}
	x, err = res.RowsAffected()
	if err != nil {
		return err
	}
	if x != 1 {
		return fmt.Errorf("failed to insert alliance data")
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// Populates an Alliance struct from the database by ID.  Gets all alliance data, with a.Data[0] being the latest
func (a *Alliance) GetById(id string) error {
	err := db.Connection.QueryRowx("SELECT * FROM alliance WHERE id=?", id).StructScan(a)
	if err != nil {
		return err
	}

	rows, err := db.Connection.Queryx("SELECT * FROM alliance_data WHERE alliance_id=? ORDER BY date DESC", id)

	if err != nil {
		return err
	}
	for rows.Next() {
		var d Data
		err = rows.StructScan(&d)
		if err != nil {
			return err
		}
		a.Data = append(a.Data, d)
	}

	return nil
}

// Populates and Alliance structure, setting a.Data[0] to data from the specified date.  The rest of the a.Data slice will be empty
func (a *Alliance) GetDataByDate(d string) error {
	return nil
}

// Finds the (first) alliance with the given tag, and returns the alliance object with all alliance data
func (a *Alliance) GetByTag(t string) error {
	var id string

	err := db.Connection.QueryRowx("SELECT alliance_id FROM alliance_data WHERE tag=? ORDER BY date DESC LIMIT 1", t).Scan(&id)
	if err != nil {
		return err
	}

	err = a.GetById(id)
	if err != nil {
		return err
	}

	return nil
}

func (a *Alliance) AllianceToJSON() ([]byte, error) {
	return json.MarshalIndent(a, "", "\t")
}

func (a *Alliance) AllianceToYAML() ([]byte, error) {
	return yaml.Marshal(a)
}
