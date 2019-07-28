// step.go
package models

/*
import (
	// "encoding/json"
	"fmt"
	"regexp"
	// "github.com/json-iterator/go"
)

// var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Step struct {
	Base
	Extend
	Path          string        `json:"path" yaml:"path" toml:"path"`
	Method        string        `json:"method" yaml:"method" toml:"method"`
	Params        Params        `json:"params,omitempty"  yaml:"params" toml:"params"`
	Headers       Header        `json:"headers,omitempty"  yaml:"headers" toml:"headers"`
	MultipartFile MultipartFile `json:"multifiles,omitempty"   yaml:"multifiles" toml:"multifiles"`
	Validate      []Validator   `json:"validate"  yaml:"validate" toml:"validate"`
}

func (step *Step) GetName() string {
	return step.Name
}

func (step *Step) Json() string {
	jsonStr, err := json.Marshal(step)
	if err != nil {
		fmt.Println(`Step to json failed:`, err.Error())
		return `{}`
	}
	return string(jsonStr)
}

func (step *Step) GetType() int {
	return TYPE_STEP
}

func (step *Step) RefTag() string {
	if step.Ref == `` {
		return regexp.MustCompile(`\s+`).ReplaceAllString(step.Name, `_`)
	}
	return step.Ref
}
*/
