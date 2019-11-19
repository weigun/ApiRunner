package parser

import (
	"strings"
	// "errors"
	"fmt"
	"reflect"

	// "strings"
	"ApiRunner/business/refs_tree"
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
			case refs.Node:
				//如果引用树,需要将下一个节点替换掉tmpdata
				// refs.user1.login.ret_email
				refs := tmpData[fieldName].(refs.Node)
				if refs.Len() > 0 {
					for i := 0; i < refs.Len(); i++ {
						ref := refs.ChildAt(i)
						tmpData[ref.Name()] = ref
					}
				} else {
					//叶子节点了
					fieldName := n.subNodes[i+1].TranslateFrom(nil, nil)
					return refs.ValueOf(fieldName).(string) //默认当做字符串来处理
				}

			default:
				return convertValue(tmpData[fieldName])
			}
			// tmpData = tmpData[fieldName].(map[string]interface{})
			log.Debug(`field name:`, fieldName, tmpData)
			i++
		}
		return convertValue(tmpData)
	case lexer.TokenFuncName:
		funName := n.subNodes[0].TranslateFrom(nil, nil) //func value
		fun := execFuncs.(map[string]reflect.Value)[funName]
		numIn := len(n.subNodes) - 1
		args := make([]reflect.Value, 0, 1)
		for i := 1; i <= numIn; i++ {
			// v := n.subNodes[i].TranslateFrom(data, execFuncs)
			nv := n.subNodes[i].(*paramNode)
			v := nv.ValueFrom(data)
			log.Debug(v)
			args = append(args, v)
		}
		for _, v := range args {
			log.Debug("type:%T,kind:%s,val:%v\n", v, v.Kind(), v)
		}
		return convertValue(fun.Call(args))
	default:
		log.Debug(`unknow node in container`)
		return ``
	}
}

func (n *containerNode) Expand(t Node) {
	switch t.Type() {
	case lexer.TokenField, lexer.TokenFuncName, lexer.TokenRawParam, lexer.TokenVarParam:
		n.subNodes = append(n.subNodes, t)
	default:
		log.Debug(`not accept sub node`)
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
	log.Debug(n.Val)
	return n.Val
}

type paramNode struct {
	*lexer.Token
}

func (n *paramNode) Type() int {
	return n.Typ
}

func (n *paramNode) TranslateFrom(data interface{}, execFuncs interface{}) string {
	// if x := strings.Index(n.Val, `$`); x >= 0 {
	// 	//需要对变量进行翻译
	// 	fmt.Println(n.Val[x+1:])
	// 	val := data.(map[string]interface{})[n.Val[x+1:]]
	// 	fmt.Println(`val:`, val)
	// 	fmt.Printf("->type:%T,kind:%s,val:%v\n", val, reflect.ValueOf(val).Kind(), val)
	// 	return convertValue(val)
	// }
	return n.Val
}

func (n *paramNode) ValueFrom(data interface{}) reflect.Value {
	if x := strings.Index(n.Val, `$`); x >= 0 {
		val := data.(map[string]interface{})[n.Val[x+1:]]
		return reflect.ValueOf(val)
	}
	return reflect.ValueOf(n.Val)
}

type varNode struct {
	*lexer.Token
}

func (n *varNode) Type() int {
	return n.Typ
}

func (n *varNode) TranslateFrom(data interface{}, execFuncs interface{}) string {
	x := strings.Index(n.Val, `$`)
	log.Debug(n.Val[x+1:])
	val := data.(map[string]interface{})[n.Val[x+1:]]
	log.Debug("%T,val:%v\n", val, val)
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
		log.Debug("%T\n", data)
	}
	return fmt.Sprintf(`%v`, data)

}
