package elemental_test

import (
	_ "embed"
	"strings"
	"testing"
	"text/scanner"

	"github.com/robbyriverside/brief/elemental"
)

//go:embed tests/test0.brf
var test0 string

func TestDecoder(t *testing.T) {
	t.Log(test0)
	node, err := elemental.Decode(strings.NewReader(test0))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s", node)
}
func TestDecodeAll(t *testing.T) {
	t.Log(test0)
	var text elemental.Scanner
	text.Init(strings.NewReader(test0), 4)
	for c := text.Scan(); c != scanner.EOF; c = text.Scan() {
		token := text.TokenText()
		t.Logf("%q: %s indent: %d start: %t line: %d\n", token, scanner.TokenString(c), text.Indent, text.LineStart, text.Pos().Line)
	}
}
