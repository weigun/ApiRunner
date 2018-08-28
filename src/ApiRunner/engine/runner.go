package engine

import (
	testcase "ApiRunner/case"
	_ "fmt"
	"log"
	_ "strings"
)

type _IStata interface {
	start(*engine.Engine)
	stop()
}

type Runner struct {
	caseset *caseset
	ready   chan bool //runner状态
}

func NewRunner(cs *caseset) *Runner {
	return &Runner{caseset, make(chan bool, 1)}
}

func (this *Runner) start(eng *Engine) {
	//start test the testcase set
	<-this.ready
	log.Println("start test the testcase set")
	log.Println(this.cacaseset.name)
	for i, ci := range this.caseset.getCases() {
		//顺序执行用例
		req := ci.buildRequest() //构造请求体
		resp := eng.do(req)
		log.Println(resp)
		//验证结果
		validate(resp, ci.getConditions())

	}
}

func (this *Runner) stop() {

}
