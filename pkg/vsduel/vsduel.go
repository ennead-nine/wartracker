package vsduel

import (
	"archive/zip"
	"bytes"
	"database/sql"
	"fmt"
	"image"
	"io"
	"sort"
	"strconv"
	"strings"
	"time"

	"wartracker/pkg/alliance"
	"wartracker/pkg/commander"
	"wartracker/pkg/db"
	"wartracker/pkg/scanner"
	"wartracker/pkg/wtid"
)

type VsDuel struct {
	Id           string  `json:"id" yaml:"id" db:"id"`
	Date         string  `json:"date" yaml:"date" db:"date"`
	LeagueLevel  string  `json:"league-level" yaml:"leagueLevel" db:"league_level"`
	LeagueId     string  `json:"league-id" yaml:"leagueId" db:"league_id"`
	TournamentId string  `json:"tournament-id" yaml:"tournamentId" db:"tournament_id"`
	Weeks        WeekMap `json:"weeks" yaml:"weeks"`
}

type Week struct {
	Id          string   `json:"id" yaml:"id" db:"id"`
	WeekNumber  int      `json:"week-number" yaml:"weekNumber" db:"vsweek_number"`
	AllianceIds []string `json:"alliance-ids" yaml:"allianceIds"`
	VsDuelId    string   `json:"vsduel-id" yaml:"duelId" db:"vsduel_id"`
	Data        DataMap  `json:"vsduel-data" yaml:"vsDuelData"`
}
type Day struct {
	Name         string `json:"name" yaml:"name" db:"name"`
	ShortName    string `json:"short-name" yaml:"shortName" db:"short_name"`
	DayOfWeek    string `json:"day-of-week" yaml:"dayOfWeek" db:"day_of_week"`
	VsDuelPoints int    `json:"vsduel-points" yaml:"vsDuelPoints" db:"vsduel_points"`
	DayNumber    int    `json:"day-number" yaml:"dayNumber" db:"day_number"`
}

type Data struct {
	Id            string           `json:"id" yaml:"id" db:"id"`
	WeekId        string           `json:"vsduel-week-id" yaml:"vsDuelWeekId" db:"vsduel_week_id"`
	DayOfWeek     string           `json:"day-of-week" yaml:"dayOfWeek" db:"day_of_week"`
	Push          bool             `json:"push" yaml:"push" db:"push"`
	Save          bool             `json:"save" yaml:"save" db:"save"`
	AllianceData  AllianceDataMap  `json:"vsduel-alliance-data" yaml:"vsDuelAllianceData"`
	CommanderData CommanderDataMap `json:"vsduel-commander-data" yaml:"vsDuelCommanderData"`
}

type AllianceData struct {
	Points       int    `json:"points" yaml:"points" db:"points"`
	VsDuelPoints int    `json:"vsduel-points" yaml:"vsDuelPoints" db:"vsduel_points"`
	AllianceId   string `json:"alliance-id" yaml:"allianceId" db:"alliance_id"`
	VsDuelDataId string `json:"vsduel-data-id" yaml:"vsDuelDataId" db:"vsduel_data_id"`
}

type CommanderData struct {
	Points       int    `json:"points" yaml:"points" db:"points"`
	Rank         int    `json:"rank" yaml:"rank" db:"rank"`
	New          bool   `json:"new" yaml:"new" db:"new"`
	Name         string `json:"name" yaml:"name" db:"name"`
	Vacation     bool   `json:"vacation" yaml:"vacation" db:"vacation"`
	AllianceId   string `json:"alliance-id" yaml:"allianceid" db:"alliance_id"`
	CommanderId  string `json:"commander-id" yaml:"commanderId" db:"commander_id"`
	VsDuelDataId string `json:"vsduel-data-id" yaml:"vsDuelDataId" db:"vsduel_data_id"`
}

// map[WeekNumber]Week
type WeekMap map[int]Week

// map[DayOfWeek]Day
type DayMap map[string]Day

