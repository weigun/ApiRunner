package manager

import (
	"flag"
	"fmt"
	_ "strings"
)

type cmdArgs struct {
	runCase string
	web     bool
}

func ParseCmd() *cmdArgs {
	args := cmdArgs{}
	flag.StringVar(&args.runCase, "run-case", "", "run cases,use ,to split,eg:logig,userinfo")
	flag.BoolVar(&args.web, "web", false, "web mode")
	flag.Parse()
	fmt.Println("get Args:")
	for _, v := range flag.Args() {
		fmt.Println(v)
	}
	return &args
}
