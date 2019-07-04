// record.go
package report

import (
	// "fmt"
	"time"

	. "ApiRunner/models"
)

const (
	SUCCESS = iota
	FAILED
	ERROR
	SKIP
)

type DateTime = time.Time

type Report struct {
	Summary Summary
	Details []Detail
}

type Summary struct {
	StartAt  DateTime
	Duration int64
	Status   []Status
}

type Status struct {
	// Type                         string
	Success, Failed, Error, Skip int64
}

type Detail struct {
	Title   string
	Status  Status
	Records []Record
}

type Record struct {
	Stat       int64
	Desc       string
	Elapsed    int64 //ms
	Request    DataMap
	Response   DataMap
	Validators []Validator
}

type DataMap = map[string]interface{}
