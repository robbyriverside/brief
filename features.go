package brief

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/scanner"
)

func (dec *Decoder) handleFeature() {
	dec.Feature = strings.ToLower(dec.Feature)
	switch dec.Feature {
	case "include":
		switch dec.ScanType {
		case scanner.String, scanner.RawString:
			dec.trimContentToken()
			dec.includeFile(dec.Token)
		}
	default:
		dec.Errorf("unknown brief feature %s", dec.Feature)
	}
}

func (dec *Decoder) trimContentToken() {
	switch dec.Token[0] {
	case '"':
		dec.Token = strings.Trim(dec.Token, "\"")
	case '`':
		dec.Token = strings.Trim(dec.Token, "`")
	}
}

func (dec *Decoder) includeFile(filename string) {
	if !filepath.IsAbs(filename) {
		filename = filepath.Join(dec.Dir, filename)
	}
	if dec.Debug {
		fmt.Println("*** include", filename)
	}
	file, err := os.Open(filename)
	if err != nil {
		dec.Error(err.Error())
		return
	}
	defer file.Close()
	dir := filepath.Dir(filename)
	idec := NewDecoder(file, dec.Text.TabCount, dir)
	idec.Debug = dec.Debug
	idec.Padding = dec.indent()
	nodes, err := idec.Decode()
	if err != nil {
		dec.Error(err.Error())
		return
	}
	size := len(nodes)
	if size == 0 {
		if dec.Debug {
			fmt.Println("*** include file was empty")
		}
		return
	}
	parent := dec.findParent(nodes[size-1].Indent)
	if parent != nil {
		parent.Body = append(parent.Body, nodes...)
		return
	}
	dec.Roots = append(dec.Roots, nodes...)
}