// map[VsDuelDayId]Data
type DataMap map[string]Data

// map[CommanderId]CommanderData
type CommanderDataMap map[string]CommanderData

// map[AllianceId]AllianceData
type AllianceDataMap map[string]AllianceData

var DayFile string

var (
	ErrDuelDataInsert = fmt.Errorf("failed to insert duel data")
	ErrNumDays        = fmt.Errorf("number of versus days is not 6")
)

var Days DayMap

func InitDays() error {
	ds := DayMap{
		"Monday": {
			Name:         "Radar Training",
			ShortName:    "Radar",
			DayOfWeek:    "Monday",
			VsDuelPoints: 1,
			DayNumber:    0,
		},
		"Tuesday": {
			Name:         "Base Expansion",
			ShortName:    "Build",
			DayOfWeek:    "Tuesday",
			VsDuelPoints: 2,
			DayNumber:    1,
		},
		"Wednesday": {
			Name:         "Age of Science",
			ShortName:    "Tech",
			DayOfWeek:    "Wednesday",
			VsDuelPoints: 2,
			DayNumber:    2,
		},
		"Thursday": {
			Name:         "Train Heros",
			ShortName:    "Hero",
			DayOfWeek:    "Thursday",
			VsDuelPoints: 2,
			DayNumber:    3,
		},
		"Friday": {
			Name:         "Total Mobilization",
			ShortName:    "Units",
			DayOfWeek:    "Friday",
			VsDuelPoints: 2,
			DayNumber:    4,
		},
		"Saturday": {
			Name:         "Enemy Buster",
			ShortName:    "Kill",
			DayOfWeek:    "Saturday",
			VsDuelPoints: 4,
			DayNumber:    5,
		},
		// Totals
		"Sunday": {
			Name:         "Totals",
			ShortName:    "Totals",
			DayOfWeek:    "Sunday",
			VsDuelPoints: 0,
			DayNumber:    6,
		},
	}

	for _, d := range ds {
		tx, err := db.Connection.Begin()
		if err != nil {
			return err
		}
		res, err := tx.Exec("INSERT INTO vsduel_day (name, short_name, day_of_week, vsduel_points, day_number) VALUES (?, ?, ?, ?, ?)",
			d.Name,
			d.ShortName,
			d.DayOfWeek,
			d.VsDuelPoints,
			d.DayNumber)
		if err != nil {
			return err
		}
		x, err := res.RowsAffected()
		if err != nil {
			tx.Rollback()
			return err
		}
		if x != 1 {
			tx.Rollback()
			return ErrDuelDataInsert
		}
		err = tx.Commit()
		if err != nil {
			return err
		}
	}

	return nil
}

func GetDays(retry ...bool) (DayMap, error) {
	ds := make(DayMap)

	rows, err := db.Connection.Queryx("SELECT * FROM vsduel_day")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var d Day
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

	res, err := tx.Exec("INSERT INTO vsduel (id, date, league_level, league_id, tournament_id) VALUES (?, ?, ?, ?, ?)",
		v.Id,
		v.Date,
		v.LeagueLevel,
		v.LeagueId,
		v.TournamentId)
	if err != nil {
		return err
	}
	x, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}
	if x != 1 {
		tx.Rollback()
		return fmt.Errorf("failed to create vsduel for %s: %w", v.Date, err)
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	_, err = GetDays()
	if err != nil && err == sql.ErrNoRows {
		err = InitDays()
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}

	return nil
}

