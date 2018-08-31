package validation

import (
	//	"time"
	//	report "ApiRunner/report"
	testcase "ApiRunner/case"
	utils "ApiRunner/utils"
	"log"
	"regexp"
	"strconv"
	"sync"
)

type PIresponseInterface interface {
	GetCode() int
	GetContent() string
	GetErrMsg() string
	GetElapsed() int64
}

type Response struct {
	// TODO 需要加入更多的字段，用于报告生成
	Code    int
	Content string
	ErrMsg  string
	Elapsed int64
}

func (this Response) GetCode() int {
	return this.Code
}

func (this Response) GetContent() string {
	return this.Content
}

func (this Response) GetErrMsg() string {
	return this.ErrMsg
}

func (this Response) GetElapsed() int64 {
	return this.Elapsed
}

type ResultItem struct {
	Tsp testcase.PIparserInsterface
	Res PIresponseInterface
}

type validation struct {
	//对应response的结构，方便进行引用
	Body map[string]interface{}
}

func (this validation) GetData() map[string]interface{} {
	//实现dataInterface接口
	return map[string]interface{}{"body": this.Body}
}

type PIassertInterface interface {
	//断言接口
	Equal(string) bool
	GreaterThan(string) bool
	LessThan(string) bool
	Regx(string) bool
}
type Assert struct {
	//	rawData     string      //存放待比较的原始数据
	compareData interface{} //存放模板转译后的数据
}

func (this *Assert) Equal(b string) bool {
	//TODO 可能需要处理编码的问题
	return this.compareData.(string) == b
}

func (this *Assert) GreaterThan(b string) bool {
	//TODO 可能需要处理编码的问题
	a := utils.ToNumber(this.compareData)
	switch a.(type) {
	case int, int64, int32:
		b, _ := strconv.ParseInt(b, 10, 64)
		return a.(int64) > b
	case float64, float32:
		b, _ := strconv.ParseFloat(b, 64)
		return a.(float64) > b
	default:
		log.Printf("not support type to compared,type is %T\n", this.compareData)
		return false
	}
}

func (this *Assert) LessThan(b string) bool {
	//TODO 可能需要处理编码的问题
	a := utils.ToNumber(this.compareData)
	switch a.(type) {
	case int, int64, int32:
		b, _ := strconv.ParseInt(b, 10, 64)
		return a.(int64) < b
	case float64, float32:
		b, _ := strconv.ParseFloat(b, 64)
		return a.(float64) < b
	default:
		log.Printf("not support type to compared,type is %T\n", this.compareData)
		return false
	}
}

func (this *Assert) Regx(b string) bool {
	//TODO 可能需要处理编码的问题
	regx := regexp.MustCompile(this.compareData.(string))
	res := regx.FindStringIndex(b)
	return res != nil

}

func NewAssert(rawData string) *Assert {
	return &Assert{compareData: rawData}
}

type ResultPool struct {
	//	testcaseChan chan testcase.PIparserInsterface
	resultChan chan ResultItem
	doneChan   chan bool //just for local mode,not web
}

const (
	//比较操作枚举
	EQ   = "eq"
	NE   = "ne"
	GT   = "gt"
	LT   = "lt"
	REGX = "regx" //正则
)

var once sync.Once
var resPool *ResultPool

func NewResultPool() *ResultPool {
	once.Do(func() {
		log.Println("ResultPool init")
		resPool = &ResultPool{make(chan ResultItem, 100), make(chan bool)}
		//用一个新的go程来专门处理结果
		go func() {
			resPool.handleResult()
		}()
	})
	return resPool
}

func (this *ResultPool) Push(r ResultItem) {
	this.resultChan <- r
}

func (this *ResultPool) Shift() ResultItem {
	return <-this.resultChan
}

func (this *ResultPool) handleResult() {
	//TODO 需要将处理结果转交给report模块
	for resItem := range this.resultChan {
		handle(resItem)
	}
	log.Println("all results came out,ready to done")
	//	time.Sleep(time.Duration(10) * time.Second)
	this.doneChan <- true
}

func (this *ResultPool) WaitForDone() {
	close(this.resultChan)
	log.Println("resultChan closed,WaitForDone")
	<-this.doneChan
}

/*
私有函数
主要是作用是对模板进行转译
*/
func handle(resItem ResultItem) {
	log.Println("resItem:", resItem)
	res := resItem.Res //response obj
	tsp := resItem.Tsp //testCaseParser obj
	contentMap := utils.Json2Map([]byte(res.GetContent()))
	log.Printf("contentMap is %T\n", contentMap)
	vali := validation{contentMap}
	tmpl := utils.GetTemplate(nil)
	for _, caseItem := range tsp.GetCaseset().GetCases() {
		for _, cond := range caseItem.Validate {
			log.Println("----------", cond)
			compareData := utils.Translate(tmpl, cond.Source, vali)
			log.Println(compareData)
			assert := NewAssert(compareData)
			switch cond.Op {
			case EQ:
				ret := assert.Equal(cond.Verified)
				log.Println(cond, ret, compareData, cond.Verified)
			case GT:
				ret := assert.GreaterThan(cond.Verified)
				log.Println(cond, ret, compareData, cond.Verified)
			case LT:
				ret := assert.LessThan(cond.Verified)
				log.Println(cond, ret, compareData, cond.Verified)
			case NE:
				ret := !assert.Equal(cond.Verified)
				log.Println(cond, ret, compareData, cond.Verified)
			case REGX:
				ret := !assert.Regx(cond.Verified)
				log.Println(cond, ret, compareData, cond.Verified)
			}
		}
	}
}
