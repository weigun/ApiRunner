package main

import (
	_ "fmt"
	"log"
	_ "strings"
)

type _IStata interface {
	start()
	stop()
}

type runner struct {
	caseset *caseset
	ready   chan bool //runner状态
}

func NewRunner(cs *caseset) *runner {
	return &runner{caseset, make(chan bool, 1)}
}

func (this *runner) start() {
	//start test the testcase set
	<-this.ready
	log.Println("start test the testcase set")
	log.Println(this.cacaseset.name)
	for i, ci := range this.caseset.getCases() {
		//顺序执行用例

	}
}

func prepare() {

}
