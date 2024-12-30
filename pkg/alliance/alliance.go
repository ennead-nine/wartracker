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
	"image"
	"strconv"
	"strings"
	"time"
	"wartracker/cmd/wartracker-cli/cmd"
	"wartracker/pkg/db"
	"wartracker/pkg/scanner"
	"wartracker/pkg/wtid"

	"gopkg.in/yaml.v3"
)

type Alliance struct {
	Id      string `json:"id" yaml:"id" db:"id"`
	Server  int    `json:"server" yaml:"server" db:"server"`
	Tag     string `json:"tag" yaml:"tag" db:"tag"`
	DataMap `json:"data" yaml:"data"`
}

type Data struct {
	Date        string `json:"date" yaml:"date" db:"date"`
	Name        string `json:"name" yaml:"name" db:"name"`
	Power       int    `json:"power" yaml:"power" db:"power"`
	GiftLevel   int    `json:"gift-level" yaml:"giftLevel" db:"gift_level"`
	MemberCount int    `json:"member-count" yaml:"memberCount" db:"member_count"`
	R5Id        string `json:"r5-id" yaml:"r5Id" db:"r5_id"`
	AllianceId  string `json:"alliance-id" yaml:"allianceId" db:"alliance_id"`
}

type DataMap map[string]Data
type AllianceMap map[string]Alliance

var (
	ErrAllianceInsert  = errors.New("alliance: failed to insert alliance")
	ErrAllianceUpdate  = errors.New("alliance: failed to update alliance")
	ErrAllianceAddData = errors.New("alliance: failed to add alliance data")
	ErrInvalidArg      = errors.New("invalid argument")
	ErrInvalidMapKey   = errors.New("invalid key in map configuration")
)

// ScanAlliance pre-processes the given image file and scans it with tessaract
// into an alliance.Alliance struct
func (a *Alliance) ScanAlliance(img image.Image, imm scanner.ImageMaps) error {
	var d Data
	var err error

	for k, im := range imm {
		switch k {
		case "tag":
			a.Tag, err = im.ProcessImageText(img)
		case "name":
			d.Name, err = im.ProcessImageText(img)
		case "power":
			d.Power, err = im.ProcessImageInt(img)
		case "giftlevel":
			d.GiftLevel, err = im.ProcessImageInt(img)
		case "membercount":
			d.MemberCount, err = im.ProcessImageInt(img)
		default:
			return ErrInvalidMapKey
		}
		if err != nil {
			return err
		}
	}

	d.Date = time.Now().Format(time.DateOnly)

	err = a.GetByTag(a.Tag)
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
		fmt.Printf("A new alliance will need to be created from this data.  Please run 'wartracker-cli alliance new -o [output file]' after verifying the data\n")
	} else {
		fmt.Printf("This alliance already exists. To add the new data run 'wartracker-cli alliance add -o [output file]' to add the new data.\n")
	}

	a.DataMap[d.Date] = d

	return nil
}

// ScanAlliance pre-processes the given image file and scans it with tessaract
// into an alliance.Alliance struct
func (a *Alliance) ScanMainAlliance(img image.Image, imm scanner.ImageMaps) error {
	var d Data
	var err error

	// Setup alliance
	for k, im := range imm {
		switch k {
		case "tag":
			if cmd.Debug {
				fmt.Printf("scanning %s...\n", k)
			}
			a.Tag, err = im.ProcessImageText(img)
		case "name":
			if cmd.Debug {
				fmt.Printf("scanning %s...\n", k)
			}
			d.Name, err = im.ProcessImageText(img)
		case "power":
			if cmd.Debug {
				fmt.Printf("scanning %s...\n", k)
			}
			d.Power, err = im.ProcessImageInt(img)
		case "membercount":
			if cmd.Debug {
				fmt.Printf("scanning %s...\n", k)
			}
			tmc, err := im.ProcessImageText(img)
			if err != nil {
				return err
			}
			tmc = strings.Split(tmc, "/")[0]
			mc, err := strconv.Atoi(tmc)
			if err != nil {
				return err
			}
			d.MemberCount = mc
		default:
			return ErrInvalidMapKey
		}
		if err != nil {
			return err
		}
	}
	d.Date = time.Now().Format(time.DateOnly)

	err = a.GetByTag(a.Tag)
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
		fmt.Printf("A new alliance will need to be created from this data.  Please run 'wartracker-cli alliance new -o output' after verifying the data\n")
	} else {
		fmt.Printf("This alliance already exists. To add the new data run 'wartracker-cli alliance add -o output' to add the new data.\n")
		d.AllianceId = a.Id
	}

	a.DataMap[d.Date] = d

	return nil
}

