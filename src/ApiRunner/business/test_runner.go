// test_runner.go
package business

import (
	"context"
	"fmt"
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
	"ApiRunner/utils/helper"
)

const (
	Queuing = iota
	Running
	Passed
	Failed
	Canceled
)

type TestRunner struct {
	ID        string
	PipeObj   models.Executable
	Reporter  *models.Report
	canceler  context.CancelFunc
	render    *renderer
	refs      refNode.Node //当前的引用
	rootRefs  refNode.Node //根引用
	mementoes *mementoMgr
	Status    int
}

func NewTestRunner(id string, pipeObj models.Executable) *TestRunner {
	t := &TestRunner{
		ID:        id,
		PipeObj:   pipeObj,
		refs:      refNode.New(`root`),
		render:    newRenderer(),
		mementoes: NewMementoMgr(),
	}
	t.refs.SetParent(t.refs)
	t.rootRefs = t.refs
	t.Reporter = models.NewReport()
	return t
}

func (r *TestRunner) Start() {
	log.Info("testrunner started")
	ctx, cancel := context.WithCancel(context.Background())
	valueCtx := context.WithValue(ctx, `status`, r.Status)
	r.canceler = cancel
	go func(ctx context.Context) {
		//TODO 需要保存堆栈
		go execute(r) //用例执行
		<-ctx.Done()
		log.Info("runner done")
		// TODO set status

	}(valueCtx)

}

func (r *TestRunner) Stop() {
	log.Info("testrunner stopping")
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
	detail := models.NewResultTree()
	detail.SetParent(detail)
	report.SetDetails(detail)
	detail.Title = r.PipeObj.GetName()
	executePipeline(r)
	sum.Duration = time.Now().Sub(sum.StartAt).Nanoseconds() / 1e6
	//统计status
	report.StatusCount()
	report.SetSummary(*sum)
	// spew.Dump(report)
	log.Info(report.Json())

}

func executePipeline(r *TestRunner) bool {
	isSucc := true
	// log.Info(`dump def:`, spew.Sdump(r.PipeObj.(*models.Pipeline).Def))
	var newDef models.Variables
	bDef := r.render.fillData(toJson(r.PipeObj.(*models.Pipeline).Def), nil)
	json.Unmarshal(bDef, &newDef)
	r.PipeObj.(*models.Pipeline).Def = newDef

	// spew.Dump(r.PipeObj.(*models.Pipeline).Def)
	log.Info(`dump newDef:`, spew.Sdump(r.PipeObj.(*models.Pipeline).Def))

	//backup self so can revert after recursive
	r.mementoes.SaveMemento(&memento{r.PipeObj.(*models.Pipeline)})
	r.mementoes.SaveMemento(&memento{r.refs.Parent()})
	r.mementoes.SaveMemento(&memento{r.Reporter.Details.Parent()})
	log.Info(`cur refs:`, r.refs)
	log.Info(`parentRef:`, r.refs.Parent())
	log.Info(`cur detail:`, r.Reporter.Details)
	log.Info(`parent detail:`, r.Reporter.Details.Parent())

	for index, execNode := range r.PipeObj.(*models.Pipeline).Steps {
		stepSuccCounter := 0
		if r.Status == Canceled {
			// 如果runner的已经取消了，就没必要再去执行下一个用例了
			log.Info(`executor stopping,because runner is canceled `)
			return false
		}
		node := refNode.New(execNode.RefTag())
		r.refs.AddChild(node)

		detail := models.NewResultTree()
		detail.Title = execNode.Desc
		r.Reporter.Details.Append(detail)
		//开始执行step
		log.Info(`--------------step begin--------------------`)
		retryTimes := execNode.Retry
		repeatTimes := execNode.Repeat
		log.Info(fmt.Sprintf("retryTimes:%d\trepeatTimes:%d\n", retryTimes, repeatTimes))
		// backupPipeObj := r.PipeObj.(*models.Pipeline)
		r.mementoes.SaveMemento(&memento{r.PipeObj.(*models.Pipeline)})
		helper.Apply(repeatTimes+1, func(curTime int) {
			log.Info(fmt.Sprintf(`cur loop is %d`, curTime+1))
			log.Info(`>>>>cur refs:`, r.refs)
			log.Info(`>>>>cur pipeline:`, r.PipeObj)
			if curTime > 0 {
				r.PipeObj = r.mementoes.PopMementoWith(r.PipeObj).GetState().(*models.Pipeline)
				r.mementoes.SaveMemento(&memento{r.PipeObj.(*models.Pipeline)})
			}
			if !executeStep(&execNode, r, index) {
				//if step failed
				// stepSuccCounter -= 1
				for j := 0; j < retryTimes; j++ {
					log.Info(fmt.Sprintf(`retry times %d,cur is %d`, retryTimes, j))
					r.PipeObj = r.mementoes.PopMementoWith(r.PipeObj).GetState().(*models.Pipeline)
					r.mementoes.SaveMemento(&memento{r.PipeObj.(*models.Pipeline)})
					if executeStep(&execNode, r, index) {
						stepSuccCounter += 1
						break
					}
				}
			} else {
				stepSuccCounter += 1
			}
		})
		log.Info(`--------------step end----------------------`)
		if stepSuccCounter != repeatTimes+1 {
			isSucc = false
		}

	}
	r.Reporter.Details = r.mementoes.PopMementoWith(r.Reporter.Details).GetState().(*models.ResultTree)
	r.refs = r.mementoes.PopMementoWith(r.refs).GetState().(refNode.Node)
	r.PipeObj = r.mementoes.PopMementoWith(r.PipeObj).GetState().(*models.Pipeline) //revert self
	log.Info(`after revert,cur refs:`, r.refs)
	log.Info(`after revert,cur detail:`, r.Reporter.Details)
	log.Info(`executePipeline result:`, isSucc)
	return isSucc
}

