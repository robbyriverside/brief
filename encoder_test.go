package brief_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/robbyriverside/brief"
)

func TestEncoder(t *testing.T) {
	t.Logf("\n%s", test0)
	nodes, err := brief.Decode(strings.NewReader(test0))
	if err != nil {
		t.Fatal(err)
	}
	before := nodes[0]
	out := before.Encode()
	t.Logf("\n%s", string(out))
	nodes, err = brief.Decode(strings.NewReader(string(out)))
	if err != nil {
		t.Fatal(err)
	}
	after := nodes[0]
	if !reflect.DeepEqual(before, after) {
		t.Fatal("before and after not the same")
	}
}
