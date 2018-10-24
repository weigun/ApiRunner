package testcase

import (
	service "ApiRunner/engine/services"
	"strconv"
	"text/template"
)

var randService = service.NewRandomService()

func randUser() string {
	return randService.GetAccount()

}

func randRange(min int64, max int64) string {
	return strconv.FormatInt(randService.GetRand(min, max), 10)
}

func getFuncMap() *template.FuncMap {
	_func := template.FuncMap{
		"randUser":  randUser,
		"randRange": randRange,
	}
	return &_func
}
