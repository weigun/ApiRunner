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
