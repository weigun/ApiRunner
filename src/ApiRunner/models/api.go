// api.go
package models

import (
	// "encoding/json"
	"fmt"

	"github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type API struct {
	Name          string        `json:"name" yaml:"name" toml:"name"`
	Variables     Variables     `json:"variables" yaml:"variables" toml:"variables"`
	Path          string        `json:"path" yaml:"path" toml:"path"`
	Method        string        `json:"method" yaml:"method" toml:"method"`
	Headers       Header        `json:"headers,omitempty"  yaml:"headers" toml:"headers"`
	Params        Params        `json:"params,omitempty"  yaml:"params" toml:"params"`
	Export        Variables     `json:"export"  yaml:"export" toml:"export"`
	MultipartFile MultipartFile `json:"multifiles,omitempty"   yaml:"multifiles" toml:"multifiles"`
	Validate      []Validator   `json:"validate"  yaml:"validate" toml:"validate"`
	BeforeRequest string        `json:"beforeRequest"  yaml:"beforeRequest" toml:"beforeRequest"`
	AfterResponse string        `json:"afterResponse"  yaml:"afterResponse" toml:"afterResponse"`
}

func (api *API) GetName() string {
	return api.Name
}

func (api *API) Json() string {
	jsonStr, err := json.Marshal(api)
	if err != nil {
		fmt.Println(`testcase to json failed:`, err.Error())
		return `{}`
	}
	return string(jsonStr)
}

func (api *API) GetType() int {
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
