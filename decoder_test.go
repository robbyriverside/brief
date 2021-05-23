package brief_test

import (
	_ "embed"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"text/scanner"

	"github.com/robbyriverside/brief"
)

//go:embed tests/test0.brf
var test0 string

//go:embed tests/test1.brf
var test1 string

//go:embed tests/test2.brf
var test2 string

//go:embed tests/test3.brf
var test3 string

func TestDecoder1(t *testing.T) {
	t.Log(test1)
	dec := brief.NewDecoder(strings.NewReader(test1), 4, "tests")
	nodes, err := dec.Decode()
	if err != nil {
		t.Fatal(err)
	}
	if len(nodes) != 1 {
		t.Errorf("fail %d != 1", len(nodes))
	}
	if nodes[0].Type != "html" {
		t.Errorf("failed html != %s", nodes[0].Type)
	}
	if nodes[0].Name != "foo" {
		t.Errorf("failed name foo != %s", nodes[0].Name)
	}
	for i, node := range nodes {
		t.Logf("%d> \n%s", i, node)
	}
}
func TestDecoderInclude(t *testing.T) {
	// FIXME: validate results
	t.Log(test2)
	dec := brief.NewDecoder(strings.NewReader(test2), 4, "tests")
	nodes, err := dec.Decode()
	if err != nil {
		t.Fatal(err)
	}
	if len(nodes) != 1 {
		t.Errorf("fail %d != 1", len(nodes))
	}
	if nodes[0].Type != "pages" {
		t.Errorf("failed pages != %s", nodes[0].Type)
	}
	for i, node := range nodes {
		t.Logf("%d> \n%s", i, node)
	}
}
func TestDecoderMultiple(t *testing.T) {
	t.Log(test3)
	dec := brief.NewDecoder(strings.NewReader(test3), 4, "tests")
	nodes, err := dec.Decode()
	if err != nil {
		t.Fatal(err)
	}
	tests := []func(node *brief.Node) bool{
		func(node *brief.Node) bool {
			return node.Type == "html" && node.Name == "foo"
		},
		func(node *brief.Node) bool {
			return node.Type == "pages"
		},
		func(node *brief.Node) bool {
			val, ok := node.Keys["class"]
			return node.Type == "html" && ok && val == "foo"
		},
	}
	for i, node := range nodes {
		if !tests[i](node) {
			t.Errorf("%d> fail", i)
		}
		t.Logf("%d> \n%s", i, node)
	}
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
		"\"X/Y = 2\"= String indent: 12 start: false line: 8",
		"`the quick brown fox\njumped over the moon and ran into a cow`= RawString indent: 12 start: false line: 9",
	}
	t.Log(test0)
	var text brief.Scanner
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

func TestShowScanner(t *testing.T) {
	fmt.Println(test1)
	var text brief.Scanner
	text.Init(strings.NewReader(test0), 4)
	for c := text.Scan(); c != scanner.EOF; c = text.Scan() {
		token := text.TokenText()
		fmt.Printf("%s= %s indent: %d start: %t line: %d\n", token, scanner.TokenString(c), text.Indent, text.LineStart, text.Pos().Line)
	}
}

func TestScannerCases(t *testing.T) {
	tests := []struct {
		Line   string
		Indent int
		Tokens []string
	}{
		{
			Line:   "       \"-42.0\":foo",
			Indent: 7,
			Tokens: []string{"\"-42.0\"", ":", "foo"},
		},
		{
			Line:   "       foo:-42.0",
			Indent: 7,
			Tokens: []string{"foo", ":", "-", "42.0"},
		},
		{
			Line:   "       +42.0:foo",
			Indent: 7,
			Tokens: []string{"+", "42.0", ":", "foo"},
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
		var text brief.Scanner
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
			t.Errorf("%d> invalid test tokens %q != %q", i, test.Tokens, tokens)
		}
	}
}
