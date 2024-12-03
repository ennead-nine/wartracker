package alliance

import (
	"context"
	"fmt"
	"time"
	"wartracker/pkg/db"
	"wartracker/pkg/wtid"
)

type Alliance struct {
	Id     string `json:"id" yaml:"id"`
	Server int64  `json:"server" yaml:"server"`
	Data
}

type data struct {
	Date        string `json:"date" yaml:"date"`
	Name        string `json:"name" yaml:"name"`
	Tag         string `json:"tag" yaml:"tag"`
	Power       int64  `json:"power" yaml:"power"`
	GiftLevel   int64  `json:"gift-level" yaml:"giftLevel"`
	MemberCount int64  `json:"member-count" yaml:"memberCount"`
	R5Id        string `json:"r5-id" yaml:"r5Id"`
}

type Data []data

// Adds and alliance to the database
func (a *Alliance) Add(server int64) error {
	var w wtid.WTID
	w.New("wartracker", "alliance", server)
	a.Id = w.Id

	query := fmt.Sprintf("INSERT INTO alliance (id, server) VALUES (\"%s\", %d);\n",
		a.Id,
		a.Server)

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	_, err := db.Connection.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	query = fmt.Sprintf("INSERT INTO alliance_data ("+
		"name, tag, date, power, gift-level, member-count, r5-id, alliance-id"+
		") VALUES ("+
		"\"%s\", \"%s\", \"%s\", %d, %d, %d, \"%s\", \"%s\");\n",
		a.Data[0].Name,
		a.Data[0].Tag,
		a.Data[0].Date,
		a.Data[0].Power,
		a.Data[0].GiftLevel,
		a.Data[0].MemberCount,
		a.Data[0].R5Id,
		a.Id)

	ctx, cancelfunc = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	_, err = db.Connection.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

// Populates an Alliance struct from the database by ID.  Gets all alliance data, with a.Data[0] being the latest
func (a *Alliance) GetById(id string) error {
	return nil
}

// Populates and Alliance structure, setting a.Data[0] to data from the specified date.  The rest of the a.Data slice will be empty
func (a *Alliance) GetDataByDate(d string) error {
	return nil
}

// Finds the (first) alliance with the given tag, and returns the alliance object with all alliance data
func (a *Alliance) GetDataByTag(n string) error {
	return nil
}
