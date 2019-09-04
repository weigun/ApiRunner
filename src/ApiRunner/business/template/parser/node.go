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
	subNodes []Node
}

func (n *fieldNode) Type() int {
	return n.Typ
}

func (n *fieldNode) String() string {
	return ``
}

func (n *fieldNode) expand(t Node) {
	if t.Type() == lexer.TokenField {
		n.subNodes = append(n.subNodes, t)
	}
}

type funcNode struct {
	*lexer.Token
	subNodes []Node
}

func (n *funcNode) Type() int {
	return n.Typ
}

func (n *funcNode) String() string {
	return ``
}

func (n *funcNode) expand(t Node) {
	switch t.Type() {
	case lexer.TokenRawParam, lexer.TokenVarParam:
		n.subNodes = append(n.subNodes, t)
	default:
		fmt.Println(`can not accept token:`, t)
	}
}

type paramNode struct {
	*lexer.Token
}

func (n *paramNode) Type() int {
	return n.Typ
}

func (n *paramNode) String() string {
	return ``
}
