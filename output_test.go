package brief_test

import (
	"os"
	"strings"
	"testing"

	"github.com/robbyriverside/brief"
)

func TestXMLOut(t *testing.T) {
	t.Logf("\n%s", test1)
	nodes, err := brief.Decode(strings.NewReader(test1))
	if err != nil {
		t.Fatal(err)
	}
	node := nodes[0]
	if err = node.WriteXML(os.Stdout); err != nil {
		t.Fatal(err)
	}
}
