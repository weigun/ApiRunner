package parser

import (
	"strings"
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
	if data == nil {
		return `${nil}`
	}
	data = data.(map[string]interface{})
	switch n.Type() {
	case lexer.TokenField:
		// var tmpData interface{}
		tmpData := data.(map[string]interface{})
		i := 0
		for i < len(n.subNodes) {
			fieldName := n.subNodes[i].TranslateFrom(nil, nil)
			switch tmpData[fieldName].(type) {
			case map[string]interface{}:
				tmpData = tmpData[fieldName].(map[string]interface{})
			default:
				return convertValue(tmpData[fieldName])
			}
			// tmpData = tmpData[fieldName].(map[string]interface{})
			fmt.Println(`field name:`, fieldName, tmpData)
			i++
		}
		return convertValue(tmpData)
	case lexer.TokenFuncName:
		funName := n.subNodes[0].TranslateFrom(nil, nil) //func value
		fun := execFuncs.(map[string]reflect.Value)[funName]
		numIn := len(n.subNodes) - 1
		args := make([]reflect.Value, 0, 1)
		for i := 1; i <= numIn; i++ {
			// args := []reflect.Value{reflect.ValueOf("wudebao"), reflect.ValueOf(30)}
			v := n.subNodes[i].TranslateFrom(data, execFuncs)
			fmt.Println(v)
			args = append(args, reflect.ValueOf(v))
		}
		for _, v := range args {
			fmt.Printf("type:%T,kind:%s,val:%v\n", v, v.Kind(), v)
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
	fmt.Println(n.Val)
	return n.Val
}

type paramNode struct {
	*lexer.Token
}

func (n *paramNode) Type() int {
	return n.Typ
}

func (n *paramNode) TranslateFrom(data interface{}, execFuncs interface{}) string {
	if x := strings.Index(n.Val, `$`); x >= 0 {
		//需要对变量进行翻译
		fmt.Println(n.Val[x+1:])
		val := data.(map[string]interface{})[n.Val[x+1:]]
		fmt.Println(`val:`, val)
		fmt.Printf("->type:%T,kind:%s,val:%v\n", val, reflect.ValueOf(val).Kind(), val)
		return convertValue(val)
	}
	return n.Val
}

type varNode struct {
	*lexer.Token
}

func (n *varNode) Type() int {
	return n.Typ
}

func (n *varNode) TranslateFrom(data interface{}, execFuncs interface{}) string {
	x := strings.Index(n.Val, `$`)
	fmt.Println(n.Val[x+1:])
	val := data.(map[string]interface{})[n.Val[x+1:]]
	fmt.Printf("%T,val:%v\n", val, val)
	return convertValue(val)
}

func convertValue(data interface{}) string {
	//将数据转换为字符串
	switch data.(type) {
	case []reflect.Value:
		data := data.([]reflect.Value)
		if len(data) == 0 {
			return ``
		}
		if len(data) > 2 || (len(data) == 2 && data[1].Type() != reflect.TypeOf((*error)(nil)).Elem()) {
			panic(`func not support return more than 2 value,and second must be error`)
		}
		switch data[0].Kind() {
		case reflect.String:
			return data[0].Interface().(string)
		case reflect.Int:
			return fmt.Sprintf(`%d`, data[0].Interface().(int))
		default:
			return fmt.Sprintf(`%s`, data[0])
		}
	case string:
		return data.(string)
	default:
		fmt.Printf("%T\n", data)
	}
	return fmt.Sprintf(`%v`, data)

}
