// testsuites
package models

import (
	"encoding/json"
	"fmt"
)

type CaseItem = map[string]TestCase

type TestSuites struct {
	Config   CaseConfig `json:"config"`
	CaseList []CaseItem `json:"case_list"`
}

func (ts *TestSuites) GetName() string {
	return ts.Config.Name
}

func (ts *TestSuites) Json() string {
	jsonStr, err := json.Marshal(ts)
	if err != nil {
		fmt.Println(`testcase to json failed:`, err.Error())
		return `{}`
	}
	return string(jsonStr)
}
