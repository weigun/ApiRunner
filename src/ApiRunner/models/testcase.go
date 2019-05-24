package models

import (
	"encoding/json"
	"fmt"
)

type ICaseObj interface {
	GetName() string
	Json() string
}

type TestCase struct {
	Config CaseConfig `json:"config"`
	APIS   []API      `json:"apis"`
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
