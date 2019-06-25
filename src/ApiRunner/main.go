package main

import (
	"ApiRunner/business"
	// "ApiRunner/services"
	"time"
)

func main() {
	// initGlobalComponents()
	// bootstrap()
	caseObj := business.ParseTestCase(`signup_case.yaml`)
	runner := business.NewTestRunner(`signup`, caseObj)
	runner.Start()
	time.Sleep(1000 * time.Second)
}
