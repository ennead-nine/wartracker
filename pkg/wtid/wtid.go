package wtid

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/google/uuid"
)

// War Tracker ID types and functions.
type WTID struct {
	Id       Id
	Org      string
	Resource string
	Server   int
	UUID     string
}

type Id string

func (id Id) String() string {
	return string(id)
}

func (wtid *WTID) New(org string, resource string, server ...int) {
	wtid.Org = org
	wtid.Resource = resource
	if len(server) > 0 {
		wtid.Server = server[0]
	} else {
		wtid.Server = 0
	}
	wtid.UUID = uuid.NewString()

	wtid.Id = Id(wtid.Org + "." + wtid.Resource + "." + strconv.Itoa(int(wtid.Server)) + "." + wtid.UUID)
}

func (wtid *WTID) Parse(id Id) error {
	re := regexp.MustCompile(`^(.*)\.(.*)\.(.*)\.(.*)$`)
	fields := re.FindAllString(string(id), 4)

	wtid.Id = Id(id)
	wtid.Org = fields[0]
	wtid.Resource = fields[1]
	si, err := strconv.Atoi(fields[2])
	if err != nil {
		return err
	}
	wtid.Server = si
	wtid.UUID = fields[3]

	return nil
}

func Validate(id Id) error {
	var w WTID

	err := w.Parse(id)
	if err != nil {
		return err
	}
	if w.Org != "wartracker" {
		return fmt.Errorf("invalid org in id")
	}
	_, err = uuid.Parse(string(w.Id))
	if err != nil {
		return err
	}
	switch w.Resource {
	case "alliance":
		return nil
	case "commander":
		return nil
	case "vsday":
		return nil
	case "vsduel":
		return nil
	default:
		return fmt.Errorf("invalid rresource in id")
	}
}
