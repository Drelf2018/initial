package fullpath_test

import (
	"testing"

	"github.com/Drelf2018/initial/fullpath"
)

type Path struct {
	Root    string
	Views   string `join:"Root"`
	Public  string `join:"Root"`
	Posts   string `join:"Public"`
	Users   string `join:"Root"`
	Log     string `join:"Root"`
	Index   string `join:"Views"`
	Version string `join:"Views"`
}

func TestFullpath(t *testing.T) {
	path := Path{
		Root:    "resource",
		Views:   "views",
		Public:  "public",
		Posts:   "posts.db",
		Users:   "users.db",
		Log:     ".log",
		Index:   "index.html",
		Version: ".version",
	}

	newPath, err := fullpath.New(path)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("\n%v\n%v", path, newPath)

	err = fullpath.Join(&path)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("\n", path)
}
