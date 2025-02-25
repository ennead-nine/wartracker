package commander

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"wartracker/pkg/db"
	"wartracker/pkg/warzone"
	"wartracker/pkg/wtid"

	"gopkg.in/yaml.v3"
)

type Commander struct {
	Id        string  `json:"id" yaml:"id" db:"id"`
	Name      string  `json:"name" yaml:"name" db:"name"`
	WarzoneId string  `json:"warzone-id" yaml:"warzone-id" db:"warzone_id"`
	Data      DataMap `json:"data" yaml:"data"`
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
	AllianceRank    int    `json:"alliance-rank" yaml:"allianceRank" db:"alliance_rank"`
	AllianceId      string `json:"alliance-id" yaml:"allianceId" db:"alliance_id"`
	CommanderId     string `json:"commander-id" yaml:"commanderId" db:"commander_id"`
}

type Alias struct {
	Alias       string `json:"alias" yaml:"alias" db:"alias"`
	Tag         string `json:"tag" yaml:"tag" db:"tag"`
	Preferred   bool   `json:"preferred" yaml:"preferred" db:"preferred"`
	CommanderId string `json:"commander-id" yaml:"commanderId" db:"commander_id"`
}

type DataMap map[string]Data
type CommanderMap map[string]Commander

func (c *Commander) Create() error {
	var z warzone.Warzone
	z.Id = c.WarzoneId
	err := z.Get()
	if err != nil {
		return fmt.Errorf("unable to lookup warzone %s: %w", c.WarzoneId, err)
	}

	var w wtid.WTID
	w.New("wartracker", "commander", z.Server)
	c.Id = string(w.Id)

	tx, err := db.Connection.Begin()
	if err != nil {
		return err
	}
	res, err := tx.Exec("INSERT INTO commander (id, name, warzone_id) VALUES (?, ?, ?)",
		c.Id,
		c.Name,
		c.WarzoneId)
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
		return fmt.Errorf("failed to insert commander: unknown error")
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return c.AddAlias(c.Name, true)
}

