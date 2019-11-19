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
	"path/filepath"

	// toml "github.com/BurntSushi/toml"
	// "github.com/davecgh/go-spew/spew"
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
	// log.Info(`*************************`)
	// spew.Dump(m)
	// spew.Dump(pipeObj)
	// spew.Dump(filePath)
	return
}

func _parsePipe(pipeMap strfacemap) *models.Pipeline {
	var pipeObj models.Pipeline
	pipeObj.Name = pipeMap[`name`].(string)
	module := pipeMap[`module`].(strfacemap)
	if host, ok := module[`host`].(string); ok {
		pipeObj.Host = host
	}
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
		// log.Info(`raw execnode:`, spew.Sdump(node))
		var nodeObj models.ExecNode
		config.Result = &nodeObj
		decoder, err := mapstructure.NewDecoder(config)
		if err != nil {
			log.Fatal(err.Error())
		}
		decoder.Decode(node)
		//require后，所有的值类型都是string，所以这里需要对一些number的字段进行转换
		if utils.ToNumber(node[`retry`]) == nil {
			nodeObj.Retry = 0
		} else {
			nodeObj.Retry = int(utils.ToNumber(node[`retry`]).(int64))
		}
		if utils.ToNumber(node[`repeat`]) == nil {
			nodeObj.Repeat = 0
		} else {
			nodeObj.Repeat = int(utils.ToNumber(node[`repeat`]).(int64))
		}

		m := require(node[`exec`].(string))
		if _, ok := m[`module`]; ok {
			var out *models.Pipeline
			out = _parsePipe(m)
			nodeObj.Exec = models.Executable(out)
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

	log.Info("ReadFile:", casePath)
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
		log.Fatal(err.Error())
	}
	decoder.Decode(in)
}
