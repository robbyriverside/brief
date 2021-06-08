package brief

import (
	"fmt"
	"strings"
)

// Spec for node Type:Name
type Spec struct {
	Type, Name string
	NoName     bool
}

// NewSpec from spec of Type:Name or just Type
func NewSpec(spec string) *Spec {
	pos := strings.IndexRune(spec, ':')
	if pos < 0 {
		return &Spec{Type: spec, NoName: true}
	}
	return &Spec{Type: spec[:pos], Name: spec[pos+1:]}
}

func (s *Spec) String() string {
	if len(s.Name) > 0 {
		return fmt.Sprintf("%s:%s", s.Type, s.Name)
	}
	return s.Type
}

// Match spec for node
func (s *Spec) Match(node *Node) bool {
	if s.Type != node.Type {
		return false
	}
	if s.NoName {
		return true
	}
	if s.Name != node.Name {
		return false
	}

	return true
}

// Find searches for a node matching name in the body of this node
// The name is a node type or a type:name pair
func (node *Node) Find(name string) *Node {
	return NewSpec(name).Find(node)
}

// Find looks for a specific node in the body that matches spec
func (s *Spec) Find(node *Node) *Node {
	for _, sub := range node.Body {
		if s.Match(sub) {
			return sub
		}
	}
	for _, sub := range node.Body {
		if found := s.Find(sub); found != nil {
			return found
		}
	}
	return nil
}

// Child follow a path to a specific node in the body
// path elements are node type or type:name pair
func (node *Node) Child(path ...string) *Node {
	at := node
	for _, name := range path {
		spec := NewSpec(name)
		var found bool
		for _, next := range at.Body {
			if spec.Match(next) {
				found = true
				at = next
				break
			}
		}
		if !found {
			return nil
		}
	}
	return at
}

func (node *Node) FindAll(name string) []*Node {
	return NewSpec(name).FindAll(node)
}

func (s *Spec) FindAll(node *Node) []*Node {
	result := make([]*Node, 0)
	for _, sub := range node.Body {
		if s.Match(sub) {
			result = append(result, sub)
		}
	}
	for _, sub := range node.Body {
		found := s.FindAll(sub)
		result = append(result, found...)
	}
	return result
}
