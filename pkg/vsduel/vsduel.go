package vsduel

import (
	"archive/zip"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"image"
	"io"
	"sort"
	"strconv"
	"strings"
	"wartracker/pkg/alliance"
	"wartracker/pkg/commander"
	"wartracker/pkg/db"
	"wartracker/pkg/scanner"
	"wartracker/pkg/wtid"

	"gopkg.in/yaml.v3"
)

type VsDuel struct {
	Id            string `json:"id" yaml:"id" db:"id"`
	Date          string `json:"date" yaml:"date" db:"date"`
	LeagueLevel   string `json:"league-level" yaml:"leagueLevel" db:"league_level"`
	LeagueId      string `json:"league-id" yaml:"leagueId" db:"league_id"`
	TournamentId  string `json:"tournament-id" yaml:"tournamentId" db:"tournament_id"`
	VsDuelDataMap `json:"vsduel-data" yaml:"vsDuelData"`
}

type VSDuelWeek struct {
	VsDuelID    string `json:"vsduel-id" yaml:"duelId" db:"vsduel_id"`
	Week        int    `json:"week" yaml:"week" db:"week"`
	Alliance1Id string `json:"allliance1-id" yaml:"alliance1Id" db:"alliance1_id"`
	Alliance2Id string `json:"allliance2-id" yaml:"alliance2Id" db:"alliance2_id"`
}
type VsDay struct {
	Id        string `json:"id" yaml:"id" db:"id"`
	Name      string `json:"name" yaml:"name" db:"name"`
	ShortName string `json:"short-name" yaml:"shortName" db:"short_name"`
	DayOfWeek string `json:"day-of-week" yaml:"dayOfWeek" db:"day_of_week"`
}

type VsDuelData struct {
	Id                 string `json:"id" yaml:"id" db:"id"`
	VsDuelWeekId       string `json:"vsduel-id" yaml:"vsDuelId" db:"vsduel_id"`
	VsDuelDayId        string `json:"vsduel-day-id" yaml:"vsDuelDayId" db:"vsduel-day-id"`
	VsAllianceDataMap  `json:"vsduel-alliance-data" yaml:"vsDuelAllianceData"`
	VsCommanderDataMap `json:"vsduel-commander-data" yaml:"vsDuelCommanderData"`
}

type VsAllianceData struct {
	Id           string `json:"id" yaml:"id" db:"id"`
	Points       int    `json:"points" yaml:"points" db:"points"`
	Tag          string `json:"tag" yaml:"tag" db:"tag"`
	AllianceId   string `json:"alliance-id" yaml:"allianceId" db:"alliance_id"`
	VsDuelDataId string `json:"vsduel-data-id" yaml:"vsDuelDataId" db:"vsduel-data-id"`
}

