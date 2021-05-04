package brief

import (
	_ "embed"
	"io"
	"text/template"

	"github.com/Masterminds/sprig"
)

//go:embed templates/xmlout.tmpl
var xmlout string

func (node *Node) WriteXML(out io.Writer) error {
	tmpl, err := template.New("xmlout").Funcs(sprig.TxtFuncMap()).Parse(xmlout)
	if err != nil {
		return err
	}
	return tmpl.Execute(out, node)
}
