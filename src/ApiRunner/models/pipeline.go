// pipeline
package models

import (
	// "encoding/json"
	"fmt"
	"regexp"
)

type Pipeline struct {
	Base
	Host          string  `json:"host"  yaml:"host" toml:"host"`
	Status        int     `json:"status"  yaml:"status"`
	Stages        []Stage `json:"stages"  yaml:"stages"`
	Notifications Action  `json:"notifications"  yaml:"notifications"`
}

func (pl *Pipeline) GetName() string {
	return pl.Name
}

func (pl *Pipeline) Json() string {
	jsonStr, err := json.Marshal(pl)
	if err != nil {
		fmt.Println(`testcase to json failed:`, err.Error())
		return `{}`
	}
	return string(jsonStr)
}

func (pl *Pipeline) GetType() int {
	return TYPE_PIPELINE
}

func (pl *Pipeline) RefTag() string {
	if pl.Ref == `` {
		return regexp.MustCompile(`\s+`).ReplaceAllString(pl.Name, `_`)
	}
	return pl.Ref
}
