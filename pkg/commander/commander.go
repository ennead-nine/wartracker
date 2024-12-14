package commander

import (
	"encoding/json"
	"fmt"
	"wartracker/pkg/db"
	"wartracker/pkg/wtid"

	"gopkg.in/yaml.v3"
)

type Commander struct {
	Id       string `json:"id" yaml:"id" db:"id"`
	NoteName string `json:"note-name" yaml:"noteName" db:"note_name"`
	Data     []Data `json:"data" yaml:"data"`
}

type Data struct {
	Date            string `json:"date" yaml:"date" db:"date"`
	PFP             []byte `json:"pfp" yaml:"pfp" db:"pfp"`
	HQLevel         int64  `json:"hq-level" yaml:"hqLevel" db:"hq_level"`
	Likes           int64  `json:"likes" yaml:"likes" db:"likes"`
	HQPower         int64  `json:"hq-power" yaml:"HqPower" db:"hq_power"`
	Kills           int64  `json:"kills" yaml:"kills" db:"kills"`
	ProfessionLevel int64  `json:"profession-level" yaml:"professionLevel" db:"profession_level"`
	TotalHeroPower  int64  `json:"total-hero-power" yaml:"totalHeroPower" db:"total_hero_power"`
	AllianceID      string `json:"alliance-id" yaml:"allianceId" db:"alliance_id"`
	CommanderID     string `json:"commander-id" yaml:"commanderId" db:"commander_id"`
}

func (c *Commander) Create(server int64) error {
	var w wtid.WTID
	w.New("wartracker", "commander", server)
	c.Id = w.Id
	c.Data[0].CommanderID = w.Id

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
		alliance_id,
		commander_id
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		c.Data[0].Date,
		c.Data[0].HQLevel,
		c.Data[0].Likes,
		c.Data[0].HQPower,
		c.Data[0].Kills,
		c.Data[0].ProfessionLevel,
		c.Data[0].TotalHeroPower,
		c.Data[0].AllianceID,
		c.Data[0].CommanderID)
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
	tx, err := db.Connection.Begin()
	if err != nil {
		return err
	}

	res, err := tx.Exec("INSERT INTO commander_data (date, pfp, hq_level, likes, hq_power, kills, profession_level, total_hero_power, alliance_id, commander_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		c.Data[0].Date,
		c.Data[0].PFP,
		c.Data[0].HQLevel,
		c.Data[0].Likes,
		c.Data[0].HQPower,
		c.Data[0].Kills,
		c.Data[0].ProfessionLevel,
		c.Data[0].TotalHeroPower,
		c.Data[0].AllianceID,
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
		c.Data = append(c.Data, d)
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
