package brief

import (
	"fmt"
	"strings"
)

// Node in a brief hierarchy
type Node struct {
	Type, Name string
	Keys       map[string]string
	Body       []*Node
	Parent     *Node
	Content    string
	Indent     int
}

// NewNode create a new Node
func NewNode(elemType string, indent int) *Node {
	return &Node{
		Type:   elemType,
		Body:   []*Node{},
		Keys:   map[string]string{},
		Indent: indent,
	}
}

func (node *Node) String() string {
	var body string
	for _, sub := range node.Body {
		body += fmt.Sprintf("\n%s", sub)
	}
	var parent string
	if node.Parent != nil {
		parent = node.Parent.Type
		if node.Parent.Name != "" {
			parent += ":" + node.Parent.Name
		}
	}
	return fmt.Sprintf("%sn(%s, %q, P(%s) %q = %q)%s", strings.Repeat(" ", node.Indent), node.Type, node.Name, parent, node.Content, node.Keys, body)
}

// Key get key value from node or return {unknown key}
func (node *Node) Key(name string) string {
	val, ok := node.Keys[name]
	if !ok {
		return "noKey"
	}
	return val
}

// Lookup a value from the above context elements
// spec can be a single name or dotted pair
// single name, returns the Name of the context
// a dotted pair returns a key value from the context {context}.{key}
func (node *Node) Lookup(spec string) string {
	return NewValueSpec(spec).Lookup(node)
}

// Slice calls Lookup on each spec and returns the slice
func (node *Node) Slice(specs ...string) []string {
	found := []string{}
	for _, spec := range specs {
		found = append(found, node.Lookup(spec))
	}
	return found
}

// Join calls Lookup on each spec and Joins them using sep
func (node *Node) Join(sep string, specs ...string) string {
	return strings.Join(node.Slice(specs...), sep)
}

// Printf calls Lookup on each spec and prints them using format
func (node *Node) Printf(format string, specs ...string) string {
	found := make([]interface{}, 0)
	for _, spec := range specs {
		found = append(found, node.Lookup(spec))
	}
	return fmt.Sprintf(format, found...)
}

// Context is a surrounding element found by node spec
// name is either a type or a type:name pair
func (node *Node) Context(name string) *Node {
	spec := NewSpec(name)
	parent := node.Parent
	for parent != nil {
		if spec.Match(parent) {
			return parent
		}
		parent = parent.Parent
	}
	return nil
}

// IndentString return a blank string width of indent.
func (node *Node) IndentString() string {
	return strings.Repeat(" ", node.Indent)
}

// NoBody true if the Node has no body
func (node *Node) NoBody() bool {
	return node.Body == nil || len(node.Body) == 0
}

// HasName true if Node has a name
func (node *Node) HasName() bool {
	return len(node.Name) > 0
}

// HasContent true if Node has content
func (node *Node) HasContent() bool {
	return len(node.Content) > 0
}

// HasKeys true if Node has keys
func (node *Node) HasKeys() bool {
	return node.Keys == nil || len(node.Keys) == 0
}

// ContentOnly true if Node only has content
func (node *Node) ContentOnly() bool {
	return node.NoBody() && !node.HasKeys() && node.HasContent()
}

// Put the value of a key
func (node *Node) Put(key, value string) {
	node.Keys[key] = value
}

// Compile adds name and content only body Nodes to the keys
func (node *Node) Compile() {
	if node.NoBody() {
		return
	}
	if node.HasName() {
		node.Put("name", node.Name)
	}
	for i, n := range node.Body {
		if n.ContentOnly() {
			if n.HasName() {
				node.Put(n.Name, n.Content)
			} else {
				node.Put(n.Type, n.Content)
			}
			node.Body = append(node.Body[:i], node.Body[i+1:]...)
		}
	}
}
