package parser

import (
	// "errors"
	// "fmt"
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

}

func (n *fieldNode) expand(t *lexer.Token) {
	n.subNodes = append(n.subNodes, t)
}
