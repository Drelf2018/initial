package initial_test

import (
	"fmt"
	"testing"

	"github.com/Drelf2018/initial"
)

type Path struct {
	Root    string `default:"resource"`
	Views   string `default:"views" abs:"Root"`
	Public  string `default:"public" abs:"Root"`
	Posts   string `default:"posts.db" abs:"Public"`
	Users   string `default:"users.db" abs:"Root"`
	Log     string `default:".log" abs:"Root"`
	Index   string `default:"index.html" abs:"Views"`
	Version string `default:".version" abs:"Views"`
	Full    *Path  `default:"Init;Abs;Self;Parent"`

	Test struct {
		T1 string  `default:"t1"`
		T2 bool    `default:"true"`
		T3 float64 `default:"3.14"`
		T4 int64   `default:"114"`
	} `default:"Default"`
}

func (p *Path) Init(_ any) {
	p.Views = "pages"
}

func (p *Path) Self(_ *Path) {
	fmt.Printf("p: %v\n", p)
}

func (*Path) Parent(parent *Path) {
	fmt.Printf("parent: %v\n", parent)
}

func TestPath(t *testing.T) {
	result := initial.Default(&Path{})
	fmt.Printf("result: %v\n", result)
}