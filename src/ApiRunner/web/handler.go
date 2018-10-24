package web

import (
	tsParser "ApiRunner/case/parser"
	engine "ApiRunner/engine"
	runner "ApiRunner/runner"
	utils "ApiRunner/utils"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	//	"strings"
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello Wrold!") //这个写入到w的是输出到客户端的
}

func runCases(w http.ResponseWriter, r *http.Request) {
	//需要加一个缓冲队列，避免多个同时执行
	log.Println("method:", r.Method)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(`get request data err `, err.Error())
		return
	}
	// should be a json,like {cases:[demo,cnfig,login]}
	eng := engine.NewEngine()
	m := utils.Json2Map(body)
	curFolder := utils.GetCwd()
	for _, v := range m[`cases`].([]interface{}) {
		//每个case一个go程跑runner
		casePath := filepath.Join(curFolder, "testcase", "conf", v.(string)+".conf")
		caseParser := tsParser.NewCaseParser(casePath)
		rn := runner.NewRunner(caseParser)
		go func(rn runner.Runner) {
			rn.Ready <- false //缓冲chan
			//			rn.Start()        //TODO 需要转为安全模式
			eng.SafeRun(rn)
		}(rn) //copy r
	}
	log.Println(`one request handle end`)
}

func SetupHandlers() {
	log.Println(`SetupHandlers`)
	http.HandleFunc(`/`, index)            // mainPage
	http.HandleFunc(`/run/case`, runCases) // runCases
}
