package commander

import (
	"encoding/json"
	"fmt"
	"time"
	"wartracker/pkg/db"
	"wartracker/pkg/wtid"

	"gopkg.in/yaml.v3"
)

type Commander struct {
	Id       string  `json:"id" yaml:"id" db:"id"`
	NoteName string  `json:"note-name" yaml:"noteName" db:"note_name"`
	Data     DataMap `json:"data" yaml:"data"`
}

type Data struct {
	Date            string `json:"date" yaml:"date" db:"date"`
	PFP             []byte `json:"pfp" yaml:"pfp" db:"pfp"`
	HQLevel         int    `json:"hq-level" yaml:"hqLevel" db:"hq_level"`
	Likes           int    `json:"likes" yaml:"likes" db:"likes"`
	HQPower         int    `json:"hq-power" yaml:"HqPower" db:"hq_power"`
	Kills           int    `json:"kills" yaml:"kills" db:"kills"`
	ProfessionLevel int    `json:"profession-level" yaml:"professionLevel" db:"profession_level"`
	TotalHeroPower  int    `json:"total-hero-power" yaml:"totalHeroPower" db:"total_hero_power"`
	AllianceId      string `json:"alliance-id" yaml:"allianceId" db:"alliance_id"`
	CommanderId     string `json:"commander-id" yaml:"commanderId" db:"commander_id"`
}

type DataMap map[string]Data

func (c *Commander) Create(server int) error {
	var w wtid.WTID
	w.New("wartracker", "commander", server)
	c.Id = string(w.Id)

	date := time.Now().Format(time.DateOnly)

	if d, ok := c.Data[date]; ok {
		d.CommanderId = c.Id
		d.Date = date
		c.Data[date] = d
	}

	tx, err := db.Connection.Begin()
	if err != nil {
		return err
	}

	res, err := tx.Exec("INSERT INTO commander (id, note_name) VALUES (?, ?)",
		c.Id,
		c.NoteName)
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
	res, err = tx.Exec(`INSERT INTO commander_data (
		date, 
		hq_level, 
		likes, 
		hq_power, 
		kills, 
		profession_level, 
		total_hero_power, 
		commander_id,
		alliance_id
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		c.Data[date].Date,
		c.Data[date].HQLevel,
		c.Data[date].Likes,
		c.Data[date].HQPower,
		c.Data[date].Kills,
		c.Data[date].ProfessionLevel,
		c.Data[date].TotalHeroPower,
		c.Data[date].CommanderId,
		c.Data[date].AllianceId)
	if err != nil {
		return err
	}
	x, err = res.RowsAffected()
	if err != nil {
		return err
	}
	if x != 1 {
		return fmt.Errorf("failed to insert commander data")
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (c *Commander) Update() error {
	if len(c.Data) < 1 {
		return fmt.Errorf("data for commander [%s] is empty", c.NoteName)
	}

	date := time.Now().Format("2006-01-02")

	if d, ok := c.Data[date]; ok {
		d.CommanderId = c.Id
		d.Date = date
		c.Data[date] = d
	}

	tx, err := db.Connection.Begin()
	if err != nil {
		return err
	}

	res, err := tx.Exec("INSERT INTO commander_data (date, hq_level, likes, hq_power, kills, profession_level, total_hero_power, alliance_id, commander_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		c.Data[date].Date,
		c.Data[date].HQLevel,
		c.Data[date].Likes,
		c.Data[date].HQPower,
		c.Data[date].Kills,
		c.Data[date].ProfessionLevel,
		c.Data[date].TotalHeroPower,
		c.Data[date].AllianceId,
		c.Id)
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
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (c *Commander) GetById(id string) error {
	err := db.Connection.QueryRowx("SELECT * FROM commander WHERE id=?", id).StructScan(c)
	if err != nil {
		return err
	}

	rows, err := db.Connection.Queryx("SELECT * FROM commander_data WHERE commander_id=? ORDER BY date DESC", id)
	if err != nil {
		return err
	}
	for rows.Next() {
		var d Data
		err = rows.StructScan(&d)
		if err != nil {
			return err
		}
		c.Data[d.Date] = d
	}

	return nil
}

func (c *Commander) GetByNoteName(n string) error {
	var id string

	err := db.Connection.QueryRowx("SELECT id FROM commander WHERE note_name=?", n).Scan(&id)
	if err != nil {
		return err
	}

	err = c.GetById(id)
	if err != nil {
		return err
	}

	return nil
}

func (c *Commander) CommanderToJSON() ([]byte, error) {
	return json.MarshalIndent(c, "", "\t")
}

func (c *Commander) CommanderToYAML() ([]byte, error) {
	return yaml.Marshal(c)
}
