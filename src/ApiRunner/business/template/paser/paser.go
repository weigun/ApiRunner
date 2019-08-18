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
	Vars   map[string]interface{}
	funcs  map[string]interface{} //interface as params

}

func Parse(input string) (string, error) {
	l := lexer.NewLexer(`test`, input)
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
		}

	}
}
