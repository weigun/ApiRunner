// requests.go
package business

import (
	"bytes"
	"crypto/tls"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	Url "net/url"
	"sync"
	"time"

	"ApiRunner/models"
	"ApiRunner/models/young"
	"ApiRunner/utils"
)

const (
	TIMEOUT          = 60
	MAX_CONNEECTIONS = 100 //连接池数
)

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

func (this *requests) Get(url string, params models.Params) *young.Response {
	request := this.BuildRequest(url, "GET", params, nil)
	return this.doRequest(request, ``, ``)
}

func (this *requests) Post(url string, params models.Params) *young.Response {
	request := this.BuildRequest(url, "POST", params, nil)
	return this.doRequest(request, ``, ``)
}

func (this *requests) doRequest(request *http.Request, beforeReqHook, afterRespHook string) *young.Response {
	//执行请求
	// log.Info("-------before request:", request)
	resp := young.Response{}
	//beforeRequest hook
	if beforeReqHook == `` {
		beforeReqHook = `beforeRequest`
	} else {
		if _, ok := hookMap[beforeReqHook]; !ok {
			log.Info(`hook not found [`, beforeReqHook, `],use default`)
			beforeReqHook = `beforeRequest`
		}
	}
	if afterRespHook == `` {
		afterRespHook = `afterResponse`
	} else {
		if _, ok := hookMap[afterRespHook]; !ok {
			log.Info(`hook not found [`, afterRespHook, `],use default`)
			afterRespHook = `afterResponse`
		}
	}
	hook := hooks{hookMap[beforeReqHook], hookMap[afterRespHook]}
	request = hook.beforeRequest(request).(*http.Request)
	response, err := this.client.Do(request)
	response = hook.afterResponse(response).(*http.Response)
	// resp.RAW = response
	if err != nil {
		resp.ErrMsg = err.Error()
		if response != nil {
			io.Copy(ioutil.Discard, response.Body)
			response.Body.Close()
		}
	} else {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Info(err.Error())
		} else {
			//			log.Info(string(body))
		}
		resp.Code = response.StatusCode
		resp.Header = response.Header
		resp.Cookies = response.Cookies()
		resp.Content = string(body)
		response.Body.Close()
		io.Copy(ioutil.Discard, response.Body)
	}
	return &resp
}

func (this *requests) BuildRequest(url, method string, params models.Params, header models.Header) *http.Request {
	//构造请求体
	// TODO 加上文件上传
	log.Info("BuildRequest:", params)
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
		log.Info(`BuildRequest error `, err.Error())
	}
	if method == `GET` {
		req.URL.RawQuery = data
	}
	// add header
	for k, v := range header {
		req.Header.Add(k, v.(string))
	}
	return req
}

func (this *requests) BuildPostFiles(url string, mpf models.MultipartFile, header models.Header) *http.Request {
	//构造文件上传的请求体
	log.Info("BuildPostFiles:", mpf)
	bodyBuf := bytes.NewBufferString(``)
	bodyWriter := multipart.NewWriter(bodyBuf)
	// boundary默认会提供一组随机数，也可以自己设置。
	bodyWriter.SetBoundary("75FE3EACA32442E1522CE5C98C6DC891")

	//数据域
	for k, v := range mpf.Params {
		bodyWriter.WriteField(k, v.(string))
	}

	//文件域
	for k, v := range mpf.Files {
		filePath := v.(string)
		_, err := bodyWriter.CreateFormFile(k, filePath)
		if err != nil {
			log.Fatal("create post file[", filePath, "] error", err.Error())
		}
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Fatal("ReadFile [", filePath, "] error", err.Error())
		}
		bodyBuf.Write(content)

	}
	bodyWriter.Close()

	reqReader := io.MultiReader(bodyBuf)
	req, err := http.NewRequest("POST", url, reqReader)
	if err != nil {
		log.Fatal(`BuildPostFiles error `, err.Error())

	}

	// add header
	for k, v := range header {
		if k == `Content-Type` {
			req.Header.Set("Content-Type", bodyWriter.FormDataContentType())
		} else {
			req.Header.Add(k, v.(string))
		}

	}
	// req.Header.Set("Content-Type", bodyWriter.FormDataContentType())
	req.ContentLength = int64(bodyBuf.Len())
	log.Info("request len:", req.ContentLength)
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
