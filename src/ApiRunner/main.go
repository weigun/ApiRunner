package main

import (
	tsParser "ApiRunner/case/parser"
	cmd "ApiRunner/cmd"
	engine "ApiRunner/engine"
	mgr "ApiRunner/manager"
	runner "ApiRunner/runner"
	utils "ApiRunner/utils"
	_ "fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
)

func appStart() {
	err := http.ListenAndServe("127.0.0.1:9090", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func runLocalTasks(tasks []string) {
	//单次执行用例
	eng := engine.NewEngine()
	_ = eng
	curFolder := utils.GetCwd()
	var wg sync.WaitGroup
	for _, v := range tasks {
		casePath := filepath.Join(curFolder, "testcase", v+".json")
		caseParser := tsParser.NewCaseParser(casePath)
		rn := runner.NewRunner(caseParser)
		wg.Add(1)
		go func(rn runner.Runner) {
			defer wg.Done()   //TODO 可能需要放在safeRun下一行
			rn.Ready <- false //缓冲chan
			rn.Start()        //TODO 需要转为安全模式
			//			eng.SafeRun(rn)
		}(rn) //copy rn
	}
	wg.Wait()
	//	TODO 生成报告
	log.Println("test done")
}

func main() {
	cmdArgs := cmd.ParseCmd()
	if cmdArgs.Web {
		mgr.SetupHandlers()
		appStart()
	} else {
		runLocalTasks(strings.Split(",", cmdArgs.RunCase))
	}
}
