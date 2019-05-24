// api.go
package models

import (
	"encoding/json"
	"fmt"
)

type API struct {
	Name          string        `json:"name"`
	Variables     []Variables   `json:"variables"`
	Path          string        `json:"path"`
	Method        string        `json:"method"`
	Headers       []Header      `json:"headers,omitempty"`
	Params        Params        `json:"params,omitempty"`
	Export        []Variables   `json:"export"`
	MultipartFile MultipartFile `json:"files,omitempty"`
	Validate      []Validator   `json:"validate"`
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

type Variables struct {
	Name string
	Val  interface{}
}

type Header struct {
	Key, Val string
}

type Params = map[string]interface{}

type Validator struct {
	Op       string
	Source   interface{}
	Verified interface{}
}

type MultipartFile struct {
	Params Params //上传的数据
	Files  Params //文件列表
}
