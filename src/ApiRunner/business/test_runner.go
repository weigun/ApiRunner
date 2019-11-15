// test_runner.go
package business

import (
	"context"
	"fmt"
	"log"
	_ "net/url"

	"regexp"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"

	refNode "ApiRunner/business/refs_tree"
	// "ApiRunner/business/template"
	"ApiRunner/models"
	"ApiRunner/models/young"

	// "ApiRunner/services"
	"ApiRunner/utils"
)

const (
	Queuing = iota
	Running
	Passed
	Failed
	Canceled
)

type TestRunner struct {
	ID       string
	PipeObj  models.Executable
	Reporter *models.Report
	canceler context.CancelFunc
	render   *renderer
	refs     refNode.Node //当前的引用
	rootRefs refNode.Node //根引用
	Status   int
}

func NewTestRunner(id string, pipeObj models.Executable) *TestRunner {
	t := &TestRunner{
		ID:      id,
		PipeObj: pipeObj,
		refs:    refNode.New(`root`),
		render:  newRenderer(),
	}
	t.refs.SetParent(t.refs)
	t.rootRefs = t.refs
	t.Reporter = models.NewReport()
	return t
}

func (r *TestRunner) Start() {
	log.Println("testrunner started")
	ctx, cancel := context.WithCancel(context.Background())
	valueCtx := context.WithValue(ctx, `status`, r.Status)
	r.canceler = cancel
	go func(ctx context.Context) {
		//TODO 需要保存堆栈
		go execute(r) //用例执行
		<-ctx.Done()
		log.Println("runner done")
		r.Status = ctx.Value(`status`).(int) //TODO 大丈夫？

	}(valueCtx)

}

func (r *TestRunner) Stop() {
	log.Println("testrunner stopping")
	r.Status = Canceled
	r.canceler()
	// TODO 干掉eventbus的channel
}

func execute(r *TestRunner) {
	//具体执行用例的实体函数
	// pipeObj := r.PipeObj.(*models.Pipeline)
	report := r.Reporter

	//顺序执行用例
	//开始计时
	sum := models.NewSummary()
	sum.StartAt = time.Now()
	report.SetSummary(*sum)
	// render := newRenderer(r.ID)
	executePipeline(r)
	sum.Duration = time.Now().Sub(sum.StartAt).Nanoseconds() / 1e6
	//统计status
	for _, dt := range report.Details {
		for _, record := range dt.Records {
			// dt.Status.Count(record.Stat)
			sum.Status[1].Count(record.Stat)
		}
		// report.Details[index] = dt
		if dt.Status.Error > 0 || dt.Status.Failed > 0 {
			sum.Status[0].Count(models.FAILED)
		} else {
			sum.Status[0].Count(models.SUCCESS)
		}

	}

	report.SetSummary(*sum)
	log.Println(report.Json())
	spew.Dump(r.refs)
}

func executePipeline(r *TestRunner) {
	log.Println(`dump def:`, spew.Sdump(r.PipeObj.(*models.Pipeline).Def))
	var newDef models.Variables
	bDef := r.render.fillData(toJson(r.PipeObj.(*models.Pipeline).Def), nil)
	json.Unmarshal(bDef, &newDef)
	r.PipeObj.(*models.Pipeline).Def = newDef

	// spew.Dump(r.PipeObj.(*models.Pipeline).Def)
	log.Println(`dump newDef:`, spew.Sdump(r.PipeObj.(*models.Pipeline).Def))

	//backup self so can revert after recursive
	parentPipe := r.PipeObj.(*models.Pipeline)
	parentRef := r.refs.Parent()
	log.Println(`cur refs:`, r.refs)
	log.Println(`parentRef:`, parentRef)

	for index, execNode := range r.PipeObj.(*models.Pipeline).Steps {
		if r.Status == Canceled {
			// 如果runner的已经取消了，就没必要再去执行下一个用例了
			log.Println(`executor stopping,because runner is canceled `)
			return
		}
		node := refNode.New(execNode.RefTag())
		r.refs.AddChild(node)
		//开始执行step
		log.Println(`--------------step begin--------------------`)
		executeStep(&execNode, r, index)
		log.Println(`--------------step end----------------------`)
	}
	r.refs = parentRef
	r.PipeObj = parentPipe //revert self
	log.Println(`after revert,cur refs:`, r.refs)
}

