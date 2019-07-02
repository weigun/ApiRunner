// hook_funcMap.go
package business

import (
	"log"
	// "net/http"
	// "github.com/davecgh/go-spew/spew"
)

var hookMap = make(map[string]hookFunc)

func beforeRequest(req interface{}) interface{} {
	log.Println(`hook beforeRequest trigger`)
	// spew.Dump(req)
	return req
}

func afterResponse(resp interface{}) interface{} {
	log.Println(`hook afterResponse trigger`)
	log.Println(resp)
	// spew.Dump(resp)
	return resp
}

//导入hook函数
func init() {
	//导入
	hookMap[`beforeRequest`] = beforeRequest
	hookMap[`afterResponse`] = afterResponse
}
