package testcase

import (
	mgr "ApiRunner/manager"
	"strconv"
)

var randService = NewRandomService()

func randUser() string {
	return randService.getAccount()

}

func randRange(min int64, max int64) string {
	return strconv.FormatInt(randService.getRand(min, max), 10)
}
