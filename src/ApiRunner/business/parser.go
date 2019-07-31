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

type strfacemap = map[string]interface{}

func ParsePipe(filePath string) (pipeObj models.Executable) {
	//只有pipeline才可以解析
	m := require(filePath)
	// spew.Dump(m)
	pipeObj = _parsePipe(m)
	// log.Println(`*************************`)
	// spew.Dump(m)
	spew.Dump(pipeObj)
	spew.Dump(filePath)
	return
}

func _parsePipe(pipeMap strfacemap) *models.Pipeline {
	var pipeObj models.Pipeline
	pipeObj.Name = pipeMap[`name`].(string)
	if host, ok := pipeMap[`host`].(string); ok {
		pipeObj.Host = host
	}
	module := pipeMap[`module`].(strfacemap)
	pipeObj.Def = module[`def`].(models.Variables)
	nodes := []interface{}{}
	if _, ok := module[`parallel`]; ok {
		pipeObj.Parallel = true
		nodes = module[`parallel`].([]interface{})
	} else {
		nodes = module[`steps`].([]interface{})
	}
	config := &mapstructure.DecoderConfig{
		TagName: "json",
	}
	for _, node := range nodes {
		node := node.(strfacemap)[`node`].(strfacemap)
		var nodeObj models.ExecNode
		config.Result = &nodeObj
		decoder, err := mapstructure.NewDecoder(config)
		if err != nil {
			log.Panic(err.Error())
		}
		decoder.Decode(node)

		m := require(node[`exec`].(string))
		if _, ok := m[`module`]; ok {
			var out models.Pipeline
			mustToObj(m, &out)
			nodeObj.Exec = models.Executable(&out)
		} else {
			var out models.API
			mustToObj(m, &out)
			nodeObj.Exec = models.Executable(&out)
		}

		pipeObj.Steps = append(pipeObj.Steps, nodeObj)
	}
	return &pipeObj
}

/*
func parsePipeline(pipeMap strfacemap) {
	//pipeline没有继承
	for index, stage := range pipeMap[`stages`].([]interface{}) {
		stage := stage.(strfacemap)
		//处理stage的继承
		mergeExtends(stage)
		//处理step的继承
		for stepIndex, step := range stage[`steps`].([]interface{}) {
			step := step.(strfacemap)
			mergeExtends(step)
			stage[`steps`].([]interface{})[stepIndex] = step
		}
		// TODO 处理hook
		pipeMap[`stages`].([]interface{})[index] = stage
	}

}

func parsePipeGroup(pipeMap strfacemap) {
	//pipegroup没有继承
	for index, pipeline := range pipeMap[`pipelines`].([]interface{}) {
		pipeline := pipeline.(strfacemap)
		parsePipeline(pipeline)
		pipeMap[`pipelines`].([]interface{})[index] = pipeline
	}
}

func mergeExtends(class strfacemap) {
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

			case strfacemap:
				//合并map
				if baseClass[k] == nil {
					//没有定义
					baseClass[k] = v.(strfacemap)
				} else {
					itemPtr := baseClass[k].(strfacemap)
					for sk, sv := range v.(strfacemap) {
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
*/
func require(casePath string) strfacemap {

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
	raw := make(strfacemap)
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

func mustToObj(in strfacemap, out interface{}) {
	config := &mapstructure.DecoderConfig{
		TagName: "json",
	}
	config.Result = out
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		log.Panic(err.Error())
	}
	decoder.Decode(in)
}
