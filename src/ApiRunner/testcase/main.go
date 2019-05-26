// main.go
package main

import (
	"log"

	"encoding/json"
	// "errors"
	"fmt"
	"io/ioutil"

	// "os"
	"path/filepath"

	toml "github.com/BurntSushi/toml"
	"gopkg.in/yaml.v2"
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

type CaseConfig struct {
	Name      string      `json:"name"`
	Host      string      `json:"host"`
	Variables []Variables `json:"variables"`
}

type ICaseObj interface {
	GetName() string
	Json() string
}

type TestCase struct {
	Config CaseConfig `json:"config"`
	APIS   []API      `json:"apis"`
}

func (tc *TestCase) GetName() string {
	return tc.Config.Name
}

func (tc *TestCase) Json() string {
	jsonStr, err := json.Marshal(tc)
	if err != nil {
		fmt.Println(`testcase to json failed:`, err.Error())
		return `{}`
	}
	return string(jsonStr)
}

type CaseItem = map[string]TestCase

type TestSuites struct {
	Config   CaseConfig `json:"config"`
	CaseList []CaseItem `json:"case_list"`
}

func (ts *TestSuites) GetName() string {
	return ts.Config.Name
}

func (ts *TestSuites) Json() string {
	jsonStr, err := json.Marshal(ts)
	if err != nil {
		fmt.Println(`testcase to json failed:`, err.Error())
		return `{}`
	}
	return string(jsonStr)
}

func ParseTestCase(filePath string) (caseObj ICaseObj, err error) {
	// if !utils.Exists(filePath) {
	// 	err = errors.New(fmt.Sprintf(`testcase file not found,%s`, filePath))
	// 	return
	// }
	content, e := ioutil.ReadFile(filePath)
	if e != nil {
		err = e
		return
	}
	switch filepath.Ext(filePath) {
	case `.yaml`, `yml`:
		return parseYamlCase(content)
	case `json`:
		return parseJsonCase(content)
	case `toml`, `tml`:
		return parseTomlCase(content)
	default:
		return
	}
}

func parseYamlCase(content []byte) (caseObj ICaseObj, err error) {
	// m := []map[string]interface{}{}
	m := []map[string]interface{}{}
	err = yaml.Unmarshal(content, &m)
	if err != nil {
		log.Printf("parse yaml error: %v", err)
		return
	}
	fmt.Printf("--- m:\n%v\n\n", m)
	caseObj, err = _parseTestCase(m)
	return
}

func parseJsonCase(content []byte) (caseObj ICaseObj, err error) {
	m := []map[string]interface{}{}
	err = json.Unmarshal(content, &m)
	if err != nil {
		log.Printf("parse json error: %v", err)
		return
	}
	fmt.Printf("--- m:\n%v\n\n", m)
	caseObj, err = _parseTestCase(m)
	return
}

func parseTomlCase(content []byte) (caseObj ICaseObj, err error) {
	m := []map[string]interface{}{}
	err = toml.Unmarshal(content, &m)
	if err != nil {
		log.Printf("parse toml error: %v", err)
		return
	}
	fmt.Printf("--- m:\n%v\n\n", m)
	caseObj, err = _parseTestCase(m)
	return
}

func _parseTestCase(ts []map[string]interface{}) (caseObj ICaseObj, err error) {
	return
}

func main() {
	fmt.Println("Hello World!")
	ParseTestCase(`case.yaml`)
}
