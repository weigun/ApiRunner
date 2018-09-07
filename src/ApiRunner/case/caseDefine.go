package testcase

import (
	varsMgr "ApiRunner/manager/bucket/caseVariables"
	utils "ApiRunner/utils"
	"bytes"
	"log"
	"net/http"
	Url "net/url"
	"strings"
)

type Header struct {
	Key, Val string
}

type Variables struct {
	Name string
	Val  interface{}
}

func (this *Variables) Conver() string {
	tmpl := utils.GetTemplate(getFuncMap())
	return utils.Translate(tmpl, this.Val.(string), nil)
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
	Export   []Variables
	Validate []Condition
}
type Caseset struct {
	Conf  CasesetConf
	Cases []CaseItem
}

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

func (this *Params) Conver(uid uint32, method string) string {
	//翻译转换为可请求的字符串格式
	tmpl := utils.GetTemplate(getFuncMap())
	if method == "GET" {
		for k, v := range this.Params {
			if strings.Index(v.(string), `$`) != -1 {
				// 存在变量引用
				val := strings.Replace(v.(string), `$`, `.`, -1)
				this.Params[k] = utils.Translate(tmpl, val, varsMgr.VarMap{uid})
			} else {
				// 没有变量引用，直接翻译即可
				this.Params[k] = utils.Translate(tmpl, v.(string), nil)
			}
		}
		return this.Encode()
	} else {
		rawData := this.ToJson()
		if strings.Index(rawData, `$`) != -1 {
			// 存在变量引用
			return utils.Translate(tmpl, rawData, varsMgr.VarMap{uid})
		} else {
			return utils.Translate(tmpl, rawData, nil)
		}
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

func (this *CaseItem) AddExportVar(v Variables) {
	this.Export = append(this.Export, v)
}

func (this *CaseItem) GetExportVars() []Variables {
	return this.Export
}

func (this *CaseItem) BuildRequest(uid uint32) *http.Request {
	//构造请求体
	//	fmt.Println("CaseItem:", this)
	api := this.Api
	method := this.Method
	var data string
	if len(this.Params.Params) == 0 {
		data = ""
	} else {
		data = this.Params.Conver(uid, this.Method)
	}
	log.Println("BuildRequest:", data)
	bodyData := bytes.NewBuffer([]byte(data)) //get方法默认是空字符串
	req, err := http.NewRequest(method, api, bodyData)
	if err != nil {
		panic(err.Error())
	}
	//	fmt.Println(req)
	for _, h := range this.Headers {
		//		fmt.Println("-------------", h)
		req.Header.Add(h.Key, h.Conver(uid))
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
func (this *Header) Conver(uid uint32) string {
	tmpl := utils.GetTemplate(getFuncMap())
	if strings.Index(this.Val, `$`) != -1 {
		// 存在变量引用
		val := strings.Replace(this.Val, `$`, `.`, -1)
		log.Println("#################:", val, varsMgr.VarMap{uid}.GetData())
		return utils.Translate(tmpl, val, varsMgr.VarMap{uid})
	} else {
		// 没有变量引用，直接翻译即可
		return utils.Translate(tmpl, this.Val, nil)
	}
}

//TODO 需要实现translate接口
//TODO 可能需要将string转为rune
//TODO 用例分层