func executeStep(execNode *models.ExecNode, r *TestRunner, rindex int) {
	log.Println(`exec step:`, execNode.Desc)
	pipeObj := r.PipeObj.(*models.Pipeline)
	report := r.Reporter
	requestor := NewRequestor()
	ref := r.refs.ChildAt(rindex) //user1 signup login
	detail := models.NewDetail()
	detail.Title = pipeObj.Name
	if execNode.Host == `` {
		execNode.Host = pipeObj.Host
	}
	// TODO:
	// 模板翻译-done
	// 拦截器-done
	// MultipartFile-done

	//MultipartFile比普通post请求优先级要高
	var header models.Header
	var resp *young.Response

	record := models.NewRecord()
	if execNode.Exec.GetType() == models.TYPE_API {
		//只是执行接口
		apiObj := execNode.Exec.(*models.API)
		url := fmt.Sprintf(`%s/%s`, execNode.Host, apiObj.Path)
		//先合并参数
		for k, v := range execNode.Args {
			apiObj.Params[k] = v
		}
		bheader := r.render.fillData(toJson(apiObj.Headers), pipeObj.Def)
		log.Println(`pipeObj.Def:`, pipeObj.Def)
		json.Unmarshal(bheader, &header)
		// render.renderObj(toJson(apiObj.Headers), true, &header)
		var startTime time.Time
		var fnBeforeRequest, fnAfterResponse string
		if execNode.Hooks[`BeforeRequest`] != nil {
			fnBeforeRequest = execNode.Hooks[`BeforeRequest`].(string)
		}
		if execNode.Hooks[`AfterResponse`] != nil {
			fnAfterResponse = execNode.Hooks[`AfterResponse`].(string)
		}

		if apiObj.MultipartFile.IsEnabled() {
			var mpf models.MultipartFile
			log.Println(`parse MultipartFile:`, apiObj.MultipartFile.Json())
			bMpf := r.render.fillData(apiObj.MultipartFile.Json(), pipeObj.Def)
			json.Unmarshal(bMpf, &mpf)
			req := requestor.BuildPostFiles(url, mpf, header)
			record.Request = req
			startTime = time.Now()
			resp = requestor.doRequest(req, fnBeforeRequest, fnAfterResponse)
		} else {
			var params models.Params
			log.Println(`parse Params:`, toJson(apiObj.Params))
			bParams := r.render.fillData(toJson(apiObj.Params), pipeObj.Def)
			json.Unmarshal(bParams, &params)

			bMethod := r.render.fillData(apiObj.Method, pipeObj.Def)
			req := requestor.BuildRequest(url, string(bMethod), params, header)
			// req := requestor.BuildRequest(url, render.renderValue(apiObj.Method, true), params, header)
			record.Request = req
			startTime = time.Now()
			resp = requestor.doRequest(req, fnBeforeRequest, fnAfterResponse)
		}
		elapsed := time.Now().Sub(startTime).Nanoseconds()
		log.Println(`got response:`, resp)
		record.Desc = execNode.Desc
		record.Elapsed = elapsed / 1e6
		record.Response = resp
		data := make(map[string]interface{})
		data[`StatusCode`] = resp.Code
		//导出变量，如token等
		if resp.ErrMsg == `` {
			//没有错误的时候才能导出变量
			//TODO assert code??
			contentMap := utils.Json2Map([]byte(resp.Content))
			if len(contentMap) == 0 {
				//如果resp.Content返回是{}或者非json
				//则直接将正文给body
				data[`body`] = resp.Content
			} else {
				data[`body`] = contentMap
			}
			for ek, ev := range execNode.Export {
				v := ev.(string)
				if strings.Index(v, `${`) != -1 && strings.Index(v, `}`) != -1 {
					//返回json则需要提取变量。这里非常简单的判断
					bindVal := string(r.render.fillData(v, data))
					ref.AddPairs(ek, bindVal)
					log.Println("add ExportVars:", ek, bindVal)
				} else {
					//plain text
					regx := regexp.MustCompile(v)
					match := regx.FindStringSubmatch(resp.Content)
					if match != nil {
						//目前暂不支持切片，如果是匹配多个值，只能是先合拼，到需要用的时候，自己再转换成字符串切片
						// services.VarsMgr.Add(fmt.Sprintf(`%s:%s`, render.tag, ek), strings.Join(match, `,`))
						ref.AddPairs(ek, strings.Join(match, `,`))
						// varsMgr.SetVar(this.Testcase.GetUid(), v.Name, strings.Join(match, `,`))
						log.Println("add ExportVars:", ek, strings.Join(match, `,`))
					}
				}
			}
		}

		//比较结果
		allPassed := true
		for _, validator := range execNode.Validate {
			// TODO 渲染变量时，适配各种数据类型
			validator.Check = validator.Actual.(string)
			compare := getAssertByOp(validator.Op)
			log.Println(`parse validator.Actual:`, validator.Actual.(string))
			actual := string(r.render.fillData(validator.Actual.(string), data))

			log.Println(`parse validator.Expected:`, validator.Expected.(string))
			expected := string(r.render.fillData(validator.Expected.(string), pipeObj.Def))
			isPassed := So(actual, compare, expected)
			if !isPassed {
				allPassed = false
			}
			log.Printf(`Actual:%v,Expected:%v,So %v`, actual, expected, isPassed)
			validator.Actual = actual
			validator.Expected = expected
			record.AddValidator(validator)
		}
		// TODO error and skip
		if allPassed {
			record.Stat = models.SUCCESS
		} else {
			record.Stat = models.FAILED
		}

		detail.AddRecord(*record)
		detail.Status.Count(record.Stat)
	} else {
		subPipeObj := execNode.Exec.(*models.Pipeline)
		//先合并参数
		log.Println(`execNode.Args:`, spew.Sdump(execNode.Args))
		for k, v := range execNode.Args {
			subPipeObj.Def[k] = v
		}
		var newDef models.Variables

		//需要合并def和refs，来作为数据源
		dataSource := make(map[string]interface{})
		for k, v := range pipeObj.Def {
			dataSource[k] = v
		}
		dataSource[`refs`] = r.rootRefs
		json.Unmarshal(r.render.fillData(toJson(subPipeObj.Def), dataSource), &newDef)
		subPipeObj.Def = newDef
		r.PipeObj = subPipeObj
		r.refs = ref
		log.Println(`replace def:`, subPipeObj.Def)
		executePipeline(r)
	}

	report.AddDetail(*detail)
}
