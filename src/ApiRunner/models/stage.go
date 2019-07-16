package models

import (
	// "encoding/json"
	"fmt"
)

const (
	TYPE_API = iota
	TYPE_TESTCASE
	TYPE_TESTSUITS
)

type ICaseObj interface {
	GetName() string
	Json() string
	GetType() int
}

type Stage struct {
	Config CaseConfig `json:"config"   yaml:"config"`
	Ref    string     `json:"ref" yaml:"ref"`
	Func   Variables  `json:"func,omitempty"  yaml:"func"`
	Return Variables  `json:"return,omitempty"  yaml:"return"`
	Jobs   []Job      `json:"jobs"  yaml:"jobs"`
}

func (st *Stage) GetName() string {
	return st.Config.Name
}

func (st *Stage) Json() string {
	jsonStr, err := json.Marshal(st)
	if err != nil {
		fmt.Println(`tesstase to json failed:`, err.Error())
		return `{}`
	}
	return string(jsonStr)
}

func (st *Stage) GetType() int {
	return TYPE_TESTCASE
}
