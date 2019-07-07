// record.go
package report

import (
	"io/ioutil"
	//std
	// "fmt"
	"net/http"
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
	Request    *http.Request
	Response   *http.Response
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

func (rc *Record) SetValidators(vds []Validator) {
	rc.Validators = vds
}

func (rc *Record) AddValidator(vd Validator) {
	rc.Validators = append(rc.Validators, vd)
}

func (rc *Record) AddValidator(vd Validator) {
	rc.Validators = append(rc.Validators, vd)
}

func (rc *Record) MarshalJSON() ([]byte, error) {
	//自定义编组过程
	dict := make(DataMap{})
	req := make(DataMap{})
	resp := make(DataMap{})
	dict[`stat`] = rc.Stat
	dict[`desc`] = rc.Desc
	dict[`elapsed`] = rc.Elapsed
	req[`url`] = rc.Request.URL.String()
	req[`method`] = rc.Request.Method
	req[`header`] = rc.Request.Header
	//handle request playload
	body, err := rc.Request.GetBody()
	if err != nil {
		panic(err)
	}
	req[`body`] = string(ioutil.ReadAll(body))
	dict[`request`] = req

	resp[`url`] = rc.Response.Request.URL.String()
	resp[`statusCode`] = rc.Response.StatusCode
	resp[`cookies`] = rc.Response.Cookies() //maybe each cookie call string()
	resp[`header`] = rc.Response.Header
	resp[`body`] = string(ioutil.ReadAll(rc.Response.Body))
	dict[`response`] = resp
	dict[`validators`] = rc.Validators
	return json.Marshal(dict)
}