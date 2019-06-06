// testcase_render.go
package business

import (
	"ApiRunner/models"
	"strings"

	// "ApiRunner/services"
	"fmt"

	"bytes"
	"log"
	"text/template"

	"github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

//TODO 模板缓存

func render(source string, renderVars bool) []byte {
	/*
		渲染测试用例
		将用例中的模板格式全部转换为具体内容
		当renderVars=true时，表示需要渲染自定义变量，否则不渲染
	*/
	tmpl, err := template.New("testcase").Funcs(funcMap).Parse(source)
	if err != nil {
		log.Fatalln(err.Error())
	}
	wr := bytes.NewBufferString(``)
	if renderVars {
		// services.VarsMgr.Get(`key`)
		if strings.Index(source, `{{$`) != -1 {
			// 存在变量引用
			source := strings.Replace(source, `$`, `.`, -1)

		} else {

		}
		tmpl.Execute(wr, nil)
	} else {
		tmpl.Execute(wr, nil)
	}
	return wr.Bytes()

}

func RenderObj(source string, renderVars bool, modelPtr interface{}) error {
	objStr := render(source, renderVars)
	switch modelPtr.(type) {
	case models.CaseConfig:
		return json.Unmarshal(objStr, modelPtr.(models.CaseConfig))
	case models.ICaseObj:
		return json.Unmarshal(objStr, modelPtr.(models.ICaseObj))
	default:
		log.Fatalln(fmt.Sprintf(`unknow model %T`, modelPtr))
	}
	return nil
}

func RenderValue(val string, renderVars bool) string {
	return string(render(val, renderVars))
}
