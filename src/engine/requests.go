package main

import (
	"bytes"
	_ "crypto/tls"
	//	"encoding/json"
	_ "fmt"
	_ "io/ioutil"
	_ "log"
	//	"math/rand"
	"net/http"
	_ "regexp"
	//	"strconv"
	//	"strings"
	_ "io"
	_ "time"
)

func NewRequest(api, method string, data string) *http.Request {
	bodyData = bytes.NewBuffer([]byte(data)) //get方法默认是空字符串
	req, err := http.NewRequest(api, method, bodyData)
	if err != nil {
		panic(err.Error())
	}
	return req
}

/////////////////////////////////////////////////////
/*
type baseRequests interface {
	send(method, api, data string) Response
}

type httpRequests interface {
	baseRequests
	Get(api, data string) Response
	Post(api, data string) Response
	AddHeader(key, value string)
}

type header struct {
	key   string
	value string
}

type webClient struct {
	header
	instance *http.Client
}

type Response struct {
	Code    int
	Content string
	ErrMsg  string
}

func NewRequest() *webClient {
	makeClient(client)
	o := webClient{instance: client}
	return &o
}

func (this *webClient) AddHeader(key, value string) {
	//	主要是设置token
	this.header.key = key
	this.header.value = value
}

func (this *webClient) Error(api string, code int) string {
	return fmt.Sprintf(`%s failed,StatusCode is %d`, api, code)
}

func (this *webClient) Get(api, data string) Response {
	//	fmt.Println(api, data)
	request, err := http.NewRequest("GET", api, nil)
	if err != nil {
		return Response{ErrMsg: err.Error()}
	}
	this.setToken(request)
	if data != "" {
		request.URL.RawQuery = data
	}
	return this.send(request)
}

func (this *webClient) Post(api, data string) Response {
	body_data := bytes.NewBuffer([]byte(data))
	request, err := http.NewRequest("POST", api, body_data)
	if err != nil {
		return Response{ErrMsg: err.Error()}
	}
	this.setToken(request)
	return this.send(request)
}

func (this *webClient) setToken(request *http.Request) {
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set(this.header.key, this.header.value)
}

func (this *webClient) send(request *http.Request) Response {
	startTime := time.Now().Unix()
	response, err := this.instance.Do(request)
	elapsed := time.Now().Unix() - startTime
	resp := Response{}
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
			fmt.Println(request.Method, api, elapsed, this.Error(api, response.StatusCode))
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
*/
