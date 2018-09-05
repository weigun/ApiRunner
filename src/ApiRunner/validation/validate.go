package validation

import (
	"io/ioutil"
	"time"
	_ "time"
	//	"time"
	//	report "ApiRunner/report"
	testcase "ApiRunner/case"
	report "ApiRunner/report"
	utils "ApiRunner/utils"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"sync"
)

type PIresponseInterface interface {
	GetCode() int
	GetContent() string
	GetErrMsg() string
	GetElapsed() int64
	GetHeader() http.Header
}

type Response struct {
	// TODO 需要加入更多的字段，用于报告生成
	Code    int
	Content string
	ErrMsg  string
	Elapsed int64
	Header  http.Header
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

func (this Response) GetHeader() http.Header {
	return this.Header
}

type ResultItem struct {
	Tsp   testcase.PIparserInsterface
	Index int64
	Res   PIresponseInterface
}

type validator struct {
	//对应response的结构，方便进行引用
	Body map[string]interface{}
}

func (this validator) GetData() map[string]interface{} {
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
	doneChan   chan uint32 //just for local mode,not web
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
		resPool = &ResultPool{make(chan ResultItem, 100), make(chan uint32, 64)}
		//用一个新的go程来专门处理结果
		go func() {
			for {
				select {
				case resItem := <-resPool.resultChan:
					log.Println("handle item---------------")
					handle(resItem)
				case uid := <-resPool.doneChan:
					log.Println("all case done")
					//收到用例跑完消息，2s后再告诉报表开始生成，暂时避免异步的问题
					time.AfterFunc(time.Duration(2)*time.Second, func() {
						log.Println("generateing report")
						var footer report.Footer
						footer.Add2Cache(uid)
					})
				}
			}
			//			resPool.handleResult()
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
	for resItem := range this.resultChan {
		handle(resItem)
	}
	log.Println("all results came out,ready to done")
	log.Println("generate report")

	//	this.doneChan <- true
}

func (this *ResultPool) Done(uid uint32) {
	this.doneChan <- uid
	//	close(this.resultChan)
	//	log.Println("resultChan closed,WaitForDone")
	//	<-this.doneChan
	//	log.Println("add footer:", uid)
	//	var footer report.Footer
	//	_ = time.Second
	//	//	time.Sleep(time.Duration(3) * time.Second)
	//	footer.Add2Cache(uid)
	//	// TODO 因为是异步，所以有概率会出现，先收到footer再收到record的情况，导致报表内容为空的问题，需要修改
}

/*
私有函数
*/
func handle(resItem ResultItem) {
	//结果处理函数
	//TODO 各种log需要集中到log中心，因为在报表性需要查看log信息
	//	log.Println("resItem:", resItem)
	res := resItem.Res     //response obj
	tsp := resItem.Tsp     //testCaseParser obj
	index := resItem.Index //用例的索引
	contentMap := utils.Json2Map([]byte(res.GetContent()))
	log.Printf("contentMap is %T\n", contentMap)
	vali := validator{contentMap}
	tmpl := utils.GetTemplate(nil)
	cn := report.CaseNum{}
	caseItem := tsp.GetCaseset().GetCases()[index]
	//TODO Status不应该这样判断
	record := report.Record{Status: res.GetCode() == 200, Api: caseItem.Api, Elapsed: res.GetElapsed(), Traceback: "coming soon"}
	makeDetail(&record, &caseItem, res, tsp.GetUid())
	//	detail.Add2Cache(tsp.GetUid())
	allPassed := true
	for _, cond := range caseItem.GetConditions() {
		//handle Validators
		//	type Validator struct {
		//	Check      bool
		//	Comparator string
		//	Expect     string
		//	Actual     string
		//		log.Println("----------", cond)
		compareData := utils.Translate(tmpl, cond.Source, vali)
		//		log.Println(compareData)
		assert := NewAssert(compareData)
		var ret bool
		switch cond.Op {
		case EQ:
			ret = assert.Equal(cond.Verified)
			//			log.Println(cond, ret, compareData, cond.Verified)
		case GT:
			ret = assert.GreaterThan(cond.Verified)
			//			log.Println(cond, ret, compareData, cond.Verified)
		case LT:
			ret = assert.LessThan(cond.Verified)
			//			log.Println(cond, ret, compareData, cond.Verified)
		case NE:
			ret = !assert.Equal(cond.Verified)
			//			log.Println(cond, ret, compareData, cond.Verified)
		case REGX:
			ret = !assert.Regx(cond.Verified)
			//			log.Println(cond, ret, compareData, cond.Verified)
		}
		if !ret {
			allPassed = false
		}
		valtor := report.Validator{ret, cond.Op, cond.Verified, compareData}
		//		valtor.Add2Cache(tsp.GetUid())
		record.Detail.Validators = append(record.Detail.Validators, valtor)
	}
	if allPassed {
		cn.Successes++
	} else {
		cn.Failures++
	}
	cn.TotalCases = cn.Successes + cn.Failures
	//	log.Println("cn:", cn)
	cn.Add2Cache(tsp.GetUid())
	record.Add2Cache(tsp.GetUid())
}

func makeDetail(record *report.Record, ci testcase.PIrequest, res PIresponseInterface, uid uint32) {
	detail := report.ExecuteDetail{}
	req := ci.BuildRequest()
	bs, _ := ioutil.ReadAll(req.Body)
	s := string(bs)
	//	log.Println("requestBody:", s)
	//s可以是json，也可以是&参数的方式，最后一种是纯字符串
	//都需要转为map结构
	detail.RequestData = make(map[string]string)
	detail.ResponseData = make(map[string]string)
	if req.Method == "GET" {
		m := utils.Params2Map(s)
		if len(m) == 0 {
			//纯字符串
			detail.RequestData["body"] = s
		} else {
			detail.RequestData["body"] = utils.Map2Json(m)
		}
	} else {
		//post方法
		detail.RequestData["body"] = s
	}
	detail.RequestData["header"] = utils.Header2Json(res.GetHeader())
	//request data finished
	//handle response
	detail.ResponseData["header"] = utils.Header2Json(res.GetHeader())
	detail.ResponseData["body"] = res.GetContent()
	log.Println("----------detail:", detail)
	record.Detail = detail

}
