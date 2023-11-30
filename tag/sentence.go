package tag

type Sentence struct {
	Token  `json:"token" yaml:"token"`
	Parent *Sentence   `json:"-" yaml:"-"`
	Body   []*Sentence `json:"body" yaml:"body"`
}

func (s *Sentence) Last() *Sentence {
	return s.Body[len(s.Body)-1]
}

func (s *Sentence) Append(t Token) (r *Sentence) {
	r = NewSentence(t, s)
	s.Body = append(s.Body, r)
	return
}

func NewSentence(token Token, parent *Sentence) *Sentence {
	return &Sentence{
		Token:  token,
		Parent: parent,
		Body:   make([]*Sentence, 0),
	}
}

func NewRoot() *Sentence {
	return NewSentence(NewToken(), nil)
}
