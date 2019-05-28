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
	"github.com/davecgh/go-spew/spew"
	"gopkg.in/yaml.v2"
)

type API struct {
	Name          string        `json:"name"`
	Variables     Variables     `json:"variables"`
	Path          string        `json:"path"`
	Method        string        `json:"method"`
	Headers       Header        `json:"headers,omitempty"`
	Params        Params        `json:"params,omitempty"`
	Export        Variables     `json:"export"`
	MultipartFile MultipartFile `json:"multifiles,omitempty"`
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
	Op       string      `json:"op"`
	Source   interface{} `json:"source"`
	Verified interface{} `json:"verified"`
}

type MultipartFile struct {
	Params Params `json:"params"` //上传的数据
	Files  Params `json:"files"`  //文件列表
}

type CaseConfig struct {
	Name      string    `json:"name"`
	Host      string    `json:"host"`
	Variables Variables `json:"variables"`
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
	CaseList []CaseItem `json:"testcases"`
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
	m := require(filePath)
	caseMap := _parseTestCase(m)
	_ = caseMap
	spew.Dump(caseMap)
	return
}

// func _parseTestCase(ts map[string]interface{}) (caseObj ICaseObj, err error) {
func _parseTestCase(ts map[string]interface{}) (caseMap map[string]interface{}) {
	// caseConf := ts[`config`]
	if _, ok := ts[`apis`]; ok {
		// testcase here
		parseSingleCase(ts)
		caseMap = ts

	} else if _, ok := ts[`testcases`]; ok {
		// testsuits here
		for index, caseInfo := range ts[`testcases`].([]interface{}) {
			caseInfo := caseInfo.(map[interface{}]interface{})
			for caseDesc, caseFile := range caseInfo {
				m := require(caseFile.(string))
				log.Println(`testsuit`)
				oneCase := _parseTestCase(m) //递归解析用例集
				//oneCase是一个测试集中一个用例的map
				caseInfo[caseDesc] = oneCase
			}
			ts[`testcases`].([]interface{})[index] = caseInfo
		}
		caseMap = ts
	}
	// spew.Dump(caseMap)
	return
}

func parseSingleCase(ts map[string]interface{}) {
	for index, apiItem := range ts[`apis`].([]interface{}) {
		//遍历接口列表，对rawApi的成员进行替换

		apiItem := apiItem.(map[interface{}]interface{})
		rawApi := require(apiItem[`api`].(string))
		for key, val := range apiItem {
			/*
				替换规则：
				非列表、map等结构直接替换
				列表和map则进行合并处理
			*/
			key := key.(string)
			if _, ok := rawApi[key]; ok {
				//只替换存在的字段
				switch val.(type) {
				case int, int8, int16, int32, int64, float32, float64, string, bool:
					//直接替换
					rawApi[key] = val
				case []interface{}:
					//合并列表
					itemListPtr := rawApi[key].([]interface{}) // 方便引用
					for _, ele := range val.([]interface{}) {
						itemListPtr = append(itemListPtr, ele)
					}
					rawApi[key] = itemListPtr
				case map[interface{}]interface{}:
					//合并map
					itemPtr := rawApi[key].(map[interface{}]interface{}) //因为多级map不可寻址，需要先提取整个val出来才能引用
					for k, v := range val.(map[interface{}]interface{}) {
						itemPtr[k.(string)] = v
					}
					rawApi[key] = itemPtr //需要回写才能更新
				default:
					log.Fatal(fmt.Sprintf(`_parseTestCase,unsupport type in case type:%T,val:%v`, val, val))
				}
			}
		}
		ts[`apis`].([]interface{})[index] = rawApi
	}
}

func require(casePath string) map[string]interface{} {

	// 将依赖的用例或者接口包含进当前用例

	log.Printf("ReadFile: %v", casePath)
	// TODO 需要设置根目录
	content, err := ioutil.ReadFile(casePath)
	if err != nil {
		log.Fatal("ReadFile error:", err)
	}
	raw := make(map[string]interface{})
	switch filepath.Ext(casePath) {
	case `.yaml`, `yml`:
		err = yaml.Unmarshal(content, &raw)
		if err != nil {
			log.Fatal("parse yaml error:", err.Error())
		}
	case `.json`, `.conf`:
		err = json.Unmarshal(content, &raw)
		if err != nil {
			log.Fatal("parse json error:", err.Error())
		}
	case `.toml`, `.tml`:
		err = toml.Unmarshal(content, &raw)
		if err != nil {
			log.Fatal("parse toml error:", err.Error())
		}
	default:
		log.Fatal("not support case format:", filepath.Ext(casePath))
	}
	return raw
}

func main() {
	fmt.Println("Hello World!")
	ParseTestCase(`signup.conf`)
	ParseTestCase(`signup.yaml`)
	ParseTestCase(`signup_case.yaml`)
	ParseTestCase(`all.yaml`)

}
