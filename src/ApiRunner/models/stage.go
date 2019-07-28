package models

/*
import (
	"regexp"
	// "encoding/json"
	"fmt"
)

const (
	TYPE_STEP = iota
	TYPE_STAGE
	TYPE_PIPELINE
	TYPE_PIPEGROUP
)

type IPipe interface {
	GetName() string
	Json() string
	GetType() int
	RefTag() string
}

type Stage struct {
	Base
	Extend
	// Config CaseConfig `json:"config"   yaml:"config"`
	Env           Variables   `json:"env"   yaml:"env"`
	MergeMode     string      `json:"mergeMode"   yaml:"mergeMode"`
	Steps         []Step      `json:"steps"  yaml:"steps"`
	Notifications StageAction `json:"notifications"  yaml:"notifications"`
}

func (st *Stage) GetName() string {
	return st.Name
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
	return TYPE_STAGE
}

func (st *Stage) RefTag() string {
	if st.Ref == `` {
		return regexp.MustCompile(`\s+`).ReplaceAllString(st.Name, `_`)
	}
	return st.Ref
}
*/
