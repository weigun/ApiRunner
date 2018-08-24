package main

import (
	"fmt"
	Url "net/url"
	_ "strings"
)

type header struct {
	key, val string
}

type variables struct {
	name string
	val  interface{}
}

type casesetConf struct {
	name       string
	host       string
	headers    []header
	globalVars []variables
}

type params struct {
	params map[string]interface{}
}

type condition struct {
	operation string
	source    string
	verified  string
}

type caseItem struct {
	name    string
	api     string
	method  string
	headers []header
	params
	validate []condition
}
type caseset struct {
	conf  casesetConf
	cases []caseItem
}

const (
	//比较操作枚举
	eq   = "eq"
	ne   = "ne"
	gt   = "gt"
	lt   = "lt"
	regx = "regx" //正则
)

func (this *params) encode() string {
	//编码查询参数
	query := Url.Values{}
	for k, v := range this.params {
		query.Add(k, v.(string))
	}
	return query.Encode()
}

func (this *params) toJson() string {
	//转json，用于post方法
	return map2Json(this.params)
}

func NewCaseset() *caseset {
	return &caseset{}
}

func (this *caseItem) hasHeader(key string) int {
	for i, v := range this.headers {
		if v.key == key {
			return i
		}
	}
	return -1

}

func (this *caseItem) addHeader(h header) {
	this.headers = append(this.headers, h)
}

func (this *caseItem) addCondition(c condition) {
	this.validate = append(this.validate, c)
}

func (this *caseset) addCaseItem(ci caseItem) {
	this.cases = append(this.cases, ci)
}

func (this *caseset) getCases() []caseItem {
	return this.cases
}

///////////////////////////////
//实现translate接口
func (this *header) cover() string {

}

//TODO 需要实现translate接口
//TODO 可能需要将string转为rune
