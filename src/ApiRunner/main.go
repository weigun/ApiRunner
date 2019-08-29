package main

import (
	"ApiRunner/business"
	// "ApiRunner/utils"

	// "fmt"

	// "ApiRunner/services"
	// "path/filepath"
	"ApiRunner/business/template/parser"
	"time"

	"github.com/davecgh/go-spew/spew"
)

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
	input := `my email is ${refs.user1.email},my luckly number is ${get_luckly_from_name($name)},and ${age} years old`
	// o, _ := parser.Parse(input)
	t := parser.Tree{}
	t.Parse(input)
	spew.Dump(t)
	time.Sleep(1000 * time.Second)
	pipObj := business.ParsePipe(`testcase\components\suits.yaml`)
	runner := business.NewTestRunner(`signup`, pipObj)
	runner.Start()
	time.Sleep(10 * time.Second)
}
