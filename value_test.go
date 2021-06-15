package brief_test

import (
	"strings"
	"testing"

	"github.com/robbyriverside/brief"
)

func TestValueSpec(t *testing.T) {
	tests := []struct {
		Spec       string
		Elem, Name string
		HasKey     bool
	}{
		{
			Spec:   "foo.bar",
			Elem:   "foo",
			Name:   "bar",
			HasKey: true,
		},
		{
			Spec:   "foo",
			Elem:   "foo",
			Name:   brief.NoKey,
			HasKey: false,
		},
		{
			Spec:   "foo:bar",
			Elem:   "foo:bar",
			Name:   brief.NoKey,
			HasKey: false,
		},
	}

	for i, test := range tests {
		spec := brief.NewValueSpec(test.Spec)
		if spec.Elem != test.Elem || spec.Name != test.Name || spec.HasKey != test.HasKey {
			t.Errorf("%d> failed Elem: %s != %s, Name: %s != %s, HasKey: %t != %t, ", i, spec.Elem, test.Elem, spec.Name, test.Name, spec.HasKey, test.HasKey)
		}
	}
}

func TestLookupValue(t *testing.T) {
	tests := []struct {
		Spec  string
		Value string
	}{
		{
			Spec:  "div",
			Value: "main",
		},
		{
			Spec:  "div.class",
			Value: "myblock",
		},
		{
			Spec:  "body.class",
			Value: "mybody",
		},
	}
	t.Logf("test0:\n%s\n", test0)
	nodes, err := brief.Decode(strings.NewReader(test0))
	if err != nil {
		t.Fatal(err)
	}
	pnode := nodes[0].Find("p")
	t.Logf("p: \n%s", pnode)
	for i, test := range tests {
		res := pnode.Lookup(test.Spec)
		if res != test.Value {
			t.Errorf("%d> failed %s != %s", i, test.Value, res)
		}
	}
}
