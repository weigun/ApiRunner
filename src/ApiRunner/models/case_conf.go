// case_conf.go
package models

import (
	"encoding/json"
	"fmt"
)

type CaseConfig struct {
	Name      string    `json:"name"  yaml:"name" toml:"name"`
	Host      string    `json:"host"  yaml:"host" toml:"host"`
	Variables Variables `json:"variables"  yaml:"variables" toml:"variables"`
}

func (cc *CaseConfig) Json() string {
	jsonStr, err := json.Marshal(cc)
	if err != nil {
		fmt.Println(`caseconfig to json failed:`, err.Error())
		return `{}`
	}
	return string(jsonStr)
}
