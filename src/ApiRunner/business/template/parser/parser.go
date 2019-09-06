package parser

import (
	// "errors"
	"fmt"

	// "strings"

	"ApiRunner/business/template/lexer"
)

/*
${email}  //var
${gen_email()}  //function
${gen_email(4)}  //function with one args
${gen_email(4,12)}  //function with args
${gen_email($min,$max)}  //function with args
${gen_email(4,$max)}  //function with mixed
${refs.user1.email}  //function with args
has ${getnum()} items,${num2} records
has num items,num2 .{[(}records
 	//null
*/
type Tree struct {
	lex      *lexer.Lexer
	fields   [][]string //二维数组来存放所有的refs
	vars     []string
	funcs    []map[string]interface{} //interface as params
	preToken *lexer.Token
	curToken *lexer.Token
	nodeList []Node
}

func (t *Tree) GetNodeList() []Node {
	return t.nodeList
}

func (t *Tree) init(input string) {
	if t.lex == nil {
		t.lex = lexer.NewLexer(input)
		t.fields = append(t.fields, []string{})
		t.funcs = append(t.funcs, map[string]interface{}{})
		// t.nodeList = []Node{}
	}
}

func (t *Tree) Parse(input string) {
	t.init(input)
	for state := startParse; state != nil; {
		state = state(t)
	}
}

func (t *Tree) getToken() *lexer.Token {
	token := t.lex.NextToken()
	t.preToken = t.curToken
	t.curToken = &token
	fmt.Println(`got token:`, token.Typ, token.Val)
	return &token
}

func (t *Tree) addNode(n Node) {
	switch n.Type() {
	case lexer.TokenField:
		nodeIndex := len(t.nodeList) - 1
		if nodeIndex < 0 || t.nodeList[nodeIndex].Type() != n.Type() {
			t.nodeList = append(t.nodeList, n)
		} else {
			// 非首个field node
			nodeIndex--
			for nodeIndex >= 0 {
				preNodeTyp := t.nodeList[nodeIndex].Type()
				if preNodeTyp != n.Type() {
					//当前nodeIndex + 1 就是祖先了
					ancestor := t.nodeList[nodeIndex+1].(*fieldNode)
					ancestor.expand(n)
					break
				}
				nodeIndex--
			}
		}

	case lexer.TokenFuncName:
		t.nodeList = append(t.nodeList, n)
	case lexer.TokenRawParam, lexer.TokenVarParam:
		nodeIndex := len(t.nodeList) - 1
		for nodeIndex >= 0 {
			preNodeTyp := t.nodeList[nodeIndex].Type()
			if preNodeTyp == lexer.TokenFuncName {
				obj := t.nodeList[nodeIndex].(*funcNode)
				obj.expand(n)
				break
			}
			nodeIndex--
		}
	default:
		t.nodeList = append(t.nodeList, n)
	}
}

func (t *Tree) ignore() {
	t.getToken()
}

type parseFn func(*Tree) parseFn

func startParse(t *Tree) parseFn {
	_token := t.getToken()
	switch _token.Typ {
	case lexer.TokenError:
		return parseError
	case lexer.TokenEOF:
		return parseEOF
	case lexer.TokenText:
		textNodeObj := &textNode{_token}
		t.addNode(textNodeObj)
		return startParse
	default:
		//must be 3 of 1,when begin
		return parseLeftDelim
	}
}

func parseLeftDelim(t *Tree) parseFn {
	t.ignore()
	return parseToken
}

func parseToken(t *Tree) parseFn {
	switch t.curToken.Typ {
	case lexer.TokenField:
		return parseField
	case lexer.TokenVariable:
		return parseVariable
	case lexer.TokenFuncName:
		return parseFuncName
	case lexer.TokenRawParam, lexer.TokenVarParam:
		return parseParams
	case lexer.TokenRightDelim:
		return parseRightDelim
	default:
		fmt.Println(`ignore token `, t.curToken)
		t.ignore()
		return parseToken
	}
}

func parseField(t *Tree) parseFn {
	index := len(t.fields) - 1
	fmt.Println(`parseField index:`, index)
	fieldNodeObj := &fieldNode{t.curToken, []Node{}}
	t.addNode(fieldNodeObj)
	t.fields[index] = append(t.fields[index], t.curToken.Val)
	t.getToken()
	return parseToken
}

func parseVariable(t *Tree) parseFn {
	t.vars = append(t.vars, t.curToken.Val)
	t.getToken()
	return parseToken
}

func parseFuncName(t *Tree) parseFn {
	// funcNode[_token.Val] = []interface{}{}
	index := len(t.funcs) - 1
	m := t.funcs[index]
	m[t.curToken.Val] = []interface{}{}
	t.funcs[index] = m
	funcNodeObj := &funcNode{t.curToken, []Node{}}
	t.addNode(funcNodeObj)
	t.getToken()
	return parseToken
}

func parseParams(t *Tree) parseFn {
	index := len(t.funcs) - 1
	m := t.funcs[index]
	for k, v := range m {
		v := v.([]interface{})
		v = append(v, t.curToken.Val)
		m[k] = v
		break
	}
	t.funcs[index] = m
	funcNodeObj := &funcNode{t.curToken, []Node{}}
	t.addNode(funcNodeObj)
	t.getToken()
	return parseToken
}

func parseRightDelim(t *Tree) parseFn {
	//一次循环结束
	switch t.preToken.Typ {
	case lexer.TokenField:
		t.fields = append(t.fields, []string{}) //插入一个新的，下一轮循环使用
	case lexer.TokenRightParen:
		//带参数的函数调用
		t.funcs = append(t.funcs, map[string]interface{}{})
	default:
		fmt.Println(`not handle token `, t.preToken)
	}
	return startParse
}

func parseError(t *Tree) parseFn {
	return parseEOF
}

func parseEOF(t *Tree) parseFn {
	fmt.Print(t.curToken.Val)
	return nil
}
