package cmd

import (
	"flag"
	"fmt"
	_ "strings"
)

type CmdArgs struct {
	RunCase string
	Web     bool
}

func ParseCmd() *CmdArgs {
	args := CmdArgs{}
	flag.StringVar(&args.RunCase, "run-case", "", "run cases,use ,to split,eg:logig,userinfo")
	flag.BoolVar(&args.Web, "web", false, "web mode")
	flag.Parse()
	fmt.Println("get Args:", args)
	for _, v := range flag.Args() {
		fmt.Println(v)
	}
	return &args
}
