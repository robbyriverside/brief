package brief

import (
	"fmt"
	"strings"
)

func (node *Node) Encode() []byte {
	var out strings.Builder

	node.writeLine(&out)

	return []byte(out.String())
}

func (node *Node) writeLine(out *strings.Builder) {
	indent := strings.Repeat(" ", node.Indent)
	out.WriteString(indent + node.Type)
	if len(node.Name) > 0 {
		out.WriteString(":" + node.Name)
	}
	for key, val := range node.Keys {
		if strings.ContainsAny(val, ": \"") {
			out.WriteString(fmt.Sprintf(" %s:%q", key, val))
			continue
		}
		out.WriteString(fmt.Sprintf(" %s:%s", key, val))
	}
	if len(node.Content) > 0 {
		out.WriteString(fmt.Sprintf(" `%s`", node.Content))
	}
	out.WriteString("\n")

	for _, sub := range node.Body {
		sub.writeLine(out)
	}
}
