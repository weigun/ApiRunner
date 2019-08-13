package lexer

import (
	"strings"
	"unicode"
)

/*
${email}  //var
${gen_email()}  //function
${gen_email(4,12)}  //function with args
${gen_email($min,$max)}  //function with args
${gen_email(4,$max)}  //function with mixed
${refs.user1.email}  //function with args
has ${num} items,${num2} records
has num items,num2 .{[(}records
*/
func LexBegin(l *Lexer) stateFn {
	l.SkipSpace()

	if strings.HasPrefix(l.InputToEnd(), LEFT_DLIM) {
		return LexLeftDelim
	} else {
		return LexText
	}
}

func LexText(l *Lexer) stateFn {
	if x := strings.Index(l.InputToEnd(), LEFT_DLIM); x >= 0 {
		//如果存在起止符，则将pos设置到起止符处,并且跳转到LexLeftDelim
		l.Pos += Pos(x)
		l.Ignore()
		return LexLeftDelim
	}
	for {
		if strings.HasPrefix(l.InputToEnd(), RIGHT_DLIM) {
			//should var
			l.Emit(TokenVariable)
			return LexRightDelim
		} else if strings.HasPrefix(l.InputToEnd(), LEFT_PAREN) {
			//function
			l.Emit(TokenFuncName)
			return LexFuncName
		} else if strings.HasPrefix(l.InputToEnd(), DOLLAR) {
			//refs
			l.Emit(TokenField)
			return LexField
		}
		l.Inc()
	}
	/*
		else if x := strings.Index(l.InputToEnd(), RIGHT_PAREN); x >= 0 {
			/*
				如果是前文是(，则应该找到)，没找到应该跳到eof
				如果找到，则开始找函数参数了，有4种情况：
				1.无参数
				2.只有明文参数
				3.只有变量参数
				4.混合参数
			* /
			if strings.HasPrefix(l.InputToEnd(), RIGHT_PAREN) {
				//无参数的情况
				return LexRightParen
			} else if strings.HasPrefix(l.InputToEnd(), DOLLAR) {
				//变量参数
				l.Ignore() //忽略$
				return LexDollar
			} else {
				//明文参数
				return LexRawParam
			}
		}*/
	//如果没有找到起止符，则应该要结束了，没必要进行下去
	l.Pos += Pos(len(l.Input))
	l.Emit(TokenEOF)
	return nil
}

func LexLeftDelim(l *Lexer) stateFn {
	l.Pos += Pos(len(LEFT_DLIM))
	l.Emit(TokenLeftDelim)
	subInput := l.InputToEnd()
	return LexText
	/*
		// l.Ignore()
		//TODO 这里可能有问题，如果input是has ${num} items,${get_records()} records
		if strings.Index(subInput, LEFT_PAREN) != -1 {
			//如果含有(，那么只能是纯函数调用(并且只有一个函数)，不能混合其他，函数参数除外
			return LexFuncName
		} else if strings.Index(subInput, DOT) != -1 {
			//如果含有.，那么只能是引用树的调用了，函数参数不支持引用树
			return TokenField
		} else {
			//只剩下变量了
			return LexVariable
		}
	*/
}

func LexVariable(l *Lexer) stateFn {
	// ${email}  //var
	if l.IsEOF() {
		//reached eof
		l.Pos += Pos(len(l.Input))
		l.Emit(TokenEOF)
		return nil
	}
	//变量名只能由字母数字以及下划线组成,且必须是字母开头
	//TODO maybe not loop
	for {
		//找到}，就可以确认变量名了
		if strings.HasPrefix(l.InputToEnd(), RIGHT_DLIM) {
			varName := l.CurrebInput()
			if !isVarNameVerified(varName) {
				//违反了命名规则
				return l.Errorf(`Variables can only consist of alphanumeric and underscores and must start with a letter`)
			}
			l.Emit(TokenVariable)
			return LexRightDelim
		}
		l.Inc()
	}
}

func LexRightDelim(l *Lexer) stateFn {
	l.Pos += Pos(len(RIGHT_DLIM))
	l.Emit(TokenRightDelim)
	return LexBegin //又一轮循环
}

func LexFuncName(l *Lexer) stateFn {
	// ${gen_email()}  //function
	if l.IsEOF() {
		//reached eof
		l.Pos += Pos(len(l.Input))
		l.Emit(TokenEOF)
		return nil
	}
	//函数名只能由字母数字以及下划线组成,且必须是字母开头
	//函数名与(之间不能存在空格
	//TODO maybe not loop
	for {
		//找到(，就可以确认函数名了
		if strings.HasPrefix(l.InputToEnd(), LEFT_PAREN) {
			fnName := l.CurrebInput()
			if !isVarNameVerified(fnName) {
				//违反了命名规则
				return l.Errorf(`function name can only consist of alphanumeric and underscores and must start with a letter`)
			}
			l.Emit(TokenFuncName)
			return LexLeftParen
		}
		l.Inc()
	}
}

func LexLeftParen(l *Lexer) stateFn {
	l.Pos += Pos(len(LEFT_PAREN))
	l.Emit(TokenLeftParen)
	return LexText
}

func LexRightParen(l *Lexer) stateFn {
	l.Pos += Pos(len((RIGHT_PAREN)))
	l.Emit(TokenRightParen) //emit wrong
}

func LexRawParam(l *Lexer) stateFn {
	/*
		明文参数
		${gen_email(4,12)}  //function with args
		如果是数字，直接解析为数字,否则直接为字符串
		需要处理多个参数的情况
	*/
	//是否多参数
	rpn := strings.Index(l.InputToEnd(), RIGHT_PAREN)
	counter := strings.Count(l.Input[l.Pos:l.Pos+rpn], COMMA)
	switch counter {
	case 0:
		l.Pos += Pos(rpn)
		l.Emit(TokenRawParam)
		return LexText
	case 1:
		for {
			if strings.HasPrefix(l.InputToEnd(), COMMA) {
				//找到一个参数
				// edit here
				// l.Ignore()
				l.Emit(TokenRawParam) //ignore?
				return LexText
			}
			l.Inc()
		}
	default:
		//2参数以上
		for {
			if strings.HasPrefix(l.InputToEnd(), COMMA) {
				//找到一个参数
				// edit here
				l.Emit(TokenRawParam) //ignore?
				return LexText
			}
			l.Inc()
		}
	}
}

// ---------------------------
func isVarNameVerified(fnName string) bool {
	runeName := []rune(fnName)
	for i, ascii := range runeName {
		if i == 0 {
			if ascii == '_' || unicode.IsLetter(ascii) {
				continue
			}
			return false
		}
		if !isAlphaNumeric(ascii) {
			return false
		}
	}
	return true
}

func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}
