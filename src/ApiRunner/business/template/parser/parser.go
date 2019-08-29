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
}

func (t *Tree) init(input string) {
	if t.lex == nil {
		t.lex = lexer.NewLexer(input)
		t.fields = append(t.fields, []string{})
		t.funcs = append(t.funcs, map[string]interface{}{})
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

/*
func Parse(input string) (*Bucket, error) {
	//TODO 优化：递归分散成函数?
	//这里是解析的入口函数
	l := lexer.NewLexer(input)
	bucketPtr := &Bucket{}
	fieldNode := []string{} //refs.user1.email分别存放refs user email
	funcNode := map[string]interface{}{}
	var preTokenType lexer.TokenType
	for {
		_token := l.NextToken()
		switch _token.Typ {
		case lexer.TokenError:
			//error
			return bucketPtr, errors.New(_token.Val)
		case lexer.TokenEOF:
			//reach inputend
			// TODO
		case lexer.TokenField:
			//refs
			fieldNode = append(fieldNode, _token.Val)
		case lexer.TokenRightDelim:
			//一个模板结束
			//根据前一个token来决定什么逻辑
			switch preTokenType {
			case lexer.TokenField:
				//应该是refs
				fields := bucketPtr.Fields
				fields = append(fields, fieldNode)
				bucketPtr.Fields = fields
				fieldNode = []string{} //reset
			case lexer.TokenRightParen: //lexer.TokenRawParam, lexer.TokenVarParam:
				//应该是带参数的函数调用
				fmt.Println(`should enter`)
				funcs := bucketPtr.Funcs
				funcs = append(funcs, funcNode) //funcNode=[print]{1,2,3}
				bucketPtr.Funcs = funcs
				funcNode = map[string]interface{}{}

			default:
				fmt.Println(`not handle token `, _token)
			}

		case lexer.TokenVariable:
			//vars
			bucketPtr.Vars = append(bucketPtr.Vars, _token.Val)
		case lexer.TokenFuncName:
			//function
			fmt.Println(`TokenFuncName `, _token.Val)
			funcNode[_token.Val] = []interface{}{}
		case lexer.TokenRawParam, lexer.TokenVarParam:
			//function params
			fmt.Println(`params `, _token.Val)
			if len(funcNode) != 1 {
				panic(`more than one func in one token`)
			}
			for k, v := range funcNode {
				v := v.([]interface{})
				v = append(v, _token.Val)
				funcNode[k] = v
				break
			}
			fmt.Println(funcNode)
		default:
			fmt.Println(`ignore token `, _token)
		}
		preTokenType = _token.Typ
	}
	return bucketPtr, nil
}
*/
