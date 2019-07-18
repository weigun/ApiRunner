// pipeline
package models

import (
	// "encoding/json"
	"fmt"
)

type Pipeline struct {
	Base
	Status        int     `json:"status"  yaml:"status"`
	Stages        []Stage `json:"stages"  yaml:"stages"`
	Notifications Event   `json:"notifications"  yaml:"notifications"`
}

func (pl *Pipeline) GetName() string {
	return pl.Config.Name
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
	return TYPE_TESTSUITS
}
