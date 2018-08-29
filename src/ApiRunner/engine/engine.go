package engine

import (
	testcase "ApiRunner/case"
	runner "ApiRunner/runner"
	_ "bytes"
	_ "fmt"
	_ "io"
	_ "io/ioutil"
	_ "log"
	_ "regexp"
	"runtime/debug"
	"sync"
	_ "time"
)

type engine struct {
	testcaseChan chan testcase.PIparserInsterface
}

var once sync.Once
var eng *engine

func NewEngine() *engine {
	once.Do(func() {
		eng = &engine{testcaseChan: make(chan testcase.PIparserInsterface, 50)}
	})
	return eng
}

func (this *engine) SafeRun(r runner.Runner) {
	defer func() {
		// don't panic
		err := recover()
		if err != nil {
			debug.PrintStack()
		}
	}()
	r.Start()
}

func (this *engine) SpawnRunner() {
	for tsp := range this.testcaseChan {
		go func(tsp testcase.PIparserInsterface) {
			r := runner.NewRunner(tsp)
			this.SafeRun(r)
		}(tsp)
	}
}
