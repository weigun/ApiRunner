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

type TestCase struct {
	Config CaseConfig `json:"config"   yaml:"config" toml:"config"`
	APIS   []API      `json:"apis"  yaml:"apis" toml:"apis"`
}

func (tc *TestCase) GetName() string {
	return tc.Config.Name
}

func (tc *TestCase) Json() string {
	jsonStr, err := json.Marshal(tc)
	if err != nil {
		fmt.Println(`testcase to json failed:`, err.Error())
		return `{}`
	}
	return string(jsonStr)
}

func (tc *TestCase) GetType() int {
	return TYPE_TESTCASE
}
