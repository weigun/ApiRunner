// test_runner.go
package business

import (
	"context"
	"fmt"
	_ "fmt"
	"log"
	_ "net/url"
	"regexp"
	"strings"
	"time"

	// "github.com/davecgh/go-spew/spew"

	refNode "ApiRunner/business/refs_tree"
	"ApiRunner/models"
	"ApiRunner/models/young"
	"ApiRunner/services"
	"ApiRunner/utils"
)

const (
	Queuing = iota
	Running
	Passed
	Failed
	Cancel
)

type TestRunner struct {
	ID       string
	PipeObj  models.IPipe
	canceler context.CancelFunc
	refs     refNode.Node
	Status   int
}

func NewTestRunner(id string, pipeObj models.IPipe) *TestRunner {
	return &TestRunner{
		ID:      id,
		PipeObj: pipeObj,
		refs:    refNode.New(`root`),
	}
}

func (r *TestRunner) Start() {
	log.Println("testrunner started")
	ctx, cancel := context.WithCancel(context.Background())
	valueCtx := context.WithValue(ctx, `status`, r.Status)
	r.canceler = cancel
	//new report
	report := models.NewReport()
	go func(ctx context.Context) {
		//TODO 需要保存堆栈
		go execute(r, report) //用例执行
		// select {
		// case <-ctx.Done():
		// 	log.Println("runner done")
		// 	r.Status = ctx.Value(`status`).(int)
		// 	return
		// }
		<-ctx.Done()
		log.Println("runner done")
		r.Status = ctx.Value(`status`).(int) //TODO 大丈夫？

	}(valueCtx)

}

func (r *TestRunner) Stop() {
	log.Println("testrunner stopping")
	r.Status = Cancel
	r.canceler()
	// TODO 干掉eventbus的channel
}

func execute(r *TestRunner, report *models.Report) {
	//具体执行用例的实体函数
	pipeObj := r.PipeObj
	switch r.PipeObj.(type) {
	case *models.Pipeline:
		// pipeObj = r.PipeObj.(*models.Pipeline)
	case *models.PipeGroup:
		// pipeObj = r.PipeObj.(*models.PipeGroup)
	default:
		log.Printf(`unknow caseobj type:%T,stop runner`, r.PipeObj)
		r.canceler()
		return
	}

	//顺序执行用例
	//开始计时
	sum := models.NewSummary()
	sum.StartAt = time.Now()
	report.SetSummary(*sum)
	render := newRenderer(r.ID)
	_type := pipeObj.GetType()
	if _type == models.TYPE_PIPELINE {
		pipeObj := r.PipeObj.(*models.Pipeline)
		executePipeline(render, pipeObj, r, report)
		// executeTestCase(render, pipeObj, r, report)
	} else {
		pipeObj := r.PipeObj.(*models.PipeGroup)
		// spew.Dump(pipeObj)
		// var caseConf models.CaseConfig
		// err := render.renderObj(pipeObj.Config.Json(), true, &caseConf)
		// if err != nil {
		// 	log.Println(`renderObj error:`, err.Error())
		// 	return
		// }
		// pipeObj.Config = caseConf
		//将全局变量同步到变量服务
		// for varName, varVal := range pipeObj.Variables {
		// 	services.VarsMgr.Add(fmt.Sprintf(`%s:%s`, render.tag, varName), varVal)
		// }
		for _, pipeline := range pipeObj.Pipelines {
			log.Println(`executeTestCase:`, pipeline.Name)
			executePipeline(render, &pipeline, r, report)
		}
	}
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
}

func executePipeline(render *renderer, pipeObj *models.Pipeline, r *TestRunner, report *models.Report) {
	requestor := NewRequestor()
	detail := models.NewDetail()
	detail.Title = pipeObj.Name
	for index, stage := range pipeObj.Stages {
		if r.Status == Cancel {
			// 如果runner的已经取消了，就没必要再去执行下一个用例了
			log.Println(`executor stopping,because runner is canceled `)
			return
		}
		//将接口的局部变量同步到变量服务
		log.Println(`render Variables`)
		var localVars models.Variables
		err := render.renderObj(utils.Map2Json(stage.Env), true, &localVars)
		if err != nil {
			log.Println(`renderObj Variables error:`, err.Error())
		}
		stage.Env = localVars
		node := refNode.New(stage.RefTag())
		r.refs.AddChild(node)
		for varName, varVal := range stage.Env {
			services.VarsMgr.Add(fmt.Sprintf(`%s:%s`, render.tag, varName), varVal)
			node.AddPairs(varName, varVal)
		}

		//开始执行step
		executeStage(render, &stage, r, report, index)
	}
}

