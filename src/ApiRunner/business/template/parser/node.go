package parser

import (
	// "errors"
	"fmt"
	"reflect"

	// "strings"
	"ApiRunner/business/template/lexer"
)

const CONTAINER int = -1

type Node interface {
	Type() int
	TranslateFrom(interface{}, interface{}) string //first interface{} must be a map[string]interface{}
}

type CompNode interface {
	Expand(Node)
}

//plain text
type textNode struct {
	*lexer.Token
}

func (n *textNode) Type() int {
	return n.Typ
}

func (n *textNode) TranslateFrom(data interface{}, execFuncs interface{}) string {
	return n.Val
}

//container node
type containerNode struct {
	subNodes []Node
}

func (n *containerNode) Type() int {
	if len(n.subNodes) <= 0 {
		return CONTAINER
	}
	return n.subNodes[0].Type()
}

func (n *containerNode) TranslateFrom(data interface{}, execFuncs interface{}) string {
	data = data.(map[string]interface{})
	switch n.Type() {
	case lexer.TokenField:
		// var tmpData interface{}
		tmpData := data.(map[string]interface{})
		i := 0
		for i < len(n.subNodes) {
			fieldName := n.subNodes[i].TranslateFrom(nil, nil)
			tmpData = tmpData[fieldName].(map[string]interface{})
			i++
		}
		return convertValue(tmpData)
	case lexer.TokenFuncName:
		funName := n.subNodes[0].TranslateFrom(nil, nil) //func value
		fun := execFuncs.(map[string]reflect.Value)[funName]
		numIn := len(n.subNodes) - 1
		args := make([]reflect.Value, numIn)
		for i := 1; i < numIn; i++ {
			// args := []reflect.Value{reflect.ValueOf("wudebao"), reflect.ValueOf(30)}
			v := n.subNodes[i].TranslateFrom(data, execFuncs)
			args = append(args, reflect.ValueOf(v))
		}
		return convertValue(fun.Call(args))
	default:
		fmt.Println(`unknow node in container`)
		return ``
	}
}

func (n *containerNode) Expand(t Node) {
	switch t.Type() {
	case lexer.TokenField, lexer.TokenFuncName, lexer.TokenRawParam, lexer.TokenVarParam:
		n.subNodes = append(n.subNodes, t)
	default:
		fmt.Println(`not accept sub node`)
	}
}

type fieldNode struct {
	*lexer.Token
}

func (n *fieldNode) Type() int {
	return n.Typ
}

func (n *fieldNode) TranslateFrom(data interface{}, execFuncs interface{}) string {
	return n.Val

}

type funcNode struct {
	*lexer.Token
}

func (n *funcNode) Type() int {
	return n.Typ
}

func (n *funcNode) TranslateFrom(data interface{}, execFuncs interface{}) string {
	return n.Val
}

type paramNode struct {
	*lexer.Token
}

func (n *paramNode) Type() int {
	return n.Typ
}

func (n *paramNode) TranslateFrom(data interface{}, execFuncs interface{}) string {
	return ``
}

type varNode struct {
	*lexer.Token
}

func (n *varNode) Type() int {
	return n.Typ
}

func (n *varNode) TranslateFrom(data interface{}, execFuncs interface{}) string {
	return ``
}

func convertValue(data interface{}) string {
	//将数据转换为字符串
	return fmt.Sprintf(`%v`, data)
}
