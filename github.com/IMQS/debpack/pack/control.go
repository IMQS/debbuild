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

	remain := func(str string, index int) int {
		if len(str) < index {
			return len(str) - 1
		}
		return index
	}

	length := len(c.LongDescription)
	start := 0
	stop := start + remain(c.LongDescription, linelen)
	for stop < length {
		stop = stop + strings.Index(c.LongDescription[stop:], " ")
		result = append(result, fmt.Sprintf(" %s", c.LongDescription[start:stop]))
		start = stop
		stop := start + remain(c.LongDescription[start:], linelen)

	}

	for i := 0; i < length; i += linelen {
		stop := i + strings.Index(c.LongDescription[i:], " ")
		if stop > 59 {
			stop = 59
		}
		line := fmt.Sprintf(" %s", c.LongDescription[i:i+stop])
		result = append(result, line)
	}

	return []byte(strings.Join(result, "\n"))
}
