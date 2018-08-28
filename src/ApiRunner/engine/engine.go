package engine

import (
	utils "ApiRunner/utils"
	_ "bytes"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	_ "log"
	"net/http"
	_ "regexp"
	"runtime/debug"
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

type Engine struct {
	//like a webclient
	Core *http.Client
}

type Response struct {
	// TODO 需要加入更多的字段，用于报告生成
	Code    int
	Content string
	ErrMsg  string
	elapsed int64
}

func NewEngine() *Engine {
	makeClient(client)
	return &Engine{Core: client}
}

func (this *Engine) safeRun(r Runner) {
	defer func() {
		// don't panic
		err := recover()
		if err != nil {
			debug.PrintStack()
		}
	}()
	r.start(this)
}

func (this *Engine) do(request *http.Request) Response {
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
