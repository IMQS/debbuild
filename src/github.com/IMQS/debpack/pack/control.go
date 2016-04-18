package pack

import (
	"fmt"
	"reflect"
	"strings"
)

const linelen = 60

type Control struct {
	Package         string `json:"package"`
	Version         string `json:"version"`
	Section         string `json:"section"`
	Priority        string `json:"priority"`
	Architecture    string `json:"architecture"`
	Depends         string `json:"depends"`
	Maintainer      string `json:"maintainer"`
	Description     string `json:"description"`
	LongDescription string `json:"longdescription"`
}

func (c *Control) Bytes() []byte {
	r := reflect.ValueOf(c).Elem()
	typeOfT := r.Type()
	result := make([]string, 0)
	for i := 0; i < r.NumField(); i++ {
		name := typeOfT.Field(i).Name
		if name == "LongDescription" {
			continue
		}
		result = append(result, fmt.Sprintf("%s: %s", name, r.Field(i).Interface()))
	}

	remain := func(str string) int {
		length := len(str)
		if length < linelen {
			return length
		}
		return linelen
	}
	start := 0
	stop := start + remain(c.LongDescription)
	for {
		spaceindex := strings.Index(c.LongDescription[stop:], " ")
		if spaceindex > 0 {
			stop = stop + spaceindex
		} else {
			stop = len(c.LongDescription)
		}
		line := strings.TrimSpace(c.LongDescription[start:stop])
		result = append(result, fmt.Sprintf(" %s", line))
		start = stop + 1
		if start >= len(c.LongDescription) {
			break
		}
		stop = start + remain(c.LongDescription[start:])
	}

	return []byte(strings.Join(result, "\n"))
}
