package brief

import (
	_ "embed"
	"io"
	"text/template"
)

//go:embed templates/xmlout.tmpl
var xmlout string

func (node *Node) WriteXML(out io.Writer) error {
	tmpl, err := template.New("xmlout").Parse(xmlout)
	if err != nil {
		return err
	}
	return tmpl.Execute(out, node)
}
