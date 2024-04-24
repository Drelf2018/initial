package initial_test

import (
	"testing"

	"github.com/Drelf2018/initial"
	"github.com/Drelf2018/initial/fullpath"
)

type File struct {
	Name string `default:"initial.go"`
}

func (f *File) BeforeDefault() error {
	println("BeforeInitial:", f.Name)
	return nil
}

func (f File) AfterDefault() error {
	println("AfterInitial: ", f.Name)
	return nil
}

type Files []File

func (f Files) BeforePathFiles(p *Path) {
	p.Files = append(f, File{}, File{p.Full.Posts})
	p.FileMap = []map[*File]*File{{{"key"}: {"value"}}}
}

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

	Datas struct {
		D1 string  `default:"t1"`
		D2 bool    `default:"true"`
		D3 float64 `default:"3.14"`
		D4 int64   `default:"114"`
	}

	Files   Files
	FileMap []map[*File]*File
}

func (*Path) BeforePathFull(p *Path) (err error) {
	p.Full, err = fullpath.New(*p)
	if err != nil {
		return
	}
	return initial.ErrBreak
}

func (p *Path) BeforeDefault() error {
	p.Files = Files{{"default.go"}}
	return nil
}

func TestPath(t *testing.T) {
	result, err := initial.New[Path]()
	if err != nil {
		t.Fatal(err)
	}
	if result.Full.Version != `resource\views\.version` {
		t.Fail()
	}
}