type VsCommanderData struct {
	Id           string `json:"id" yaml:"id" db:"id"`
	Points       int    `json:"points" yaml:"points" db:"points"`
	Rank         int    `json:"rank" yaml:"rank" db:"rank"`
	Name         string `json:"name" yaml:"name" db:"name"`
	AllianceId   string `json:"alliance-id" yaml:"allianceid" db:"alliance_id"`
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
	ds := VsDays{
		"Monday": {
			Name:      "Radar Training",
			ShortName: "Radar",
			DayOfWeek: "Monday",
		},
		"Tuesday": {
			Name:      "Base Expansion",
			ShortName: "Construction",
			DayOfWeek: "Tuesday",
		},
		"Wednesday": {
			Name:      "Age of Science",
			ShortName: "Tech",
			DayOfWeek: "Wednesday",
		},
		"Thursday": {
			Name:      "Train Heros",
			ShortName: "Hero",
			DayOfWeek: "Thursday",
		},
		"Friday": {
			Name:      "Total Mobilization",
			ShortName: "Units",
			DayOfWeek: "Friday",
		},
		"Saturday": {
			Name:      "Enemy Buster",
			ShortName: "Kill",
			DayOfWeek: "Saturday",
		},
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

func GetDays(retry ...bool) (VsDays, error) {
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

	res, err := tx.Exec("INSERT INTO vsduel (id, date, league_level, league_id, tournamet_id) VALUES (?, ?, ?, ?)",
		v.Id,
		v.Date,
		v.LeagueLevel,
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
		res, err := tx.Exec("INSERT INTO vsduel_commander (points, rank, name, alliance_id, commander_id, vsduel_data_id) VALUES (?, ?, ?, ?, ?, ?)",
			d.Points,
			d.Rank,
			d.Name,
			d.AllianceId,
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

// Takes a zipfile containing a single day of versus points ranking screen shots and adds the data to the duel.
func (d *VsDuel) ScanPointsRanking(z []byte, did string) (VsCommanderDataMap, error) {
	cd := make(VsCommanderDataMap)

	files, err := unzipSS(z)
	if err != nil {
		return nil, fmt.Errorf("unable to unzip sreen shots: %w", err)
	}

	for _, f := range files {
		zf, err := f.Open()
		if err != nil {
			return nil, err
		}
		defer zf.Close()

		b, err := io.ReadAll(zf)
		if err != nil {
			return nil, err
		}

		// TODO: Crop points should be config/args
		ns, err := scanner.NewScanner(b, false, true, image.Point{447, 715}, image.Point{552, 1725})
		if err != nil {
			return nil, fmt.Errorf("could not create names scanner: %w", err)
		}
		// TODO: Crop points should be config/args
		ps, err := scanner.NewScanner(b, false, true, image.Point{993, 715}, image.Point{312, 1725})
		if err != nil {
			return nil, fmt.Errorf("could not create points scanner: %w", err)
		}

		var names [7]string
		var alliances [7]string
		var points [7]int

		// Scan Names/Alliaces
		t, err := ns.ScanImage()
		if err != nil {
			return nil, fmt.Errorf("error scanning image: %w", err)
		}
		c := 0
		for _, n := range t {
			if strings.HasPrefix(n, "[") {
				alliances[c] = n
				c++
			} else {
				names[c] += n
			}
		}

		// Scan Points
		t, err = ps.ScanImage()
		if err != nil {
			return nil, fmt.Errorf("error scanning image: %w", err)
		}
		c = 0
		for _, p := range t {
			pt := strings.Replace(p, ",", "", -1)
			points[c], err = strconv.Atoi(pt)
			if err != nil {
				return nil, fmt.Errorf("could not convert %s to integer: %w", p, err)
			}
			c++
		}

		for i := 0; i < 7; i++ {
			var x VsCommanderData
			x, err = processCommander(alliances[i], names[i], points[i])
			if err != nil {
				return nil, fmt.Errorf("error processing commander %s - %s: %w", alliances[i], names[i], err)
			}
			cd[names[i]] = x
		}
	}

	cd.SetRank()

	return cd, nil
}

func processCommander(a, c string, p int) (VsCommanderData, error) {
	var cd VsCommanderData
	var cc commander.Commander
	var aa alliance.Alliance

	t, err := alliance.SplitTagName(a)
	if err != nil {
		return cd, err
	}
	err = aa.GetByTag(t[0])
	if err == nil {
		cd.AllianceId = aa.Id
	} else if err == sql.ErrNoRows {
		aa.Tag = t[0]
		aa.Server = 0
		err := aa.Create()
		if err != nil {
			return cd, fmt.Errorf("could not create alliance [%s] %s: %w", t[0], t[1], err)
		}
		cd.AllianceId = aa.Id
	} else {
		return cd, fmt.Errorf("error getting alliance data for [%s] %s: %w", t[0], t[1], err)
	}

	cd.Name = c
	err = cc.GetByAlias(cd.Name)
	if err == nil {
		cd.CommanderId = cc.Id
	} else if err == sql.ErrNoRows {
		cc.NoteName = cd.Name
		cc.Server = aa.Server
		cc.Tag = aa.Tag
		err := cc.Create()
		if err != nil {
			return cd, fmt.Errorf("could not create commander [%s] %s: %w", cc.Tag, cc.NoteName, err)
		}
		cd.CommanderId = cc.Id
	} else {
		return cd, fmt.Errorf("error getting  commander [%s] %s: %w", aa.Tag, cd.Name, err)
	}

	cd.Points = p

	return cd, nil
}

func unzipSS(z []byte) ([]*zip.File, error) {
	a, err := zip.NewReader(bytes.NewReader(z), int64(len(z)))
	if err != nil {
		return nil, fmt.Errorf("failed to create zip reader: %w", err)
	}

	return a.File, nil
}

func (cd VsCommanderDataMap) SetRank() {
	type kv struct {
		Key   string
		Value int
	}

	var ss []kv
	for k, v := range cd {
		ss = append(ss, kv{k, v.Points})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	for i, kv := range ss {
		c := cd[kv.Key]
		c.Rank = i + 1
		cd[kv.Key] = c
	}
}
