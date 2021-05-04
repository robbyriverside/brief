package brief_test

import (
	"os"
	"strings"
	"testing"

	"github.com/robbyriverside/brief"
)

func TestXMLOut(t *testing.T) {
	node, err := brief.Decode(strings.NewReader(test1))
	if err != nil {
		t.Fatal(err)
	}
	if err = node.WriteXML(os.Stdout); err != nil {
		t.Fatal(err)
	}
}
