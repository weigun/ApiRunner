// pipeline
package models

import (
	// "encoding/json"
	"fmt"
)

type StageItem = map[string]Stage

type Pipeline struct {
	Config CaseConfig  `json:"config"  yaml:"config" toml:"config"`
	Stages []StageItem `json:"stages"  yaml:"stages" toml:"testcases"`
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
