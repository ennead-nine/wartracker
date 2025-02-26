package warzone

import (
	"database/sql"
	"fmt"
	"wartracker/pkg/db"
	"wartracker/pkg/wtid"
)

type Warzone struct {
	Id     string `json:"id" yaml:"id" db:"id"`
	Server int    `json:"server" yaml:"server" db:"server"`
}

type WarzoneMap map[string]Warzone

func (z *Warzone) Create(s int) error {
	err := z.GetByServer(s)
	if err == nil {
		return fmt.Errorf("warzone for server %d already exists", s)
	} else if err != sql.ErrNoRows {
		return fmt.Errorf("error checking if warzone for server %d exists: %w", s, err)
	}

	z.Server = s
	var i wtid.WTID
	i.New("wartracker", "warzone", z.Server)
	z.Id = string(i.Id)

	tx, err := db.Connection.Begin()
	if err != nil {
		return err
	}
	res, err := tx.Exec("INSERT INTO warzone (id, server) VALUES (?, ?)",
		z.Id,
		z.Server)
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

	return nil
}

func (w *Warzone) Get() error {
	err := db.Connection.QueryRowx("SELECT * FROM warzone WHERE id=?", w.Id).StructScan(w)
	if err != nil {
		return err
	}

	return nil
}

func (w *Warzone) GetByServer(s int) error {
	err := db.Connection.QueryRowx("SELECT * FROM warzone WHERE server=?", s).StructScan(w)
	if err != nil {
		return err
	}

	return nil
}

func List() (WarzoneMap, error) {
	var zs = make(WarzoneMap)

	rows, err := db.Connection.Queryx("SELECT * FROM warzone")

	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var z Warzone
		err = rows.StructScan(&z)
		if err != nil {
			return nil, err
		}
		zs[z.Id] = z
	}

	return zs, nil
}
