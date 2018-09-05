package runner

import (
	testcase "ApiRunner/case"
	report "ApiRunner/report"
	utils "ApiRunner/utils"
	validation "ApiRunner/validation"
	_ "bytes"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	_ "strings"
	"time"
)

const (
	timeout = 15
	maxConn = 100 //连接池数
)

var client *http.Client

func makeClient(_client *http.Client) {
	//TODO 可能非线程安全，需要改为once.DO方式
	if _client != nil {
		return
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		MaxIdleConnsPerHost: maxConn,
		MaxIdleConns:        maxConn,
		DisableCompression:  true,
		DisableKeepAlives:   false,
	}
	client = &http.Client{
		Transport: tr,
		Timeout:   time.Duration(timeout) * time.Second,
	}
}

type Runner struct {
	Core     *http.Client
	Testcase testcase.PIparserInsterface
	Ready    chan bool //runner状态,如果true，则表明是web模式，false则是本地模式
}

func getErr(api string, code int) string {
	return fmt.Sprintf(`%s failed,StatusCode is %d`, api, code)
}

func NewRunner(cs testcase.PIparserInsterface) Runner {
	makeClient(client)
	return Runner{client, cs, make(chan bool, 1)}
}

func (this *Runner) Start() {
	//start test the testcase set
	status := <-this.Ready
	log.Println("start test the testcase set")
	s := report.Summary{Title: utils.GetDateTime(), StartTime: time.Now().Unix()}
	s.Add2Cache(this.Testcase.GetUid())
	caseName := this.Testcase.GetCaseset().Conf.Name
	log.Println(caseName)
	info := report.Info{caseName, this.Testcase.GetCaseset().Conf.Host}
	info.Add2Cache(this.Testcase.GetUid())
	resPool := validation.NewResultPool()
	for i, ci := range this.Testcase.GetCaseset().GetCases() {
		//顺序执行用例
		log.Printf("testcase %d\n", i)
		req := ci.BuildRequest() //构造请求体
		resp := this.doRequest(req)
		ts := this.Testcase
		log.Println(resp)
		resPool.Push(validation.ResultItem{ts, int64(i), resp}) //推送到结果池进行验证
		//TODO 各种log需要集中到log中心，因为在报表性需要查看log信息

	}
	resPool.Done(this.Testcase.GetUid())
	if !status {
		log.Println("testcase", caseName, "finished ", "waiting for result to finish")
		//		resPool.WaitForDone(this.Testcase.GetUid())
	}
}

func (this *Runner) stop() {

}

func (this *Runner) doRequest(request *http.Request) validation.Response {
	//执行请求
	startTime := utils.Now4ms()
	response, err := this.Core.Do(request)
	elapsed := utils.Now4ms() - startTime
	resp := validation.Response{Elapsed: elapsed, Header: response.Header}
	api := request.URL.String()
	if err != nil {
		resp.ErrMsg = err.Error()
		if response != nil {
			io.Copy(ioutil.Discard, response.Body)
			response.Body.Close()
		}
	} else {
		if response.StatusCode == http.StatusOK {
			log.Println(request.Method, api, elapsed, response.ContentLength)
		} else {
			log.Println(request.Method, api, elapsed, getErr(api, response.StatusCode))
		}
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			//			log.Printf("%v\n", err)
		} else {
			//			log.Println(string(body))
		}
		resp.Code = response.StatusCode
		resp.Content = string(body)
		response.Body.Close()
		io.Copy(ioutil.Discard, response.Body)
	}
	return resp
}
