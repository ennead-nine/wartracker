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
	Tag      string  `json:"tag" yaml:"tag" db:"tag"`
	Server   int     `json:"server" yaml:"server" db:"server"`
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

type Alias struct {
	Alias       string `json:"alias" yaml:"alias" db:"alias"`
	Tag         string `json:"tag" yaml:"tag" db:"tag"`
	Server      int    `json:"server" yaml:"server" db:"server"`
	Preferred   bool   `json:"preferred" yaml:"preferred" db:"preferred"`
	CommanderId string `json:"commander-id" yaml:"commanderId" db:"commander_id"`
}

type DataMap map[string]Data

func (c *Commander) Create() error {
	var w wtid.WTID
	w.New("wartracker", "commander", c.Server)
	c.Id = string(w.Id)

	tx, err := db.Connection.Begin()
	if err != nil {
		return err
	}

	res, err := tx.Exec("INSERT INTO commander (id, note_name, tag, server) VALUES (?, ?, ?, ?)",
		c.Id,
		c.NoteName,
		c.Tag,
		c.Server)
	if err != nil {
		return err
	}
	x, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if x != 1 {
		return fmt.Errorf("failed to insert commander")
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	tx, err = db.Connection.Begin()
	if err != nil {
		return err
	}
	res, err = tx.Exec(`INSERT INTO commander_aliases (
		alias, 
		tag, 
		server, 
		commander_id 
		) VALUES (?, ?, ?, ?)`,
		c.NoteName, c.Tag, c.Server, c.Id)
	if err != nil {
		return err
	}
	x, err = res.RowsAffected()
	if err != nil {
		return err
	}
	if x != 1 {
		return fmt.Errorf("failed to insert alias")
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (c *Commander) AddData() error {
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

func (c *Commander) GetById(id string, latest ...bool) error {
	err := db.Connection.QueryRowx("SELECT * FROM commander WHERE id=?", id).StructScan(c)
	if err != nil {
		return err
	}

	q := "SELECT * FROM commander_data WHERE commander_id=? ORDER BY date DESC"
	if len(latest) > 0 {
		if latest[0] {
			q += " LIMIT 1"
		}
	}

	rows, err := db.Connection.Queryx(q, id)
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

func (c *Commander) GetByAlias(n string, latest ...bool) error {
	var id string

	err := db.Connection.QueryRowx("SELECT commander_id FROM commander_aliases WHERE alias=?", n).Scan(&id)
	if err != nil {
		return err
	}

	if len(latest) > 0 {
		if latest[0] {
			err = c.GetById(id, latest[0])
		} else {
			err = c.GetById(id)
		}
	}
	if err != nil {
		return err
	}

	return nil
}

// AddAlias adds an alias to the character's name list
func (c *Commander) AddAlias(a string) error {
	return fmt.Errorf("not implemented")
}

// SetNoteName sets the commander's note name (main name) to n.  If the note
// name is not equal to the new name, the old note name is saved as alias if
// needed.
func (c *Commander) SetNoteName(n string) error {
	return fmt.Errorf("not implemented")
}

// Merge will copy data and aliases from s to p and delete s, and add s.NoteName as an alias for p
func Merge(p, s Commander) error {
	return fmt.Errorf("not implemented")
}

func List() ([]Commander, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *Commander) CommanderToJSON() ([]byte, error) {
	return json.MarshalIndent(c, "", "\t")
}

func (c *Commander) CommanderToYAML() ([]byte, error) {
	return yaml.Marshal(c)
}
