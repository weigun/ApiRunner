// record.go
package report

import (
	"fmt"
	"time"
)

const (
	SUCCESS = iota
	FAILED
	ERROR
	SKIP
)

type DateTime = time.Time

type Summary struct {
	StartAt  DateTime
	Duration int64
	Status   []Status
}

type Status struct {
	Type                         string
	Success, Failed, Error, Skip int64
}
