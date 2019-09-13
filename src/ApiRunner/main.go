package main

import (
	"bytes"
	"time"

	"fmt"
	// "path/filepath"

	"ApiRunner/business"
	// "ApiRunner/utils"
	// "ApiRunner/services"
	"ApiRunner/business/template"
	// "ApiRunner/business/template/lexer"

	"github.com/davecgh/go-spew/spew"
)

func get_luckly_from_name(name string) string {
	switch name {
	case `weigun`:
		return `666`
	default:
		return `233`
	}
}

func main() {
	// initGlobalComponents()
	// bootstrap()
	// caseObj := business.ParseTestCase(filepath.Join(utils.GetCwd(), `testcase`, `signup_case.yaml`))
	/*
		caseObj := business.ParseTestCase(`testcase\signup_case.yaml`)
		spew.Dump(caseObj)
		runner := business.NewTestRunner(`signup`, caseObj)
		runner.Start()

		// caseObj2 := business.ParseTestCase(filepath.Join(utils.GetCwd(), `testcase`, `all.yaml`))
		caseObj2 := business.ParseTestCase(`testcase\all.yaml`)
		spew.Dump(caseObj2)
		runner2 := business.NewTestRunner(`all`, caseObj2)
		runner2.Start()
		time.Sleep(1000 * time.Second)
	*/
	fnMap := make(template.FuncMap)
	fnMap[`get_luckly_from_name`] = get_luckly_from_name
	input := `my email is ${refs.user1.email},my luckly number is ${get_luckly_from_name($name)},and ${age} years old`
	t := template.New().Funcs(fnMap)
	t.Parse(input)
	spew.Dump(t)
	wr := bytes.NewBufferString(``)
	data := make(map[string]interface{})
	subData := make(map[string]interface{})
	subData[`user1`] = map[string]interface{}{`email`: `283257958@qq.com`}
	data[`refs`] = subData
	data[`age`] = 6
	data[`name`] = `weigun`
	t.Execute(wr, data)
	fmt.Println(wr.String())
	time.Sleep(1000 * time.Second)
	pipObj := business.ParsePipe(`testcase\components\suits.yaml`)
	runner := business.NewTestRunner(`signup`, pipObj)
	runner.Start()
	time.Sleep(10 * time.Second)
}
