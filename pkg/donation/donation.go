package donation

import (
	"fmt"
	"io"
	"path/filepath"
	"strconv"
	"strings"

	"wartracker/pkg/commander"
	"wartracker/pkg/common"
	"wartracker/pkg/db"
	"wartracker/pkg/scanner"
	"wartracker/pkg/warzone"
	"wartracker/pkg/wtid"
)

type Donation struct {
	Id          string `json:"id" yaml:"id" db:"id"`
	Amount      int    `json:"amount" yaml:"amount" db:"amount"`
	Date        string `json:"date" yaml:"date" db:"date"`
	Name        string `json:"name" yaml:"name" db:"name"`
	CommanderId string `json:"commander-id" yaml:"commanderId" db:"commander_id"`
}

type DonationMap map[string]Donation

var CropNames = common.CropRatios{
	Rpx: .38,
	Rpy: .21,
	Rsx: .39,
	Rsy: .58,
}
var CropPoints = common.CropRatios{
	Rpx: .75,
	Rpy: .21,
	Rsx: .22,
	Rsy: .58,
}

func AddDonation(d *Donation) error {
	var c commander.Commander
	err := c.GetByName(d.Name)
	if err != nil {
		return fmt.Errorf("unable to lookup commander %s: %w", c.Name, err)
	}
	d.CommanderId = c.Id

	var z warzone.Warzone
	z.Id = c.WarzoneId
	err = z.Get()
	if err != nil {
		return fmt.Errorf("unable to lookup warzone %s: %w", c.WarzoneId, err)
	}

	var w wtid.WTID
	w.New("wartracker", "commander", z.Server)
	d.Id = string(w.Id)

	tx, err := db.Connection.Begin()
	if err != nil {
		return err
	}
	res, err := tx.Exec("INSERT INTO donation (id, amount, date, name, commander_id) VALUES (?, ?, ?, ?, ?)",
		d.Id,
		d.Amount,
		d.Date,
		c.Name,
		d.CommanderId)
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

func (dm DonationMap) ScanDonations(z []byte, date string) ([]string, error) {
	var badimg []string

	files, err := common.UnzipSS(z)
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

		if filepath.Ext(f.Name) == "jpg" {
			b, err = common.ConvertJpgToPng(b)
			if err != nil {
				return nil, err
			}
		}

		nb, err := common.CropImage(b, CropNames)
		if err != nil {
			return nil, err
		}
		pb, err := common.CropImage(b, CropPoints)
		if err != nil {
			return nil, err
		}

		// TODO: Crop points should be config/args
		var ns *scanner.Scanner
		var ps *scanner.Scanner
		ns, err = scanner.NewScanner(nb, true)
		if err != nil {
			return nil, fmt.Errorf("could not create names scanner: %w", err)
		}
		ps, err = scanner.NewScanner(pb, true)
		if err != nil {
			return nil, fmt.Errorf("could not create points scanner: %w", err)
		}

		var names []string
		var points []int

		t, err := ns.ScanImage()
		if err != nil {
			return nil, fmt.Errorf("error scanning image: %s: %w", f.Name, err)
		}
		names = append(names, t...)

		t, err = ps.ScanImage()
		if err != nil {
			return nil, fmt.Errorf("error scanning image: %w", err)
		}
		for _, p := range t {
			pt := strings.Replace(p, ",", "", -1)
			var pi int
			pi, err = strconv.Atoi(pt)
			points = append(points, pi)
		}
		if err != nil {
			badimg = append(badimg, f.Name)
			continue
		}

		if len(names) != 7 || len(points) != 7 {
			badimg = append(badimg, f.Name)
			continue
		}

		for i, n := range names {
			var d Donation
			d.Name = n
			d.Date = date
			d.Amount = points[i]
			dm[d.Name] = d
		}

		for _, d := range dm {
			err = AddDonation(&d)
			if err != nil {
				return badimg, fmt.Errorf("failed to add donation %#v: %w", d, err)
			}
		}
	}

	return badimg, nil
}
