package utils

import (
	//validation "ApiRunner/validation"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

func Map2Json(m map[string]interface{}) string {
	jsonStr, err := json.Marshal(m)
	if err != nil {
		return "{}"
	}
	return string(jsonStr)
}

func Json2Map(js []byte) map[string]interface{} {
	var mapResult = make(map[string]interface{})
	//使用 json.Unmarshal(data []byte, v interface{})进行转换,返回 error 信息
	if err := json.Unmarshal(js, &mapResult); err != nil {
		return mapResult
	}
	return mapResult
}

func GetCwd() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0])) //返回绝对路径  filepath.Dir(os.Args[0])去除最后一个元素的路径
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

func GetErr(api string, code int) string {
	//请求错误码格式化
	return fmt.Sprintf(`%s failed,StatusCode is %d`, api, code)
}

func GetTemplate() *template.Template {
	_func := template.FuncMap{
		"randUser":  randUser,
		"randRange": randRange,
	}

	t := template.New("conf")
	t.Funcs(_func)
	return t
}

func Translate(data string) string {
	//将模板翻译
	tmpl := GetTemplate() //不能放到全局或者通过闭包的方式，因为这个是携程不安全的
	wr := bytes.NewBufferString("")
	tmpl, err := tmpl.Parse(data)
	if err != nil {
		log.Fatalln(err)
	}
	tmpl.Execute(wr, nil)
	return wr.String()
}

func TranslateValidata(data string, resp validation) string {
	//将模板翻译
	tmpl := getTemplate() //不能放到全局或者通过闭包的方式，因为这个是携程不安全的
	wr := bytes.NewBufferString("")
	tmpl, err := tmpl.Parse(data)
	if err != nil {
		log.Fatalln(err)
	}
	switch resp.Body.(type) {
	case respBodyMap:
		tmpl.Execute(wr, resp.Body.(respBodyMap))
	case respBodySlice:
		tmpl.Execute(wr, resp.Body.(respBodySlice))
	}
	return wr.String()
}
