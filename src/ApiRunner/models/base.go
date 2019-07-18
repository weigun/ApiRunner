// step.go
package models

import (
	// "encoding/json"
	"fmt"

	"github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Base struct {
	Name string `json:"name" yaml:"name"`
	Ref  string `json:"ref" yaml:"ref"`
}

type Extend struct {
	Extends string `json:"extends" yaml:"extends"`
	Status  int    `json:"status" yaml:"status"`
}

type Variables = map[string]interface{}

type Header = Variables

type Params = Variables

type Validator struct {
	Check    string      `json:"check"  yaml:"check" toml:"check"`
	Op       string      `json:"op"  yaml:"op" toml:"op"`
	Actual   interface{} `json:"actual"  yaml:"actual" toml:"actual"`
	Expected interface{} `json:"expected"  yaml:"expected" toml:"expected"`
}

type MultipartFile struct {
	Params Params `json:"params"  yaml:"params" toml:"params"` //上传的数据
	Files  Params `json:"files"   yaml:"files" toml:"files"`   //文件列表
}

func (mf *MultipartFile) IsEnabled() bool {
	return len(mf.Params) > 0 || len(mf.Files) > 0
}

func (mf *MultipartFile) Json() string {
	jsonStr, err := json.Marshal(mf)
	if err != nil {
		fmt.Println(`MultipartFile to json failed:`, err.Error())
		return `{}`
	}
	return string(jsonStr)
}

type callbacks = []string

type Event struct {
	onSuccess callbacks `json:"onSuccess"  yaml:"onSuccess"`
	onFailure callbacks `json:"onFailure"  yaml:"onFailure"`
}

type StageEvent struct {
	Event
	onStepFailure callbacks `json:"onStepFailure"  yaml:"onStepFailure"`
	onStepSuccess callbacks `json:"onStepSuccess"  yaml:"onStepSuccess"`
}
