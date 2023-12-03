# initial

依赖注入初始化

### 使用

```golang
package initial_test

import (
	"fmt"
	"testing"

	"github.com/Drelf2018/initial"
)

type File struct {
	Name string `default:"initial.go"`
}

func (f *File) BeforeInitial(parent any) {
	fmt.Printf("BeforeInitial: %v\n", f)
}

func (f *File) AfterInitial(parent any) {
	fmt.Printf("AfterInitial: %v\n", f)
}

func (f *File) Info(*Path) {
	fmt.Println(f.Name)
}

type Files []File

func (f *Files) Add(p *Path) {
	*f = append(*f, File{p.Full.Posts}, File{}, File{p.Full.Index})
	p.FileMap = []map[*File]*File{{{"twice"}: {"once"}}}
}

type Path struct {
	Root    string `default:"resource"`
	Views   string `default:"views" abs:"Root"`
	Public  string `default:"public" abs:"Root"`
	Posts   string `default:"posts.db" abs:"Public"`
	Users   string `default:"users.db" abs:"Root"`
	Log     string `default:".log" abs:"Root"`
	Index   string `default:"index.html" abs:"Views"`
	Version string `default:".version" abs:"Views"`
	Full    *Path  `default:"-,Init,new,initial.Abs,Self,Parent"`

	Test struct {
		T1  string  `default:"t1"`
		T2  bool    `default:"true"`
		T3  float64 `default:"3.14"`
		T4  int64   `default:"114"`
		New *Path   `default:"-"`
	}

	Null string `default:""`

	Files Files `default:"Add,[Info]"`

	FileMap []map[*File]*File `default:"[{Info,Info:Info}]"`
}

func (p *Path) Init(_ any) {
	p.Views = "pages"
}

func (p *Path) Self(_ *Path) {
	fmt.Printf("self: %v\n", p)

}

func (*Path) Parent(parent *Path) error {
	fmt.Printf("parent: %v\n", parent)
	return initial.ErrBreak
}

func NewPath(self *Path) {
	fmt.Printf("new: %v\n", self)
}

func init() {
	initial.Register("new", NewPath, initial.SELF)
}

func (p *Path) BeforeDefault() {
	fmt.Printf("BeforeDefault\n")
}

func (p *Path) AfterDefault() {
	fmt.Printf("AfterDefault\n")
}

func TestPath(t *testing.T) {
	result := initial.Default(&Path{})
	fmt.Printf("result: %v\n", result)

	if result.Full.Version != `resource\pages\.version` {
		t.Fail()
	}
}
```

#### 控制台

```
BeforeDefault
new: &{ pages       <nil> { false 0 0 <nil>}  [] []}
self: &{resource resource\pages resource\public resource\public\posts.db resource\users.db resource\.log resource\pages\index.html resource\pages\.version <nil> { false 0 0 <nil>}  [] []}
parent: &{resource views public posts.db users.db .log index.html .version 0xc000216100 { false 0 0 <nil>}  [] []}
BeforeInitial: &{resource\public\posts.db}
resource\public\posts.db
AfterInitial: &{resource\public\posts.db}
BeforeInitial: &{}
initial.go
AfterInitial: &{initial.go}
BeforeInitial: &{resource\pages\index.html}
resource\pages\index.html
AfterInitial: &{resource\pages\index.html}
BeforeInitial: &{twice}
twice
twice
AfterInitial: &{twice}
BeforeInitial: &{once}
once
AfterInitial: &{once}
AfterDefault
result: &{resource views public posts.db users.db .log index.html .version 0xc000216100 {t1 true 3.14 114 <nil>}  [{resource\public\posts.db} {initial.go} {resource\pages\index.html}] [map[0xc00020e660:0xc00020e650]]}
```