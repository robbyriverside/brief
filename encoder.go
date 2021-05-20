package brief

import (
	"fmt"
	"strings"
	"unicode"
)

// Encode converts a node into brief format
func (node *Node) Encode() []byte {
	var out strings.Builder
	body := node.write(&out)
	for len(body) > 0 {
		next := body[0]
		body = body[1:]
		more := next.write(&out)
		body = append(more, body...)
	}

	return []byte(out.String())
}

func isIdentRune(ch rune, i int) bool {
	return ch == '_' || unicode.IsLetter(ch) || unicode.IsDigit(ch) && i > 0
}

func isSymbol(token string) bool {
	for i, ch := range token {
		if !isIdentRune(ch, i) {
			return false
		}
	}
	return true
}

func (node *Node) write(out *strings.Builder) []*Node {
	indent := strings.Repeat(" ", node.Indent)
	out.WriteString(indent + node.Type)
	if len(node.Name) > 0 {
		out.WriteString(":" + node.Name)
	}
	for key, val := range node.Keys {
		if isSymbol(val) {
			out.WriteString(fmt.Sprintf(" %s:%s", key, val))
			continue
		}
		out.WriteString(fmt.Sprintf(" %s:%q", key, val))
	}
	if len(node.Content) > 0 {
		out.WriteString(fmt.Sprintf(" `%s`", node.Content))
	}
	out.WriteString("\n")
	return node.Body
}
