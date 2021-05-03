package brief_test

import (
	"strings"
	"testing"

	"github.com/robbyriverside/brief"
)

func TestEncoder(t *testing.T) {
	node, err := brief.Decode(strings.NewReader(test0))
	if err != nil {
		t.Fatal(err)
	}
	out := node.Encode()
	t.Logf("%s", string(out))
}
