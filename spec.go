package brief

import (
	"fmt"
	"strings"
)

// Spec for node Type:Name
type Spec struct {
	Type, Name string
}

// NewSpec from spec of Type:Name or just Type
func NewSpec(spec string) *Spec {
	pos := strings.IndexRune(spec, ':')
	if pos < 0 {
		return &Spec{Type: spec}
	}
	return &Spec{Type: spec[:pos], Name: spec[pos+1:]}
}

func (s *Spec) String() string {
	if len(s.Name) > 0 {
		return fmt.Sprintf("%s:%s", s.Type, s.Name)
	}
	return s.Type
}

// Match spec to node
// If spec.Name is empty then match Type only
func (s *Spec) Match(node *Node) bool {
	if s.Type != node.Type {
		return false
	}
	if len(s.Name) > 0 && s.Name != node.Name {
		return false
	}
	return true
}

// Same spec for node
func (s *Spec) Same(node *Node) bool {
	if s.Type != node.Type {
		return false
	}
	if s.Name != node.Name {
		return false
	}
	return true
}

// FindNode a subnode by spec
// spec values can be Type:Name  or just Type
// if not found returns nil
func (node *Node) FindNode(path ...string) *Node {
	result := node
	for _, spec := range path {
		s := NewSpec(spec)
		var found bool
		for _, n := range result.Body {
			if s.Match(n) {
				found = true
				result = n
				break
			}
		}
		if !found {
			return nil
		}
	}
	return result
}

// GetNode a subnode by spec
// spec values can be Type:Name  or just Type
// spec must match exactly
// if not found returns nil
func (node *Node) GetNode(path ...string) *Node {
	result := node
	for _, spec := range path {
		s := NewSpec(spec)
		var found bool
		for _, n := range result.Body {
			if s.Same(n) {
				found = true
				result = n
				break
			}
		}
		if !found {
			return nil
		}
	}
	return result
}
