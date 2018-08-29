package engine

import (
	testcase "ApiRunner/case"
	runner "ApiRunner/runner"
)

type PImailBoxInterface interface {
	GetTestcase() testcase.PIparserInsterface
}

type PIreportInterface interface {
	PImailBoxInterface
}
