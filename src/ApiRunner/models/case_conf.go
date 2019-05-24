// case_conf.go
package models

type CaseConfig struct {
	Name      string      `json:"name"`
	Host      string      `json:"host"`
	Variables []Variables `json:"variables"`
}
