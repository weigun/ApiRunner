package parser

import (
	"errors"
	"fmt"

	// "strings"

	"ApiRunner/business/template/lexer"
)

func isEOF(token lexer.Token) bool {
	return token.Typ == lexer.TokenEOF
}

type Bucket struct {
	Fields [][]string //二维数组来存放所有的refs
	Vars   []string
	Funcs  []map[string]interface{} //interface as params
}

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
func Parse(input string) (*Bucket, error) {
	//TODO 优化：递归分散成函数?
	l := lexer.NewLexer(`test`, input)
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
