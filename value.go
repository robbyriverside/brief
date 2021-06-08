package brief

import "strings"

// ValueSpec states
const (
	NoKey  = "noKey"
	NoVal  = "noVal"
	NoCTX  = "noCTX"
	NoName = "noName"
)

// ValueSpec <elem>.<key> or <elem>.Name
type ValueSpec struct {
	Elem, Name string
	HasKey     bool
}

// NewValueSpec spec for a value Name or Key-value
// <elem>.<key> or <elem>.Name
func NewValueSpec(spec string) *ValueSpec {
	elem, name, hasKey := ParseValueSpec(spec)
	return &ValueSpec{
		Elem:   elem,
		Name:   name,
		HasKey: hasKey,
	}
}

// Value from node if matching elem and has key
func (val *ValueSpec) Value(node *Node) (string, bool) {
	if val.Elem != node.Type {
		return "", false
	}
	if val.HasKey {
		val, ok := node.Keys[val.Name]
		return val, ok
	}
	if len(node.Name) == 0 {
		return "", false
	}
	return node.Name, true
}

// ParseValueSpec type.key or type.Name
// return type, value and hasKey
func ParseValueSpec(spec string) (string, string, bool) {
	values := strings.Split(spec, ".")
	hasKey := len(values) > 1
	elem := spec
	name := NoKey
	if hasKey {
		elem = values[0]
		name = values[1]
	}
	return elem, name, hasKey
}

// Lookup value spec in parents
func (val *ValueSpec) Lookup(node *Node) string {
	ctx := node.Context(val.Elem)
	if ctx == nil {
		return NoCTX
	}
	if val.HasKey {
		kval, ok := ctx.Keys[val.Name]
		if !ok {
			return NoKey
		}
		if len(kval) == 0 {
			return NoVal
		}
	}
	if len(ctx.Name) == 0 {
		return NoName
	}
	return ctx.Name
}

// Collect value specs up the hierarchy into a Slice
func (node *Node) Collect(names ...string) []string {
	at := node
	specs := make([]*ValueSpec, len(names))
	for i, name := range names {
		specs[i] = NewValueSpec(name)
	}
	result := make([]string, 0)
	for {
		if at == nil {
			break
		}
		for i, spec := range specs {
			val, ok := spec.Value(at)
			if !ok {
				continue
			}
			result = append(result, val)
			specs = specs[i:]
			break
		}
		at = at.Parent
	}
	return result
}
