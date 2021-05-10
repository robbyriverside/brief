package brief_test

import (
	"strings"
	"testing"

	"github.com/robbyriverside/brief"
)

func TestFindSpec(t *testing.T) {
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
			Path:  []string{"body", "div"},
			Found: true,
		},
		{
			Path:  []string{"body", "div", "p"},
			Found: true,
		},
	}
	t.Log(test1)
	nodes, err := brief.Decode(strings.NewReader(test1))
	if err != nil {
		t.Fatal(err)
	}
	node := nodes[0]
	for _, test := range tests {
		found := node.FindNode(test.Path...)
		if found == nil && test.Found {
			t.Error("failed:", test.Found, test.Path)
		}
		if found != nil && !test.Found {
			t.Error("failed:", test.Found, test.Path)
		}
		t.Logf("%s --> %s", test.Path, found)
	}
}

func TestGetSpec(t *testing.T) {
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
			Path:  []string{"body", "div"},
			Found: false,
		},
		{
			Path:  []string{"body", "div:main", "p"},
			Found: true,
		},
	}
	t.Log(test1)
	nodes, err := brief.Decode(strings.NewReader(test1))
	if err != nil {
		t.Fatal(err)
	}
	node := nodes[0]
	for _, test := range tests {
		found := node.GetNode(test.Path...)
		if found == nil && test.Found {
			t.Error("failed:", test.Found, test.Path)
		}
		if found != nil && !test.Found {
			t.Error("failed:", test.Found, test.Path)
		}
		t.Logf("%s --> %s", test.Path, found)
	}
}
