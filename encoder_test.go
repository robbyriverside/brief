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

func TestNoQuote(t *testing.T) {
	tests := []struct {
		value string
		valid bool
	}{
		{
			value: "one",
			valid: true,
		},
		{
			value: "-one",
			valid: false,
		},
		{
			value: "-",
			valid: false,
		},
		{
			value: "24",
			valid: true,
		},
		{
			value: "33.0e4",
			valid: true,
		},
		{
			value: "-33.0e4",
			valid: true,
		},
		{
			value: "-22",
			valid: true,
		},
		{
			value: "- 22",
			valid: false,
		},
		{
			value: "24a",
			valid: false,
		},
		{
			value: "33.0e-4",
			valid: true,
		},
		{
			value: "x-22",
			valid: false,
		},
		{
			value: "- - 33",
			valid: false,
		},
		{
			value: "go-flags",
			valid: false,
		},
		{
			value: "go flags",
			valid: false,
		},
	}
	for i, test := range tests {
		if brief.NoQuote(test.value) != test.valid {
			t.Errorf("%d> failed %s", i, test.value)
		}
	}
}