// Create adds and alliance resource to the database
func (a *Alliance) Create(server int) error {
	var w wtid.WTID
	w.New("wartracker", "alliance", server)
	a.Id = string(w.Id)

	tx, err := db.Connection.Begin()
	if err != nil {
		return err
	}
	res, err := tx.Exec("INSERT INTO alliance (id, tag, server) VALUES (?, ?, ?)",
		a.Id,
		a.Tag,
		a.Server)
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
		return ErrAllianceInsert
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// Update updates an alliance resource in the database
func (a *Alliance) Update() error {
	date := time.Now().Format("2006-01-02")

	if d, ok := a.DataMap[date]; ok {
		d.AllianceId = a.Id
		d.Date = date
		a.DataMap[date] = d
	}

	tx, err := db.Connection.Begin()
	if err != nil {
		return err
	}
	res, err := tx.Exec("UPDATE alliance SET tag=?, server=? WHERE id=?",
		a.Tag,
		a.Server,
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
		return ErrAllianceUpdate
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// AddData adds data to an alliance resource in the database
func (a *Alliance) AddData(date string) error {
	tx, err := db.Connection.Begin()
	if err != nil {
		return err
	}
	res, err := tx.Exec("INSERT INTO alliance_data (name, date, power, gift_level, member_count, r5_id, alliance_id) VALUES (?, ?, ?, ?, ?, ?, ?)",
		a.DataMap[date].Name,
		a.DataMap[date].Date,
		a.DataMap[date].Power,
		a.DataMap[date].GiftLevel,
		a.DataMap[date].MemberCount,
		a.DataMap[date].R5Id,
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
		return ErrAllianceAddData
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// List geta all alliances
func List(withData ...bool) (AllianceMap, error) {
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

	if len(withData) > 0 {
		if withData[0] {
			for _, a := range as {
				a.GetLatestData()
			}
		}
	}

	return as, nil
}

// Get gets an alliance resource from the database
func (a *Alliance) Get(withData ...bool) error {
	a.DataMap = make(DataMap)

	err := db.Connection.QueryRowx("SELECT * FROM alliance WHERE id=?", a.Id).StructScan(a)
	if err != nil {
		return err
	}

	if len(withData) > 0 {
		if withData[0] {
			a.GetLatestData()
		}
	}

	return nil
}

// GetData retrives the data sets related to an alliance resource
func (a *Alliance) GetData() error {
	a.DataMap = make(DataMap)

	rows, err := db.Connection.Queryx("SELECT * FROM alliance_data WHERE alliance_id=? ORDER BY date DESC", a.Id)

	if err != nil {
		return err
	}
	for rows.Next() {
		var d Data
		err = rows.StructScan(&d)
		if err != nil {
			return err
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
		return err
	}
	a.DataMap[d.Date] = d

	return nil
}

// Finds the (first) alliance with the given tag, and returns the alliance object with all alliance data
func (a *Alliance) GetByTag(t string) error {
	err := db.Connection.QueryRowx("SELECT id FROM alliance WHERE tag=?", t).Scan(&a.Id)
	if err != nil {
		return err
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
