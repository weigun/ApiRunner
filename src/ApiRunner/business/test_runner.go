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

	// "github.com/davecgh/go-spew/spew"

	"ApiRunner/models"
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
	CaseObj  models.ICaseObj
	canceler context.CancelFunc
	Status   int
}

func NewTestRunner(id string, caseObj models.ICaseObj) *TestRunner {
	return &TestRunner{
		ID:      id,
		CaseObj: caseObj,
	}
}

func (r *TestRunner) Start() {
	log.Println("testrunner started")
	ctx, cancel := context.WithCancel(context.Background())
	valueCtx := context.WithValue(ctx, `status`, r.Status)
	r.canceler = cancel
	go func(ctx context.Context) {
		//TODO 需要保存堆栈
		go execute(r) //用例执行
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
}

func execute(r *TestRunner) {
	//具体执行用例的实体函数
	caseObj := r.CaseObj
	switch r.CaseObj.(type) {
	case *models.TestCase:
		// caseObj = r.CaseObj.(*models.TestCase)
	case *models.TestSuites:
		// caseObj = r.CaseObj.(*models.TestSuites)
	default:
		log.Printf(`unknow caseobj type:%T,stop runner`, r.CaseObj)
		r.canceler()
		return
	}

	//顺序执行用例
	render := newRenderer(r.ID)
	// requestor := NewRequestor()
	_type := caseObj.GetType()
	if _type == models.TYPE_TESTCASE {
		caseObj := r.CaseObj.(*models.TestCase)
		executeTestCase(render, caseObj, r)
	} else {
		caseObj := r.CaseObj.(*models.TestSuites)
		// spew.Dump(caseObj)
		var caseConf models.CaseConfig
		err := render.renderObj(caseObj.Config.Json(), true, &caseConf)
		if err != nil {
			log.Println(`renderObj error:`, err.Error())
			return
		}
		caseObj.Config = caseConf
		//将全局变量同步到变量服务
		for varName, varVal := range caseConf.Variables {
			services.VarsMgr.Add(fmt.Sprintf(`%s:%s`, render.tag, varName), varVal)
		}
		for _, caseItem := range caseObj.CaseList {
			for caseName, ts := range caseItem {
				log.Println(`executeTestCase:`, caseName)
				executeTestCase(render, &ts, r)
			}
		}
	}
}

func executeTestCase(render *renderer, caseObj *models.TestCase, r *TestRunner) {
	// spew.Dump(caseObj)
	requestor := NewRequestor()
	//caseConfStr := renderTestCase(caseObj.Config.Json(), true)
	//caseConf := json.Unmarshal([]byte(caseConfStr), &models.CaseConfig{})
	var caseConf models.CaseConfig
	log.Println(`render config`)
	err := render.renderObj(caseObj.Config.Json(), true, &caseConf)
	if err != nil {
		log.Println(`renderObj error:`, err.Error())
		return
	}
	caseObj.Config = caseConf
	//将全局变量同步到变量服务
	for varName, varVal := range caseConf.Variables {
		services.VarsMgr.Add(fmt.Sprintf(`%s:%s`, render.tag, varName), varVal)
	}
	for _, api := range caseObj.APIS {
		if r.Status == Cancel {
			// 如果runner的已经取消了，就没必要再去执行下一个用例了
			log.Println(`executor stopping,because runner is canceled `)
			return
		}
		//将接口的局部变量同步到变量服务
		log.Println(`render Variables`)
		var localVars models.Variables
		err := render.renderObj(utils.Map2Json(api.Variables), true, &localVars)
		if err != nil {
			log.Println(`renderObj Variables error:`, err.Error())
		}
		api.Variables = localVars
		for varName, varVal := range api.Variables {
			services.VarsMgr.Add(fmt.Sprintf(`%s:%s`, render.tag, varName), varVal)
		}

		url := fmt.Sprintf(`%s/%s`, render.renderValue(caseObj.Config.Host, true), render.renderValue(api.Path, true))
		// TODO:
		// 模板翻译-done
		// 拦截器-done
		// MultipartFile-done

		//MultipartFile比普通post请求优先级要高
		var header models.Header
		var resp *Response

		render.renderObj(toJson(api.Headers), true, &header)
		if api.MultipartFile.IsEnabled() {
			var mpf models.MultipartFile
			render.renderObj(api.MultipartFile.Json(), true, &mpf)
			req := requestor.BuildPostFiles(url, mpf, header)
			resp = requestor.doRequest(req, api.BeforeRequest, api.AfterResponse)
		} else {
			var params models.Params
			render.renderObj(toJson(api.Params), true, &params)
			req := requestor.BuildRequest(url, render.renderValue(api.Method, true), params, header)
			resp = requestor.doRequest(req, api.BeforeRequest, api.AfterResponse)
		}
		log.Println(resp)
		data := make(map[string]interface{})
		data[`StatusCode`] = resp.Code
		//导出变量，如token等
		if resp.ErrMsg == `` {
			//没有错误的时候才能导出变量
			//TODO assert code??
			for ek, ev := range api.Export {
				v := ev.(string)
				if strings.Index(v, `{{`) != -1 && strings.Index(v, `}}`) != -1 {
					//返回json则需要提取变量
					contentMap := utils.Json2Map([]byte(resp.Content))
					data[`body`] = contentMap
					bindVal := render.renderWithData(v, data)
					services.VarsMgr.Add(fmt.Sprintf(`%s:%s`, render.tag, ek), bindVal)
					log.Println("add ExportVars:", render.tag, ek, bindVal)
				} else {
					//plain text
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
		for _, validator := range api.Validate {
			// TODO 渲染变量时，适配各种数据类型
			validator.Check = validator.Actual.(string)
			compare := getAssertByOp(validator.Op)
			actual := render.renderWithData(validator.Actual.(string), data)
			expected := render.renderValue(validator.Expected.(string), true)
			isPassed := So(actual, compare, expected)
			log.Printf(`Actual:%v,Expected:%v,So %v`, actual, expected, isPassed)
		}

	}
}
