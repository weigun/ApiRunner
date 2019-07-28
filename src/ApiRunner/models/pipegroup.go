// pipegroup
package models

/*
import (
	"fmt"
	"regexp"
)

type PipeGroup struct {
	Base
	Pipelines []Pipeline `json:"pipelines"  yaml:"pipelines"`
}

func (pg *PipeGroup) GetName() string {
	return pg.Name
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
	return TYPE_PIPEGROUP
}

func (pg *PipeGroup) RefTag() string {
	if pg.Ref == `` {
		return regexp.MustCompile(`\s+`).ReplaceAllString(pg.Name, `_`)
	}
	return pg.Ref
}
*/
