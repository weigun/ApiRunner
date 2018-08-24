package main

import (
	"flag"
	"fmt"
	_ "strings"
)

type cmdArgs struct {
	runCase string
	web     bool
}

func parseCmd() {
	args := cmdArgs{}
	flag.Var(&args.runCase, "run-case", "run cases,use ,to split,eg:logig,userinfo")
	flag.Var(&args.web, "web", "web mode")
	flag.Parse()
	fmt.Println("get Args:")
	for _, v := range flag.Args() {
		fmt.Println(v)
	}
}
