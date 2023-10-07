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
	Name string
}

func (f *File) Info(*Path) {
	fmt.Println(f.Name)
}

type Files []File

func (f *Files) Add(p *Path) {
	*f = append(*f, File{p.Full.Posts}, File{p.Full.Index})
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
	Full    *Path  `default:"Init;new;initial.Abs;Self;Parent"`

	Test struct {
		T1 string  `default:"t1"`
		T2 bool    `default:"true"`
		T3 float64 `default:"3.14"`
		T4 int64   `default:"114"`
	} `default:"initial.Default"`

	Files Files `default:"Add;range.Info"`
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

func NewPath(self *Path) {
	fmt.Printf("self: %v\n", self)
}

func init() {
	initial.Register("new", NewPath, "self")
}

func TestPath(t *testing.T) {
	result := initial.Default(&Path{})
	fmt.Printf("result: %v\n", result)
}
```

#### 控制台

```
self: &{ pages       <nil> { false 0 0} []}
p: &{resource resource\pages resource\public resource\public\posts.db resource\users.db resource\.log resource\pages\index.html resource\pages\.version <nil> { false 0 0} []}
parent: &{resource views public posts.db users.db .log index.html .version 0xc00001e8f0 { false 0 0} []}
resource\public\posts.db
resource\pages\index.html
result: &{resource views public posts.db users.db .log index.html .version 0xc00001e8f0 {t1 true 3.14 114} [{resource\public\posts.db} {resource\pages\index.html}]}
```