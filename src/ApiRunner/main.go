package main

import (
	"ApiRunner/business"
	"ApiRunner/utils"

	// "fmt"

	// "ApiRunner/services"
	"path/filepath"
	"time"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	// initGlobalComponents()
	// bootstrap()
	// caseObj := business.ParseTestCase(`signup_case.yaml`)
	// spew.Dump(caseObj)
	// runner := business.NewTestRunner(`signup`, caseObj)
	// runner.Start()

	caseObj2 := business.ParseTestCase(filepath.Join(utils.GetCwd(), `testcase`, `all.yaml`))
	spew.Dump(caseObj2)
	runner2 := business.NewTestRunner(`signup`, caseObj2)
	runner2.Start()
	time.Sleep(1000 * time.Second)
}
