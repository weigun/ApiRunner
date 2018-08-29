package runner

import (
	testcase "ApiRunner/case"
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
	Ready    chan bool //runner状态
}
type Response struct {
	// TODO 需要加入更多的字段，用于报告生成
	Code    int
	Content string
	ErrMsg  string
	elapsed int64
}

func (this *Response) GetCode() int {
	return this.Code
}

func (this *Response) GetContent() string {
	return this.Content
}

func (this *Response) GetErrMsg() string {
	return this.ErrMsg
}

func (this *Response) GetElapsed() int64 {
	return this.elapsed
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
	<-this.Ready
	log.Println("start test the testcase set")
	caseName := this.Testcase.GetCaseset().Conf.Name
	log.Println(caseName)
	for i, ci := range this.Testcase.GetCaseset().GetCases() {
		//顺序执行用例
		fmt.Printf("testcase %d\n", i)
		req := ci.BuildRequest() //构造请求体
		resp := this.doRequest(req)
		log.Println(resp)
		//验证结果
		//		validate(resp, ci.GetConditions())

	}
}

func (this *Runner) stop() {

}

func (this *Runner) doRequest(request *http.Request) Response {
	//执行请求
	startTime := time.Now().Unix()
	response, err := this.Core.Do(request)
	elapsed := time.Now().Unix() - startTime
	resp := Response{elapsed: elapsed}
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
