package main

import (
	_ "fmt"
	"log"
	_ "strings"
)

type _IStata interface {
	start(*engine)
	stop()
}

type runner struct {
	caseset *caseset
	ready   chan bool //runner状态
}

func NewRunner(cs *caseset) *runner {
	return &runner{caseset, make(chan bool, 1)}
}

func (this *runner) start(eng *engine) {
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

func (this *runner) stop() {

}
