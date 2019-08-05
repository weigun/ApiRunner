package lexer

import (
	"strings"
)

/*
${email}  //var
${gen_email()}  //function
${gen_email(4,12)}  //function with args
${gen_email($min,$max)}  //function with args
${refs.user1.email}  //function with args
has ${num} items
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
	//如果没有找到起止符，则应该要结束了，没必要进行下去
	l.Pos += Pos(len(l.Input))
	l.Emit(TokenEOF)
	return nil
}

func LexLeftDelim(l *Lexer) stateFn {
	l.Pos += Pos(len(LEFT_DLIM))
	l.Emit(TokenLeftDelim)
	subInput := l.InputToEnd()
	l.Ignore()
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
			if !isVarNameVerified(&varName) {
				//违反了命名规则
				return l.Errorf(`Variables can only consist of alphanumeric and underscores and must start with a letter`)
			}
			l.Emit(TokenVariable)
			return LexRightDelim
		}
		l.Inc()
	}
}
