package main

import (
	"fmt"
)

type respBody struct {
	Code int
	Msg  string
	Data interface{}
}

type validation struct {
	//验证体
	//定义验证数据的结构，可通过模板翻译，用来做结果对比
	Code int
	Body respBody
}

func validate(resp Response, conds []condition) {
	//TODO 结果需要给报告
	_ := resp
	for i, cond := range conds {
		fmt.Println(cond)
	}

}
