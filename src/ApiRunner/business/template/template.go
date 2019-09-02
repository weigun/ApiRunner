package template

import (
	"io"
	"reflect"

	"ApiRunner/business/template/parser"
)

/*
TODO:
1.需要捕获panic，并且进行恢复
2.同一个模板可以重复解析
3.Template允许并行化

*/
type FuncMap = map[string]interface{}

type Template struct {
	*parser.Tree
	execFuncs map[string]reflect.Value //TODO 这个结构可以抽出来在不同的模板之间共享
}

func (t *Template) Parse(text string) (*Template, error) {
	t.Tree.Parse(text)
	return t, nil
}

func (t *Template) Funcs(funcMap FuncMap) *Template {
	for name, fn := range funcMap {
		v := reflect.ValueOf(fn)
		if v.Kind() != reflect.Func {
			panic("value for " + name + " not a function")
		}
		t.execFuncs[name] = v
	}
	return t
}

func (t *Template) Execute(wr io.Writer, data interface{}) {

}

func New() *Template {
	t := &Template{&parser.Tree{}, make(map[string]reflect.Value)}
	return t
}
