package elemental_test

import (
	_ "embed"
	"fmt"
	"reflect"
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

func TestScanner(t *testing.T) {
	tests := []string{
		"html= Ident indent: 0 start: true line: 1",
		"head= Ident indent: 4 start: true line: 2",
		"title= Ident indent: 8 start: true line: 3",
		"`My Web Page`= RawString indent: 8 start: false line: 3",
		"body= Ident indent: 4 start: true line: 4",
		"class= Ident indent: 4 start: false line: 4",
		":= \":\" indent: 4 start: false line: 4",
		"mybody= Ident indent: 4 start: false line: 4",
		"h1= Ident indent: 8 start: true line: 5",
		"`My Web Page`= RawString indent: 8 start: false line: 5",
		"div= Ident indent: 8 start: true line: 7",
		":= \":\" indent: 8 start: false line: 7",
		"main= Ident indent: 8 start: false line: 7",
		"class= Ident indent: 8 start: false line: 7",
		":= \":\" indent: 8 start: false line: 7",
		"myblock= Ident indent: 8 start: false line: 7",
		"p= Ident indent: 12 start: true line: 8",
		"id= Ident indent: 12 start: false line: 8",
		":= \":\" indent: 12 start: false line: 8",
		"\"X:Y = 2\"= String indent: 12 start: false line: 8",
		"`the quick brown fox\njumped over the moon and ran into a cow`= RawString indent: 12 start: false line: 9",
	}
	t.Log(test0)
	var text elemental.Scanner
	text.Init(strings.NewReader(test0), 4)
	i := 0
	for c := text.Scan(); c != scanner.EOF; c = text.Scan() {
		token := text.TokenText()
		scan := fmt.Sprintf("%s= %s indent: %d start: %t line: %d", token, scanner.TokenString(c), text.Indent, text.LineStart, text.Pos().Line)
		if scan != tests[i] {
			t.Errorf("%d> compare failed: %q != \n           %q", i, scan, tests[i])
		}
		i++
	}
}

func TestScannerCases(t *testing.T) {
	tests := []struct {
		Line   string
		Indent int
		Tokens []string
	}{
		{
			Line:   "       elem:foo",
			Indent: 7,
			Tokens: []string{"elem", ":", "foo"},
		},
		{
			Line:   "   \nelem `some text`",
			Indent: 0,
			Tokens: []string{"elem", "`some text`"},
		},
		{
			Line: "			elem id:\"some text\"", // using tabs
			Indent: 12,
			Tokens: []string{"elem", "id", ":", "\"some text\""},
		},
	}

	for i, test := range tests {
		var text elemental.Scanner
		text.Init(strings.NewReader(test.Line), 4)
		tokens := []string{}
		for c := text.Scan(); c != scanner.EOF; c = text.Scan() {
			token := text.TokenText()
			tokens = append(tokens, token)
			if text.LineStart {
				if test.Indent != text.Indent {
					t.Errorf("%d> invalid test indent %d != %d", i, test.Indent, text.Indent)
				}
			}
		}
		if !reflect.DeepEqual(test.Tokens, tokens) {
			t.Errorf("%d> invalid test tokens %s != %s", i, test.Tokens, tokens)
		}
	}
}