func (v *VsDuel) StartWeek(week int, as []string) error {
	var w wtid.WTID
	w.New("wartracker", "vsduelweek", 0)

	var k Week
	k.Id = string(w.Id)
	k.WeekNumber = week
	k.AllianceIds = append(k.AllianceIds, as...)
	k.VsDuelId = v.Id

	tx, err := db.Connection.Begin()
	if err != nil {
		return err
	}

	res, err := tx.Exec("INSERT INTO vsduel_week (id, vsweek_number, vsduel_id) VALUES (?, ?, ?)",
		k.Id,
		k.WeekNumber,
		k.VsDuelId)
	if err != nil {
		return err
	}
	x, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to start week %d for %s: %w", k.WeekNumber, k.VsDuelId, err)
	}
	if x != 1 {
		tx.Rollback()
		return fmt.Errorf("failed to start week %d for %s: %w", k.WeekNumber, k.VsDuelId, db.ErrDbErrorUnknown)
	}

	for _, a := range k.AllianceIds {
		res, err = tx.Exec("INSERT INTO vsduel_alliance (alliance_id, vsduel_week_id) VALUES (?, ?)", a, k.Id)
		if err != nil {
			tx.Rollback()
			return err
		}
		x, err := res.RowsAffected()
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to start week %d for %s: %w", k.WeekNumber, k.VsDuelId, err)
		}
		if x != 1 {
			tx.Rollback()
			return fmt.Errorf("failed to start week %d for %s: %w", k.WeekNumber, k.VsDuelId, db.ErrDbErrorUnknown)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to start week %d for %s: %w", k.WeekNumber, k.VsDuelId, err)
	}

	ks := make(WeekMap)
	ks[k.WeekNumber] = k
	v.Weeks = ks

	return nil
}

func Get(vid string) (*VsDuel, error) {
	var v VsDuel

	err := db.Connection.QueryRowx("SELECT * FROM vsduel WHERE id=?", vid).StructScan(&v)
	if err != nil {
		return nil, err
	}

	return &v, nil
}

func (v *VsDuel) GetWeeks() error {
	rows, err := db.Connection.Queryx("SELECT * FROM vsduel_week WHERE vsduel_id=?", v.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no weeks found for %s", v.Id)
		} else {
			return err
		}
	}
	ks := make(WeekMap)
	for rows.Next() {
		var k Week
		err = rows.StructScan(&k)
		if err != nil {
			return err
		}
		err = k.GetAlliances()
		if err != nil {
			return err
		}
		ks[k.WeekNumber] = k
	}

	v.Weeks = ks

	return nil
}

func GetWeek(kid string) (*Week, error) {
	var k Week
	err := db.Connection.QueryRowx("SELECT * FROM vsduel_week WHERE id=?", kid).StructScan(&k)
	if err != nil {
		return nil, err
	}

	err = k.GetAlliances()
	if err != nil {
		return nil, err
	}

	return &k, nil
}

func (k *Week) GetAlliances() error {
	rows, err := db.Connection.Queryx("SELECT alliance_id FROM vsduel_alliance WHERE vsduel_week_id=?", k.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no alliances found for week %s", k.Id)
		} else {
			return err
		}
	}
	for rows.Next() {
		var aid string
		err = rows.Scan(&aid)
		if err != nil {
			return err
		}
		k.AllianceIds = append(k.AllianceIds, aid)
	}

	return nil
}

func (k *Week) StartDay(dow string) error {
	var w wtid.WTID
	w.New("wartracker", "vsdueldata", 0)

	var d Data
	d.Id = string(w.Id)
	d.DayOfWeek = dow
	d.WeekId = k.Id

	tx, err := db.Connection.Begin()
	if err != nil {
		return err
	}
	res, err := tx.Exec("INSERT INTO vsduel_data (id, day_of_week, vsduel_week_id) VALUES (?, ?, ?)",
		d.Id,
		d.DayOfWeek,
		d.WeekId)
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
		return fmt.Errorf("failed to insert duel data")
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (k *Week) GetDays() error {
	dm := make(DataMap)
	rows, err := db.Connection.Queryx("SELECT * FROM vsduel_data WHERE vsduel_week_id=?", k.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no days found for %s", k.Id)
		} else {
			return err
		}
	}
	for rows.Next() {
		var d Data
		err = rows.StructScan(&d)
		if err != nil {
			return err
		}
		dm[d.DayOfWeek] = d
	}

	k.Data = dm
	return nil
}

