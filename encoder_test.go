package brief_test

import (
	"strings"
	"testing"

	"github.com/robbyriverside/brief"
)

func TestEncoder(t *testing.T) {
	// FIXME: validate results
	nodes, err := brief.Decode(strings.NewReader(test0))
	if err != nil {
		t.Fatal(err)
	}
	out := nodes[0].Encode()
	t.Logf("%s", string(out))
}