func executeStage(render *renderer, stageObj *models.Stage, r *TestRunner, report *models.Report, rindex int) {
	requestor := NewRequestor()
	ref := r.refs.ChildAt(rindex)

	for _, step := range stageObj.Steps {
		if r.Status == Cancel {
			// 如果runner的已经取消了，就没必要再去执行下一个用例了
			log.Println(`executor stopping,because runner is canceled `)
			return
		}
		//将接口的局部变量同步到变量服务
		// log.Println(`render Variables`)
		// var localVars models.Variables
		// err := render.renderObj(utils.Map2Json(stage.Env), true, &localVars)
		// if err != nil {
		// 	log.Println(`renderObj Variables error:`, err.Error())
		// }
		// stage.Env = localVars
		node := refNode.New(step.RefTag())
		ref.AddChild(node)
		for varName, varVal := range step.Params {
			services.VarsMgr.Add(fmt.Sprintf(`%s:%s`, render.tag, varName), varVal)
			node.AddPairs(varName, varVal)
		}
		//write to here,stop
		url := fmt.Sprintf(`%s/%s`, render.renderValue(pipeObj.Host, true), render.renderValue(stage.Path, true))
		// TODO:
		// 模板翻译-done
		// 拦截器-done
		// MultipartFile-done

		//MultipartFile比普通post请求优先级要高
		var header models.Header
		var resp *young.Response

		record := models.NewRecord()

		render.renderObj(toJson(stage.Headers), true, &header)
		var startTime time.Time
		if stage.MultipartFile.IsEnabled() {
			var mpf models.MultipartFile
			render.renderObj(stage.MultipartFile.Json(), true, &mpf)
			req := requestor.BuildPostFiles(url, mpf, header)
			record.Request = req
			startTime = time.Now()
			resp = requestor.doRequest(req, stage.BeforeRequest, stage.AfterResponse)
		} else {
			var params models.Params
			render.renderObj(toJson(stage.Params), true, &params)
			req := requestor.BuildRequest(url, render.renderValue(stage.Method, true), params, header)
			record.Request = req
			startTime = time.Now()
			resp = requestor.doRequest(req, stage.BeforeRequest, stage.AfterResponse)
		}
		elapsed := time.Now().Sub(startTime).Nanoseconds()
		log.Println(resp)
		record.Desc = stage.Name
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
			for ek, ev := range stage.Export {
				v := ev.(string)
				if strings.Index(v, `{{`) != -1 && strings.Index(v, `}}`) != -1 {
					//返回json则需要提取变量
					// contentMap := utils.Json2Map([]byte(resp.Content))
					// data[`body`] = contentMap
					bindVal := render.renderWithData(v, data)
					services.VarsMgr.Add(fmt.Sprintf(`%s:%s`, render.tag, ek), bindVal)
					log.Println("add ExportVars:", render.tag, ek, bindVal)
				} else {
					//plain text
					// data[`body`] = resp.Content
					regx := regexp.MustCompile(v)
					match := regx.FindStringSubmatch(resp.Content)
					if match != nil {
						//目前暂不支持切片，如果是匹配多个值，只能是先合拼，到需要用的时候，自己再转换成字符串切片
						services.VarsMgr.Add(fmt.Sprintf(`%s:%s`, render.tag, ek), strings.Join(match, `,`))
						// varsMgr.SetVar(this.Testcase.GetUid(), v.Name, strings.Join(match, `,`))
						log.Println("add ExportVars:", render.tag, ek, strings.Join(match, `,`))
					}
				}
			}
		}

		//比较结果
		allPassed := true
		for _, validator := range stage.Validate {
			// TODO 渲染变量时，适配各种数据类型
			validator.Check = validator.Actual.(string)
			compare := getAssertByOp(validator.Op)
			actual := render.renderWithData(validator.Actual.(string), data)
			expected := render.renderValue(validator.Expected.(string), true)
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
	}
	report.AddDetail(*detail)
}

