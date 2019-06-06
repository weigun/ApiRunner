// testcase_parser.go
package business

import (
	"ApiRunner/models"
	"log"

	// "encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	// toml "github.com/BurntSushi/toml"
	"github.com/davecgh/go-spew/spew"
	// "github.com/json-iterator/go"
	"gopkg.in/yaml.v2"
)

// var json = jsoniter.ConfigCompatibleWithStandardLibrary

func ParseTestCase(filePath string) (caseObj models.ICaseObj) {
	m := require(filePath)
	caseMap := _parseTestCase(m)
	if len(caseMap) == 0 {
		return
	}
	caseObj = toObj(caseMap, filepath.Ext(filePath))
	spew.Dump(caseObj)
	return
}

func _parseTestCase(ts map[string]interface{}) (caseMap map[string]interface{}) {
	if _, ok := ts[`apis`]; ok {
		// testcase here
		parseSingleCase(ts)

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
		// caseMap = ts
	}
	caseMap = ts
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
	default:
		log.Fatal("not support case format:", filepath.Ext(casePath))
	}
	return raw
}

func toObj(caseMap map[string]interface{}, ext string) models.ICaseObj {
	var isTestSuits bool
	if _, ok := caseMap[`testcases`]; ok {
		isTestSuits = true
	}
	switch ext {
	case `.yaml`, `.yml`:
		byteCaseMap, _ := yaml.Marshal(caseMap)
		if isTestSuits {
			//用例集
			var ts models.TestSuites
			yaml.Unmarshal(byteCaseMap, &ts)
			return &ts
		} else {
			//单个用例
			var ts models.TestCase
			yaml.Unmarshal(byteCaseMap, &ts)
			return &ts
		}
	case `.json`, `.conf`:
		byteCaseMap, _ := json.Marshal(caseMap)
		if isTestSuits {
			//用例集
			var ts models.TestSuites
			json.Unmarshal(byteCaseMap, &ts)
			return &ts
		} else {
			//单个用例
			var ts models.TestCase
			json.Unmarshal(byteCaseMap, &ts)
			return &ts
		}
	default:
		log.Fatal(`toObj failed,not support case format `, ext)
	}
	return &models.TestCase{}

}
