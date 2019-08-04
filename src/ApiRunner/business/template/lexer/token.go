package lexer

import (
	"fmt"
	"strings"
)

/*
${email}  //var
${gen_email()}  //function
${gen_email(4,12)}  //function with args
${gen_email($min,$max)}  //function with args
${refs.user1.email}  //function with args





*/
type Token struct {
}

type TokenType int

const (
	TokenError TokenType = iota
	TokenEOF
	TokenLeftDelim
	TokenVariable
	TokenRightDelim
	TokenFuncName
	TokenLeftParen
	TokenRawParam
	TokenVarParam
	TokenComma
	TokenRightParen
	TokenField
	TokenDot
)
