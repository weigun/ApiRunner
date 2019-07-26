// testcase_parser.go
package business

import (
	"ApiRunner/models"
	"ApiRunner/utils"
	"ApiRunner/utils/yaml"

	// "os"

	// "encoding/json"
	// "fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	// toml "github.com/BurntSushi/toml"
	"github.com/davecgh/go-spew/spew"
	// "github.com/json-iterator/go"

	// "gopkg.in/yaml.v2"
	"github.com/mitchellh/mapstructure"
)

// var json = jsoniter.ConfigCompatibleWithStandardLibrary

func ParsePipe(filePath string) (pipeObj models.IPipe) {
	//只解析pipeline和pipegroup
	m := require(filePath)
	// spew.Dump(m)
	pipeMap := _parsePipe(m)
	// log.Println(`*************************`)
	// spew.Dump(m)
	// spew.Dump(caseMap)
	if len(pipeMap) == 0 {
		return
	}
	pipeObj = toObj(pipeMap)
	spew.Dump(filePath)
	return
}

func _parsePipe(pipeMap map[string]interface{}) map[string]interface{} {
	if stages, ok := pipeMap[`stages`]; ok {
		parsePipeline(pipeMap)
	} else if pipelines, ok := pipeMap[`pipelines`]; ok {
		parsePipeGroup(pipeMap)
	}
	return pipeMap
}

func parsePipeline(pipeMap map[string]interface{}) {
	//pipeline没有继承
	for index, stage := range pipeMap[`stages`].([]interface{}) {
		stage := stage.(map[string]interface{})
		//处理stage的继承
		mergeExtends(stage)
		//处理step的继承
		for stepIndex, step := range stage[`steps`].([]interface{}) {
			step := step.(map[string]interface{})
			mergeExtends(step)
			stage[`steps`].([]interface{})[stepIndex] = step
		}
		// TODO 处理hook
		pipeMap[`stages`].([]interface{})[index] = stage
	}

}

func parsePipeGroup(pipeMap map[string]interface{}) {
	//pipegroup没有继承
	for index, pipeline := range pipeMap[`pipelines`].([]interface{}) {
		pipeline := pipeline.(map[string]interface{})
		parsePipeline(pipeline)
		pipeMap[`pipelines`].([]interface{})[index] = pipeline
	}
}

func mergeExtends(class map[string]interface{}) {
	if filePath, ok := class[`extends`]; ok {
		baseClass := require(filePath.(string))
		mergeExtends(baseClass)
		for k, v := range class {
			log.Println(`merge `, k)
			switch v.(type) {
			case []interface{}:
				//合并列表
				//stage的合并需要指定是头合并还是尾合并,默认是尾合并
				var mergeMode string
				if mergeMode, ok := class[`mergeMode`]; !ok {
					mergeMode = `tail`
				}
				itemListPtr := baseClass[k].([]interface{}) // 方便引用
				tmp := make([]interface{}, 0, len(itemListPtr)+len(v.([]interface{})))
				for _, ele := range v.([]interface{}) {
					if mergeMode == `tail` {
						itemListPtr = append(itemListPtr, ele)
					} else {
						tmp = append(tmp, ele)
					}
				}
				if mergeMode != `tail` {
					tmp = append(tmp, itemListPtr...)
					baseClass[k] = tmp
				} else {
					baseClass[k] = itemListPtr
				}

			case map[string]interface{}:
				//合并map
				if baseClass[k] == nil {
					//没有定义
					baseClass[k] = v.(map[string]interface{})
				} else {
					itemPtr := baseClass[k].(map[string]interface{})
					for sk, sv := range v.(map[string]interface{}) {
						itemPtr[sk] = sv
					}
					baseClass[k] = itemPtr //需要回写才能更新
				}
			default:
				baseClass[k] = v
			}
		}
		class = baseClass
	}

}

func require(casePath string) map[string]interface{} {

	// 将依赖的用例或者接口包含进当前用例

	log.Printf("ReadFile: %v", casePath)
	// TODO 需要设置根目录
	pathList := []string{utils.GetCwd()}
	pathList = append(pathList, casePath)
	casePathAbs := filepath.Join(pathList...)
	content, err := ioutil.ReadFile(casePathAbs)
	if err != nil {
		log.Fatal("ReadFile error:", err)
	}
	raw := make(map[string]interface{})
	raw = yaml.UnmarshalToMapStr(content)
	if err != nil {
		log.Fatal("parse yaml error:", err.Error())
	}
	switch filepath.Ext(casePath) {
	case `.yaml`, `yml`:
		raw = yaml.UnmarshalToMapStr(content)
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

func toObj(caseMap map[string]interface{}) models.IPipe {
	var isPipeGroup bool
	if _, ok := caseMap[`pipelines`]; ok {
		isPipeGroup = true
	}
	config := &mapstructure.DecoderConfig{
		TagName: "json",
	}
	if isPipeGroup {
		//用例集
		var pg models.PipeGroup
		config.Result = &pg
		decoder, err := mapstructure.NewDecoder(config)
		if err != nil {
			log.Panic(err.Error())
		}
		decoder.Decode(caseMap)
		// spew.Dump(pg)
		return &pg
	} else {
		//单个用例
		var pl models.Pipeline
		config.Result = &pl
		decoder, err := mapstructure.NewDecoder(config)
		if err != nil {
			log.Panic(err.Error())
		}
		decoder.Decode(caseMap)
		return &pl
	}
	return &models.Pipeline{}
}
