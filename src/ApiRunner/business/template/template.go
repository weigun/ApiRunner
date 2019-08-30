package template

import (
	"ApiRunner/business/template/parser"
)

type Template struct {
	*parser.Tree
}

func New() *Template {
	t := &Template{}
	t.init()
	return t
}
