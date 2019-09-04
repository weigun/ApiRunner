package parser

import (
	// "errors"
	"fmt"
	// "strings"
	"ApiRunner/business/template/lexer"
)

type Node interface {
	Type() int
	String() string
}

//plain text
type textNode struct {
	*lexer.Token
}

func (n *textNode) Type() int {
	return n.Typ
}

func (n *textNode) String() string {
	return n.Val
}

type fieldNode struct {
	*lexer.Token
	subNodes []*lexer.Token
}

func (n *fieldNode) Type() int {
	return n.Typ
}

func (n *fieldNode) String() string {
	return ``
}

func (n *fieldNode) expand(t *lexer.Token) {
	if t.Typ == lexer.TokenField {
		n.subNodes = append(n.subNodes, t)
	}
}

type funcNode struct {
	*lexer.Token
	subNodes []*lexer.Token
}

func (n *funcNode) Type() int {
	return n.Typ
}

func (n *funcNode) String() string {
	return ``
}

func (n *funcNode) expand(t *lexer.Token) {
	switch t.Typ {
	case lexer.TokenRawParam, lexer.TokenVarParam:
		n.subNodes = append(n.subNodes, t)
	default:
		fmt.Println(`can not accept token:`, t)
	}
}