func executeTestCase(render *renderer, pipeObj *models.Pipeline, r *TestRunner, report *models.Report) {
	// spew.Dump(pipeObj)
	requestor := NewRequestor()
	detail := models.NewDetail()
	// var caseConf models.CaseConfig
	// log.Println(`render config`)
	// err := render.renderObj(pipeObj.Config.Json(), true, &caseConf)
	// if err != nil {
	// 	log.Println(`renderObj error:`, err.Error())
	// 	return
	// }
	// pipeObj.Config = caseConf
	detail.Title = pipeObj.Name
	//将全局变量同步到变量服务
	// for varName, varVal := range caseConf.Variables {
	// 	services.VarsMgr.Add(fmt.Sprintf(`%s:%s`, render.tag, varName), varVal)
	// }

	for _, stage := range pipeObj.Stages {
		if r.Status == Cancel {
			// 如果runner的已经取消了，就没必要再去执行下一个用例了
			log.Println(`executor stopping,because runner is canceled `)
			return
		}
		//将接口的局部变量同步到变量服务
		log.Println(`render Variables`)
		var localVars models.Variables
		err := render.renderObj(utils.Map2Json(stage.Env), true, &localVars)
		if err != nil {
			log.Println(`renderObj Variables error:`, err.Error())
		}
		stage.Env = localVars
		node := refNode.New(stage.RefTag())
		node.SetParent(refs)
		refs.AddChild(node)
		for varName, varVal := range stage.Env {
			services.VarsMgr.Add(fmt.Sprintf(`%s:%s`, render.tag, varName), varVal)
			node.AddPairs(varName, varVal)
		}

		url := fmt.Sprintf(`%s/%s`, render.renderValue(pipeObj.Host, true), render.renderValue(stage.Path, true))
		// TODO:
		// 模板翻译-done
		// 拦截器-done
		// MultipartFile-done

		//MultipartFile比普通post请求优先级要高
		var header models.Header
		var resp *young.Response

		record := models.NewRecord()

		render.renderObj(toJson(stage.Headers), true, &header)
		var startTime time.Time
		if stage.MultipartFile.IsEnabled() {
			var mpf models.MultipartFile
			render.renderObj(stage.MultipartFile.Json(), true, &mpf)
			req := requestor.BuildPostFiles(url, mpf, header)
			record.Request = req
			startTime = time.Now()
			resp = requestor.doRequest(req, stage.BeforeRequest, stage.AfterResponse)
		} else {
			var params models.Params
			render.renderObj(toJson(stage.Params), true, &params)
			req := requestor.BuildRequest(url, render.renderValue(stage.Method, true), params, header)
			record.Request = req
			startTime = time.Now()
			resp = requestor.doRequest(req, stage.BeforeRequest, stage.AfterResponse)
		}
		elapsed := time.Now().Sub(startTime).Nanoseconds()
		log.Println(resp)
		record.Desc = stage.Name
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
			for ek, ev := range stage.Export {
				v := ev.(string)
				if strings.Index(v, `{{`) != -1 && strings.Index(v, `}}`) != -1 {
					//返回json则需要提取变量
					// contentMap := utils.Json2Map([]byte(resp.Content))
					// data[`body`] = contentMap
					bindVal := render.renderWithData(v, data)
					services.VarsMgr.Add(fmt.Sprintf(`%s:%s`, render.tag, ek), bindVal)
					log.Println("add ExportVars:", render.tag, ek, bindVal)
				} else {
					//plain text
					// data[`body`] = resp.Content
					regx := regexp.MustCompile(v)
					match := regx.FindStringSubmatch(resp.Content)
					if match != nil {
						//目前暂不支持切片，如果是匹配多个值，只能是先合拼，到需要用的时候，自己再转换成字符串切片
						services.VarsMgr.Add(fmt.Sprintf(`%s:%s`, render.tag, ek), strings.Join(match, `,`))
						// varsMgr.SetVar(this.Testcase.GetUid(), v.Name, strings.Join(match, `,`))
						log.Println("add ExportVars:", render.tag, ek, strings.Join(match, `,`))
					}
				}
			}
		}

		//比较结果
		allPassed := true
		for _, validator := range stage.Validate {
			// TODO 渲染变量时，适配各种数据类型
			validator.Check = validator.Actual.(string)
			compare := getAssertByOp(validator.Op)
			actual := render.renderWithData(validator.Actual.(string), data)
			expected := render.renderValue(validator.Expected.(string), true)
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
	}
	report.AddDetail(*detail)
}
