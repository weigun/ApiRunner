// hook_funcMap.go
package business

import (
	"fmt"
	"net/http"
)

var hookMap = make(map[string]func(interface{}) interface{})

func beforeRequest(req interface{}) interface{} {
	fmt.Println("Hello World!")
	return nil
}

//导入hook函数
func init() {
	//导入
	funcMap[`beforeRequest`] = beforeRequest
}
