// component.go
package models

import (
	"fmt"
)

type Component struct {
	Base
	API        API         `json:"api"  yaml:"api"`
	Props      Variables   `json:"props"  yaml:"props"`
	Components []Component `json:"components"  yaml:"components"`
	Validate   []Validator `json:"validate"  yaml:"validate" toml:"validate"`
	Event      Variables   `json:"event"  yaml:"event"`
}
