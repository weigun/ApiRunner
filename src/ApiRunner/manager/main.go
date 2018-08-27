package main

import (
	_ "fmt"
	"log"
	"net/http"
	"strings"
)

func setupHandlers() {
	for path, handler := range handlerMap {
		http.HandleFunc(path, handler) //设置访问的路由
	}
}

func appStart() {
	err := http.ListenAndServe("127.0.0.1:9090", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func main() {
	cmdArgs := parseCmd()
	if cmdArgs.web {
		setupHandlers()
		appStart()
	} else {
		runLocalTasks(strings.Split(",", cmdArgs.runCase))
	}
}
