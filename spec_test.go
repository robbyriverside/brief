package brief_test

import (
	"strings"
	"testing"

	"github.com/robbyriverside/brief"
)

func TestChild(t *testing.T) {
	tests := []struct {
		Path  []string
		Found bool
	}{
		{
			Path:  []string{"head", "title"},
			Found: true,
		},
		{
			Path:  []string{"body", "div:main"},
			Found: true,
		},
		{
			Path:  []string{"body", "div"}, // div name is ignored
			Found: true,
		},
		{
			Path:  []string{"body", "div:main", "p"},
			Found: true,
		},
		{
			Path:  []string{"body", "div:"}, // div does have a name
			Found: false,
		},
		{
			Path:  []string{"body", "div", "p:foo"}, // p does NOT have a name
			Found: false,
		},
	}
	t.Log(test1)
	nodes, err := brief.Decode(strings.NewReader(test1))
	if err != nil {
		t.Fatal(err)
	}
	node := nodes[0]
	for _, test := range tests {
		found := node.Child(test.Path...)
		if found == nil && test.Found {
			t.Error("failed:", test.Found, test.Path)
		}
		if found != nil && !test.Found {
			t.Error("failed:", test.Found, test.Path)
		}
		t.Logf("%s %t --> %s", test.Path, test.Found, found)
	}
}

func TestFind(t *testing.T) {
	tests := []struct {
		Name  string
		Found bool
	}{
		{
			Name:  "head",
			Found: true,
		},
		{
			Name:  "title",
			Found: true,
		},
		{
			Name:  "p",
			Found: true,
		},
		{
			Name:  "foo", // not found
			Found: false,
		},
		{
			Name:  "div:", // div has a name
			Found: false,
		},
		{
			Name:  "p:foo", // p has no name
			Found: false,
		},
	}
	t.Log(test1)
	nodes, err := brief.Decode(strings.NewReader(test1))
	if err != nil {
		t.Fatal(err)
	}
	node := nodes[0]
	for _, test := range tests {
		found := node.Find(test.Name)
		if found == nil && test.Found {
			t.Error("failed:", test.Found, test.Name)
		}
		if found != nil && !test.Found {
			t.Error("failed:", test.Found, test.Name)
		}
		t.Logf("%s %t --> %s", test.Name, test.Found, found)
	}
}
