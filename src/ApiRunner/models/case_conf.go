// case_conf.go
package models

type CaseConfig struct {
	Name      string    `json:"name"  yaml:"name" toml:"name"`
	Host      string    `json:"host"  yaml:"host" toml:"host"`
	Variables Variables `json:"variables"  yaml:"variables" toml:"variables"`
}