func (k *Week) EndWeek() error {
	return fmt.Errorf("not yet implemented")
}

func (k *Week) AllianceTotal(aid string) (int, error) {
	var points int

	for dow, d := range k.Data {
		var p int
		did := k.Data[dow].Id
		aid := d.AllianceData[aid].AllianceId

		err := db.Connection.QueryRowx("SELECT SUM(points) FROM vsduel_alliance_data WHERE vsduel_data_id=? AND alliance_id=?", did, aid).Scan(&p)
		if err != nil {
			return 0, err
		}

		points += p
	}

	return points, nil
}

func (k *Week) CommanderTotal(cid string) (int, error) {
	var points int

	for _, d := range k.Data {
		var p int
		did := d.Id
		cid := d.CommanderData[cid].CommanderId

		err := db.Connection.QueryRowx("SELECT SUM(points) FROM vsduel_commander_data WHERE vsduel_data_id=? AND commander_id=?", did, cid).Scan(&p)
		if err != nil {
			return 0, err
		}

		points += p
	}

	return points, nil
}

func GetData(did string) (*Data, error) {
	var d Data

	err := db.Connection.QueryRowx("SELECT * FROM vsduel_data WHERE id=?", did).StructScan(&d)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

func (d *Data) GetCommanderData() error {
	cdm := make(CommanderDataMap)

	rows, err := db.Connection.Queryx("SELECT * FROM vsduel_commander_data WHERE vsduel_data_id=?", d.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no commander data found for %s: %w", d.Id, err)
		} else {
			return fmt.Errorf("failed to get commander data for %s: %w", d.Id, err)
		}
	}
	for rows.Next() {
		var cd CommanderData
		err = rows.StructScan(&cd)
		if err != nil {
			return err
		}
		cdm[cd.CommanderId] = cd
	}

	d.CommanderData = cdm

	return nil
}

func (d *Data) UpsertAllianceData(ad AllianceDataMap) error {
	tx, err := db.Connection.Begin()
	if err != nil {
		return err
	}

	if len(ad) < 1 {
		return fmt.Errorf("upsertAllianceData: no data present")
	}

	var vp int
	for _, a := range ad {
		res, err := tx.Exec("INSERT INTO vsduel_alliance_data (points, vsduel_points, alliance_id, vsduel_data_id) VALUES (?, ?, ?, ?)",
			a.Points,
			a.VsDuelPoints,
			a.AllianceId,
			d.Id)
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
			return fmt.Errorf("failed to insert duel data")
		}
		vp += a.VsDuelPoints
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	if vp == 0 {
		return fmt.Errorf("upsertAllianceData: vs points are zero: %v", ad)
	}

	d.AllianceData = ad

	return nil
}

func (d *Data) UpsertCommanderData(cd CommanderDataMap) error {
	tx, err := db.Connection.Begin()
	if err != nil {
		return err
	}
	for _, c := range cd {
		res, err := tx.Exec("INSERT INTO vsduel_commander_data (points, rank, new, name, alliance_id, commander_id, vsduel_data_id) VALUES (?, ?, ?, ?, ?, ?, ?)",
			c.Points,
			c.Rank,
			c.New,
			c.Name,
			c.AllianceId,
			c.CommanderId,
			d.Id)
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
			return fmt.Errorf("failed to insert duel data")
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	d.CommanderData = cd

	return nil
}

func (cdm *CommanderDataMap) Exists(cid string) (bool, error) {
	var n string
	err := db.Connection.QueryRowx("SELECT name FROM vsduel_commander_data WHERE id=?", cid).Scan(&n)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("error checking if commander exists: %w", err)
	}

	return true, nil
}

