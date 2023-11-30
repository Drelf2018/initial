package tag

const (
	COMMA      = iota + 4 // ,
	COLON                 // :
	LBRACKET              // [
	RBRACKET              // ]
	LBRACE                // {
	RBRACE                // }
	LGROUP                // (
	RGROUP                // )
	IDENTIFIER            // identifier
)

func getKind(value string) int {
	switch value {
	case ",", ";":
		return COMMA
	case ":":
		return COLON
	case "[":
		return LBRACKET
	case "]":
		return RBRACKET
	case "{":
		return LBRACE
	case "}":
		return RBRACE
	case "(":
		return LGROUP
	case ")":
		return RGROUP
	default:
		return IDENTIFIER
	}
}

// 最小词语单元
type Token struct {
	Kind  int    `json:"kind" yaml:"kind"`
	Value string `json:"value" yaml:"value"`
}

// 判断合法性
func (t *Token) IsVaild() bool {
	return t.Kind > 0
}

// 判断是否包含下级
func (t *Token) IsBegin() bool {
	switch t.Kind {
	case LBRACKET, LBRACE, LGROUP:
		return true
	default:
		return false
	}
}

// 判断是否结束包含
func (t *Token) IsEnd() bool {
	switch t.Kind {
	case RBRACKET, RBRACE, RGROUP:
		return true
	default:
		return false
	}
}

// 新建
func (t *Token) New(kind int, value string) {
	t.Kind = kind
	t.Value = value
}

// 自动推断类型
func (t *Token) Set(value string) {
	t.New(getKind(value), value)
}

// 设置标识符
func (t *Token) Var(value string) {
	t.New(IDENTIFIER, value)
}

// 切换
func (t *Token) Shift(n *Token) {
	t.Kind, t.Value, n.Kind = n.Kind, n.Value, -1
}

// 新建列表单元
func NewToken() (t Token) {
	t.Set("[")
	return
}
