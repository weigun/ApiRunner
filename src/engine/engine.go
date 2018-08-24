package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	_ "io"
	_ "io/ioutil"
	_ "log"
	"net/http"
	_ "regexp"
	_ "time"
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

type engine struct {
	//like a webclient
	core *http.Client
}

func NewEngine() *engine {
	makeClient(client)
	return &engine{core: client}
}

func (this *engine) safeRun(r runner) {
	defer func() {
		// don't panic
		err := recover()
		if err != nil {
			debug.PrintStack()
		}
	}()
	r.start()
}

func (this *engine) runnerPool() {

}
