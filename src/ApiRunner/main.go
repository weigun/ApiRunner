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
	runner := &business.TestRunner{ID: `signup`, CaseObj: caseObj}
	runner.Start()
	time.Sleep(1000 * time.Second)
}
