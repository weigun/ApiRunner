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
	Base
	Extend
	// Config CaseConfig `json:"config"   yaml:"config"`
	Env           Variables  `json:"env"   yaml:"env"`
	MergeMode     string     `json:"mergeMode"   yaml:"mergeMode"`
	Steps         []Step     `json:"steps"  yaml:"steps"`
	Notifications StageEvent `json:"notifications"  yaml:"notifications"`
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
