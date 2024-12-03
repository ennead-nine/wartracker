package wtid

import (
	"regexp"
	"strconv"

	"github.com/google/uuid"
)

// War Tracker ID types and functions.
type WTID struct {
	Id       string
	Org      string
	Resource string
	Server   int64
	UUID     string
}

func (wtid *WTID) New(org string, resource string, server ...int64) {
	wtid.Org = org
	wtid.Resource = resource
	if len(server) > 0 {
		wtid.Server = server[0]
	} else {
		wtid.Server = 0
	}
	wtid.UUID = uuid.NewString()

	wtid.Id = wtid.Org + "." + wtid.Resource + "." + strconv.Itoa(int(wtid.Server)) + "." + wtid.UUID
}

func (wtid *WTID) Parse(id string) error {
	re := regexp.MustCompile(`^(.*)\.(.*)\.(.*)\.(.*)$`)
	fields := re.FindAllString(id, 4)

	wtid.Id = id
	wtid.Org = fields[0]
	wtid.Resource = fields[1]
	si, err := strconv.Atoi(fields[2])
	if err != nil {
		return err
	}
	wtid.Server = int64(si)
	wtid.UUID = fields[3]

	return nil
}