func executeStep(execNode *models.ExecNode, r *TestRunner, rindex int) bool {
	var isSucc bool
	log.Info(`exec step:`, execNode.Desc, ` rindex:`, rindex)
	pipeObj := r.PipeObj.(*models.Pipeline)
	// report := r.Reporter
	requestor := NewRequestor()
	ref := r.refs.ChildAt(rindex) //user1 signup login
	detail := r.Reporter.Details.ChildAt(rindex)
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

	if execNode.Exec.GetType() == models.TYPE_API {
		//只是执行接口
		record := models.NewRecord()
		apiObj := execNode.Exec.(*models.API)
		url := fmt.Sprintf(`%s/%s`, execNode.Host, apiObj.Path)
		//先合并参数
		for k, v := range execNode.Args {
			apiObj.Params[k] = v
		}
		bheader := r.render.fillData(toJson(apiObj.Headers), pipeObj.Def)
		log.Info(`pipeObj.Def:`, pipeObj.Def)
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
			log.Info(`parse MultipartFile:`, apiObj.MultipartFile.Json())
			bMpf := r.render.fillData(apiObj.MultipartFile.Json(), pipeObj.Def)
			json.Unmarshal(bMpf, &mpf)
			req := requestor.BuildPostFiles(url, mpf, header)
			record.Request = req
			startTime = time.Now()
			resp = requestor.doRequest(req, fnBeforeRequest, fnAfterResponse)
		} else {
			var params models.Params
			log.Info(`parse Params:`, toJson(apiObj.Params))
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
		log.Info(`got response:`, resp)
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
					log.Info("add ExportVars:", ek, bindVal)
				} else {
					//plain text
					regx := regexp.MustCompile(v)
					match := regx.FindStringSubmatch(resp.Content)
					if match != nil {
						//目前暂不支持切片，如果是匹配多个值，只能是先合拼，到需要用的时候，自己再转换成字符串切片
						// services.VarsMgr.Add(fmt.Sprintf(`%s:%s`, render.tag, ek), strings.Join(match, `,`))
						ref.AddPairs(ek, strings.Join(match, `,`))
						// varsMgr.SetVar(this.Testcase.GetUid(), v.Name, strings.Join(match, `,`))
						log.Info("add ExportVars:", ek, strings.Join(match, `,`))
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
			log.Info(`parse validator.Actual:`, validator.Actual.(string))
			actual := string(r.render.fillData(validator.Actual.(string), data))

			log.Info(`parse validator.Expected:`, validator.Expected.(string))
			expected := string(r.render.fillData(validator.Expected.(string), pipeObj.Def))
			isPassed := So(actual, compare, expected)
			if !isPassed {
				allPassed = false
			}
			log.Info(fmt.Sprintf(`Actual:%v,Expected:%v,So %v`, actual, expected, isPassed))
			validator.Actual = actual
			validator.Expected = expected
			record.AddValidator(validator)
		}
		// TODO error and skip
		if allPassed {
			record.Stat = models.SUCCESS
			isSucc = true
		} else {
			record.Stat = models.FAILED
		}
		detail.SetRecord(record)
		// detail.Status.Count(record.Stat)
		// detail.Status = record.Stat
		execNode.Retry -= 1
		execNode.Repeat -= 1
	} else {
		subPipeObj := execNode.Exec.(*models.Pipeline)
		//先合并参数
		log.Info(`execNode.Args:`, spew.Sdump(execNode.Args))
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
		r.Reporter.Details = detail
		log.Info(`replace def:`, subPipeObj.Def)
		isSucc = executePipeline(r)
	}
	if isSucc || execNode.Retry <= 0 {
		if execNode.Repeat <= 0 {
			// report.AddDetail(*detail)
			// detail.
		}
	}
	return isSucc
}
