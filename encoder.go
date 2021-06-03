package brief

import (
	"fmt"
	"strings"
	"text/scanner"
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

// NoQuote tests if the value is an identifier or number
func NoQuote(value string) bool {
	var s scanner.Scanner
	s.Init(strings.NewReader(value))
	tok := s.Scan()
	var minus bool
	for {
		switch tok {
		case scanner.Ident:
			return !minus && s.TokenText() == value
		case scanner.Float, scanner.Int:
			numval := s.TokenText()
			if minus {
				numval = "-" + numval
			}
			return numval == value
		case '-':
			if minus {
				return false
			}
			minus = true
			tok = s.Scan()
		default:
			return false
		}
	}

}

func (node *Node) write(out *strings.Builder) []*Node {
	indent := strings.Repeat(" ", node.Indent)
	out.WriteString(indent + node.Type)
	if len(node.Name) > 0 {
		if NoQuote(node.Name) {
			out.WriteString(":" + node.Name)
		} else {
			out.WriteString(fmt.Sprintf(":%q", node.Name))
		}
	}
	for key, val := range node.Keys {
		if NoQuote(val) {
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
