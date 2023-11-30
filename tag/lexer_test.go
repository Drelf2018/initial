package tag_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/Drelf2018/initial/tag"
)

func TestLexer(t *testing.T) {
	p := tag.NewParser("Add,[initial.Default,Info(parent,p)],{A:B},Test")
	b, _ := json.MarshalIndent(p.Sentence, "", "  ")
	fmt.Printf("%v\n", string(b))
}
