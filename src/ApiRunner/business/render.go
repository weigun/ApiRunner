// testcase_render.go
package business

import (
	"ApiRunner/business/template"
	// "ApiRunner/models"

	"bytes"
	// "fmt"
	// "log"
	// "strings"
	"sync"

	"github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type renderer struct {
	sync.Mutex
	buf *bytes.Buffer
}

func newRenderer() *renderer {
	return &renderer{buf: bytes.NewBufferString(``)}
}

func (r *renderer) fillData(src string, data interface{}) []byte {
	r.Lock()
	defer r.Unlock()
	t := template.New().Funcs(funcMap)
	t.Parse(src)
	if data == nil {
		data = make(map[string]interface{})
	}
	t.Execute(r.buf, data)
	val := r.buf.Bytes()
	r.buf.Reset()
	return val
}

// func (r *renderer) render(source string, renderVars bool) []byte {
// 	/*
// 		渲染测试用例
// 		将用例中的模板格式全部转换为具体内容
// 		当renderVars=true时，表示需要渲染自定义变量，否则不渲染
// 	*/
// 	tmpl := template.New("testcase").Funcs(funcMap)
// 	wr := bytes.NewBufferString(``)
// 	if renderVars {
// 		// services.VarsMgr.Get(`key`)
// 		if strings.Index(source, `{{$`) != -1 {
// 			// 存在变量引用
// 			key := fmt.Sprintf(`.%s.`, r.tag)
// 			source := strings.Replace(source, `$`, key, -1)
// 			tmpl, err := tmpl.Parse(source)
// 			if err != nil {
// 				log.Fatalln(err.Error())
// 			}
// 			//从变量服务中取出需要的变量
// 			//TODO 如果没有找到变量，则懒加载？
// 			log.Println(r.tag)
// 			m := services.VarsMgr.GetByGroup(r.tag)
// 			varsMap := make(map[string]map[string]string)
// 			varsMap[r.tag] = m
// 			// spew.Dump(varsMap)
// 			tmpl.Execute(wr, varsMap)
// 		} else {
// 			tmpl, err := tmpl.Parse(source)
// 			if err != nil {
// 				log.Fatalln(err.Error())
// 			}
// 			tmpl.Execute(wr, nil)
// 		}
// 	} else {
// 		tmpl, err := tmpl.Parse(source)
// 		if err != nil {
// 			log.Fatalln(err.Error())
// 		}
// 		tmpl.Execute(wr, nil)
// 	}
// 	return wr.Bytes()

// }

// func (r *renderer) renderObj(source string, renderVars bool, modelPtr interface{}) error {
// 	objStr := r.render(source, renderVars)
// 	switch modelPtr.(type) {
// 	case *models.Executable:
// 		return json.Unmarshal(objStr, modelPtr.(*models.Executable))
// 	// case *models.Params:
// 	// return json.Unmarshal(objStr, modelPtr.(*models.Params))
// 	// case *models.Variables:
// 	case *map[string]interface{}:
// 		log.Println(`----------`, string(objStr))
// 		return json.Unmarshal(objStr, modelPtr.(*map[string]interface{}))
// 	case *models.MultipartFile:
// 		return json.Unmarshal(objStr, modelPtr.(*models.MultipartFile))
// 	default:
// 		log.Fatalln(fmt.Sprintf(`unknow model %T`, modelPtr))
// 	}
// 	return nil
// }

// func (r *renderer) renderValue(val string, renderVars bool) string {
// 	return string(r.render(val, renderVars))
// }

// func (r *renderer) renderWithData(source string, data interface{}) string {
// 	tmpl, err := template.New("testcase").Funcs(funcMap).Parse(source)
// 	if err != nil {
// 		log.Fatalln(err.Error())
// 	}
// 	wr := bytes.NewBufferString(``)
// 	tmpl.Execute(wr, data)
// 	return wr.String()
// }

// func fillData() {

// }
