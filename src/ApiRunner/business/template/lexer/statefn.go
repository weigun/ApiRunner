package lexer

import (
	// "fmt"
	// "fmt"
	"strings"
	"unicode"
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
func LexBegin(l *Lexer) stateFn {
	// l.SkipSpace()
	l.InAction = ACTION_VAR
	for {
		if l.IsEOF() {
			//reached eof
			// l.Pos += Pos(len(l.Input))
			l.Emit(TokenEOF)
			return nil
		}
		if x := strings.Index(l.InputToEnd(), LEFT_DLIM); x >= 0 {
			l.Pos += Pos(x)
			l.Emit(TokenText)
			return LexLeftDelim
		}
		l.Pos = Pos(len(l.Input))
		l.Emit(TokenText)
		return LexBegin
		// if strings.HasPrefix(l.InputToEnd(), LEFT_DLIM) {
		// 	l.Emit(TokenText)
		// 	return LexLeftDelim
		// }
		// l.Inc()
	}
}

func LexText(l *Lexer) stateFn {
	for {
		if strings.HasPrefix(l.InputToEnd(), RIGHT_DLIM) {
			//TODO 应该要加上当前模式的判定
			switch l.InAction {
			case ACTION_REFS:
				l.Emit(TokenField)
			case ACTION_FUNC, ACTION_RAW_TEXT:
				// 无需处理
			default:
				l.Emit(TokenVariable)
			}
			return LexRightDelim

		} else if strings.HasPrefix(l.InputToEnd(), LEFT_PAREN) {
			//function
			l.InAction = ACTION_FUNC
			l.Emit(TokenFuncName)
			return LexLeftParen
		} else if strings.HasPrefix(l.InputToEnd(), DOT) {
			//refs
			//先检查字段名是否合法
			fieldName := l.CurrebInput()
			if !isVarNameVerified(fieldName) {
				//违反了命名规则
				return l.Errorf(`field can only consist of alphanumeric and underscores and must start with a letter`)
			}
			l.InAction = ACTION_REFS
			l.Emit(TokenField)
			return LexDot
		} else if strings.HasPrefix(l.InputToEnd(), RIGHT_PAREN) {
			// param end
			// 需要判断有无参数
			tmpPos := l.Pos
			paramsStr := l.CurrebFnArgs() //(),(4),(4,12) etc
			log.Debug(`paramsStr:`, paramsStr)
			if paramsStr == `` {
				//无参数
				return LexRightParen
			}
			params := strings.Split(paramsStr, COMMA) //分割参数
			if len(params) == 1 {
				//1个参数
				//判断是明文参数还是变量参数
				if strings.Index(params[0], DOLLAR) == -1 {
					//明文参数
					// l.Pos += len(params[0]) + 1
					l.Emit(TokenRawParam)
				} else {
					//变量参数
					l.Emit(TokenVarParam)
				}
				return LexRightParen
			} else {
				//多个参数
				for _, v := range params {
					start := l.Start
					for {
						log.Debug(l.Input[start:])
						if strings.HasPrefix(l.Input[start:], v) {
							l.Pos = start + len(v)
							l.Start = start
							break
						}
						start++
					}
					// l.Pos = l.Start + len(v)
					log.Debug(`888888888888888`, l.CurrebInput(), v)
					//判断是明文参数还是变量参数
					if strings.Index(v, DOLLAR) == -1 {
						//明文参数
						l.Emit(TokenRawParam)
					} else {
						//变量参数
						l.Emit(TokenVarParam)
					}
				}
				l.Pos = tmpPos //还原位置
				return LexRightParen
			}
		}
		l.Inc()
	}
	//如果没有找到起止符，则应该要结束了，没必要进行下去
	l.Pos += Pos(len(l.Input))
	l.Emit(TokenEOF)
	return nil
}

func LexLeftDelim(l *Lexer) stateFn {
	l.Pos += Pos(len(LEFT_DLIM))
	l.Emit(TokenLeftDelim)
	return LexText
}

func LexVariable(l *Lexer) stateFn {
	// ${email}  //var
	if l.IsEOF() {
		//reached eof
		// l.Pos += Pos(len(l.Input))
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
		// l.Pos += Pos(len(l.Input))
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
	l.FnStart = l.Pos
	l.Pos += Pos(len(LEFT_PAREN))
	l.Emit(TokenLeftParen)
	return LexText
}

func LexRightParen(l *Lexer) stateFn {
	l.Pos += Pos(len((RIGHT_PAREN)))
	l.Emit(TokenRightParen) //emit wrong
	return LexRightDelim
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

func LexDot(l *Lexer) stateFn {
	l.Pos += Pos(len((DOT)))
	l.Emit(TokenDot)
	return LexText
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
