package paser

import (
	"fmt"
	"strings"

	"ApiRunner/business/template/lexer"
)

func isEOF(token lexer.Token) bool {
	return token.Typ == lexer.TokenEOF
}

type Bucket struct {
	Fields [][]string //二维数组来存放所有的refs
	Vars   []string
	funcs  []map[string]interface{} //interface as params

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
func Parse(input string) (string, error) {
	l := lexer.NewLexer(`test`, input)
	bucketPtr := &Bucket{}
	fieldNode := []string{} //refs.user1.email分别存放refs user email
	funcNode := make(map[string]interface{})
	var preTokenType lexer.TokenType
	for {
		_token := l.NextToken()
		switch _token.Typ {
		case lexer.TokenError:
			//error
			return input, _token.Val
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
			default:
				fmt.Println(`not handle token `, _token)
			}

		case lexer.TokenVariable:
			//vars
			bucketPtr.Vars = append(bucketPtr.Vars, _token.Val)
		case lexer.TokenFuncName:
			//function
			funcNode[_token.Val] = []string{}
		case lexer.TokenRawParam:
			//function params
			if len(funcNode) != 1 {
				panic(`more than one func in one token`)
			}
			for k, v := range funcNode {
				v = append(v, _token.Val)
				funcNode[k] = v
				break
			}
			//here
		}
		preTokenType = _token.Typ

	}
}
