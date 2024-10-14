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

func TestParse(t *testing.T) {
	values := initial.ParseValues(reflect.TypeOf(Path{}))
	for _, val := range values {
		t.Log(val)
	}
}

func TestInitial(t *testing.T) {
	result, err := initial.New[Path]()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)
}

func TestError(t *testing.T) {
	err := initial.Initial(&Path{Files: []File{{"error.log"}}})
	if err == nil {
		t.Fatal()
	}
	t.Log(err)
}
