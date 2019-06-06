// runner_group.go
package business

import (
	_ "context"
	_ "fmt"
	_ "log"
	"sync"
)

/*
管理全局的testrunner，方便终止执行
*/

type RunnerGroup struct {
	sync.RWMutex
	group map[string]*TestRunner
}

func (rg *RunnerGroup) Add(rid string, rPtr *TestRunner) {
	rg.Lock()
	defer rg.Unlock()
	rg.group[rid] = rPtr
}

func (rg *RunnerGroup) Remove(rid string) {
	rg.Lock()
	defer rg.Unlock()
	delete(rg.group, rid)
}

func (rg *RunnerGroup) Get(rid string) *TestRunner {
	rg.RLock()
	defer rg.RUnlock()
	return rg.group[rid]
}

func (rg *RunnerGroup) StopRunner(rid string) {
	runner := rg.Get(rid)
	runner.Stop()
}
