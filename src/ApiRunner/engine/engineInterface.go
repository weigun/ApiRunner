package engine

import (
	testcase "ApiRunner/case"
)

type PImailBoxInterface interface {
	GetTestcase() testcase.PIparserInsterface
}

type PIreportInterface interface {
	PImailBoxInterface
}
