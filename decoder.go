package brief

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/scanner"
)

type DecoderState int

const (
	Unknown    DecoderState = iota
	NewLine                 // LineStart
	KeyElem                 // Key set to elem
	KeyValue                // Key set to Key
	KeyEmpty                // ready for next key or content  key is empty
	OnName                  // Set Name
	OnValue                 // Put key-value
	OnFeature               // Exec Feature
	FeatureSet              // Feature value is set
)

type Decoder struct {
	Err            error
	Roots, Nesting []*Node
	Text           Scanner
	ScanType       rune
	Token          string
	State          DecoderState
	Key, Feature   string
	Padding        int
}

func NewDecoder(reader io.Reader, tabsize int) *Decoder {
	var decoder Decoder
	decoder.Text.Init(reader, tabsize)
	decoder.Roots = make([]*Node, 0)
	decoder.Nesting = make([]*Node, 0)
	return &decoder
}

func (dec *Decoder) TopLevel() bool {
	return len(dec.Nesting) == 0
}

func (dec *Decoder) Indent() int {
	return dec.Text.Indent + dec.Padding
}

func (dec *Decoder) Parent() *Node {
	size := len(dec.Nesting)
	if size > 0 {
		return dec.Nesting[size-1]
	}
	return nil
}

func (dec *Decoder) next() bool {
	dec.ScanType = dec.Text.Scan()
	dec.Token = dec.Text.TokenText()
	return dec.ScanType != scanner.EOF
}

func (dec *Decoder) setName() {
	parent := dec.Parent()
	if parent == nil {
		dec.Error("SetName parent not found")
		return
	}
	parent.Name = strings.Trim(dec.Token, "\"")
}

func (dec *Decoder) Error(msg string) error {
	pos := dec.Text.Pos()
	dec.Err = fmt.Errorf("%s at %q on pos %d:%d", msg, dec.Token, pos.Line, pos.Offset)
	return dec.Err
}

func (dec *Decoder) setValue() {
	parent := dec.Parent()
	if parent == nil {
		dec.Error("SetValue parent not found")
		return
	}
	if len(dec.Key) == 0 {
		dec.Error("SetValue no key")
	}
	parent.Put(dec.Key, strings.Trim(dec.Token, "\""))
}

func (dec *Decoder) setContent() {
	parent := dec.Parent()
	if parent == nil {
		dec.Error("SetContent parent not found")
		return
	}
	parent.Content = strings.Trim(dec.Token, "`")
}

func (dec *Decoder) findParent(indent int) *Node {
	for size := len(dec.Nesting); size > 0; size = len(dec.Nesting) {
		last := size - 1
		parent := dec.Nesting[last]
		if indent > parent.Indent {
			return parent
		}
		dec.Nesting = dec.Nesting[:last]
	}
	return nil
}

func (dec *Decoder) addNode() {
	node := NewNode(dec.Token, dec.Indent())
	parent := dec.findParent(node.Indent)
	if parent != nil {
		parent.Body = append(parent.Body, node)
	} else {
		dec.Roots = append(dec.Roots, node)
	}
	dec.Nesting = append(dec.Nesting, node)
}

// Decode creates a Node by parsing brief format from reader
func Decode(reader io.Reader) ([]*Node, error) {
	dec := NewDecoder(reader, 4)
	return dec.Decode()
}

// Decode creates a Node by parsing brief format from reader
func (dec *Decoder) Decode() ([]*Node, error) {
	for dec.next() {
		if dec.Err != nil {
			return nil, dec.Err
		}
		if dec.Text.LineStart {
			dec.State = NewLine
		}
		switch dec.ScanType {
		case scanner.Comment: // skip comments
		case scanner.Ident:
			switch dec.State {
			case NewLine:
				dec.addNode()
				dec.Key = dec.Token
				dec.State = KeyElem
			case KeyElem: // no colon after elem
				dec.Key = dec.Token
				dec.State = KeyValue
			case KeyEmpty:
				dec.Key = dec.Token
				dec.State = KeyValue
			case OnName:
				dec.setName()
				dec.Key = ""
				dec.State = KeyEmpty
			case OnValue:
				dec.setValue()
				dec.Key = ""
				dec.State = KeyEmpty
			case OnFeature:
				dec.Feature = dec.Token
				dec.State = FeatureSet
			default:
				return nil, dec.Error("invalid identifier found")
			}
		case scanner.String, scanner.Int, scanner.Float:
			switch dec.State {
			case OnName:
				dec.setName()
				dec.Key = ""
				dec.State = KeyEmpty
			case OnValue:
				dec.setValue()
				dec.Key = ""
				dec.State = KeyEmpty
			default:
				return nil, dec.Error("invalid value found")
			}
		case scanner.RawString:
			switch dec.State {
			case KeyElem, KeyEmpty:
				dec.Key = ""
				dec.setContent()
				dec.State = KeyEmpty
			case FeatureSet:
				dec.contentFeature()
			default:
				return nil, dec.Error("invalid content found")
			}
		case ':':
			switch dec.State {
			case KeyElem:
				dec.State = OnName
			case KeyValue:
				dec.State = OnValue
			default:
				return nil, dec.Error("invalid syntax ':'")
			}
		case '+':
			switch dec.State {
			case NewLine:
				dec.State = KeyEmpty
			default:
				return nil, dec.Error("invalid syntax '+'")
			}
		case '#':
			switch dec.State {
			case NewLine:
				dec.State = OnFeature
			default:
				return nil, dec.Error("invalid syntax '#'")
			}
		}
	}
	return dec.Roots, nil
}

func (dec *Decoder) contentFeature() {
	content := strings.Trim(dec.Token, "`")
	feature := strings.ToLower(dec.Feature)
	if len(feature) == 0 {
		dec.Err = dec.Error("empty feature")
	}
	if len(content) == 0 {
		dec.Err = dec.Error("empty feature content")
	}
	switch feature {
	case "include":
		dec.Err = dec.includeFile(content)
	}
}

func (dec *Decoder) includeFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	idec := NewDecoder(file, dec.Text.TabCount)
	idec.Padding = dec.Indent()
	nodes, err := idec.Decode()
	if err != nil {
		return err
	}
	size := len(nodes)
	if size == 0 {
		return nil
	}
	parent := dec.findParent(nodes[size-1].Indent)
	if parent != nil {
		parent.Body = append(parent.Body, nodes...)
		return nil
	}
	dec.Roots = append(dec.Roots, nodes...)
	return nil
}
