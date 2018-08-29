package engine

import (
	testcase "ApiRunner/case"
	runner "ApiRunner/runner"
)

type PImailBoxInterface interface {
	GetTestcase() testcase.PIparserInsterface
	GetResult() runner.PIresponseInterface
}

type PIreportInterface interface {
	PImailBoxInterface
}
