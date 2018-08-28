package testcase

import (
	utils "ApiRunner/utils"
	_ "fmt"
	_ "log"
	Url "net/url"
	_ "strings"
)

type Header struct {
	key, val string
}

type Variables struct {
	name string
	val  interface{}
}

type CasesetConf struct {
	name       string
	host       string
	headers    []Header
	globalVars []Variables
}

type Params struct {
	params map[string]interface{}
}

type Condition struct {
	operation string
	source    string
	verified  string
}

type CaseItem struct {
	name    string
	api     string
	method  string
	headers []Header
	params
	validate []Condition
}
type Caseset struct {
	conf  casesetConf
	cases []caseItem
}

const (
	//比较操作枚举
	EQ   = "eq"
	NE   = "ne"
	GT   = "gt"
	LT   = "lt"
	REGX = "regx" //正则
)

func (this *params) Encode() string {
	//编码查询参数
	query := Url.Values{}
	for k, v := range this.params {
		query.Add(k, v.(string))
	}
	return query.Encode()
}

func (this *params) ToJson() string {
	//转json，用于post方法
	return utils.Map2Json(this.params)
}

func (this *params) Conver(method string) string {
	//翻译转换为可请求的字符串格式
	if method == "GET" {
		for k, v := range this.params {
			this.params[k] = utils.Translate(v.(string))
		}
		return this.Encode()
	} else {
		rawData := this.ToJson()
		return utils.Translate(rawData)
	}
}

func NewCaseset() *Caseset {
	return &Caseset{}
}

func (this *CaseItem) HasHeader(key string) int {
	for i, v := range this.headers {
		if v.key == key {
			return i
		}
	}
	return -1

}

func (this *CaseItem) Addheader(h Header) {
	this.headers = append(this.headers, h)
}

func (this *CaseItem) Addcondition(c Condition) {
	this.validate = append(this.validate, c)
}

func (this *CaseItem) Getconditions() []Condition {
	return this.validate
}

func (this *CaseItem) Cover() {
	//TODO 将整个ci翻译？
}

func (this *CaseItem) Buildrequest() *http.Request {
	//构造请求体
	api := this.api
	method := this.method
	var data string
	if this.params == nil {
		data = ""
	} else {
		data = this.params.Conver(this.method)
	}
	return NewRequest(api, method, data)
}

func (this *Caseset) Addcaseitem(ci Caseitem) {
	this.cases = append(this.cases, ci)
}

func (this *Caseset) Getcases() []Caseitem {
	return this.cases
}

///////////////////////////////
//实现translate接口
func (this *Header) Cover() string {

}

//TODO 需要实现translate接口
//TODO 可能需要将string转为rune
//TODO 用例分层
