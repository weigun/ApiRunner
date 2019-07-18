// pipegroup
package models

import (
	"fmt"
)

type PipeGroup struct {
	Base
	Pipelines []Pipeline `json:"Pipelines"  yaml:"Pipelines"`
}

func (pg *PipeGroup) GetName() string {
	return pg.Config.Name
}

func (pg *PipeGroup) Json() string {
	jsonStr, err := json.Marshal(pg)
	if err != nil {
		fmt.Println(`testcase to json failed:`, err.Error())
		return `{}`
	}
	return string(jsonStr)
}

func (pg *PipeGroup) GetType() int {
	return TYPE_TESTSUITS
}
