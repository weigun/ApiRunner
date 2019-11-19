// hook_funcMap.go
package business

// "log"
// "net/http"
// "github.com/davecgh/go-spew/spew"

var hookMap = make(map[string]hookFunc)

func beforeRequest(req interface{}) interface{} {
	log.Info(`hook beforeRequest trigger`)
	// spew.Dump(req)
	return req
}

func afterResponse(resp interface{}) interface{} {
	log.Info(`hook afterResponse trigger`)
	log.Info(resp)
	// spew.Dump(resp)
	return resp
}

//导入hook函数
func init() {
	//导入
	hookMap[`beforeRequest`] = beforeRequest
	hookMap[`afterResponse`] = afterResponse
}
