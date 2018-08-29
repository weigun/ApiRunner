package testcase

import (
	utils "ApiRunner/utils"
	"bytes"
	"fmt"
	_ "fmt"
	_ "log"
	"net/http"
	Url "net/url"
	_ "strings"
)

type Header struct {
	Key, Val string
}

type Variables struct {
	Name string
	Val  interface{}
}

type CasesetConf struct {
	Name       string
	Host       string
	Headers    []Header
	GlobalVars []Variables
}

type Params struct {
	Params map[string]interface{}
}

type Condition struct {
	Op       string
	Source   string
	Verified string
}

type CaseItem struct {
	Name    string
	Api     string
	Method  string
	Headers []Header
	Params
	Validate []Condition
}
type Caseset struct {
	Conf  CasesetConf
	Cases []CaseItem
}

const (
	//比较操作枚举
	EQ   = "eq"
	NE   = "ne"
	GT   = "gt"
	LT   = "lt"
	REGX = "regx" //正则
)

func (this *Params) Encode() string {
	//编码查询参数
	query := Url.Values{}
	for k, v := range this.Params {
		query.Add(k, v.(string))
	}
	return query.Encode()
}

func (this *Params) ToJson() string {
	//转json，用于post方法
	return utils.Map2Json(this.Params)
}

func (this *Params) Conver(method string) string {
	//翻译转换为可请求的字符串格式
	tmpl := utils.GetTemplate(getFuncMap())
	if method == "GET" {
		for k, v := range this.Params {
			this.Params[k] = utils.Translate(tmpl, v.(string))
		}
		return this.Encode()
	} else {
		rawData := this.ToJson()
		return utils.Translate(tmpl, rawData)
	}
}

func NewCaseset() *Caseset {
	return &Caseset{}
}

func (this *CaseItem) HasHeader(key string) int {
	for i, v := range this.Headers {
		if v.Key == key {
			return i
		}
	}
	return -1

}

func (this *CaseItem) AddHeader(h Header) {
	this.Headers = append(this.Headers, h)
}

func (this *CaseItem) AddCondition(c Condition) {
	this.Validate = append(this.Validate, c)
}

func (this *CaseItem) GetConditions() []Condition {
	return this.Validate
}

func (this *CaseItem) cover() {
	//TODO 将整个ci翻译？
}

func (this *CaseItem) BuildRequest() *http.Request {
	//构造请求体
	fmt.Println(this)
	api := this.Api
	method := this.Method
	var data string
	if len(this.Params.Params) == 0 {
		data = ""
	} else {
		data = this.Params.Conver(this.Method)
	}
	fmt.Println(data)
	bodyData := bytes.NewBuffer([]byte(data)) //get方法默认是空字符串
	req, err := http.NewRequest(method, api, bodyData)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(req)
	for _, h := range this.Headers {
		fmt.Println("-------------", h)
		req.Header.Add(h.Key, h.Val)
	}
	return req
}

func (this *Caseset) AddCaseItem(ci CaseItem) {
	this.Cases = append(this.Cases, ci)
}

func (this *Caseset) GetCases() []CaseItem {
	return this.Cases
}

///////////////////////////////
//实现translate接口
func (this *Header) cover() string {
	return ""
}

//TODO 需要实现translate接口
//TODO 可能需要将string转为rune
//TODO 用例分层
