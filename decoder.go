package brief

import (
	"fmt"
	"io"
	"strings"
	"text/scanner"
)

// Decode creates a Node by parsing brief format from reader
func Decode(reader io.Reader) (*Node, error) {
	var root *Node
	var text Scanner
	text.Init(reader, 4)
	isFirst := true
	nesting := []*Node{}
	var isElem, addValue bool
	var key string
	for tt := text.Scan(); tt != scanner.EOF; tt = text.Scan() {
		token := text.TokenText()
		if token[0] == '/' {
			continue
		}
		if isFirst {
			if tt != scanner.Ident {
				return nil, fmt.Errorf("line %d must begin with an identifer: %q", text.Pos().Line, token)
			}
			isFirst = false
			root = NewNode(token, text.Indent)
			nesting = append(nesting, root)
			isElem = true
			continue
		}

		leaf := len(nesting) - 1
		parent := nesting[leaf]
		if text.LineStart {
			if !(tt == scanner.Ident || tt == '+') {
				return nil, fmt.Errorf("line %d must begin with an identifer or plus: %q", text.Pos().Line, token)
			}

			if text.Indent <= parent.Indent {
				var found bool
				for i := leaf - 1; i > -1; i-- {
					if text.Indent > nesting[i].Indent {
						nesting = nesting[:i+1]
						found = true
						break
					}
				}
				if !found {
					// TODO: could allow more than one tree per input
					return nil, fmt.Errorf("nesting error on line %d", text.Pos().Line)
				}
				leaf = len(nesting) - 1
				parent = nesting[leaf]
			}

			if tt == '+' {
				key = ""
				continue
			}

			node := NewNode(token, text.Indent)
			parent.Body = append(parent.Body, node)
			nesting = append(nesting, node)
			isElem = true
			key = token
			continue
		}

		if addValue {
			addValue = false
			switch tt {
			case scanner.Ident, scanner.String, scanner.Float, scanner.Int:
				if tt == scanner.String {
					token = strings.Trim(token, "\"")
				}
				if isElem {
					isElem = false
					parent.Name = token
				} else {
					parent.Put(key, token)
				}
			default:
				return nil, fmt.Errorf("invalid value %s on line %d", scanner.TokenString(tt), text.Pos().Line)
			}
			key = "" // clear the key
			continue
		}

		switch tt {
		case ':':
			addValue = true
		case scanner.Ident:
			if key != "" && !isElem {
				return nil, fmt.Errorf("key %q has no value on line %d", key, text.Pos().Line)
			}
			key = token
			isElem = false
		case scanner.RawString:
			if key != "" && !isElem {
				return nil, fmt.Errorf("key %q has no value on line %d", key, text.Pos().Line)
			}
			parent.Content += strings.Trim(token, "`")
			isElem = false
		default:
			return nil, fmt.Errorf("invalid %s token on line %d", scanner.TokenString(tt), text.Pos().Line)
		}

	}
	return root, nil
}
