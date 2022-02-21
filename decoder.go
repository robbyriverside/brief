package brief

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/scanner"
)

// DecoderState constants
type DecoderState int

// States of the decoder
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
	NegValue                // Minus sign instead of a value
	OnComment               // A comment
)

// Decoder for brief formated files
type Decoder struct {
	Err            error
	Roots, Nesting []*Node
	Text           Scanner
	ScanType       rune
	Token          string
	State          DecoderState
	Key, Feature   string
	Padding        int
	Dir            string
	Debug          bool
}

// NewDecoder from reader with tabsize and optional directory
// srcdir is used with #include files relative to this reader
func NewDecoder(reader io.Reader, tabsize int, srcdir string) *Decoder {
	dir := srcdir
	if len(dir) == 0 {
		var err error
		dir, err = os.Getwd()
		if err != nil {
			dir = "./"
		}
	}
	var decoder Decoder
	decoder.Dir = dir
	decoder.Text.Init(reader, tabsize)
	decoder.Roots = make([]*Node, 0)
	decoder.Nesting = make([]*Node, 0)
	return &decoder
}

// NewFileDecoder new decoder that reads from a filename
func NewFileDecoder(filename string) (*Decoder, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return NewDecoder(file, 4, fileDir(filename)), nil
}

func fileDir(filename string) string {
	if filepath.IsAbs(filename) {
		return filepath.Dir(filename)
	}
	abs, err := filepath.Abs(filename)
	if err != nil {
		abs = "./" + filename // because os.Getwd failed
	}
	return filepath.Dir(abs)
}

// DecodeFile into brief Nodes
func DecodeFile(filename string) ([]*Node, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	nodes, err := Decode(file, fileDir(filename))
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

// Decode creates a Node by parsing brief format from reader
func Decode(reader io.Reader, srcdir string) ([]*Node, error) {
	dec := NewDecoder(reader, 4, srcdir)
	return dec.Decode()
}

// Errorf added to decoder and returned
func (dec *Decoder) Errorf(format string, args ...interface{}) error {
	return dec.Error(fmt.Sprintf(format, args...))
}

// Error added to decoder and returned
func (dec *Decoder) Error(msg string) error {
	pos := dec.Text.Pos()
	err := fmt.Errorf("%s on %q at %d:%d", msg, dec.Token, pos.Line, pos.Column-len(dec.Token))
	if dec.Err != nil {
		err = fmt.Errorf("%s\n%s", dec.Err, err)
	}
	dec.Err = err
	return dec.Err
}

// topLevel returns true if
func (dec *Decoder) topLevel() bool {
	return len(dec.Nesting) == 0
}

// indent adds padding to indent
func (dec *Decoder) indent() int {
	return dec.Text.Indent + dec.Padding
}

// parent returns parent node, if any
// reduces nesting
func (dec *Decoder) parent() *Node {
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
	parent := dec.parent()
	if parent == nil {
		dec.Error("SetName parent not found")
		return
	}
	parent.Name = strings.Trim(dec.Token, "\"")
}

func (dec *Decoder) setValue(neg bool) {
	parent := dec.parent()
	if parent == nil {
		dec.Error("SetValue parent not found")
		return
	}
	if len(dec.Key) == 0 {
		dec.Error("SetValue no key")
	}
	if neg && dec.Token[0] != '"' {
		parent.Put(dec.Key, "-"+dec.Token)
		return
	}
	parent.Put(dec.Key, strings.Trim(dec.Token, "\""))
}

func (dec *Decoder) setContent() {
	parent := dec.parent()
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
	node := NewNode(dec.Token, dec.indent())
	parent := dec.findParent(node.Indent)
	if parent != nil {
		node.Parent = parent
		parent.Body = append(parent.Body, node)
	} else {
		dec.Roots = append(dec.Roots, node)
	}
	dec.Nesting = append(dec.Nesting, node)
}

// Decode creates a Node by parsing brief format from reader
func (dec *Decoder) Decode() ([]*Node, error) {
	dec.State = KeyEmpty
	for dec.next() {
		if dec.Err != nil {
			return nil, dec.Err
		}
		if dec.Text.LineStart {
			switch dec.State {
			case KeyElem, KeyEmpty, OnComment:
				dec.State = NewLine
			default:
				return nil, dec.Error("invalid stray token at end of line above")
			}
		}
		// if this is a feature use the feature handler
		if dec.State == FeatureSet {
			dec.handleFeature()
			dec.State = KeyEmpty
			continue
		}
		switch dec.ScanType {
		case scanner.Comment: // skip comments
			dec.State = OnComment
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
			case NegValue:
				return nil, dec.Error("invalid minus before symbol")
			case OnValue:
				dec.setValue(false)
				dec.Key = ""
				dec.State = KeyEmpty
			case OnFeature:
				dec.Feature = dec.Token
				dec.State = FeatureSet
			default:
				return nil, dec.Error("invalid identifier found")
			}
		case scanner.String, scanner.Int, scanner.Float:
			if dec.State == NegValue && dec.ScanType == scanner.String {
				return nil, dec.Error("invalid minus before string")
			}
			switch dec.State {
			case OnName:
				dec.setName()
				dec.Key = ""
				dec.State = KeyEmpty
			case OnValue, NegValue:
				dec.setValue(dec.State == NegValue)
				dec.Key = ""
				dec.State = KeyEmpty
			default:
				return nil, dec.Error("invalid value found")
			}
		case scanner.RawString:
			if dec.State == NegValue {
				return nil, dec.Error("invalid minus before content")
			}
			switch dec.State {
			case KeyElem, KeyEmpty:
				dec.Key = ""
				dec.setContent()
				dec.State = KeyEmpty
			default:
				return nil, dec.Error("invalid content found")
			}
		case '-':
			switch dec.State {
			case OnValue, OnName:
				dec.State = NegValue
			default:
				return nil, dec.Error("invalid minus")
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
			case KeyElem, KeyEmpty:
				dec.Key = ""
				dec.readBlock()
				dec.State = KeyEmpty
			case NewLine:
				dec.State = OnFeature
			default:
				return nil, dec.Error("invalid syntax '#'")
			}
		}
	}
	return dec.Roots, nil
}

func (dec *Decoder) readBlock() error {
	delim := dec.Text.Next()
	if !strings.ContainsAny(string(delim), "|@$%") {
		return dec.Error("invalid block delimiter: #" + string(delim))
	}
	var build strings.Builder
	for ch := dec.Text.Next(); ch != scanner.EOF; ch = dec.Text.Next() {
		if ch == delim {
			at := dec.Text.Next()
			if at == '#' {
				dec.Token = build.String()
				dec.setContent()
				return nil
			}
			build.WriteRune(ch)
			build.WriteRune(at)
			continue
		}
		build.WriteRune(ch)
	}
	return dec.Error("Found EOF while reading block no matching " + string(delim))
}
