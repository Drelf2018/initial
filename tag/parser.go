package tag

// 语法分析器
type Parser struct {
	// 词法分析器
	Lexer
	// 根 Sentence
	Sentence *Sentence
}

func (p *Parser) build() {
	sentence := p.Sentence
	for p.Lexer.Next() {
		token := p.Lexer.Read()
		if token.Kind == COMMA {
			continue
		}
		if token.Kind == COLON {
			sentence = sentence.Parent.Append(NewToken())
			continue
		}
		if token.IsEnd() {
			sentence = sentence.Parent
			if token.Kind == RBRACE {
				sentence = sentence.Parent
			}
			continue
		}
		if token.Kind == LGROUP {
			sentence = sentence.Last()
			continue
		}
		s := sentence.Append(token)
		if token.Kind == LBRACKET {
			sentence = s
			continue
		}
		if token.Kind == LBRACE {
			sentence = s.Append(NewToken())
		}
	}
}

func NewParser(s string) (p Parser) {
	p = Parser{
		Lexer:    NewLexer(s),
		Sentence: NewRoot(),
	}
	p.build()
	return
}
