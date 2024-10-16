# initial

依赖注入初始化

### 使用

```go
package initial_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/Drelf2018/initial"
)

type File struct {
	Name string `default:"initial.go"`
}

func (f *File) BeforeInitial() {
	println("FileBeforeInitial:", f.Name)
}

func (f File) AfterInitial() error {
	if f.Name == "error.log" {
		return errors.New("test error")
	}
	println("FileAfterInitial: ", f.Name)
	return nil
}

type Files []File

type Path struct {
	ID uint16 `default:"9000"`

	Root    string `default:"resource"`
	Views   string `default:"views"      join:"Root"`
	Public  string `default:"public"     join:"Root"`
	Posts   string `default:"posts.db"   join:"Public"`
	Users   string `default:"users.db"   join:"Root"`
	Log     string `default:".log"       join:"Root"`
	Index   string `default:"index.html" join:"Views"`
	Version string `default:".version"   join:"Views"`
	Full    *Path

	Data struct {
		D1 string  `default:"d1"`
		D2 bool    `default:"true"`
		D3 float64 `default:"3.14"`
		D4 int64   `default:"114"`
	}

	Files   Files
	FileMap []map[*File]*File

	// Error error `default:"some error"` // can't work
	Error error `default:"$myError"` // works
}

func init() {
	initial.SetDefaultValue("$myError", errors.New("some error"))
}

func (p *Path) BeforeFiles() {
	p.Files = append(p.Files, File{}, File{p.Full.Posts})
	p.FileMap = []map[*File]*File{{{"key.txt"}: {"value.txt"}}}
}

func (p *Path) BeforeFull() {
	p.Full.Posts = p.Posts
}

func (p *Path) BeforeInitial() {
	p.Files = append(p.Files, File{"default.go"})
}
```

#### 测试解析结果

```go
func TestParse(t *testing.T) {
	values := initial.ParseValues(reflect.TypeOf(Path{}))
	for _, val := range values {
		t.Log(val)
	}
}
```

```
initial_test.go:76: {0 true <uint16 Value> <invalid Value> <invalid Value>}
initial_test.go:76: {1 true resource <invalid Value> <invalid Value>}
initial_test.go:76: {2 true views <invalid Value> <invalid Value>}
initial_test.go:76: {3 true public <invalid Value> <invalid Value>}
initial_test.go:76: {4 true posts.db <invalid Value> <invalid Value>}
initial_test.go:76: {5 true users.db <invalid Value> <invalid Value>}
initial_test.go:76: {6 true .log <invalid Value> <invalid Value>}
initial_test.go:76: {7 true index.html <invalid Value> <invalid Value>}
initial_test.go:76: {8 true .version <invalid Value> <invalid Value>}
initial_test.go:76: {9 false <invalid Value> <func(*initial_test.Path) Value> <invalid Value>}
initial_test.go:76: {10 true <invalid Value> <invalid Value> <invalid Value>}
initial_test.go:76: {11 true <invalid Value> <func(*initial_test.Path) Value> <invalid Value>}
initial_test.go:76: {12 true <invalid Value> <invalid Value> <invalid Value>}
initial_test.go:76: {13 true <*errors.errorString Value> <invalid Value> <invalid Value>}
```

#### 测试初始化

```go
func TestInitial(t *testing.T) {
	result, err := initial.New[Path]()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)
}
```

```
FileBeforeInitial: default.go
FileAfterInitial:  default.go
FileBeforeInitial:
FileAfterInitial:  initial.go
FileBeforeInitial: posts.db
FileAfterInitial:  posts.db
FileBeforeInitial: key.txt
FileAfterInitial:  key.txt
FileBeforeInitial: value.txt
FileAfterInitial:  value.txt
    initial_test.go:85: &{9000 resource views public posts.db users.db .log index.html .version 0xc000170200 {d1 true 3.14 114} [{default.go} {initial.go} {posts.db}] [map[0xc00010e800:0xc00010e7f0]] some error}
```

#### 测试运行时返回错误

```go
func TestError(t *testing.T) {
	err := initial.Initial(&Path{Files: []File{{"error.log"}}})
	if err == nil {
		t.Fatal()
	}
	t.Log(err)
}
```

```
FileBeforeInitial: error.log
    initial_test.go:93: test error
```