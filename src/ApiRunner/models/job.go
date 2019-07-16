// job.go
package models

import (
	// "encoding/json"
	"fmt"

	"github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Job struct {
	Name          string        `json:"name" yaml:"name" toml:"name"`
	Path          string        `json:"path" yaml:"path" toml:"path"`
	Method        string        `json:"method" yaml:"method" toml:"method"`
	Ref           string        `json:"ref" yaml:"ref"`
	Extends       string        `json:"extends" yaml:"extends"`
	Import        string        `json:"import" yaml:"import"`
	Func          Variables     `json:"func,omitempty"  yaml:"func"`
	Return        Variables     `json:"return,omitempty"  yaml:"return"`
	Headers       Header        `json:"headers,omitempty"  yaml:"headers" toml:"headers"`
	MultipartFile MultipartFile `json:"multifiles,omitempty"   yaml:"multifiles" toml:"multifiles"`
	Validate      []Validator   `json:"validate"  yaml:"validate" toml:"validate"`
	Event         Event         `json:"when,omitempty"  yaml:"when"`
}

func (job *Job) GetName() string {
	return job.Name
}

func (job *Job) Json() string {
	jsonStr, err := json.Marshal(job)
	if err != nil {
		fmt.Println(`testcase to json failed:`, err.Error())
		return `{}`
	}
	return string(jsonStr)
}

func (job *Job) GetType() int {
	return TYPE_API
}

// type Variables struct {
// 	Name string
// 	Val  interface{}
// }
type Variables = map[string]interface{}

type Header = map[string]interface{}

// type Header struct {
// 	Key, Val string
// }

type Params = map[string]interface{}

type Event = map[string]interface{}

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
