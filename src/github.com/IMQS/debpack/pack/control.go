package pack

import (
	"strings"
)

type Control struct {
	Section         string   `json:"section"`
	Priority        string   `json:"priority"`
	Architecture    string   `json:"architecture"`
	Depends         string   `json:"depends"`
	Maintainer      string   `json:"maintainer"`
	Description     string   `json:"description"`
	LongDescription []string `json:"longdescription"`
}

func (c *Control) JoinedDescription() string {
	return strings.Join(c.LongDescription, "\n")
}