// Takes a zipfile containing a single day of versus points ranking screen shots and adds the data to the duel.
func (k *Week) ScanPointsRanking(z []byte, dow string) ([]string, error) {
	var badimg []string

	err := k.GetDays()
	if err != nil {
		return nil, err
	}

	d := k.Data[dow]

	cd := make(CommanderDataMap)

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
		var ns *scanner.Scanner
		var ps *scanner.Scanner
		if dow == "Saturday" {
			ns, err = scanner.NewScanner(b, true, true, image.Point{435, 580}, image.Point{540, 1775})
		} else {
			ns, err = scanner.NewScanner(b, false, true, image.Point{447, 715}, image.Point{552, 1725})
		}
		if err != nil {
			return nil, fmt.Errorf("could not create names scanner: %w", err)
		}
		// TODO: Crop points should be config/args
		if dow == "Saturday" {
			ps, err = scanner.NewScanner(b, true, true, image.Point{934, 580}, image.Point{358, 1775})
		} else {
			ps, err = scanner.NewScanner(b, false, true, image.Point{993, 715}, image.Point{312, 1725})
		}
		if err != nil {
			return nil, fmt.Errorf("could not create points scanner: %w", err)
		}

		var names [7]string
		var points [7]int
		var alliances [7]string

		// Scan Names/Alliaces
		t, err := ns.ScanImage()
		if err != nil {
			return nil, fmt.Errorf("error scanning image: %s: %w", f.Name, err)
		}
		c := 0
		for _, n := range t {
			if c >= 7 {
				badimg = append(badimg, f.Name)
				continue
			}
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
		for i, p := range t {
			pt := strings.Replace(p, ",", "", -1)
			points[i], err = strconv.Atoi(pt)
		}
		if err != nil {
			badimg = append(badimg, f.Name)
			continue
		}

		if len(names) < 7 || len(alliances) < 7 || len(points) < 7 {
			badimg = append(badimg, f.Name)
			continue
		}
		for i := 0; i < 7; i++ {
			var x CommanderData
			err = x.processCommander(alliances[i], names[i], points[i])
			if err != nil {
				badimg = append(badimg, f.Name)
				continue
			}
			if dow == "Saturday" {
				k.fixSaturday(&x)
			}
			cd[x.CommanderId] = x
		}
	}

	cd.setRank()
	err = d.UpsertCommanderData(cd)
	d.CommanderData = cd
	if err != nil {
		return nil, err
	}

	k.Data[dow] = d

	return badimg, nil
}

func (d *Data) CalculateAllianceData() error {
	adm := make(AllianceDataMap)
	dow := d.DayOfWeek
	k, err := GetWeek(d.WeekId)
	if err != nil {
		return fmt.Errorf("failed to get week %s: %w", d.WeekId, err)
	}

	for _, aid := range k.AllianceIds {
		var ad AllianceData
		ad.AllianceId = aid
		ad.VsDuelDataId = adm[aid].VsDuelDataId
		ad.processAlliance(d.CommanderData)
		adm[aid] = ad
	}
	vp := Days[dow].VsDuelPoints
	if adm[k.AllianceIds[0]].Points >= adm[k.AllianceIds[1]].Points {
		x := adm[k.AllianceIds[0]]
		x.VsDuelPoints = vp
		adm[k.AllianceIds[0]] = x
	} else {
		x := adm[k.AllianceIds[1]]
		x.VsDuelPoints = vp
		adm[k.AllianceIds[1]] = x
	}

	if adm[k.AllianceIds[0]].VsDuelPoints == 0 && adm[k.AllianceIds[1]].VsDuelPoints == 0 {
		return fmt.Errorf("scanPointsRanking: vspoints are zero: %v: should be: %d", adm, vp)
	}
	err = d.UpsertAllianceData(adm)
	if err != nil {
		return err
	}

	return nil
}

func (ad *AllianceData) processAlliance(cdm CommanderDataMap) {
	p := 0
	for _, c := range cdm {
		if c.AllianceId == ad.AllianceId {
			p += c.Points
		}
	}
	ad.Points = p
}

