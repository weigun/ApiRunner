package lexer

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

const (
	EOF                = -1
	LEFT_DLIM   string = "${"
	RIGHT_DLIM  string = "}"
	LEFT_PAREN  string = `(`
	RIGHT_PAREN string = `)`
	COMMA       string = `,`
	DOT         string = `.`
	DOLLAR      string = `$`
	NEWLINE     string = "\n"
)

type stateFn func(*Lexer) stateFn

type Lexer struct {
	Name   string
	Input  string
	Tokens chan Token
	Start  Pos //start position of this token
	Pos    Pos // current position in the input
	Width  Pos //width of last rune read from input
}

func (l *Lexer) Next() rune {
	/*
	   Reads the next rune (character) from the input stream
	   and advances the lexer position.
	*/
	if l.Pos > utf8.RuneCountInString(l.Input) {
		//reach eof
		l.Width = 0
		return EOF
	}
	//get remain char from cur position
	r, w := utf8.DecodeRuneInString(l.Input[l.Pos:])
	l.Width = Pos(w)
	l.Pos += l.Width
	return r
}

func (l *Lexer) Backup() {
	l.Pos -= l.Width
}

func (l *Lexer) CurrebInput() string {
	return l.Input[l.Start:l.Pos]
}

func (l *Lexer) Dec() {
	l.Pos--
}

func (l *Lexer) Inc() {
	l.Pos++
	if l.Pos >= utf8.RuneCountInString(l.Input) {
		l.Emit(TokenEOF)
	}
}

func (l *Lexer) Emit(tokenTyp TokenType) {
	l.Tokens <- Token{tokenTyp, l.Start, l.Input[l.Start:l.Pos]}
	l.Start = l.Pos
}

func (l *Lexer) Errorf(format string, args ...interface{}) stateFn {
	l.Tokens <- Token{TokenError, l.Start, fmt.Sprintf(format, args...)}
	return nil
}

func (l *Lexer) Ignore() {
	l.Start = l.Pos
}

func (l *Lexer) InputToEnd() string {
	return l.Input[l.Pos:]
}

func (l *Lexer) IsEOF() bool {
	return l.Pos >= len(l.Input)
}

func (l *Lexer) IsSpace() bool {
	r, _ := utf8.DecodeRuneInString(l.Input[l.Pos:])
	return unicode.IsSpace(r)
}

func (l *Lexer) NextToken() Token {
	return <-l.Tokens
}

func (l *Lexer) Peek() rune {
	/*
		Returns the next rune in the stream, then puts the lexer
		position back. Basically reads the next rune without consuming
		it.
	*/
	r := l.Next()
	l.Backup()
	return r
}

func (l *Lexer) Run() {
	for state := LexBegin; state != nil; {
		state = state(l)
	}
	close(l.Tokens)
}

func (l *Lexer) SkipSpace() {
	for {
		r := l.Next()
		if !unicode.IsSpace(r) {
			l.Dec()
			break
		}
		if int(r) == TokenEOF {
			l.Emit(TokenEOF)
			break
		}
	}
}
