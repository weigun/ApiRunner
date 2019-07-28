// api.go
package models

import (
	// "encoding/json"
	"fmt"

	"github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type API struct {
	Name string `json:"name" yaml:"name" toml:"name"`
	// Variables     Variables     `json:"variables" yaml:"variables" toml:"variables"`
	Host          string        `json:"host"  yaml:"host" toml:"host"`
	Path          string        `json:"path" yaml:"path" toml:"path"`
	Method        string        `json:"method" yaml:"method" toml:"method"`
	Headers       Header        `json:"headers,omitempty"  yaml:"headers" toml:"headers"`
	Params        Params        `json:"params,omitempty"  yaml:"params" toml:"params"`
	MultipartFile MultipartFile `json:"multifiles,omitempty"   yaml:"multifiles" toml:"multifiles"`
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
