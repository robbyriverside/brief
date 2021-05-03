package elemental

import (
	"fmt"
	"strings"
)

type Node struct {
	Type, Name string
	Keys       map[string]string
	Body       []*Node
	Content    string
	Indent     int
}

func NewNode(elemType string) *Node {
	return &Node{
		Type: elemType,
		Body: []*Node{},
		Keys: map[string]string{},
	}
}

func (node *Node) String() string {
	var body string
	for _, sub := range node.Body {
		body += fmt.Sprintf("\n%s", sub)
	}
	return fmt.Sprintf("%sn(%s, %q, %q = %q)%s", strings.Repeat(" ", node.Indent), node.Type, node.Name, node.Content, node.Keys, body)
}

func (node *Node) NoBody() bool {
	return node.Body == nil || len(node.Body) == 0
}

func (node *Node) HasName() bool {
	return len(node.Name) > 0
}

func (node *Node) HasContent() bool {
	return len(node.Content) > 0
}

func (node *Node) HasKeys() bool {
	return node.Keys == nil || len(node.Keys) == 0
}

func (node *Node) IsBodyKey() bool {
	return node.NoBody() && !node.HasKeys() && node.HasContent()
}

func (node *Node) Get(key string) (string, bool) {
	val, ok := node.Keys[key]
	return val, ok
}

func (node *Node) Put(key, value string) {
	node.Keys[key] = value
}

func (node *Node) Compile() {
	if node.NoBody() {
		return
	}
	if node.HasName() {
		node.Put("name", node.Name)
	}
	for i, n := range node.Body {
		if n.IsBodyKey() {
			if n.HasName() {
				node.Put(n.Name, n.Content)
			} else {
				node.Put(n.Type, n.Content)
			}
			node.Body = append(node.Body[:i], node.Body[i+1:]...)
		}
	}
}
