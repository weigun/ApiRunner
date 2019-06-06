package main

import (
	"bytes"
	"log"
	"text/template"
)

//TODO 模板缓存

func render(source string, renderVars bool) string {
	/*
		渲染测试用例
		将用例中的模板格式全部转换为具体内容
		当renderVars=true时，表示需要渲染自定义变量，否则不渲染
	*/
	tmpl, err := template.New("testcase").Parse(source)
	if err != nil {
		log.Println(tmpl)
		log.Fatalln(err.Error())
	}
	wr := bytes.NewBufferString(``)
	if renderVars {
		// services.VarsMgr.Get(`key`)
		tmpl.Execute(wr, nil)
	} else {
		tmpl.Execute(wr, nil)
	}
	return wr.String()

}

func main() {
	render(`hello {{name}}`, false)
}
