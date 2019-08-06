package lexer

/*
不支持多行
${email}  //var
${gen_email()}  //function
${gen_email(4,12)}  //function with args
${gen_email($min,$max)}  //function with args
${gen_email(4,$max)}  //function with mixed
${refs.user1.email}  //function with args
has ${num} items,${num2} records
*/
type Token struct {
	Typ TokenType
	Pos Pos
	Val string
}

type TokenType int
type Pos int

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
	TokenDollar
	TokenText //just text,no var and func. e.g,has ${num} items-has,items
)
