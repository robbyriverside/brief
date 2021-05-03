package elemental_test

import (
	"strings"
	"testing"

	"github.com/robbyriverside/brief/elemental"
)

func TestEncoder(t *testing.T) {
	node, err := elemental.Decode(strings.NewReader(test0))
	if err != nil {
		t.Fatal(err)
	}
	out := node.Encode()
	t.Logf("%s", string(out))
}
