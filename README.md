# Brief Specification Format

In short, the brief format is XML with minimal syntax and indented blocks like Python.

This repo contains readers and translators for the brief format written in Go.

## Why?

The unique feature of XML having both keyword parameters and nested content body is useful when writing a specification.  However, writing XML by hand can be tedious because of all the syntax.  The primary design goal of brief to have the same structure as XML but easy to write.

Brief is not intended to be a data interchange format.  However, it can be easily converted to XML.  

Brief is the primary input format for the Brevity App Meta-Generator.

## Brief Library

### Brief Encoder

Parses the brief format and creates a Node object.

```go
type Node struct {
	Type, Name string
	Keys       map[string]string
	Body       []*Node
	Content    string
	Indent     int
}
```

This is more efficient than using reflection to map to an arbitrary structure.  

```go
var in io.Reader
rootNode, err := brief.Decode(in)
```

The Node structure can be walked to convert to a custom struct, if required.

### Brief Decoder

Writes the Node object in brief format.

```go
var node Node
var out io.Writer
err := node.Encode(out)
```

### Brief XML Output

Writes the Node object in XML format.

```go
var node Node
var out io.Writer
err := node.WriteXML(out)
```

XML output uses a go text template.  This serves as an example of using brief with a template.

Contents of templates/xmlout.tmpl:

```template
{{define "node" -}}
{{.IndentString}}<{{.Type}} {{- if .Name}} name="{{.Name}}"{{end}} {{- range $key, $val := .Keys}} {{$key}}="{{$val}}"{{end}}>
{{- .Content}}
{{.IndentString}}{{range $node := .Body}}{{ template "node" $node }}{{end -}}
{{.IndentString}}</{{.Type}}>
{{end}}
{{- template "node" . }}
```

## Brief Format

With one exception, the + operator (see below), the first token on each line is the element type.  After the element type, is a series of key-value pairs, optionally followed by a text body.  Child elements are indented on the lines below the parent element.

### Example

No better example of XML format than the widely known HTML dialect.  HTML5 has some variations, but we will skip over them for our purposes.

An HTML page contains a single top-level structure and two sub-structures:

```brief
html
    head
    body
```

Sub-structures are indented to identify a sub-block.  The first identifier on each line is an element name or type.  The back-tic contains multiline text which forms the contents of an element.  The back-tic must be after any key-value pairs.

```brief
html
    head
        title `My Web Page`
    body class:mybody
        h1 `My Web Page`

        div id:main class:myblock
            p id:X `the quick brown fox
jumped over the moon and ran into a cow`
```

Here is the equivalent in HTML.

```html
<html>
<head>
    <title>My Web Page</title>
</head>
<body class="mybody">
    <h1>My Web Page</h1>

    <div id="main" class="myblock">
        <p id="X">the quick brown fox
jumped over the moon and ran into a cow</p>
    </div>
</body>
</html>
```

### Key Values

For key-value pairs with a value that is more than just a simple token double quotes are used.  

```brief
elem key:"value of key"  <->   <elem key="value of key"/>

elem key:"my brother \"Bill\""  <->  <elem key="my brother \"Bill\"">
```

Simple tokens cannot contain brief syntactic characters:  space, colon, back-tic, double-quote.  This allows number formats to be simple tokens.

```brief
elem size:33 max:1.4e3
```

### Name Key

Because specification elements often have a "name" keyword to identify them in the document, we give them a special place.  The element type can be a keyword by adding a colon (:) to the end and so it can also have a value, which is the name.

```brief
body
    div:foo  <->  <div name="foo"/>
```

The purpose of this shorthand is to improve readability and standardize on elements having names.  Names are important in a written specification which is the primary purpose of the brief format.

### Multi-line and Content

Line cannot start with a back-tic.  Body text requires an element.  Content back-tic is last feature on a line.
A line that starts with plus '+' is a continuation of the attributes on the line above.  May be followed by a space.

```brief
elem:foo bar:zed
    + more:true range:"3 to 5" `
  more content here`
```
