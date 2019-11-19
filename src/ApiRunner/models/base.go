// step.go
package models

import (
	"github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

const (
	TYPE_API = iota
	TYPE_PIPELINE
)

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
		log.Warning(`MultipartFile to json failed:`, err.Error())
		return `{}`
	}
	return string(jsonStr)
}

// type Module struct {
// 	Host  string     `json:"host"  yaml:"host" toml:"host"`
// 	Def   Variables  `json:"def"  yaml:"def"`
// 	Steps []ExecNode `json:"steps"  yaml:"steps"`
// }

type Executable interface {
	GetName() string
	Json() string
	GetType() int
}
