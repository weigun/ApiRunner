// test_runner.go
package runner

import (
	"context"
	_ "fmt"
	"log"

	"ApiRunner/models"
)

const (
	Queuing = iota
	Running
	Passed
	Failed
	Cancel
)

type TestRunner struct {
	ID       string
	CaseObj  models.ICaseObj
	canceler context.CancelFunc
	Status   int
}

func (r *TestRunner) Start() {
	log.Println("testrunner started")
	ctx, cancel := context.WithCancel(context.Background())
	valueCtx := context.WithValue(ctx, `status`, r.Status)
	r.canceler = cancel
	go func(ctx context.Context) {
		//TODO 需要保存堆栈
		go execute(r) //用例执行
		select {
		case <-ctx.Done():
			log.Println("runner done")
			r.Status = ctx.Value(`status`).(int)
			return
		}
	}(valueCtx)

}

func (r *TestRunner) Stop() {
	log.Println("testrunner stopping")
	r.Status = Cancel
	r.canceler()
}

func execute(r *TestRunner) {
	//具体执行用例的实体函数
	caseObj := r.CaseObj
	switch r.CaseObj.(type) {
	case *models.TestCase:
		caseObj = r.CaseObj.(*models.TestCase)
	case *models.TestSuites:
		caseObj = r.CaseObj.(*models.TestSuites)
	default:
		log.Printf(`unknow caseobj type:%T,stop runner`, r.CaseObj)
		r.canceler()
		return
	}
	_ = caseObj
	//顺序执行用例
}