func (c *Commander) AddAlias(n string, p bool) error {
	if p {
		tx, err := db.Connection.Begin()
		if err != nil {
			return err
		}
		_, err = tx.Exec("UPDATE commander_alias SET preferred = false WHERE commander_id = ?", c.Id)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update preferred alias for %s: %w", c.Id, err)
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
	res, err := tx.Exec(`INSERT INTO commander_alias (
		alias,
		tag,
		preferred,
		commander_id 
		) VALUES (?, ?, ?, ?)`,
		n, "NA", p, c.Id)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to add alias for %s: %w", c.Id, err)
	}
	x, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to add alias for %s: %w", c.Id, err)
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

func (c *Commander) GetAliases() ([]Alias, error) {
	var as []Alias

	q := "SELECT * FROM commander_alias WHERE commander_id=?"
	rows, err := db.Connection.Queryx(q, c.Id)
	if err != nil {
		return nil, fmt.Errorf("unable to get aliases for %s, %w", c.Id, err)
	}
	for rows.Next() {
		var a Alias
		err = rows.StructScan(&a)
		if err != nil {
			rows.Close()
			return nil, fmt.Errorf("failed to get aliases for %s: %w", c.Id, err)
		}
		as = append(as, a)
	}

	return as, nil
}

func (c *Commander) AddData(date string, d Data) error {
	c.Data = make(DataMap)
	c.Data[date] = d

	tx, err := db.Connection.Begin()
	if err != nil {
		return err
	}

	res, err := tx.Exec("INSERT INTO commander_data (date, hq_level, likes, hq_power, kills, profession_level, total_hero_power, alliance_rank, alliance_id, commander_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		c.Data[date].Date,
		c.Data[date].HQLevel,
		c.Data[date].Likes,
		c.Data[date].HQPower,
		c.Data[date].Kills,
		c.Data[date].ProfessionLevel,
		c.Data[date].TotalHeroPower,
		c.Data[date].AllianceRank,
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

func (c *Commander) GetData() error {
	c.Data = make(DataMap)
	q := "SELECT * from commander_data WHERE commander_id=?"
	rows, err := db.Connection.Queryx(q, c.Id)
	if err != nil {
		return fmt.Errorf("unable to get data for %s, %w", c.Name, err)
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

func (c *Commander) GetLatestData() error {
	q := "SELECT * from commander_data WHERE commander_id=? ORDER BY date DESC LIMIT 1"
	rows, err := db.Connection.Queryx(q, c.Id)
	if err != nil {
		return fmt.Errorf("unable to get latest data for commander %s, %w", c.Name, err)
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

func (c *Commander) GetDataByDate(date string) error {
	q := "SELECT * from commander_data WHERE commander_id=? AND date=\"?\""
	rows, err := db.Connection.Queryx(q, c.Id, date)
	if err != nil {
		return fmt.Errorf("unable to get latest data for commander %s, %w", c.Name, err)
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

func (c *Commander) Get() error {
	err := db.Connection.QueryRowx("SELECT * FROM commander WHERE id=?", c.Id).StructScan(c)
	if err != nil {
		return err
	}

	return nil
}

func (c *Commander) Update() error {
	var z warzone.Warzone
	z.Id = c.WarzoneId
	err := z.Get()
	if err != nil {
		return fmt.Errorf("error looking up warzone %s: %w", c.WarzoneId, err)
	}

	tx, err := db.Connection.Begin()
	if err != nil {
		return err
	}
	res, err := tx.Exec("UPDATE commander SET name=?, warzone_id=? WHERE id=?",
		c.Name,
		c.WarzoneId,
		c.Id)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return c.AddAlias(c.Name, true)
}

func (c *Commander) GetByName(n string) error {
	var id string

	err := db.Connection.QueryRowx("SELECT DISTINCT commander_id FROM commander_alias WHERE alias=?", n).Scan(&id)
	if err != nil {
		return err
	}

	c.Id = id
	err = c.Get()
	if err != nil {
		return err
	}

	return nil
}

// Merge will copy data and aliases from s to p and delete s, and add s.NoteName as an alias for p
func Merge(src, dst string) error {
	var s Commander
	s.Id = src
	err := s.Get()
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return fmt.Errorf("failed to merge commander from %s to %s: %w", src, dst, err)
	}

	var p Commander
	p.Id = dst
	err = p.Get()
	if err != nil {
		return fmt.Errorf("failed to merge commander from %s to %s: %w", src, dst, err)
	}

	sas, err := s.GetAliases()
	if err != nil {
		return err
	}
	for _, sa := range sas {
		err = p.AddAlias(sa.Alias, false)
		if err != nil {
			return err
		}
	}

	err = s.GetData()
	if err != nil {
		return err
	}
	for _, sd := range s.Data {
		err = p.AddData(sd.Date, sd)
		if err != nil {
			return err
		}
	}

	return s.Delete()
}

func (c *Commander) DeleteData(date string) error {
	tx, err := db.Connection.Begin()
	if err != nil {
		return err
	}
	res, err := tx.Exec("DELETE FROM commander_data WHERE date=? AND commander_id=?",
		date,
		c.Id)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	delete(c.Data, date)
	return nil
}

func (c *Commander) DeleteAlias(alias string) error {
	tx, err := db.Connection.Begin()
	if err != nil {
		return err
	}
	res, err := tx.Exec("DELETE FROM commander_alias WHERE alias=? AND commander_id=?",
		alias,
		c.Id)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (c *Commander) Delete() error {
	as, err := c.GetAliases()
	if err != nil {
		return err
	}
	for _, a := range as {
		err = c.DeleteAlias(a.Alias)
		if err != nil {
			return err
		}
	}

	err = c.GetData()
	if err != nil {
		return err
	}
	for _, cd := range c.Data {
		err = c.DeleteData(cd.Date)
		if err != nil {
			return err
		}
	}

	tx, err := db.Connection.Begin()
	if err != nil {
		return err
	}
	res, err := tx.Exec("DELETE FROM commander WHERE id=?",
		c.Id)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func List() (CommanderMap, error) {
	var cs = make(CommanderMap)

	rows, err := db.Connection.Queryx("SELECT * FROM commander")

	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var c Commander
		err = rows.StructScan(&c)
		if err != nil {
			return nil, err
		}
		cs[c.Id] = c
	}

	return cs, nil
}

func ListByAlliance(a string) (CommanderMap, error) {
	var cs = make(CommanderMap)

	rows, err := db.Connection.Queryx("SELECT DISTINCT commander_id FROM commander_data WHERE alliance_id")

	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var c Commander
		var i string

		err = rows.Scan(&i)
		if err != nil {
			return nil, err
		}

		c.Id = i
		err = c.Get()
		if err != nil {
			return nil, fmt.Errorf("error getting commander %s: %w", i, err)
		}
		cs[c.Id] = c
	}

	return cs, nil
}

func (c *Commander) CommanderToJSON() ([]byte, error) {
	return json.MarshalIndent(c, "", "\t")
}

func (c *Commander) CommanderToYAML() ([]byte, error) {
	return yaml.Marshal(c)
}
