/*
Copyright © 2024 P4K Ennead  <ennead.tbc@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package alliance

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"wartracker/pkg/db"
	"wartracker/pkg/warzone"
	"wartracker/pkg/wtid"

	"gopkg.in/yaml.v3"
)

type Alliance struct {
	Id        string `json:"id" yaml:"id" db:"id"`
	WarzoneId string `json:"warzone-id" yaml:"warzoneId" db:"warzone_id"`
	Tag       string `json:"tag" yaml:"tag" db:"tag"`
	DataMap   `json:"data" yaml:"data"`
}

type Data struct {
	Date        string `json:"date" yaml:"date" db:"date"`
	Name        string `json:"name" yaml:"name" db:"name"`
	Power       int64  `json:"power" yaml:"power" db:"power"`
	GiftLevel   int    `json:"gift-level" yaml:"giftLevel" db:"gift_level"`
	MemberCount int    `json:"member-count" yaml:"memberCount" db:"member_count"`
	AllianceId  string `json:"alliance-id" yaml:"allianceId" db:"alliance_id"`
}

type DataMap map[string]Data
type AllianceMap map[string]Alliance

var (
	ErrAllianceExists     = errors.New("alliance already exists")
	ErrAllianceDataExists = errors.New("alliance data for date already exists")
	ErrAllianceNotFound   = errors.New("alliance not found")
)

// Create adds and alliance resource to the database
func (a *Alliance) Create() error {
	var z warzone.Warzone
	z.Id = a.WarzoneId
	err := z.Get()
	if err != nil {
		return fmt.Errorf("error looking up warzone %s: %w", a.WarzoneId, err)
	}

	var w wtid.WTID
	w.New("wartracker", "alliance", z.Server)
	a.Id = string(w.Id)

	tx, err := db.Connection.Begin()
	if err != nil {
		return err
	}
	res, err := tx.Exec("INSERT INTO alliance (id, tag, warzone_id) VALUES (?, ?, ?)",
		a.Id,
		a.Tag,
		a.WarzoneId)
	if err != nil {
		tx.Rollback()
		return err
	}
	x, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}
	if x != 1 {
		tx.Rollback()
		return fmt.Errorf("failed to insert alliance %s: %w", a.Tag, db.ErrDbErrorUnknown)
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return a.AddAlias(a.Tag, true)
}

func (a *Alliance) AddAlias(n string, p bool) error {
	if p {
		tx, err := db.Connection.Begin()
		if err != nil {
			return err
		}
		_, err = tx.Exec("UPDATE alliance_alias SET preferred = false WHERE alliance_id = ?", a.Id)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update preferred alias for %s: %w", a.Id, err)
		}
		err = tx.Commit()
		if err != nil {
			return err
		}
	}
	tx, err := db.Connection.Begin()
	if err != nil {
		return err
	}
	res, err := tx.Exec(`INSERT INTO alliance_alias (
		alias, 
		preferred,
		alliance_id 
		) VALUES (?, ?, ?)`,
		n, p, a.Id)
	if err != nil {
		tx.Rollback()
		return err
	}
	x, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}
	if x != 1 {
		return fmt.Errorf("failed to insert alias: unknown error")
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// Update updates an alliance resource in the database
func (a *Alliance) Update() error {
	var z warzone.Warzone
	z.Id = a.WarzoneId
	err := z.Get()
	if err != nil {
		return fmt.Errorf("error looking up warzone %s: %w", a.WarzoneId, err)
	}

	tx, err := db.Connection.Begin()
	if err != nil {
		return err
	}
	res, err := tx.Exec("UPDATE alliance SET tag=?, warzone_id=? WHERE id=?",
		a.Tag,
		a.WarzoneId,
		a.Id)
	if err != nil {
		tx.Rollback()
		return err
	}
	x, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}
	if x != 1 {
		tx.Rollback()
		return fmt.Errorf("unable to update alliance [%s]: %w", a.Tag, db.ErrDbErrorUnknown)
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// AddData adds data to an alliance resource in the database
func (a *Alliance) AddData(date string, d Data) error {
	a.DataMap[date] = d

	tx, err := db.Connection.Begin()
	if err != nil {
		return err
	}
	res, err := tx.Exec("INSERT INTO alliance_data (name, date, power, gift_level, member_count, alliance_id) VALUES (?, ?, ?, ?, ?, ?)",
		a.DataMap[date].Name,
		a.DataMap[date].Date,
		a.DataMap[date].Power,
		a.DataMap[date].GiftLevel,
		a.DataMap[date].MemberCount,
		a.Id)
	if err != nil {
		tx.Rollback()
		return err
	}
	x, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}
	if x != 1 {
		tx.Rollback()
		return fmt.Errorf("unable to add alliance data to alliance %s: %w", d.Name, db.ErrDbErrorUnknown)
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return a.AddAlias(a.DataMap[date].Name, false)
}

// List geta all alliances
func List() (AllianceMap, error) {
	var as = make(AllianceMap)

	rows, err := db.Connection.Queryx("SELECT * FROM alliance")

	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var a Alliance
		err = rows.StructScan(&a)
		if err != nil {
			return nil, err
		}
		as[a.Id] = a
	}

	return as, nil
}

// Get gets an alliance resource from the database
func (a *Alliance) Get() error {
	err := db.Connection.QueryRowx("SELECT * FROM alliance WHERE id=?", a.Id).StructScan(a)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrAllianceNotFound
		} else {
			return fmt.Errorf("unable to retrive alliance %s: %w", a.Id, err)
		}
	}

	return nil
}

// GetData retrives the data sets related to an alliance resource
func (a *Alliance) GetData() error {
	a.DataMap = make(DataMap)

	rows, err := db.Connection.Queryx("SELECT * FROM alliance_data WHERE alliance_id=? ORDER BY date DESC", a.Id)
	if err != nil {
		return fmt.Errorf("failed to get data for alliance %s: %w", a.Id, err)
	}
	for rows.Next() {
		var d Data
		err = rows.StructScan(&d)
		if err != nil {
			return fmt.Errorf("failed to get data for alliance %s: %w", a.Id, err)
		}
		a.DataMap[d.Date] = d
	}

	return nil
}

// GetLatestData retrives the latest data set related to an alliance resource
func (a *Alliance) GetLatestData() error {
	a.DataMap = make(DataMap)
	var d Data

	err := db.Connection.QueryRowx("SELECT * FROM alliance_data WHERE alliance_id=? ORDER BY date DESC LIMIT 1", a.Id).StructScan(&d)
	if err != nil {
		return fmt.Errorf("failed to get the latest data for alliance %s: %w", a.Id, err)
	}
	a.DataMap[d.Date] = d

	return nil
}

func (a *Alliance) GetDataByDate(date string) error {
	a.DataMap = make(DataMap)
	var d Data

	err := db.Connection.QueryRowx("SELECT * FROM alliance_data WHERE date=?", date).StructScan(&d)
	if err != nil {
		return fmt.Errorf("failed to get the alliance data for alliance %s from %s: %w", a.Id, date, err)
	}
	a.DataMap[d.Date] = d

	return nil
}

// Finds the (first) alliance with the given tag, and returns the alliance object with all alliance data
func (a *Alliance) GetByTag(t string) error {
	err := db.Connection.QueryRowx("SELECT id FROM alliance WHERE tag=?", t).Scan(&a.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrAllianceNotFound
		} else {
			return fmt.Errorf("unable to retrive alliance by tag [%s]: %w", t, err)
		}
	}

	err = a.Get()
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

func SplitTagName(s string) ([]string, error) {
	r := regexp.MustCompile(`^\[(.*)\] (.*)$`)
	m := r.FindAllStringSubmatch(s, -1)
	if m == nil {
		return nil, fmt.Errorf("no tag in string: \"%s\"", s)
	}
	if len(m[0]) != 3 {
		return nil, fmt.Errorf("could not split alliance name: %s", s)
	}
	return m[0][1:], nil
}