func (cd *CommanderData) processCommander(a, c string, p int) error {
	var cc commander.Commander
	var aa alliance.Alliance

	cd.Points = p

	t, err := alliance.SplitTagName(a)
	if err != nil {
		return err
	}
	err = aa.GetByTag(t[0])
	if err != nil {
		return err
	}
	cd.AllianceId = aa.Id

	err = cc.GetByName(c)
	if err != nil {
		if err == sql.ErrNoRows {
			cc.WarzoneId = aa.WarzoneId
			cc.Name = c
			err = cc.Create()
			if err != nil {
				return err
			}
			var ccd commander.Data
			date := time.Now().Format("2006-01-02")
			ccd.AllianceId = aa.Id
			ccd.Date = date
			err := cc.AddData(date, ccd)
			if err != nil {
				return err
			}

			cd.New = true
		} else {
			return err
		}
	}
	cd.CommanderId = cc.Id
	cd.Name = c

	return nil
}

func unzipSS(z []byte) ([]*zip.File, error) {
	a, err := zip.NewReader(bytes.NewReader(z), int64(len(z)))
	if err != nil {
		return nil, fmt.Errorf("failed to create zip reader: %w", err)
	}

	return a.File, nil
}

func (k *Week) fixSaturday(cd *CommanderData) error {
	var mft int

	dows := [6]string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday"}

	for _, dow := range dows {
		var dp int
		err := db.Connection.QueryRowx("SELECT points FROM vsduel_commander_data WHERE vsduel_data_id=? AND commander_id=?", k.Data[dow].Id, cd.CommanderId).Scan(&dp)
		if err != nil && err != sql.ErrNoRows {
			return err
		}
		mft += dp
	}

	cd.Points = cd.Points - mft

	return nil
}

func (cd CommanderDataMap) setRank() {
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

func (cd CommanderDataMap) UpdateRanks() error {
	cd.setRank()

	tx, err := db.Connection.Begin()
	if err != nil {
		return err
	}
	for _, c := range cd {
		res, err := tx.Exec("UPDATE vsduel_commander_data SET rank=? WHERE commander_id=? AND vsduel_data_id=?",
			c.Rank,
			c.CommanderId,
			c.VsDuelDataId)
		if err != nil {
			tx.Rollback()
			return err
		}
		x, err := res.RowsAffected()
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update rank for %s: %w", c.CommanderId, err)
		}
		if x != 1 {
			tx.Rollback()
			return fmt.Errorf("failed to update rank for %s to %d", c.CommanderId, c.Rank)
		}
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to update ranks: %w", err)
	}

	return nil
}

func (k *Week) MergeCommanderData(src, dst string) error {
	for _, d := range k.Data {
		var sd CommanderData
		var dc commander.Commander
		dc.Id = dst
		err := dc.Get()
		if err != nil {
			return fmt.Errorf("merge commanders: failed to get commander %s: %w", dst, err)
		}

		err = d.GetCommanderData()
		if err != nil {
			return fmt.Errorf("merge commanders: failed to get commander data: %w", err)
		}

		sd = d.CommanderData[src]
		tx, err := db.Connection.Begin()
		if err != nil {
			return err
		}
		res, err := tx.Exec("UPDATE vsduel_commander_data SET commander_id=? WHERE commmander_id=? amd vsduel_data_id=?", dst, sd.CommanderId, sd.VsDuelDataId)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute merge vsduel commander data for %s: %w", sd.CommanderId, err)
		}
		_, err = res.RowsAffected()
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to merge vsduel commander data for %s: %w", sd.CommanderId, err)
		}
		err = tx.Commit()
		if err != nil {
			return fmt.Errorf("failed to update vsuel commander data for %s: %w", sd.CommanderId, err)
		}

		sd.CommanderId = dst
		d.CommanderData[dst] = sd
		delete(d.CommanderData, src)
	}

	return commander.Merge(src, dst)
}
