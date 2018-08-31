package runner

import (
	testcase "ApiRunner/case"
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
	caseName := this.Testcase.GetCaseset().Conf.Name
	log.Println(caseName)
	resPool := validation.NewResultPool()
	for i, ci := range this.Testcase.GetCaseset().GetCases() {
		//顺序执行用例
		fmt.Printf("testcase %d\n", i)
		req := ci.BuildRequest() //构造请求体
		resp := this.doRequest(req)
		ts := this.Testcase
		log.Println(resp)
		resPool.Push(validation.ResultItem{ts, resp}) //推送到结果池进行验证

	}
	if !status {
		log.Println("testcase", caseName, "finished ", "waiting for result to finish")
		resPool.WaitForDone()
	}
}

func (this *Runner) stop() {

}

func (this *Runner) doRequest(request *http.Request) validation.Response {
	//执行请求
	startTime := time.Now().Unix()
	response, err := this.Core.Do(request)
	elapsed := time.Now().Unix() - startTime
	resp := validation.Response{Elapsed: elapsed}
	api := request.URL.String()
	if err != nil {
		resp.ErrMsg = err.Error()
		if response != nil {
			io.Copy(ioutil.Discard, response.Body)
			response.Body.Close()
		}
	} else {
		if response.StatusCode == http.StatusOK {
			fmt.Println(request.Method, api, elapsed, response.ContentLength)
		} else {
			fmt.Println(request.Method, api, elapsed, getErr(api, response.StatusCode))
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
