// record.go
package report

import (
	//std
	// "fmt"
	"time"

	//third party
	"github.com/json-iterator/go"

	. "ApiRunner/models"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

const (
	SUCCESS = iota
	FAILED
	ERROR
	SKIP

	TESTCASE = iota + 100
	TESTSUITS
)

type DateTime = time.Time

type Report struct {
	Summary Summary  `json:"summary"`
	Details []Detail `json:"details"`
}

type Summary struct {
	StartAt  DateTime `json:"startAt"`
	Duration int64    `json:"duration"`
	Status   []Status `json:"status"`
}

type Status struct {
	// Type                         string
	Success int64 `json:"success"`
	Failed  int64 `json:"failed"`
	Error   int64 `json:"error"`
	Skip    int64 `json:"skip"`
}

func (s *Status) Count(stat int64) {
	switch stat {
	case SUCCESS:
		s.Success += 1
	case FAILED:
		s.Failed += 1
	case ERROR:
		s.Error += 1
	case SKIP:
		s.Skip += 1
	}
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

////////////////////////////////////////////////////////
//report
func NewReport() *Report {
	return &Report{}
}

func (rp *Report) SetSummary(sum Summary) {
	rp.Summary = sum
}

func (rp *Report) SetDetails(details []Detail) {
	rp.Details = details
}

func (rp *Report) AddDetail(detail Detail) {
	rp.Details = append(rp.Details, detail)
}

func (rp *Report) Json() string {
	jsonStr, err := json.Marshal(rp)
	if err != nil {
		fmt.Println(`Report to json failed:`, err.Error())
		return `{}`
	}
	return string(jsonStr)
}

// Summary

func NewSummary() *Summary {
	return &Summary{Status: make([]Status{}, 2)}
}

func (sum *Summary) Counter(which int64) *Status {
	if which == TESTSUITS {
		return *sum.Status[0]
	}
	return *sum.Status[1]
}

// Details
func NewDetail() *Detail {
	return &Detail{Status: Status{}}
}

func (dt *Detail) SetRecords(records []Record) {
	dt.Records = records
}

func (dt *Detail) AddRecord(record Record) {
	dt.Records = append(dt.Records, record)
}

// Record
func NewRecord() *Record {
	return &Record{}
}

func name() {

}
