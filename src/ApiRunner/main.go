package main

import (
	testcase "ApiRunner/case"
	engine "ApiRunner/engine"
	mgr "ApiRunner/manager"
	utils "ApiRunner/utils"
	_ "fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

func appStart() {
	err := http.ListenAndServe("127.0.0.1:9090", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func RunLocalTasks(tasks []string) {
	//单次执行用例
	eng := engine.NewEngine()
	curFolder := utils.GetCwd()
	var wg sync.WaitGroup
	for i, v := range tasks {
		casePath := filepath.Join(curFolder, "testcase", v+".json")
		caseParser := testcase.NewCaseParser(casePath)
		rn := engine.NewRunner{caseParser.GetCaseset()}
		wg.Add(1)
		go func(rn runner) {
			defer wg.Done()  //TODO 可能需要放在safeRun下一行
			rn.ready <- true //缓冲chan
			eng.safeRun(rn)
		}(rn) //copy rn
	}
	wg.Wait()
	//	TODO 生成报告
	log.Println("test done")
}

func main() {
	cmdArgs := parseCmd()
	if cmdArgs.web {
		mgr.SetupHandlers()
		appStart()
	} else {
		RunLocalTasks(strings.Split(",", cmdArgs.runCase))
	}
}
