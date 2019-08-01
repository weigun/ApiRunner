// node
package models

import (
	"regexp"
)

type ExecNode struct {
	Desc     string      `json:"desc" yaml:"desc"`
	Host     string      `json:"host"  yaml:"host"`
	Ref      string      `json:"ref"  yaml:"ref"`
	Exec     Executable  `json:"exec"  yaml:"exec"`
	Args     Variables   `json:"args"  yaml:"args"`
	Export   Variables   `json:"export"  yaml:"export"`
	Validate []Validator `json:"validate"  yaml:"validate"`
	Hooks    Variables   `json:"hooks"  yaml:"hooks"`
	Retry    int         `json:"retry"  yaml:"retry"`
	Repeat   int         `json:"repeat"  yaml:"repeat"`
}

func (exec *ExecNode) RefTag() string {
	if exec.Ref == `` {
		return regexp.MustCompile(`\s+`).ReplaceAllString(exec.Desc, `_`)
	}
	return exec.Ref
}
