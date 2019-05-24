// test_report.go
package models

type CaseNum struct {
	TotalCases uint64
	Successes  uint64
	Failures   uint64
	Errors     uint64
}

type Summary struct {
	Title     string
	StartTime uint64
	Duration  uint64
	CaseNum
}

// type Validator struct {
// 	Check      bool
// 	Comparator string
// 	Expect     string
// 	Actual     string
// }

type ExecuteDetail struct {
	RequestData  map[string]string
	ResponseData map[string]string
	// Validators   []Validator
}

type Record struct {
	//一个执行结果
	//	CaseSetName string
	Status    bool
	Api       string
	Elapsed   uint64 //ms
	Detail    ExecuteDetail
	Traceback string // 捕获到的异常
}

type Info struct {
	TaskName string `json:"taskName"`
	Host     string `json:"host"`
}

type RecordSet struct {
	//执行结果的集合，即一个用例集，都是以用例集为单位的
	Info
	List []Record
	CaseNum
}

type TestReport struct {
	TaskName string `json:"taskName"`
	Host     string `json:"host"`
	// Results  []Result `json:"results"`
}
