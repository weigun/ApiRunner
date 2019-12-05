// record.go
package models

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"ApiRunner/models/young"
)

// var json = jsoniter.ConfigCompatibleWithStandardLibrary

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
	Summary Summary `json:"summary"`
	// Details []Detail `json:"details"`
	Details *ResultTree `json:"details"`
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
	default:
		log.Warning(fmt.Sprintf(`unknow stat %T,%v`, stat, stat))
	}
}

func (s *Status) Total() int64 {
	return s.Success + s.Error + s.Failed + s.Skip
}

type Detail struct {
	Title  string
	Status int64
	Record *Record
}

type Record struct {
	Stat       int64
	Desc       string
	Elapsed    int64 //ms
	Request    *http.Request
	Response   *young.Response
	Validators []Validator
}

type ResultTree struct {
	*Detail
	parent   *ResultTree
	children []*ResultTree
}

func (rt *ResultTree) Parent() *ResultTree {
	return rt.parent
}

func (rt *ResultTree) SetParent(result *ResultTree) {
	rt.parent = result
}

func (rt *ResultTree) Append(result *ResultTree) {
	result.SetParent(rt)
	rt.children = append(rt.children, result)
}

func (rt *ResultTree) ChildAt(index int) *ResultTree {
	if index > len(rt.children) {
		panic(`IndexError: list assignment index out of range`)
	}
	return rt.children[index]
}

func (rt *ResultTree) Len() int {
	return len(rt.children)
}

func (rt *ResultTree) MarshalJSON() ([]byte, error) {
	//自定义编组过程
	dict := make(DataMap)
	dict[`title`] = rt.Title
	dict[`status`] = rt.Status
	dict[`records`] = rt.Record
	dict[`children`] = rt.children
	return json.Marshal(dict)
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

func (rp *Report) SetDetails(details *ResultTree) {
	rp.Details = details
}

func (rp *Report) StatusCount() {
	statusCounter(&rp.Summary, rp.Details, true)
}

func statusCounter(sum *Summary, dt *ResultTree, rootFlag bool) {
	if dt.Record != nil {
		sum.Status[1].Count(dt.Record.Stat) //step
		//如果有任一failed的，那个父节点及其祖先也是failed的
		if dt.Record.Stat == FAILED || dt.Record.Stat == ERROR {
			if dt.Parent().Status != FAILED {
				for parent := dt.Parent(); ; {
					parent.Status = FAILED
					if parent.Parent() == parent {
						break
					}
					parent = parent.Parent()
				}
			}
		}
		dt.Status = dt.Record.Stat
	}
	if dt.Len() > 0 {
		for i := 0; i < dt.Len(); i++ {
			statusCounter(sum, dt.ChildAt(i), false)
			if rootFlag {
				sum.Status[0].Count(dt.ChildAt(i).Status)
			}
		}
	}
}

func (rp *Report) Json() string {
	jsonStr, err := json.Marshal(rp)
	if err != nil {
		log.Info(`Report to json failed:`, err.Error())
		return `{}`
	}
	return string(jsonStr)
}

func NewResultTree() *ResultTree {
	rt := &ResultTree{}
	rt.Detail = NewDetail()
	rt.SetParent(rt)
	return rt
}

// Summary

func NewSummary() *Summary {
	return &Summary{Status: make([]Status, 2)}
}

// Details
func NewDetail() *Detail {
	return &Detail{Status: 0}
}

func (dt *Detail) SetRecord(record *Record) {
	dt.Record = record
}

// func (dt *Detail) AddRecord(record Record) {
// 	dt.Record = append(dt.Record, record)
// }

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

func (rc *Record) MarshalJSON() ([]byte, error) {
	//自定义编组过程
	dict := make(DataMap)
	req := make(DataMap)
	resp := make(DataMap)
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
	bBody, err := ioutil.ReadAll(body)
	if err != nil {
		panic(err)
	}
	req[`body`] = string(bBody)
	dict[`request`] = req

	resp[`url`] = rc.Request.URL.String()
	resp[`statusCode`] = rc.Response.Code
	resp[`cookies`] = rc.Response.Cookies
	resp[`header`] = rc.Response.Header
	resp[`body`] = rc.Response.Content
	dict[`response`] = resp
	dict[`validators`] = rc.Validators
	return json.Marshal(dict)
}
