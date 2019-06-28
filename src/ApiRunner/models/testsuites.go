// testsuites
package models

import (
	// "encoding/json"
	"fmt"
)

type CaseItem = map[string]TestCase

type TestSuites struct {
	Config   CaseConfig `json:"config"  yaml:"config" toml:"config"`
	CaseList []CaseItem `json:"testcases"  yaml:"testcases" toml:"testcases"`
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

func (ts *TestSuites) GetType() int {
	return TYPE_TESTSUITS
}
