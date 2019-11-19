// initializer.go
package main

import (
	_ "ApiRunner/business"
	_ "ApiRunner/dao"
	_ "ApiRunner/models"
	"ApiRunner/utils/logger"
)

var log = logger.GetLogger(nil, `MAIN`)

func initGlobalComponents() {
	log.Info(`initGlobalComponents`)
}
