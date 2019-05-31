// testcase_render.go
package business

import (
	"text/template"

	"github.com/Masterminds/sprig"
)

func GetTemplate(_func *template.FuncMap) *template.Template {
	t := template.New("conf")
	if _func != nil {
		t.Funcs(*_func)
	}
	return t //不能放到全局或者通过闭包的方式，因为这个是携程不安全的
}
