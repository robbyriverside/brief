package brief

// M map of key pairs
type M map[string]string

var indent = 4

// SetIndent for node creation
func SetIndent(val int) {
	indent = val
}

// S create string body node
func S(elem string, keys M, content string) *Node {
	return &Node{
		Type:    elem,
		Content: content,
		Keys:    keys,
		Indent:  indent,
	}
}

// B create branch body node
func B(elem string, keys M, body ...*Node) *Node {
	return &Node{
		Type:   elem,
		Keys:   keys,
		Body:   body,
		Indent: indent,
	}
}

// Add a child to a node
func (node *Node) Add(child *Node) *Node {
	node.Body = append(node.Body, child)
	return node
}
