package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"sync"
)

func map2Json(m map[string]interface{}) string {
	jsonStr, err := json.Marshal(m)
	if err != nil {
		return "{}"
	}
	return string(jsonStr)
}

func json2Map(js []byte) map[string]interface{} {
	var mapResult = make(map[string]interface{})
	//使用 json.Unmarshal(data []byte, v interface{})进行转换,返回 error 信息
	if err := json.Unmarshal(js, &mapResult); err != nil {
		return mapResult
	}
	return mapResult
}

func getCwd() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0])) //返回绝对路径  filepath.Dir(os.Args[0])去除最后一个元素的路径
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

func runLocalTasks(tasks []string) {
	//单次执行用例
	eng := NewEngine()
	curFolder := getCwd()
	var wg sync.WaitGroup
	for i, v := range tasks {
		casePath := filepath.Join(curFolder, "testcase", v+".json")
		caseParser := NewCaseParser(casePath)
		rn := NewRunner{caseParser.getCaseset()}
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
