package tag

import (
	"strings"
)

// 词法分析器
type Lexer struct {
	Scanner
	token  *Token
	victim *Token
}

func (l *Lexer) Next() bool {
	if l.token.IsVaild() {
		return true
	}
	if l.victim.IsVaild() {
		l.token.Shift(l.victim)
		return true
	}
	l.read()
	return l.token.IsVaild()
}

func (l *Lexer) save(s string) {
	if l.Length() == 0 {
		l.token.Set(s)
		return
	}
	l.token.Var(l.Restore())
	l.victim.Set(s)
}

func (l *Lexer) read() {
	for l.Scanner.Next() {
		switch s := l.Scanner.Read(); s {
		case ",", ":", ";", "[", "]", "{", "}", "(", ")":
			l.save(s)
			return
		default:
			l.Store()
		}
	}
	if l.Length() != 0 {
		l.token.Var(l.Restore())
	}
}

func (l *Lexer) Read() (t Token) {
	t.Shift(l.token)
	return
}

func NewLexer(s string) Lexer {
	return Lexer{
		NewScanner(strings.NewReader(s)),
		new(Token),
		new(Token),
	}
}
