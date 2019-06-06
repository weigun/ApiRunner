// render_funcMap.go
package business

import (
	"github.com/Masterminds/sprig"
)

//自定义一个funcMap，以sprig这个库为基础
var funcMap = sprig.TxtFuncMap()

//可以在这里编写自定义函数
func world() string {
	return `world`
}

//导入自定义函数
func init() {
	//导入
	funcMap[`world`] = world
}
