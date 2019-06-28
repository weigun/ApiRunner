package main

import (
	"ApiRunner/business"
	// "fmt"

	// "ApiRunner/services"
	"time"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	// initGlobalComponents()
	// bootstrap()
	caseObj := business.ParseTestCase(`signup_case.yaml`)
	spew.Dump(caseObj)
	runner := business.NewTestRunner(`signup`, caseObj)
	runner.Start()
	time.Sleep(1000 * time.Second)
}
