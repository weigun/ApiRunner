// requests.go
package runner

import (
	// "bytes"
	"crypto/tls"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	Url "net/url"

	// "strings"
	"sync"
	"time"

	"ApiRunner/utils"
)

const (
	timeout = 60
	maxConn = 100 //连接池数
)

var client *http.Client
var once sync.Once

func makeClient() {
	//TODO 可能非线程安全，需要改为once.DO方式
	once.Do(func() {
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
	})
}

type httpClient struct {
	Core *http.Client
}

func GetHttpClient() *httpClient {
	makeClient()
	return &httpClient{client}
}

func (this *httpClient) Get(url string, params Params) Response {
	request := buildRequest(url, "GET", encode(params))
	return this.doRequest(request)
}

func (this *httpClient) Post(url string, params Params) Response {
	request := buildRequest(url, "POST", encode(params))
	return this.doRequest(request)
}

func (this *httpClient) doRequest(request *http.Request) Response {
	//执行请求
	// log.Println("-------before request:", request)
	resp := Response{}
	response, err := this.Core.Do(request)
	if err != nil {
		resp.ErrMsg = err.Error()
		if response != nil {
			io.Copy(ioutil.Discard, response.Body)
			response.Body.Close()
		}
	} else {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Printf("%v\n", err)
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

type Params map[string]interface{}

func encode(params Params) string {
	//编码查询参数
	query := Url.Values{}
	for k, v := range params {
		query.Add(k, v.(string))
	}
	return query.Encode()
}

func toJson(params Params) string {
	//转json，用于post方法
	return utils.Map2Json(params)
}

func buildRequest(url, method, data string) *http.Request {
	//构造请求体
	//	fmt.Println("CaseItem:", this)
	// api := this.Api
	// method := this.Method
	// var data string
	// if len(this.Params.Params) == 0 {
	// 	data = ""
	// } else {
	// 	data = this.Params.Conver(uid, this.Method)
	// }
	// log.Println("BuildRequest:", data)
	// bodyData := bytes.NewBuffer([]byte(data)) //get方法默认是空字符串
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		panic(err.Error())
	}
	req.URL.RawQuery = data
	// log.Println(req)
	// for _, h := range this.Headers {
	// 	//		fmt.Println("-------------", h)
	// 	req.Header.Add(h.Key, h.Conver(uid))
	// }
	return req
}

type Response struct {
	Code    int
	Content string
	ErrMsg  string
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
