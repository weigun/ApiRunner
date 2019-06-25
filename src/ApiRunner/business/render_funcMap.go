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

func gen_email() string {
	return `test@qq1.com`
}

//导入自定义函数
func init() {
	//导入
	funcMap[`world`] = world
	funcMap[`gen_email`] = gen_email
}
