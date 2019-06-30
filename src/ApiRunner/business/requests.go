// requests.go
package business

import (
	"bytes"
	"crypto/tls"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptrace"
	Url "net/url"

	// "strings"
	"sync"
	"time"

	"ApiRunner/models"
	"ApiRunner/utils"
)

const (
	TIMEOUT          = 60
	MAX_CONNEECTIONS = 100 //连接池数
)

type RefReq = *http.Request
type RefRsp = *http.Response

var client *http.Client
var once sync.Once

func makeClient() {
	once.Do(func() {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			MaxIdleConnsPerHost: MAX_CONNEECTIONS,
			MaxIdleConns:        MAX_CONNEECTIONS,
			DisableCompression:  true,
			DisableKeepAlives:   false,
		}
		client = &http.Client{
			Transport: tr,
			Timeout:   time.Duration(TIMEOUT) * time.Second,
		}
	})
}

type requests struct {
	client *http.Client
}

func NewRequestor() *requests {
	makeClient()
	return &requests{client}
}

func (this *requests) Get(url string, params models.Params) *Response {
	request := this.BuildRequest(url, "GET", params)
	return this.doRequest(request)
}

func (this *requests) Post(url string, params models.Params) *Response {
	request := this.BuildRequest(url, "POST", params)
	return this.doRequest(request)
}

func (this *requests) doRequest(request *http.Request) *Response {
	//执行请求
	// log.Println("-------before request:", request)
	resp := Response{}
	//hook test
	trace := &httptrace.ClientTrace{
		GotConn: func(connInfo httptrace.GotConnInfo) {
			log.Println("Got Conn")
		},
		ConnectStart: func(network, addr string) {
			log.Println("Dial start")
		},
		ConnectDone: func(network, addr string, err error) {
			log.Println("Dial done")
		},
		GotFirstResponseByte: func() {
			log.Println("First response byte!")
		},
		WroteHeaders: func() {
			log.Println("Wrote headers")
		},
		WroteRequest: func(wr httptrace.WroteRequestInfo) {
			log.Println("Wrote request", wr)
		},
	}
	request = request.WithContext(httptrace.WithClientTrace(request.Context(), trace))
	response, err := this.client.Do(request)
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
	return &resp
}

func (this *requests) BuildRequest(url, method string, params models.Params) *http.Request {
	//构造请求体
	// TODO 加上文件上传
	log.Println("BuildRequest:", params)
	var data string
	switch method {
	case `GET`:
		data = encode(params)
	case `POST`:
		data = toJson(params)
	}
	bodyData := bytes.NewBuffer([]byte(data)) //get方法默认是空字符串
	req, err := http.NewRequest(method, url, bodyData)
	if err != nil {
		log.Printf(`BuildRequest error %v\n`, err.Error())
	}
	if method == `GET` {
		req.URL.RawQuery = data
	}
	//TODO add header
	return req
}

func encode(params models.Params) string {
	//编码查询参数
	query := Url.Values{}
	for k, v := range params {
		query.Add(k, v.(string))
	}
	return query.Encode()
}

func toJson(params models.Params) string {
	//转json，用于post方法
	return utils.Map2Json(params)
}

type Response struct {
	Code    int
	Content string
	ErrMsg  string
}
