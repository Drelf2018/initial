package tag

import (
	"io"
)

type Reader interface {
	ReadRune() (rune, int, error)
}

type Scanner struct {
	Reader
	// 当前字符
	current rune
	// 读取时错误
	err error
	// 暂存的字符
	storage []rune
}

// 获取下位字符
func (s *Scanner) Next() bool {
	s.current, _, s.err = s.Reader.ReadRune()
	if s.err != nil && s.err != io.EOF {
		panic(s.err)
	}
	return s.err != io.EOF
}

// 读取 string
func (s *Scanner) Read() string {
	return string(s.current)
}

// 获取暂存字符长度
func (s *Scanner) Length() int {
	return len(s.storage)
}

// 暂存当前字符
func (s *Scanner) Store() {
	s.storage = append(s.storage, s.current)
}

// 清空暂存并以 string 返回
func (s *Scanner) Restore() (r string) {
	r = string(s.storage)
	s.storage = make([]rune, 0, 16)
	return
}

func NewScanner(r Reader) Scanner {
	return Scanner{
		Reader:  r,
		current: 0,
		err:     nil,
		storage: make([]rune, 0, 16),
	}
}
