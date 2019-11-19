package business

import (
	//third party
	//project self
	_ "ApiRunner/services"
	"ApiRunner/utils/logger"
)

var log = logger.GetLogger(nil, `business`)

func init() {
	log.Info(`init business layer`)
}
